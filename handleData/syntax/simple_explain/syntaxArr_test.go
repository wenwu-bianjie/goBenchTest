package synatx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var testData string = `{
"SYS_TRA_NO":"990771",
"dataSource":"CPUS",
"TERM_ID":"",
"ISS_RESP_CD":"90",
"MCHNT_CD":"842584073990001",
"RESV_FLD1_2":"1",
"FWD_SYS_ID":"D",
"RESV_FLD1_1":"0",
"MCHNT_TP":"7399",
"PRI_ACCT_NO_CONV":"196228481938229227475",
"CROSS_DIST_IN":"0",
"FWD_LINE_NO":"0000041274",
"TRANS_ST":"10000",
"STI_IN":"0",
"CARD_MEDIA":"2",
"CARD_CLASS":"01",
"SETTLE_DT":"20200325",
"TRANS_ID":"W20",
"RCV_INS_ID_CD":"01039200",
"ACQ_INS_ID_CD":"48429202",
"ISS_INS_ID_CD":"01030000",
"POS_ENTRY_MD_CD":"012",
"ACPT_RESP_CD":"91",
"TRANS_AT":"00000000000000000000000500000",
"TRANS_FIN_TS":"2020-03-25 09:18:49.594581",
"DB_ID":"1",
"SYS_ID":"D",
"toTs":"2020-03-25 09:18:49.535638",
"TRANS_RCV_TS":"2020-03-25 09:18:49.535638",
"MSG_TP":"9900",
"TFR_DT_TM":"0325113725",
"TRANS_CHNL":"07",
"CARD_ATTR":"01",
"TRANS_ID_CONV":"W20",
"CARD_BIN":"19622848",
"CARD_BRAND":"12",
"dataType":"DGF11",
"RCV_PROC_IN":"1",
"APP_HOST_ID":"14",
"FWD_INS_ID_CD":"48429202",
"FWD_PROC_IN":"1"
}`

func TestNewSyntaxANodes(t *testing.T) {
	var IsSwtSucc_sql string = "not(ISS_RESP_CD='90' and ACPT_RESP_CD<>'90') or (ISS_RESP_CD='93'and ACPT_RESP_CD<>'93') or (ISS_RESP_CD='A0'and ACPT_RESP_CD<>'A0') or (ISS_RESP_CD='A7'and ACPT_RESP_CD<>'A7') or (ISS_RESP_CD='A8'and ACPT_RESP_CD<>'A8') or (ISS_RESP_CD='A9'and ACPT_RESP_CD<>'A9')"
	s := strings.Replace(IsSwtSucc_sql, "=", " = ", -1)
	s = strings.Replace(s, "<>", " <> ", -1)
	s = strings.Replace(s, "'", "", -1)

	o := NewSyntaxANodes(s)
	var value map[string]interface{}
	if err := json.Unmarshal([]byte(testData), &value); err == nil {
		if o.SyntaxNodes.MatchJson(&value) != false {
			t.Error("MatchString fail")
		}
	}

	//var o1 = NewSyntaxANodes("@RCV_INS_ID_CD = 01000000")
	//
	//if o1.SyntaxNodes.MatchString("01000000") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//if o1.SyntaxNodes.MatchString("010000") != false {
	//	t.Error("MatchString fail")
	//}
	//
	//var o2 = NewSyntaxANodes("@RCV_INS_ID_CD in {01000000, 01000001}")
	//
	//if o2.SyntaxNodes.MatchString("01000000") != true || o2.SyntaxNodes.MatchString("01000001") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//if o2.SyntaxNodes.MatchString("aaaa") != false {
	//	t.Error("MatchString fail")
	//}
	//
	//var o3 = NewSyntaxANodes("@RCV_INS_ID_CD ex {01000000, 01000001}")
	//
	//if o3.SyntaxNodes.MatchString("0100000") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var o4 = NewSyntaxANodes("@requestUrl head /sso/login")
	//
	//if o4.SyntaxNodes.MatchString("/sso/login") != false || o4.SyntaxNodes.MatchString("/sso/login/") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var o5 = NewSyntaxANodes("@requestUrl tail  /sso/login")
	//
	//if o5.SyntaxNodes.MatchString("/sso/login") != false || o5.SyntaxNodes.MatchString("s/sso/login") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var o6 = NewSyntaxANodes("@requestUrl vague   /sso/login")
	//
	//if o6.SyntaxNodes.MatchString("/sso/login") != false || o6.SyntaxNodes.MatchString("s/sso/login/sfd") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var o7 = NewSyntaxANodes("@ORDER_TYPE sub[0,2] in {02}  ")
	//
	//if o7.SyntaxNodes.MatchString("02333") != true || o7.SyntaxNodes.MatchString("04") != false {
	//	t.Error("MatchString fail")
	//}
	//
	//var res = NewSyntaxANodes("@RCV_INS_ID_CD = 01000000|FWD_INS_ID_CD = 010100000")
	//if res.SyntaxNodes.MatchString("01000000") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var res1 = NewSyntaxANodes("@FWD = 00000000 | @RCV in {12345678,87654321} | (@ISS = 1222222 | @ISS = 888888883 | (@ISS = 222222 & @ISS in {222222  }  )  ) ")
	//
	//if res1.SyntaxNodes.MatchString("222222") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var res2 = NewSyntaxANodes("@FWD = 00000000 | @RCV in {12345678,87654321} | (@ISS = 1222222 | @ISS = 888888883 | (@ISS = 222222 & @ISS head 222  ) )")
	//
	//if res2.SyntaxNodes.MatchString("222222") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var res3 = NewSyntaxANodes("@FWD tail 1 | @RCV in {123456781,87654321} | (@ISS = 1222222 | @ISS = 888888883 | (@ISS = 222222 & @ISS head 222  ) )")
	//
	//if res3.SyntaxNodes.MatchString("87654321") != true {
	//	t.Error("MatchString fail")
	//}
}

func BenchmarkNewSyntaxANodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var res3 = NewSyntaxANodes("@FWD tail 1 | @RCV in {123456781,87654321} | (@ISS = 1222222 | @ISS = 888888883 | (@ISS = 222222 & @ISS head 222  ) )")

		if res3.SyntaxNodes.MatchString("87654321") != true {
			b.Error("MatchString fail")
		}
	}
}

func TestAa(t *testing.T) {
	var s interface{} = 144.3
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

	fmt.Println(reflect.TypeOf(s), s)
}
