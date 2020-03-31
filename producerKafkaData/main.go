package main

import (
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	brokerList      = flag.String("brokers", "localhost:9092", "The comma separated list of brokers in the Kafka cluster")
	topic           = flag.String("topic", "testData", "The topic to produce messages to")
	messageBodySize = flag.Int("message-body-size", 100, "The size of the message payload")
	waitForAll      = flag.Bool("wait-for-all", false, "Whether to wait for all ISR to Ack the message")
	sleep           = flag.Float64("sleep", 1, "The number of seconds to sleep between messages")
	verbose         = flag.Bool("verbose", false, "Whether to enable Sarama logging")
	statFrequency   = flag.Int("statFrequency", 1000, "How frequently (in messages) to print throughput and latency")
)

var testData string = `{
"SYS_TRA_NO":"990771",
"dataSource":"CPUS",
"TERM_ID":"",
"ISS_RESP_CD":"00",
"MCHNT_CD":"842584073990001",
"RESV_FLD1_2":"1",
"FWD_SYS_ID":"D",
"RESV_FLD1_1":"0",
"MCHNT_TP":"7399",
"PRI_ACCT_NO_CONV":"196228481938229227475",
"CROSS_DIST_IN":"0",
"FWD_LINE_NO":"0000041274",
"TRANS_ST":"10000",
"STI_IN":"0",
"CARD_MEDIA":"2",
"CARD_CLASS":"01",
"SETTLE_DT":"20200325",
"TRANS_ID":"W20",
"RCV_INS_ID_CD":"01039200",
"ACQ_INS_ID_CD":"48429202",
"ISS_INS_ID_CD":"01030000",
"POS_ENTRY_MD_CD":"012",
"ACPT_RESP_CD":"00",
"TRANS_AT":"00000000000000000000000500000",
"TRANS_FIN_TS":"2020-03-25 09:18:49.594581",
"DB_ID":"1",
"SYS_ID":"D",
"toTs":"2020-03-25 09:18:49.535638",
"TRANS_RCV_TS":"2020-03-25 09:18:49.535638",
"MSG_TP":"9900",
"TFR_DT_TM":"0325113725",
"TRANS_CHNL":"07",
"CARD_ATTR":"01",
"TRANS_ID_CONV":"W20",
"CARD_BIN":"19622848",
"CARD_BRAND":"12",
"dataType":"DGF11",
"RCV_PROC_IN":"1",
"APP_HOST_ID":"14",
"FWD_INS_ID_CD":"48429202",
"FWD_PROC_IN":"1"
}`

var num int = 1675

type MessageMetadata struct {
	EnqueuedAt time.Time
}

func (mm *MessageMetadata) Latency() time.Duration {
	return time.Since(mm.EnqueuedAt)
}

func producerConfiguration() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	if *waitForAll {
		config.Producer.RequiredAcks = sarama.WaitForAll
	} else {
		config.Producer.RequiredAcks = sarama.WaitForLocal
	}

	return config
}

func main() {
	flag.Parse()

	var (
		wg                            sync.WaitGroup
		enqueued, successes, failures int
		totalLatency                  time.Duration
	)

	if *verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	producer, err := sarama.NewAsyncProducer(strings.Split(*brokerList, ","), producerConfiguration())
	if err != nil {
		log.Fatalln("Failed to start producer:", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var (
			latency, batchDuration time.Duration
			batchStartedAt         time.Time
			rate                   float64
		)

		batchStartedAt = time.Now()
		for message := range producer.Successes() {
			totalLatency += message.Metadata.(*MessageMetadata).Latency()
			successes++

			if successes%*statFrequency == 0 {

				batchDuration = time.Since(batchStartedAt)
				rate = float64(*statFrequency) / (float64(batchDuration) / float64(time.Second))
				latency = totalLatency / time.Duration(*statFrequency)

				log.Printf("Rate: %0.2f/s; latency: %0.2fms\n", rate, float64(latency)/float64(time.Millisecond))

				totalLatency = 0
				batchStartedAt = time.Now()
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println("FAILURE:", err)
			failures++
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGTERM)

	//messageBody := sarama.ByteEncoder(make([]byte, *messageBodySize))
	messageBody := sarama.StringEncoder(testData)
ProducerLoop:
	for i := 0; i < num; i++ {
		message := &sarama.ProducerMessage{
			Topic:    *topic,
			Key:      sarama.StringEncoder(fmt.Sprintf("%d", enqueued)),
			Value:    messageBody,
			Metadata: &MessageMetadata{EnqueuedAt: time.Now()},
		}

		select {
		case <-signals:
			producer.AsyncClose()
			break ProducerLoop
		case producer.Input() <- message:
			enqueued++
		}
		fmt.Println(i)

		if *sleep > 0 {
			time.Sleep(time.Duration(*sleep) * time.Second)
		}
	}

	fmt.Println("Waiting for in flight messages to be processed...")
	wg.Wait()

	log.Println()
	log.Printf("Enqueued: %d; Produced: %d; Failed: %d.\n", enqueued, successes, failures)
}
