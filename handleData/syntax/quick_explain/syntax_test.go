package main

import (
	"testing"
)

func TestSynax_Marshal(t *testing.T) {
	var o1 = Synax{
		DataSource: "CUPS",
		Field:      "RCV_INS_IS_CD",
		Value:      "00000000",
	}

	var res1 = "CUPS_00000000_R"

	var o2 = Synax{
		DataSource: "CUPS",
		Field:      "MCHNT_CD",
		Value:      "A&B",
	}

	var res2 = "CUPS_AB_M"

	if o1.Marshal() != res1 {
		t.Error("marshal fail")
	}

	if o2.Marshal() != res2 {
		t.Error("marshal fail")
	}
}
