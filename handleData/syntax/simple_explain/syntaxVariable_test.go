package synatx

import (
	"testing"
)

func TestSyntax_MarshalWithVariable(t *testing.T) {
	var o1 = Syntax{
		Expression: "@RCV_INS_ID_CD in {%CUPS_14293410_R_LIST%}",
	}

	var data = make(map[string]interface{})

	data["CUPS_14293410_R_LIST"] = "1111111111111111"

	o1.Marshal().MarshalWithVariable().ReplaceVariableWithData(data).MarshalRegexp()
	if o1.MatchString("1111111111111111") != true {
		t.Error("MatchString fail")
	}

	//json格式变量
	var o2 = Syntax{
		Expression: "@RCV_INS_ID_CD ex {%$UPACt_ysfAPP_NA_LIST%$}",
	}

	data["UPACt"] = map[string]interface{}{
		"ysfAPP": map[string]interface{}{
			"NA": map[string]interface{}{
				"LIST": "11111111111",
			},
		},
	}

	o2.Marshal().MarshalWithJsonVariable("_").ReplaceJsonVariableWithData(data).MarshalRegexp()
	if o2.MatchString("1111111111") != true {
		t.Error("MatchString fail")
	}
}

func BenchmarkSyntax_MarshalWithJsonVariable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var data = make(map[string]interface{})
		//json格式变量
		var o2 = Syntax{
			Expression: "@RCV_INS_ID_CD ex {%$UPACt_ysfAPP_NA_LIST%$}",
		}

		data["UPACt"] = map[string]interface{}{
			"ysfAPP": map[string]interface{}{
				"NA": map[string]interface{}{
					"LIST": "11111111111",
				},
			},
		}

		o2.Marshal().MarshalWithJsonVariable("_").ReplaceJsonVariableWithData(data).MarshalRegexp()
		if o2.MatchString("1111111111") != true {
			b.Error("MatchString fail")
		}
	}
}
