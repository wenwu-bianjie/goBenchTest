package synatx

import (
	"fmt"
	"github.com/wenwu-bianjie/goBenchTest/handleData/syntax/util"
	"regexp"
	"strconv"
	"strings"
)

//监控对象识别：简单表达式识别
type symbol string

const (
	Equal   symbol = "="
	unEqual symbol = "<>"
	In      symbol = "in"
	Ex      symbol = "ex"
	Head    symbol = "head"
	Tail    symbol = "tail"
	Vague   symbol = "vague"
)

type Syntax struct {
	Expression   string           `json:"expression"`
	Symbol       symbol           `json:"symbol"`
	IsNot        bool             `json:"is_not"`
	Field        string           `json:"field"`
	Value        string           `json:"value"`
	Regs         []*regexp.Regexp `json:"regs"`
	Sub          Sub              `json:"sub"`
	Variable     Variable         `json:"variable"`
	JsonVariable JsonVariable
}

type Sub struct {
	IsSub bool `json:"is_sub"`
	From  int  `json:"from"`
	To    int  `json:"to"`
}

type Variable struct {
	HasVariable bool   `json:"has_variable"`
	Variable    string `json:"variable"`
	Value       string `json:"value"`
}

type JsonVariable struct {
	HasVariable bool     `json:"has_variable"`
	Variable    []string `json:"variable"`
	Value       string   `json:"value"`
}

func (s *Syntax) Marshal() *Syntax {
	arr := strings.SplitN(s.Expression, " ", -1)
	for _, v := range arr {
		v = util.TrimStr(v)
		if len(v) == 0 {
			continue
		}
		if len(s.Field) == 0 {
			//s.Field = util.RemoveStringFirstWord(v)
			s.Field = v
			continue
		} else {
			m := util.SubRegMatch(v)
			if len(m) == 1 {
				s.Sub.IsSub = true
				if len(m[0]) == 3 {
					var mStr = m[0]
					from, err := strconv.Atoi(string(mStr[1]))
					if err != nil {
						fmt.Printf("Marshal sub err: %v", err)
					}
					s.Sub.From = from
					to, err := strconv.Atoi(string(mStr[2]))
					if err != nil {
						fmt.Printf("Marshal sub err: %v", err)
					}
					s.Sub.To = to
				}
				continue
			}
		}

		if len(s.Symbol) == 0 {
			s.Symbol = symbol(v)
			continue
		}

		if len(s.Value) == 0 {
			s.Value = v
		} else {
			s.Value = s.Value + v
		}
	}
	return s
}

func (s *Syntax) MarshalRegexp() *Syntax {
	var (
		reg    *regexp.Regexp
		err    error
		regStr string
	)
	switch s.Symbol {
	case Equal:
		regStr = fmt.Sprintf("^%s$", s.Value)
		if reg, err = regexp.Compile(regStr); err == nil {
			s.Regs = append(s.Regs, reg)
		}
	case unEqual:
		regStr = fmt.Sprintf("^%s$", s.Value)
		if reg, err = regexp.Compile(regStr); err == nil {
			s.Regs = append(s.Regs, reg)
		}
	case In:
		values := util.FormatStringWithBraceToSlice(s.Value)
		for _, v := range values {
			regStr = fmt.Sprintf("^%s$", v)
			if reg, err = regexp.Compile(regStr); err == nil {
				s.Regs = append(s.Regs, reg)
			}
		}
	case Ex:
		values := util.FormatStringWithBraceToSlice(s.Value)
		for _, v := range values {
			regStr = fmt.Sprintf("^%s$", v)
			if reg, err = regexp.Compile(regStr); err == nil {
				s.Regs = append(s.Regs, reg)
			}
		}
	case Head:
		regStr = fmt.Sprintf("^%s%s$", s.Value, ".+")
		if reg, err = regexp.Compile(regStr); err == nil {
			s.Regs = append(s.Regs, reg)
		}
	case Tail:
		regStr = fmt.Sprintf("^%s%s$", ".+", s.Value)
		if reg, err = regexp.Compile(regStr); err == nil {
			s.Regs = append(s.Regs, reg)
		}
	case Vague:
		regStr = fmt.Sprintf("^%s%s%s$", ".+", s.Value, ".+")
		if reg, err = regexp.Compile(regStr); err == nil {
			s.Regs = append(s.Regs, reg)
		}
	}

	return s
}

func (s *Syntax) MatchString(v string) bool {
	if s.Sub.IsSub == true {
		v = v[s.Sub.From:s.Sub.To]
	}
	if len(s.Regs) == 0 {
		return true
	}
	if (s.Symbol != Ex) && (s.Symbol != unEqual) {
		for _, r := range s.Regs {
			res := r.Find([]byte(v))
			if res != nil {
				return true
			}
		}
		return false
	} else {
		for _, r := range s.Regs {
			res := r.Find([]byte(v))
			if res != nil {
				return false
			}
		}
		return true
	}
}
