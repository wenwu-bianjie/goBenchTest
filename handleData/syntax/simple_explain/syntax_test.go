package synatx

import (
	"testing"
)

func TestSyntax_Mamshal(t *testing.T) {
	var o = Syntax{
		Expression: "@RCV_INS_ID_CD <> 01000000",
	}

	o.Marshal().MarshalRegexp()

	if o.MatchString("01000000") != true {
		t.Error("MatchString fail")
	}

	//var o1 = Syntax{
	//	Expression: "@RCV_INS_ID_CD = 01000000",
	//}
	//
	//o1.Marshal().MarshalRegexp()
	//
	//if o1.MatchString("01000000") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//if o1.MatchString("010000") != false {
	//	t.Error("MatchString fail")
	//}
	//
	//var o2 = Syntax{
	//	Expression: "@RCV_INS_ID_CD in {01000000, 01000001}",
	//}
	//
	//o2.Marshal().MarshalRegexp()
	//
	//if o2.MatchString("01000000") != true || o2.MatchString("01000001") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//if o2.MatchString("aaaa") != false {
	//	t.Error("MatchString fail")
	//}
	//
	//var o3 = Syntax{
	//	Expression: "@RCV_INS_ID_CD ex {01000000, 01000001}",
	//}
	//
	//o3.Marshal().MarshalRegexp()
	//
	//if o3.MatchString("0100000") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var o4 = Syntax{
	//	Expression: "@requestUrl head /sso/login",
	//}
	//
	//o4.Marshal().MarshalRegexp()
	//
	//if o4.MatchString("/sso/login") != false || o4.MatchString("/sso/login/") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var o5 = Syntax{
	//	Expression: "@requestUrl tail  /sso/login",
	//}
	//
	//o5.Marshal().MarshalRegexp()
	//
	//if o5.MatchString("/sso/login") != false || o5.MatchString("s/sso/login") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var o6 = Syntax{
	//	Expression: "@requestUrl vague   /sso/login",
	//}
	//
	//o6.Marshal().MarshalRegexp()
	//
	//if o6.MatchString("/sso/login") != false || o6.MatchString("s/sso/login/sfd") != true {
	//	t.Error("MatchString fail")
	//}
	//
	//var o7 = Syntax{
	//	Expression: "@ORDER_TYPE sub[0,2] in {02}",
	//}
	//
	//o7.Marshal().MarshalRegexp()
	//
	//fmt.Println(o7)
	//
	//if o7.MatchString("02333") != true || o7.MatchString("04") != false {
	//	t.Error("MatchString fail")
	//}
}

func BenchmarkSyntax_Marshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var o7 = Syntax{
			Expression: "@ORDER_TYPE sub[0,2] in {02}",
		}

		o7.Marshal().MarshalRegexp()

		if o7.MatchString("02333") != true || o7.MatchString("04") != false {
			b.Error("MatchString fail")
		}
	}
}
