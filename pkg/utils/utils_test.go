package utils

import "testing"

func TestCreateID(t *testing.T) {
	uid := "u1"
	uid2 := "u2"
	id := CreateID(uid, uid2)
	if id != uid+"->"+uid2 {
		t.Fatal("create id fail")
	}
}

func TestIDGenerate(t *testing.T) {
	id := IDGenerate()
	if id == "" {
		t.Fatal("create id fail")
	}
}

func TestCheckAiKeyWord(t *testing.T) {
	ok, _ := CheckAiKeyWord("@AI")
	if ok == false {
		t.Fatal("check ai key word fail")
	}
	ok, _ = CheckAiKeyWord("@ss")
	if ok == true {
		t.Fatal("check ai key word fail")
	}
}
