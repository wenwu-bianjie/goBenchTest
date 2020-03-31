package consumer

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wenwu-bianjie/goBenchTest/handleData/config"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

var (
	//brokerList = flag.String("brokers", "localhost:9092", "The comma separated list of brokers in the Kafka cluster")
	//topic      = flag.String("topic", "testData", "The topic to consume")
	//partition  = flag.Int("partition", 0, "The partition to consume")
	//offset     = flag.String("offset", "oldest", "The offset to start with. Can be `oldest`, `newest`, or an actual offset")
	//verbose    = flag.Bool("verbose", false, "Whether to turn on sarama logging")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func ForConsumer(dataChan chan map[string]interface{}, keyChan chan []byte) {
	flag.Parse()

	if config.G_config.ConsumerVerbose {
		sarama.Logger = logger
	}

	var (
		initialOffset int64
		offsetError   error
	)
	switch config.G_config.ConsumerOffset {
	case "oldest":
		initialOffset = sarama.OffsetOldest
	case "newest":
		initialOffset = sarama.OffsetNewest
	default:
		initialOffset, offsetError = strconv.ParseInt(config.G_config.ConsumerOffset, 10, 64)
	}

	if offsetError != nil {
		logger.Fatalln("Invalid initial offset:", config.G_config.ConsumerOffset)
	}

	c, err := sarama.NewConsumer(strings.Split(config.G_config.ConsumerBrokerList, ","), nil)
	if err != nil {
		logger.Fatalln(err)
	}

	pc, err := c.ConsumePartition(config.G_config.ConsumerTopic, int32(config.G_config.ConsumerPartition), initialOffset)
	if err != nil {
		logger.Fatalln(err)
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		close(dataChan)
		close(keyChan)
		pc.AsyncClose()
	}()

	for msg := range pc.Messages() {
		//fmt.Printf("Offset: %d\n", msg.Offset)
		var value map[string]interface{}
		if err := json.Unmarshal(msg.Value, &value); err == nil {
			dataChan <- value
			keyChan <- msg.Key
		} else {
			fmt.Printf("Value:  %s\n", string(msg.Value))
		}
		if msg.Offset == config.G_config.ConsumerNumber {
			close(dataChan)
			close(keyChan)
			pc.AsyncClose()
			break
		}
	}

	if err := c.Close(); err != nil {
		fmt.Println("Failed to close consumer: ", err)
	}
}
