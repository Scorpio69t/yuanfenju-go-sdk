package yuanfenju

import (
	"encoding/json"
	"testing"
)

func TestQueryMerchantDataUnmarshalObject(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"merchant_type":"非会员","merchant_email":"a@example.com","merchant_nickname":"demo","merchant_register_time":"2026-03-19 17:29:00","merchant_expire_time":"--","merchant_remaining_call_times":"99"}}`
	var resp CommonResponse[QueryMerchantData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal object failed: %v", err)
	}

	if resp.Data.MerchantType != "非会员" {
		t.Fatalf("unexpected merchant_type: %s", resp.Data.MerchantType)
	}
	if resp.Data.MerchantEmail != "a@example.com" {
		t.Fatalf("unexpected merchant_email: %s", resp.Data.MerchantEmail)
	}
	if resp.Data.MerchantExpireTime != "--" || resp.Data.ExpireTime != "--" {
		t.Fatalf("unexpected merchant_expire_time: %s / %s", resp.Data.MerchantExpireTime, resp.Data.ExpireTime)
	}
	if resp.Data.MerchantRemainingCallTimes != "99" || resp.Data.CanUseNum != "99" {
		t.Fatalf("unexpected merchant_remaining_call_times: %s / %s", resp.Data.MerchantRemainingCallTimes, resp.Data.CanUseNum)
	}
}

func TestQueryMerchantDataUnmarshalLegacyObject(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"merchant_type":"次数会员","expire_time":"2027-01-31","can_use_num":"8"}}`
	var resp CommonResponse[QueryMerchantData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal legacy object failed: %v", err)
	}

	if resp.Data.MerchantType != "次数会员" {
		t.Fatalf("unexpected merchant_type: %s", resp.Data.MerchantType)
	}
	if resp.Data.MerchantExpireTime != "2027-01-31" || resp.Data.ExpireTime != "2027-01-31" {
		t.Fatalf("unexpected merchant_expire_time: %s / %s", resp.Data.MerchantExpireTime, resp.Data.ExpireTime)
	}
	if resp.Data.MerchantRemainingCallTimes != "8" || resp.Data.CanUseNum != "8" {
		t.Fatalf("unexpected merchant_remaining_call_times: %s / %s", resp.Data.MerchantRemainingCallTimes, resp.Data.CanUseNum)
	}
}

func TestQueryMerchantDataUnmarshalArray(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":[{"merchant_type":"包年会员","merchant_expire_time":"2027-01-31","merchant_remaining_call_times":"--"}]}`
	var resp CommonResponse[QueryMerchantData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal array failed: %v", err)
	}

	if resp.Data.MerchantType != "包年会员" {
		t.Fatalf("unexpected merchant_type: %s", resp.Data.MerchantType)
	}
	if resp.Data.MerchantExpireTime != "2027-01-31" || resp.Data.ExpireTime != "2027-01-31" {
		t.Fatalf("unexpected merchant_expire_time: %s / %s", resp.Data.MerchantExpireTime, resp.Data.ExpireTime)
	}
	if resp.Data.MerchantRemainingCallTimes != "--" || resp.Data.CanUseNum != "--" {
		t.Fatalf("unexpected merchant_remaining_call_times: %s / %s", resp.Data.MerchantRemainingCallTimes, resp.Data.CanUseNum)
	}
}

func TestQueryMerchantDataUnmarshalEmptyArray(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":[]}`
	var resp CommonResponse[QueryMerchantData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal empty array failed: %v", err)
	}
	if resp.Data != (QueryMerchantData{}) {
		t.Fatalf("expected zero value for empty data array, got %#v", resp.Data)
	}
}
