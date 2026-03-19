package yuanfenju

import (
	"encoding/json"
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
