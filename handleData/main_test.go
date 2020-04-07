package main

import (
	"fmt"
	"github.com/wenwu-bianjie/goBenchTest/handleData/config"
	"github.com/wenwu-bianjie/goBenchTest/handleData/consumer"
	synatx "github.com/wenwu-bianjie/goBenchTest/handleData/syntax/simple_explain"
	"github.com/wenwu-bianjie/goBenchTest/handleData/syntax/util"
	"strings"
	"testing"
	"time"
)


func BenchmarkSubstr(b *testing.B) {
	var err error
	// 初始化命令行参数

	// 加载配置
	if err = config.InitConfig("./config.json"); err != nil {
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
	t1 := float64(time.Now().UnixNano())
	// 消费数据，并转为map格式
	consumer.ForConsumer(isSwtSucc_sql_o, ToTsExpression_o)

	t2 := float64(time.Now().UnixNano())
	fmt.Printf("总耗时: %v ms",(t2 - t1) / 1000000)
}