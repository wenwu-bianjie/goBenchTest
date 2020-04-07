package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wenwu-bianjie/goBenchTest/handleData/config"
	"github.com/wenwu-bianjie/goBenchTest/handleData/consumer"
	"github.com/wenwu-bianjie/goBenchTest/handleData/producer"
	synatx "github.com/wenwu-bianjie/goBenchTest/handleData/syntax/simple_explain"
	"github.com/wenwu-bianjie/goBenchTest/handleData/syntax/util"
	"strconv"
	"strings"
	"time"
)

var dataChan chan map[string]interface{} = make(chan map[string]interface{}, 1)
var keyChan chan []byte = make(chan []byte, 1)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./config.json", "指定config.json")
	flag.Parse()
}

func main() {
	t1 := time.Now()
	var err error
	// 初始化命令行参数
	initArgs()

	// 加载配置
	if err = config.InitConfig(confFile); err != nil {
		fmt.Println(err)
		return
	}

	// 生成IsSwtSucc_sql语法表达式的匹配对象
	isSwtSucc_sql_s := strings.Replace(config.G_config.IsSwtSucc_sql, "=", " = ", -1)
	isSwtSucc_sql_s = strings.Replace(isSwtSucc_sql_s, "<>", " <> ", -1)
	isSwtSucc_sql_s = strings.Replace(isSwtSucc_sql_s, "'", "", -1)
	isSwtSucc_sql_o := synatx.NewSyntaxANodes(isSwtSucc_sql_s)

	// 生成监控对象的匹配对象
	var ToTsExpression_o = synatx.NewSyntaxANodes(util.RemoveStringFirstWord(config.G_config.ToTsExpression))

	// 消费数据，并转为map格式
	go consumer.ForConsumer(dataChan, keyChan)

	// 转换时间戳格式
	//nowTimeStamp := time.Now().UnixNano()
	for data := range dataChan {
		key := <-keyChan
		//if t, ok := data[config.ToTs]; ok {
		//	switch t.(type) {
		//	case string:
		//		if ts := transferTime.TransferTime(t.(string)); ts < nowTimeStamp {
		//			//continue
		//		}
		//	}
		//}

		// IsSwtSucc_sql语法表达式匹配
		res := isSwtSucc_sql_o.SyntaxNodes.MatchJson(&data)
		data[config.G_config.IsSwtSuccKey] = strconv.FormatBool(res)

		// 监控对象识别
		m := ToTsExpression_o.SyntaxNodes.MatchJson(&data)

		if m {
			// 发送给kafka
			value, err := json.Marshal(data)
			if err == nil {
				producer.ForProducer(string(key), string(value))
			}
		}
	}
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
}
