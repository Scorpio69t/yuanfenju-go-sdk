package yuanfenju

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestToolsQqRequestValidate(t *testing.T) {
	if err := (ToolsQqRequest{Qq: "323366223", Lang: "zh-cn"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	err := (ToolsQqRequest{Qq: "", Lang: "zh-cn"}).Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestToolsShoujiRequestValidate(t *testing.T) {
	if err := (ToolsShoujiRequest{Shouji: "13800138000", Lang: "en-us"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	err := (ToolsShoujiRequest{Shouji: "13800138000", Lang: "zh-tw"}).Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestToolsShuziRequestValidate(t *testing.T) {
	if err := (ToolsShuziRequest{Shuzi: "8888", Lang: "zh-cn"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	cases := []ToolsShuziRequest{
		{Shuzi: "", Lang: "zh-cn"},
		{Shuzi: "-1", Lang: "zh-cn"},
		{Shuzi: "1000000000000", Lang: "zh-cn"},
		{Shuzi: "abc", Lang: "zh-cn"},
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

func TestToolsJixiongDataUnmarshalWithNumericData(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"desc":"云开见月","xiongji":"吉","desc1":"","score":"95","data":323366223,"shuli":39}}`
	var resp CommonResponse[ToolsJixiongData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal tools jixiong data failed: %v", err)
	}

	if resp.Data.InputData != "323366223" || resp.Data.Shuli != 39 || resp.Data.Xiongji != "吉" {
		t.Fatalf("unexpected tools jixiong data: %#v", resp.Data)
	}
}

func TestToolsJixiongDataUnmarshalWithStringData(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"desc":"测试","xiongji":"凶","desc1":"d","score":88,"data":"13800138000","shuli":"18"}}`
	var resp CommonResponse[ToolsJixiongData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal tools jixiong data failed: %v", err)
	}

	if resp.Data.Score != "88" || resp.Data.InputData != "13800138000" || resp.Data.Shuli != 18 {
		t.Fatalf("unexpected tools jixiong data: %#v", resp.Data)
	}
}
