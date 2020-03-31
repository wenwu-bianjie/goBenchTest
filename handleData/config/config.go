package config

import (
	"encoding/json"
	"io/ioutil"
)

var TRANS_FIN_TS string = "TRANS_FIN_TS"
var ToTs string = "toTs"
var TRANS_RCV_TS string = "TRANS_RCV_TS"

type Config struct {
	ConsumerBrokerList string `json:"consumer_broker_list"`
	ConsumerTopic      string `json:"consumer_topic"`
	ConsumerPartition  int    `json:"consumer_partition"`
	ConsumerOffset     string `json:"consumer_offset"`  //The offset to start with. Can be `oldest`, `newest`, or an actual offset
	ConsumerVerbose    bool   `json:"consumer_verbose"` //Whether to turn on sarama logging
	ConsumerNumber     int64  `json:"consumer_number"`

	ProducerBrokerList  string `json:"producer_broker_list"`
	ProducerTopic       string `json:"producer_topic"`
	ProducerPartitioner string `json:"producer_partitioner"`
	ProducerVerbose     bool   `json:"producer_verbose"`

	IsSwtSucc_sql  string `json:"is_swt_succ_sql"`
	IsSwtSuccKey   string `json:"is_swt_succ_key"`
	ToTsExpression string `json:"to_ts_expression"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	// 1, 把配置文件读进来
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 2, 做JSON反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	// 3, 赋值单例
	G_config = &conf

	return
}
