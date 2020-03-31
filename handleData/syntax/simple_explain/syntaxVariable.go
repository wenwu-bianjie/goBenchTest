package synatx

import (
	"regexp"
	"strings"
)

var marshalWithVariableReg = regexp.MustCompile("%([^%]+)%")

var marshalWithJsonVariableReg = regexp.MustCompile("%\\$([^%]+)%\\$")

func (s *Syntax) MarshalWithVariable() *Syntax {
	m := marshalWithVariableReg.FindAllSubmatch([]byte(s.Value), -1)
	if len(m) != 0 {
		for _, v := range m {
			if len(v) == 2 {
				s.Variable.HasVariable = true
				s.Variable.Value = string(v[0])
				s.Variable.Variable = string(v[1])
			}
		}
	}

	return s
}

func (s *Syntax) ReplaceVariableWithData(data map[string]interface{}) *Syntax {
	if s.Variable.HasVariable == true {
		if v, ok := data[s.Variable.Variable]; ok {
			if value, ok := v.(string); ok {
				s.Value = strings.Replace(s.Value, s.Variable.Value, value, -1)
			}
		}
	}
	return s
}

func (s *Syntax) MarshalWithJsonVariable(sep string) *Syntax {
	m := marshalWithJsonVariableReg.FindAllSubmatch([]byte(s.Value), -1)
	if len(m) != 0 {
		for _, v := range m {
			if len(v) == 2 {
				s.JsonVariable.HasVariable = true
				s.JsonVariable.Value = string(v[0])
				strs := strings.SplitN(string(v[1]), sep, -1)
				s.JsonVariable.Variable = strs
			}
		}
	}
	return s
}

func (s *Syntax) ReplaceJsonVariableWithData(data map[string]interface{}) *Syntax {
	var (
		source map[string]interface{}
	)
	source = data
	if s.JsonVariable.HasVariable == true {
		length := len(s.JsonVariable.Variable) - 1
		for i, key := range s.JsonVariable.Variable {
			if v, ok := source[key]; ok {
				if i == length {
					if next, ok := v.(string); ok {
						s.Value = strings.Replace(s.Value, s.JsonVariable.Value, next, -1)
					}
				}
				if next, ok := v.(map[string]interface{}); ok {
					source = next
				}
			}
		}
	}
	return s
}
