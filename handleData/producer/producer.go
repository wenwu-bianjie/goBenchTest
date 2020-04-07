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

func GetProducer() (*sarama.AsyncProducer, error) {
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

	conf := sarama.NewConfig()
	conf.Producer.Partitioner = partitionerConstructor
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(strings.Split(config.G_config.ProducerBrokerList, ","), conf)
	if err != nil {
		logger.Fatalln("FAILED to open the producer:", err)
		return nil, err
	}

	return &producer, err
}
