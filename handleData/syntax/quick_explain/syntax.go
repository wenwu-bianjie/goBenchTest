package main

import (
	"fmt"
	"github.com/wenwu-bianjie/aimTest/demo3/util"
	"strings"
)

//监控对象识别：快速识别
const and = "&"

type Synax struct {
	DataSource string `json:"data_source"`
	Field      string `json:"field"`
	Value      string `json:"value"`
}

//转为监控对象的ID格式
//根据value中是否包含&字符分开处理
func (s *Synax) Marshal() (res string) {
	if strings.Contains(s.Value, and) {
		value := strings.Replace(s.Value, and, "", 1)
		res = fmt.Sprintf("%s_%s_%s", s.DataSource, value, util.SubStringFirstWord(s.Field))
	} else {
		res = fmt.Sprintf("%s_%s_%s", s.DataSource, s.Value, util.SubStringFirstWord(s.Field))
	}
	return
}
