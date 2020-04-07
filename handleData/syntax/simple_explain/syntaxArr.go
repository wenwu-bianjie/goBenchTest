package synatx

import (
	"strconv"
	"strings"
)

const (
	Not string = "not"
	//And           string = "&"
	And string = "and"
	//Or            string = "|"
	Or            string = "or"
	LeftBrackets  string = "("
	RightBrackets string = ")"
)

type SyntaxNode struct {
	Syntax      *Syntax `json:"syntax"`
	LogicSymbol string
	SyntaxRes   SyntaxRes
}

type SyntaxNodes []*SyntaxNode

func (nodes SyntaxNodes) MatchString(s string) bool {
	var m bool
	for _, v := range nodes {
		if len(v.SyntaxRes.SyntaxNodes) > 0 {
			m = v.SyntaxRes.SyntaxNodes.MatchString(s)
			if v.Syntax.IsNot {
				m = !m
			}
		} else {
			m = v.Syntax.MatchString(s)
			if v.Syntax.IsNot {
				m = !m
			}
		}
		if m == true {
			if v.LogicSymbol == Or {
				return true
			} else if v.LogicSymbol == "" {
				return true
			} else {
				continue
			}
		}
		if m == false {
			if v.LogicSymbol == And {
				return false
			} else if v.LogicSymbol == "" {
				return false
			} else {
				continue
			}
		}
	}
	return false
}

func (nodes SyntaxNodes) MatchJson(data *map[string]interface{}) bool {
	var m bool
	for _, v := range nodes {
		if len(v.SyntaxRes.SyntaxNodes) > 0 {
			m = v.SyntaxRes.SyntaxNodes.MatchJson(data)
			if v.Syntax.IsNot {
				m = !m
			}
		} else {
			if s, ok := (*data)[v.Syntax.Field]; ok {
				switch s.(type) {
				case int:
					s = strconv.Itoa(s.(int))
				case int64:
					s = strconv.FormatInt(s.(int64), 10)
				case float64:
					s = strconv.FormatInt(int64(s.(float64)), 10)
				case float32:
					s = strconv.FormatInt(int64(float64(s.(float32))), 10)
				}

				switch s.(type) {
				case string:
					m = v.Syntax.MatchString(s.(string))
					if v.Syntax.IsNot {
						m = !m
					}
				}
			}
		}
		if m == true {
			if v.LogicSymbol == Or {
				return true
			} else if v.LogicSymbol == "" {
				return true
			} else {
				continue
			}
		}
		if m == false {
			if v.LogicSymbol == And {
				return false
			} else if v.LogicSymbol == "" {
				return false
			} else {
				continue
			}
		}
	}
	return false
}

type SyntaxRes struct {
	Symbols     []string
	SyntaxNodes SyntaxNodes
}

func NewSyntaxANodes(s string) *SyntaxRes {
	var syntaxRes = &SyntaxRes{}
	syntaxRes.forNewSyntaxANodes(s, false)
	symbolsLen := len(syntaxRes.Symbols)
	for i, v := range syntaxRes.SyntaxNodes {
		if i < symbolsLen {
			v.LogicSymbol = syntaxRes.Symbols[i]
		}
		if strings.Contains(v.Syntax.Expression, And) || strings.Contains(v.Syntax.Expression, Or) {
			var expression string = v.Syntax.Expression
			if strings.Contains(v.Syntax.Expression, Not) {
				v.Syntax.IsNot = true
				expression = strings.Replace(v.Syntax.Expression, Not, "", 1)
				expression = strings.TrimLeft(expression, " ")
			}
			var lastRightBracketsIndex = len(expression) - 1
			if strings.Index(expression, LeftBrackets) == 0 && strings.LastIndex(expression, RightBrackets) == lastRightBracketsIndex {
				syntaxRes.SyntaxNodes[i].SyntaxRes = *NewSyntaxANodes(expression[1:lastRightBracketsIndex])
			}
		} else {
			v.Syntax.Marshal().MarshalRegexp()
		}
	}
	return syntaxRes
}

func (r *SyntaxRes) forNewSyntaxANodes(s string, noSplitN bool) {
	if noSplitN {
		r.SyntaxNodes = append(r.SyntaxNodes, &SyntaxNode{
			Syntax: &Syntax{
				Expression: strings.Trim(s, " "),
			},
		})
		return
	}
	var (
		hasAnd        bool
		hasOr         bool
		length        int
		i             int
		v             string
		hasBracketStr string
		andIndex      int
		orIndex       int
	)
	hasAnd = strings.Contains(s, And)
	hasOr = strings.Contains(s, Or)
	if hasAnd && hasOr {
		andIndex = strings.Index(s, And)
		orIndex = strings.Index(s, Or)

		if andIndex < orIndex {
			hasOr = false
		} else {
			hasAnd = false
		}
		if hasAnd {
			m := strings.SplitN(s, Or, -1)
			if len(m) > 0 && strings.Count(m[0], LeftBrackets) > 0 && strings.Count(m[0], LeftBrackets) == strings.Count(m[0], RightBrackets) {
				hasOr = true
				hasAnd = false
			}
		}
		if hasOr {
			m := strings.SplitN(s, And, -1)
			if len(m) > 0 && strings.Count(m[0], LeftBrackets) > 0 && strings.Count(m[0], LeftBrackets) == strings.Count(m[0], RightBrackets) {
				hasOr = false
				hasAnd = true
			}
		}
	}
	if hasAnd == true {
		m := strings.SplitN(s, And, -1)
		length = len(m) - 1
		for ; i <= length; i++ {
			v = m[i]
			if strings.Contains(v, LeftBrackets) {
				if (strings.Count(hasBracketStr, LeftBrackets) != strings.Count(hasBracketStr, RightBrackets)) && hasBracketStr != "" {
					hasBracketStr = hasBracketStr + And + v
				} else {
					hasBracketStr = v
				}
				if strings.Count(hasBracketStr, LeftBrackets) == strings.Count(hasBracketStr, RightBrackets) {
					r.forNewSyntaxANodes(hasBracketStr, true)
					hasBracketStr = ""
					r.Symbols = append(r.Symbols, And)
				}
				continue
			}
			r.forNewSyntaxANodes(v, false)
			if i < length {
				r.Symbols = append(r.Symbols, And)
			}
		}
	}

	if hasOr == true {
		m := strings.SplitN(s, Or, -1)
		length = len(m) - 1
		for ; i <= length; i++ {
			v = m[i]
			if strings.Contains(v, LeftBrackets) {
				if (strings.Count(hasBracketStr, LeftBrackets) != strings.Count(hasBracketStr, RightBrackets)) && hasBracketStr != "" {
					hasBracketStr = hasBracketStr + Or + v
				} else {
					hasBracketStr = v
				}
				if strings.Count(hasBracketStr, LeftBrackets) == strings.Count(hasBracketStr, RightBrackets) {
					r.forNewSyntaxANodes(hasBracketStr, true)
					hasBracketStr = ""
					r.Symbols = append(r.Symbols, Or)
				}
				continue
			}
			r.forNewSyntaxANodes(v, false)
			if i < length {
				r.Symbols = append(r.Symbols, Or)
			}
		}
	}

	if hasAnd == false && hasOr == false {
		r.SyntaxNodes = append(r.SyntaxNodes, &SyntaxNode{
			Syntax: &Syntax{
				Expression: strings.Trim(s, " "),
			},
		})
	}
}
