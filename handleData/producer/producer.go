package producer

import (
	"github.com/Shopify/sarama"
	"github.com/wenwu-bianjie/goBenchTest/handleData/config"
	"log"
	"os"
	"strings"
)

var (
	// The comma separated list of brokers in the Kafka cluster
	//brokerList = "localhost:9092"
	// The topic to produce to
	//topic = "test"
	// The partitioning scheme to use. Can be `hash`, or `random`
	//partitioner = "hash"
	// Whether to turn on sarama logging
	//verbose = false

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func ForProducer(k, v string) {
	if config.G_config.ProducerVerbose {
		sarama.Logger = logger
	}

	var partitionerConstructor sarama.PartitionerConstructor
	switch config.G_config.ProducerPartitioner {
	case "hash":
		partitionerConstructor = sarama.NewHashPartitioner
	case "random":
		partitionerConstructor = sarama.NewRandomPartitioner
	default:
		log.Fatalf("Partitioner %s not supported.", config.G_config.ProducerPartitioner)
	}

	var keyEncoder, valueEncoder sarama.Encoder
	//if *key != "" {
	//	keyEncoder = sarama.StringEncoder(*key)
	//}
	//if *value != "" {
	//	valueEncoder = sarama.StringEncoder(*value)
	//}

	if k != "" {
		keyEncoder = sarama.StringEncoder(k)
	}

	if v != "" {
		valueEncoder = sarama.StringEncoder(v)
	}

	conf := sarama.NewConfig()
	conf.Producer.Partitioner = partitionerConstructor
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(config.G_config.ProducerBrokerList, ","), conf)
	if err != nil {
		logger.Fatalln("FAILED to open the producer:", err)
	}
	defer producer.Close()

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: config.G_config.ProducerTopic,
		Key:   keyEncoder,
		Value: valueEncoder,
	})

	if err != nil {
		logger.Println("FAILED to produce message:", err)
	} else {
		//fmt.Printf("topic=%s\tpartition=%d\toffset=%d\n", config.G_config.ProducerTopic, partition, offset)
	}
}
