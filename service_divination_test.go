package yuanfenju

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestMeiriDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"number":343,"guaming":"留连","description":{"卦曰":"a","解曰":"b","凶吉":"c","运势":"d","财富":"e","感情":"f","事业":"g","身体":"h","神鬼":"i","行人":"j"}}}`
	var resp CommonResponse[MeiriData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal meiri data failed: %v", err)
	}

	if resp.Data.Number != 343 {
		t.Fatalf("unexpected number: %d", resp.Data.Number)
	}
	if resp.Data.GuaMing != "留连" {
		t.Fatalf("unexpected guaming: %s", resp.Data.GuaMing)
	}
	if resp.Data.Description.GuaYue != "a" || resp.Data.Description.XingRen != "j" {
		t.Fatalf("unexpected description: %#v", resp.Data.Description)
	}
}

func TestMeiriRequestValidate(t *testing.T) {
	if err := (MeiriRequest{Lang: "zh-cn"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	err := (MeiriRequest{Lang: "fr-fr"}).Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestXiaoliurenRequestValidate(t *testing.T) {
	if err := (XiaoliurenRequest{Shuzi: "0", Lang: "zh-cn"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}
	if err := (XiaoliurenRequest{Shuzi: "999999999999", Lang: "en-us"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	cases := []XiaoliurenRequest{
		{Shuzi: "", Lang: "zh-cn"},
		{Shuzi: "-1", Lang: "zh-cn"},
		{Shuzi: "1000000000000", Lang: "zh-cn"},
		{Shuzi: "abc", Lang: "zh-cn"},
		{Shuzi: "100", Lang: "fr-fr"},
	}
	for _, req := range cases {
		err := req.Validate()
		if err == nil {
			t.Fatalf("expected validation error for req=%#v", req)
		}
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation, got: %v", err)
		}
	}
}

func TestXiaoliurenDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"number":100,"guaming":"大安","description":{"卦曰":"a","解曰":"b","凶吉":"c","运势":"d","财富":"e","感情":"f","事业":"g","身体":"h","神鬼":"i","行人":"j"}}}`
	var resp CommonResponse[XiaoliurenData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal xiaoliuren data failed: %v", err)
	}
	if resp.Data.Number != 100 || resp.Data.GuaMing != "大安" {
		t.Fatalf("unexpected xiaoliuren data: %#v", resp.Data)
	}
	if resp.Data.Description.GuaYue != "a" || resp.Data.Description.XingRen != "j" {
		t.Fatalf("unexpected description: %#v", resp.Data.Description)
	}
}
