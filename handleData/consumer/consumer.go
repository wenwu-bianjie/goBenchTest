package consumer

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wenwu-bianjie/goBenchTest/handleData/config"
	"github.com/wenwu-bianjie/goBenchTest/handleData/producer"
	synatx "github.com/wenwu-bianjie/goBenchTest/handleData/syntax/simple_explain"
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
  	kafkaProducer *sarama.AsyncProducer
)

func ForConsumer(isSwtSucc_sql_o *synatx.SyntaxRes, ToTsExpression_o *synatx.SyntaxRes) {
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

	kafkaProducer, err = producer.GetProducer()

	if err != nil {
		logger.Fatalln(err)
	}
	defer (*kafkaProducer).Close()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		pc.AsyncClose()
	}()

	var i int64 = 1

	for msg := range pc.Messages() {
		//fmt.Printf("Offset: %d\n", msg.Offset)
		var value map[string]interface{}
		if err := json.Unmarshal(msg.Value, &value); err == nil {
			// IsSwtSucc_sql语法表达式匹配
			res := isSwtSucc_sql_o.SyntaxNodes.MatchJson(&value)
			value[config.G_config.IsSwtSuccKey] = strconv.FormatBool(res)

			// 监控对象识别
			m := ToTsExpression_o.SyntaxNodes.MatchJson(&value)

			if m {
				// 发送给kafka
				value, err := json.Marshal(value)
				if err == nil {
					var keyEncoder, valueEncoder sarama.Encoder
					keyEncoder = sarama.StringEncoder(msg.Key)
					valueEncoder = sarama.StringEncoder(value)

					msg := &sarama.ProducerMessage{
						Topic: config.G_config.ProducerTopic,
						Key:   keyEncoder,
						Value: valueEncoder,
					}
					(*kafkaProducer).Input() <- msg
				}
			}
		} else {
			fmt.Printf("Value:  %s\n", string(msg.Value))
		}

		if i >= config.G_config.ConsumerNumber {
			pc.AsyncClose()
			break
		}
		i++
	}
	fmt.Println(i)

	if err := c.Close(); err != nil {
		fmt.Println("Failed to close consumer: ", err)
	}
}
