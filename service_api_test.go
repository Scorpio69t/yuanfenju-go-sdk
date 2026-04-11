package yuanfenju

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func newTestClient(t *testing.T, expectedPath string, status int, body string) *Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		if r.PostFormValue("api_key") != "test_api_key" {
			t.Fatalf("api_key not injected: %q", r.PostFormValue("api_key"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(Config{
		APIKey:  "test_api_key",
		BaseURL: srv.URL,
		Timeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("new client failed: %v", err)
	}
	return c
}

func TestFreeQueryMerchant_Success(t *testing.T) {
	client := newTestClient(t, "/v1/Free/querymerchant", http.StatusOK, `{"errcode":0,"errmsg":"ok","data":{"merchant_type":"非会员","merchant_remaining_call_times":"12"}}`)
	resp, err := client.Free.QueryMerchant(context.Background())
	if err != nil {
		t.Fatalf("query merchant failed: %v", err)
	}
	if resp.Data.MerchantType != "非会员" {
		t.Fatalf("unexpected merchant_type: %s", resp.Data.MerchantType)
	}
}

func TestFreeQueryMerchant_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Free/querymerchant", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":[]}`)
	_, err := client.Free.QueryMerchant(context.Background())
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
	if apiErr.Code != -1 {
		t.Fatalf("unexpected API error code: %d", apiErr.Code)
	}
}

func TestFreeQueryMerchant_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Free/querymerchant", http.StatusInternalServerError, `server exploded`)
	_, err := client.Free.QueryMerchant(context.Background())
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 500") {
		t.Fatalf("expected http 500 error, got: %v", err)
	}
}

func TestFreeQueryTimes_Success(t *testing.T) {
	client := newTestClient(t, "/v1/Free/querytimes", http.StatusOK, `{"errcode":0,"errmsg":"ok","data":{"call_times":"7","expire_time":3600,"expire_time_message":"msg"}}`)
	resp, err := client.Free.QueryTimes(context.Background())
	if err != nil {
		t.Fatalf("query times failed: %v", err)
	}
	if resp.Data.CallTimes != "7" {
		t.Fatalf("unexpected call_times: %s", resp.Data.CallTimes)
	}
}

func TestFreeQueryTimes_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Free/querytimes", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Free.QueryTimes(context.Background())
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
	if apiErr.Code != -1 {
		t.Fatalf("unexpected API error code: %d", apiErr.Code)
	}
}

func TestFreeQueryTimes_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Free/querytimes", http.StatusBadGateway, `bad gateway`)
	_, err := client.Free.QueryTimes(context.Background())
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Fatalf("expected http 502 error, got: %v", err)
	}
}

func TestBaziPaipan_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"测试"},"bazi_info":{"kw":"戌亥"},"dayun_info":{},"start_info":{},"detail_info":{}}}`
	client := newTestClient(t, "/v1/Bazi/paipan", http.StatusOK, body)

	resp, err := client.Bazi.Paipan(context.Background(), BaziPaipanRequest{
		Sex:   "1",
		Type:  "1",
		Year:  "1990",
		Month: "01",
		Day:   "01",
		Hours: "12",
	})
	if err != nil {
		t.Fatalf("paipan failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "测试" {
		t.Fatalf("unexpected name: %s", resp.Data.BaseInfo.Name)
	}
}

func TestBaziPaipan_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/paipan", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Paipan(context.Background(), BaziPaipanRequest{
		Sex:   "1",
		Type:  "1",
		Year:  "1990",
		Month: "01",
		Day:   "01",
		Hours: "12",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
	if apiErr.Code != -1 {
		t.Fatalf("unexpected API error code: %d", apiErr.Code)
	}
}

func TestBaziPaipan_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/paipan", http.StatusServiceUnavailable, `unavailable`)
	_, err := client.Bazi.Paipan(context.Background(), BaziPaipanRequest{
		Sex:   "1",
		Type:  "1",
		Year:  "1990",
		Month: "01",
		Day:   "01",
		Hours: "12",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 503") {
		t.Fatalf("expected http 503 error, got: %v", err)
	}
}

func TestBaziJiuxing_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"张三"},"jiuxing":{"九星":"文曲星"}}}`
	client := newTestClient(t, "/v1/Bazi/jiuxing", http.StatusOK, body)

	resp, err := client.Bazi.Jiuxing(context.Background(), BaziJiuxingRequest{
		Sex:    "0",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err != nil {
		t.Fatalf("jiuxing failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.Jiuxing.Jiuxing != "文曲星" {
		t.Fatalf("unexpected jiuxing response: %#v", resp.Data)
	}
}

func TestBaziJiuxing_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/jiuxing", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Jiuxing(context.Background(), BaziJiuxingRequest{
		Sex:    "0",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziJiuxing_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/jiuxing", http.StatusBadGateway, `bad gateway`)
	_, err := client.Bazi.Jiuxing(context.Background(), BaziJiuxingRequest{
		Sex:    "0",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Fatalf("expected http 502 error, got: %v", err)
	}
}

func TestBaziHehun_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"male":{"name":"男方"},"female":{"name":"女方"},"minggong":{"score":"30"},"nianqitongzhi":{"score":"20"},"yueling":{"score":"5"},"rigan":{"score":"25"},"tiangan":{"score":"5"},"zinv":{"nannv":"男","score":"15"},"all_score":100}}`
	client := newTestClient(t, "/v1/Bazi/hehun", http.StatusOK, body)

	resp, err := client.Bazi.Hehun(context.Background(), BaziHehunRequest{
		MaleType:     "1",
		MaleYear:     "1988",
		MaleMonth:    "11",
		MaleDay:      "8",
		MaleHours:    "12",
		MaleMinute:   "20",
		FemaleType:   "1",
		FemaleYear:   "1988",
		FemaleMonth:  "11",
		FemaleDay:    "8",
		FemaleHours:  "12",
		FemaleMinute: "20",
	})
	if err != nil {
		t.Fatalf("hehun failed: %v", err)
	}
	if resp.Data.AllScore != 100 || resp.Data.Male.Name != "男方" || resp.Data.Female.Name != "女方" {
		t.Fatalf("unexpected hehun response: %#v", resp.Data)
	}
}

func TestBaziHehun_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/hehun", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Hehun(context.Background(), BaziHehunRequest{
		MaleType:     "1",
		MaleYear:     "1988",
		MaleMonth:    "11",
		MaleDay:      "8",
		MaleHours:    "12",
		MaleMinute:   "20",
		FemaleType:   "1",
		FemaleYear:   "1988",
		FemaleMonth:  "11",
		FemaleDay:    "8",
		FemaleHours:  "12",
		FemaleMinute: "20",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziHehun_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/hehun", http.StatusBadRequest, `bad request`)
	_, err := client.Bazi.Hehun(context.Background(), BaziHehunRequest{
		MaleType:     "1",
		MaleYear:     "1988",
		MaleMonth:    "11",
		MaleDay:      "8",
		MaleHours:    "12",
		MaleMinute:   "20",
		FemaleType:   "1",
		FemaleYear:   "1988",
		FemaleMonth:  "11",
		FemaleDay:    "8",
		FemaleHours:  "12",
		FemaleMinute: "20",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 400") {
		t.Fatalf("expected http 400 error, got: %v", err)
	}
}

func TestBaziHepan_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"male":{"name":"甲方"},"female":{"name":"乙方"},"minggong":{"score":"10"},"nianqitongzhi":{"score":"20"},"yueling":{"score":"5"},"rigan":{"score":"25"},"tiangan":{"score":"5"},"jiankang":{"score":"10"},"all_score":100}}`
	client := newTestClient(t, "/v1/Bazi/hepan", http.StatusOK, body)

	resp, err := client.Bazi.Hepan(context.Background(), BaziHepanRequest{
		MaleType:     "1",
		MaleYear:     "1988",
		MaleMonth:    "11",
		MaleDay:      "8",
		MaleHours:    "12",
		MaleMinute:   "20",
		FemaleType:   "1",
		FemaleYear:   "1988",
		FemaleMonth:  "11",
		FemaleDay:    "8",
		FemaleHours:  "12",
		FemaleMinute: "20",
	})
	if err != nil {
		t.Fatalf("hepan failed: %v", err)
	}
	if resp.Data.AllScore != 100 || resp.Data.Male.Name != "甲方" || resp.Data.Female.Name != "乙方" {
		t.Fatalf("unexpected hepan response: %#v", resp.Data)
	}
}

func TestBaziHepan_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/hepan", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Hepan(context.Background(), BaziHepanRequest{
		MaleType:     "1",
		MaleYear:     "1988",
		MaleMonth:    "11",
		MaleDay:      "8",
		MaleHours:    "12",
		MaleMinute:   "20",
		FemaleType:   "1",
		FemaleYear:   "1988",
		FemaleMonth:  "11",
		FemaleDay:    "8",
		FemaleHours:  "12",
		FemaleMinute: "20",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziHepan_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/hepan", http.StatusBadRequest, `bad request`)
	_, err := client.Bazi.Hepan(context.Background(), BaziHepanRequest{
		MaleType:     "1",
		MaleYear:     "1988",
		MaleMonth:    "11",
		MaleDay:      "8",
		MaleHours:    "12",
		MaleMinute:   "20",
		FemaleType:   "1",
		FemaleYear:   "1988",
		FemaleMonth:  "11",
		FemaleDay:    "8",
		FemaleHours:  "12",
		FemaleMinute: "20",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 400") {
		t.Fatalf("expected http 400 error, got: %v", err)
	}
}

func TestBaziCesuan_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"张三"},"bazi_info":{"kw":"戌亥"},"chenggu":{"total_weight":"5.6"},"wuxing":{"simple_desc":"木"},"yinyuan":{"sanshishu_yinyuan":"姻缘"},"caiyun":{"sanshishu_caiyun":{"simple_desc":"先苦后甜"}},"sizhu":{"rizhu":"日柱"},"mingyun":{"sanshishu_mingyun":"命运"},"sx":"龙","xz":"天蝎座","xiyongshen":{"shui_score":118}}}`
	client := newTestClient(t, "/v1/Bazi/cesuan", http.StatusOK, body)

	resp, err := client.Bazi.Cesuan(context.Background(), BaziCesuanRequest{
		Sex:       "0",
		Type:      "1",
		Year:      "1988",
		Month:     "11",
		Day:       "8",
		Hours:     "12",
		Minute:    "20",
		Zhen:      "3",
		Longitude: "116.46",
		Latitude:  "39.92",
	})
	if err != nil {
		t.Fatalf("cesuan failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.Caiyun.SanshishuCaiyun.SimpleDesc != "先苦后甜" || resp.Data.Xiyongshen.ShuiScore != 118 {
		t.Fatalf("unexpected cesuan response: %#v", resp.Data)
	}
}

func TestBaziCesuan_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/cesuan", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Cesuan(context.Background(), BaziCesuanRequest{
		Sex:       "0",
		Type:      "1",
		Year:      "1988",
		Month:     "11",
		Day:       "8",
		Hours:     "12",
		Minute:    "20",
		Zhen:      "3",
		Longitude: "116.46",
		Latitude:  "39.92",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziCesuan_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/cesuan", http.StatusBadGateway, `bad gateway`)
	_, err := client.Bazi.Cesuan(context.Background(), BaziCesuanRequest{
		Sex:       "0",
		Type:      "1",
		Year:      "1988",
		Month:     "11",
		Day:       "8",
		Hours:     "12",
		Minute:    "20",
		Zhen:      "3",
		Longitude: "116.46",
		Latitude:  "39.92",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Fatalf("expected http 502 error, got: %v", err)
	}
}

func TestBaziJingpan_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"张三","sex":"乾造"},"detail_info":{"dayun_info":[{"dayun_index":1}],"sizhu_info":{"year":{"tg":"甲"}}}}}`
	client := newTestClient(t, "/v1/Bazi/jingpan", http.StatusOK, body)

	resp, err := client.Bazi.Jingpan(context.Background(), BaziJingpanRequest{
		Sex:    "0",
		Type:   "1",
		Year:   "1994",
		Month:  "4",
		Day:    "30",
		Hours:  "10",
		Minute: "0",
	})
	if err != nil {
		t.Fatalf("jingpan failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.DetailInfo.DayunInfo[0].DayunIndex != 1 {
		t.Fatalf("unexpected jingpan response: %#v", resp.Data)
	}
}

func TestBaziJingpan_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/jingpan", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Jingpan(context.Background(), BaziJingpanRequest{
		Sex:    "0",
		Type:   "1",
		Year:   "1994",
		Month:  "4",
		Day:    "30",
		Hours:  "10",
		Minute: "0",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziJingpan_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/jingpan", http.StatusBadRequest, `bad request`)
	_, err := client.Bazi.Jingpan(context.Background(), BaziJingpanRequest{
		Sex:    "0",
		Type:   "1",
		Year:   "1994",
		Month:  "4",
		Day:    "30",
		Hours:  "10",
		Minute: "0",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 400") {
		t.Fatalf("expected http 400 error, got: %v", err)
	}
}

func TestBaziJingsuan_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"张三","sex":"乾造"},"detail_info":{"dayun_info":[{"dayun_index":1,"dayun_indication":{"shiye":"事业"}}],"sizhu_info":{"year":{"tg":"戊"},"sizhu_indication":{"caiyun":{"sanshishu_caiyun":{"simple_desc":"知足常乐"}}}}}}}`
	client := newTestClient(t, "/v1/Bazi/jingsuan", http.StatusOK, body)

	resp, err := client.Bazi.Jingsuan(context.Background(), BaziJingsuanRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err != nil {
		t.Fatalf("jingsuan failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.DetailInfo.SizhuInfo.SizhuIndication.Caiyun.SanshishuCaiyun.SimpleDesc != "知足常乐" {
		t.Fatalf("unexpected jingsuan response: %#v", resp.Data)
	}
}

func TestBaziJingsuan_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/jingsuan", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Jingsuan(context.Background(), BaziJingsuanRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziJingsuan_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/jingsuan", http.StatusBadGateway, `bad gateway`)
	_, err := client.Bazi.Jingsuan(context.Background(), BaziJingsuanRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Fatalf("expected http 502 error, got: %v", err)
	}
}

func TestBaziWeilai_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"张三","sex":"乾造"},"detail_info":{"sizhu_info":{"year":{"tg":"戊"}},"yunshi_year_info":{"yunshi_year":{"year":2026,"indication":{"shiye":"事业"}}},"yunshi_month_info":[{"month":"1月","yunshi_day_info":[{"day":"1号"}]}]}}}`
	client := newTestClient(t, "/v1/Bazi/weilai", http.StatusOK, body)

	resp, err := client.Bazi.Weilai(context.Background(), BaziWeilaiRequest{
		Sex:        "1",
		Type:       "1",
		Year:       "1988",
		Month:      "11",
		Day:        "8",
		Hours:      "12",
		Minute:     "20",
		YunshiYear: strconv.Itoa(time.Now().Year()),
	})
	if err != nil {
		t.Fatalf("weilai failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || len(resp.Data.DetailInfo.YunshiMonthInfo) == 0 || resp.Data.DetailInfo.YunshiMonthInfo[0].YunshiDayInfo[0].Day != "1号" {
		t.Fatalf("unexpected weilai response: %#v", resp.Data)
	}
}

func TestBaziWeilai_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/weilai", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Weilai(context.Background(), BaziWeilaiRequest{
		Sex:        "1",
		Type:       "1",
		Year:       "1988",
		Month:      "11",
		Day:        "8",
		Hours:      "12",
		Minute:     "20",
		YunshiYear: strconv.Itoa(time.Now().Year()),
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziWeilai_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/weilai", http.StatusServiceUnavailable, `unavailable`)
	_, err := client.Bazi.Weilai(context.Background(), BaziWeilaiRequest{
		Sex:        "1",
		Type:       "1",
		Year:       "1988",
		Month:      "11",
		Day:        "8",
		Hours:      "12",
		Minute:     "20",
		YunshiYear: strconv.Itoa(time.Now().Year()),
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 503") {
		t.Fatalf("expected http 503 error, got: %v", err)
	}
}

func TestBaziZwpan_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"张三","sex":"坤造","age":36},"gong_pan":[{"minggong":"疾厄宫","tianfuxing":"太阴"}]}}`
	client := newTestClient(t, "/v1/Bazi/zwpan", http.StatusOK, body)
	resp, err := client.Bazi.Zwpan(context.Background(), BaziZwpanRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err != nil {
		t.Fatalf("zwpan failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || len(resp.Data.GongPan) == 0 || resp.Data.GongPan[0].Minggong != "疾厄宫" {
		t.Fatalf("unexpected zwpan response: %#v", resp.Data)
	}
}

func TestBaziZwpan_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/zwpan", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Zwpan(context.Background(), BaziZwpanRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziZwpan_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/zwpan", http.StatusBadRequest, `bad request`)
	_, err := client.Bazi.Zwpan(context.Background(), BaziZwpanRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 400") {
		t.Fatalf("expected http 400 error, got: %v", err)
	}
}

func TestBaziYunshi_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"张三","sex":"坤造","yeargz":"庚辰"},"yunshi_info":{"lucky_number":"4、9","fortune_score":76,"jixiong_today":"中吉"}}}`
	client := newTestClient(t, "/v1/Bazi/yunshi", http.StatusOK, body)
	resp, err := client.Bazi.Yunshi(context.Background(), BaziYunshiRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "2000",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "11",
	})
	if err != nil {
		t.Fatalf("yunshi failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.YunshiInfo.FortuneScore != 76 {
		t.Fatalf("unexpected yunshi response: %#v", resp.Data)
	}
}

func TestBaziYunshi_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/yunshi", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Yunshi(context.Background(), BaziYunshiRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "2000",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "11",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziYunshi_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/yunshi", http.StatusBadGateway, `bad gateway`)
	_, err := client.Bazi.Yunshi(context.Background(), BaziYunshiRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "2000",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "11",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Fatalf("expected http 502 error, got: %v", err)
	}
}

func TestBaziCaiyunfenxi_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"name":"张三","liunian":"乙巳"},"bazi_info":{"kw":"子丑"},"caiyun_info":{"yearlyOverallFortuneScore":86,"yearlyMonthlyFortuneAndInvestmentTips":["1月建议"]}}}`
	client := newTestClient(t, "/v1/Bazi/caiyunfenxi", http.StatusOK, body)
	resp, err := client.Bazi.Caiyunfenxi(context.Background(), BaziCaiyunfenxiRequest{
		Sex:     "1",
		Type:    "1",
		Year:    "1988",
		LiuYear: "2025",
		Month:   "1",
		Day:     "8",
		Hours:   "12",
		Minute:  "20",
	})
	if err != nil {
		t.Fatalf("caiyunfenxi failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.CaiyunInfo.YearlyOverallFortuneScore != 86 {
		t.Fatalf("unexpected caiyunfenxi response: %#v", resp.Data)
	}
}

func TestBaziCaiyunfenxi_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/caiyunfenxi", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Bazi.Caiyunfenxi(context.Background(), BaziCaiyunfenxiRequest{
		Sex:     "1",
		Type:    "1",
		Year:    "1988",
		LiuYear: "2025",
		Month:   "1",
		Day:     "8",
		Hours:   "12",
		Minute:  "20",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestBaziCaiyunfenxi_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Bazi/caiyunfenxi", http.StatusServiceUnavailable, `unavailable`)
	_, err := client.Bazi.Caiyunfenxi(context.Background(), BaziCaiyunfenxiRequest{
		Sex:     "1",
		Type:    "1",
		Year:    "1988",
		LiuYear: "2025",
		Month:   "1",
		Day:     "8",
		Hours:   "12",
		Minute:  "20",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 503") {
		t.Fatalf("expected http 503 error, got: %v", err)
	}
}

func TestDivinationMeiri_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"number":777,"guaming":"大安","description":{"卦曰":"a"}}}`
	client := newTestClient(t, "/v1/Zhanbu/meiri", http.StatusOK, body)
	resp, err := client.Divination.Meiri(context.Background(), MeiriRequest{Lang: "zh-cn"})
	if err != nil {
		t.Fatalf("meiri failed: %v", err)
	}
	if resp.Data.Number != 777 || resp.Data.GuaMing != "大安" {
		t.Fatalf("unexpected meiri response: %#v", resp.Data)
	}
}

func TestDivinationMeiri_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/meiri", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Meiri(context.Background(), MeiriRequest{Lang: "zh-cn"})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
	if apiErr.Code != -1 {
		t.Fatalf("unexpected API error code: %d", apiErr.Code)
	}
}

func TestDivinationMeiri_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/meiri", http.StatusBadRequest, `bad request`)
	_, err := client.Divination.Meiri(context.Background(), MeiriRequest{Lang: "zh-cn"})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 400") {
		t.Fatalf("expected http 400 error, got: %v", err)
	}
}

func TestDivinationXiaoliuren_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"number":100,"guaming":"大安","description":{"卦曰":"a"}}}`
	client := newTestClient(t, "/v1/Zhanbu/xiaoliuren", http.StatusOK, body)
	resp, err := client.Divination.Xiaoliuren(context.Background(), XiaoliurenRequest{Shuzi: "100", Lang: "zh-cn"})
	if err != nil {
		t.Fatalf("xiaoliuren failed: %v", err)
	}
	if resp.Data.Number != 100 || resp.Data.GuaMing != "大安" {
		t.Fatalf("unexpected xiaoliuren response: %#v", resp.Data)
	}
}

func TestDivinationXiaoliuren_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/xiaoliuren", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Xiaoliuren(context.Background(), XiaoliurenRequest{Shuzi: "100", Lang: "zh-cn"})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
	if apiErr.Code != -1 {
		t.Fatalf("unexpected API error code: %d", apiErr.Code)
	}
}

func TestDivinationXiaoliuren_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/xiaoliuren", http.StatusBadGateway, `bad gateway`)
	_, err := client.Divination.Xiaoliuren(context.Background(), XiaoliurenRequest{Shuzi: "100", Lang: "zh-cn"})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Fatalf("expected http 502 error, got: %v", err)
	}
}

func TestDivinationZhiwen_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"muzhi":"箩纹","shizhi":"箩纹","zhongzhi":"箩纹","wumingzhi":"箩纹","xiaozhi":"箩纹","description":{"分析":"a"}}}`
	client := newTestClient(t, "/v1/Zhanbu/zhiwen", http.StatusOK, body)
	resp, err := client.Divination.Zhiwen(context.Background(), ZhiwenRequest{
		Muzhi:     "0",
		Shizhi:    "0",
		Zhongzhi:  "0",
		Wumingzhi: "0",
		Xiaozhi:   "0",
	})
	if err != nil {
		t.Fatalf("zhiwen failed: %v", err)
	}
	if resp.Data.Muzhi != "箩纹" || resp.Data.Description.Fenxi != "a" {
		t.Fatalf("unexpected zhiwen response: %#v", resp.Data)
	}
}

func TestDivinationZhiwen_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/zhiwen", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Zhiwen(context.Background(), ZhiwenRequest{
		Muzhi:     "0",
		Shizhi:    "0",
		Zhongzhi:  "0",
		Wumingzhi: "0",
		Xiaozhi:   "0",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestDivinationZhiwen_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/zhiwen", http.StatusServiceUnavailable, `unavailable`)
	_, err := client.Divination.Zhiwen(context.Background(), ZhiwenRequest{
		Muzhi:     "0",
		Shizhi:    "0",
		Zhongzhi:  "0",
		Wumingzhi: "0",
		Xiaozhi:   "0",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 503") {
		t.Fatalf("expected http 503 error, got: %v", err)
	}
}

func TestDivinationYaogua_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"id":21,"common_desc1":"火雷噬嗑","common_desc2":"象曰","common_desc3":"解卦","shiye":"事业","jingshang":"经商","qiuming":"求名","waichu":"外出","hunlian":"婚恋","juece":"决策","image":"https://yuanfenju.com/Public/img/zhouyi64gua/21.jpg"}}`
	client := newTestClient(t, "/v1/Zhanbu/yaogua", http.StatusOK, body)
	resp, err := client.Divination.Yaogua(context.Background(), YaoguaRequest{Lang: "zh-cn"})
	if err != nil {
		t.Fatalf("yaogua failed: %v", err)
	}
	if resp.Data.ID != 21 || resp.Data.CommonDesc1 != "火雷噬嗑" {
		t.Fatalf("unexpected yaogua response: %#v", resp.Data)
	}
}

func TestDivinationYaogua_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/yaogua", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Yaogua(context.Background(), YaoguaRequest{Lang: "zh-cn"})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestDivinationYaogua_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/yaogua", http.StatusServiceUnavailable, `unavailable`)
	_, err := client.Divination.Yaogua(context.Background(), YaoguaRequest{Lang: "zh-cn"})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 503") {
		t.Fatalf("expected http 503 error, got: %v", err)
	}
}

func TestDivinationTaluojiedu_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"cards":[{"positions_index":1,"positions_name":"选项A","card_name":"隐者","card_interpretation":{"general":"g"}}],"overall_interpretation":{"summary_message":"summary","oracle_message":"oracle"},"environment":{"calculation_time":"2025-12-18 14:58:09","time_ganzhi":"乙未","time_element":"土局"}}}`
	client := newTestClient(t, "/v1/Zhanbu/taluojiedu", http.StatusOK, body)
	resp, err := client.Divination.Taluojiedu(context.Background(), TaluojieduRequest{
		SpreadID: "2",
		TopicID:  "1",
		Lang:     "zh-cn",
	})
	if err != nil {
		t.Fatalf("taluojiedu failed: %v", err)
	}
	if len(resp.Data.Cards) != 1 || resp.Data.Cards[0].CardName != "隐者" || resp.Data.OverallInterpretation.OracleMessage != "oracle" {
		t.Fatalf("unexpected taluojiedu response: %#v", resp.Data)
	}
}

func TestDivinationTaluojiedu_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/taluojiedu", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Taluojiedu(context.Background(), TaluojieduRequest{
		SpreadID: "2",
		TopicID:  "1",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestDivinationTaluojiedu_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/taluojiedu", http.StatusBadGateway, `bad gateway`)
	_, err := client.Divination.Taluojiedu(context.Background(), TaluojieduRequest{
		SpreadID: "2",
		TopicID:  "1",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Fatalf("expected http 502 error, got: %v", err)
	}
}

func TestDivinationTaluoxipai_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"0":21,"1":10,"2":14,"image":"https://yuanfenju.com/Public/img/taluo/back.jpg"}}`
	client := newTestClient(t, "/v1/Zhanbu/taluoxipai", http.StatusOK, body)
	resp, err := client.Divination.Taluoxipai(context.Background(), TaluoxipaiRequest{TaluoSpreads: "3"})
	if err != nil {
		t.Fatalf("taluoxipai failed: %v", err)
	}
	if len(resp.Data.CardNos) != 3 || resp.Data.CardNos[0] != 21 || resp.Data.Image == "" {
		t.Fatalf("unexpected taluoxipai response: %#v", resp.Data)
	}
}

func TestDivinationTaluoxipai_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/taluoxipai", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Taluoxipai(context.Background(), TaluoxipaiRequest{TaluoSpreads: "3"})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestDivinationTaluoxipai_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/taluoxipai", http.StatusBadGateway, `bad gateway`)
	_, err := client.Divination.Taluoxipai(context.Background(), TaluoxipaiRequest{TaluoSpreads: "3"})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Fatalf("expected http 502 error, got: %v", err)
	}
}

func TestDivinationTaluozhanbu_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"牌名":"月亮 The Moon","关键字":"不安","星相":"双鱼座","四要素":"水","牌面描述":"desc","正位含义":{"基本含义":"a","建议":"b"},"逆位含义":{"基本含义":"c","建议":"d"},"含义":{"基本含义":"e","建议":"f"},"正逆":"正位","id":19,"image":"https://yuanfenju.com/Public/img/taluo/18.jpg"}}`
	client := newTestClient(t, "/v1/Zhanbu/taluozhanbu", http.StatusOK, body)
	resp, err := client.Divination.Taluozhanbu(context.Background(), TaluozhanbuRequest{
		TaluoInverse: "0",
		Lang:         "zh-cn",
	})
	if err != nil {
		t.Fatalf("taluozhanbu failed: %v", err)
	}
	if resp.Data.CardName != "月亮 The Moon" || resp.Data.Orientation != "正位" || resp.Data.Meaning.Advice != "f" {
		t.Fatalf("unexpected taluozhanbu response: %#v", resp.Data)
	}
}

func TestDivinationTaluozhanbu_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/taluozhanbu", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Taluozhanbu(context.Background(), TaluozhanbuRequest{
		TaluoInverse: "0",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestDivinationTaluozhanbu_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/taluozhanbu", http.StatusServiceUnavailable, `unavailable`)
	_, err := client.Divination.Taluozhanbu(context.Background(), TaluozhanbuRequest{
		TaluoInverse: "0",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 503") {
		t.Fatalf("expected http 503 error, got: %v", err)
	}
}

func TestDivinationTaluospreads_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":[{"position":"第1号位牌信息","image":"https://yuanfenju.com/Public/img/taluo/0.jpg","card_info":{"cart_reverse":"逆位","card_description":{"base_desc":"a","advice":"b"},"card_name":"愚人 The Fool","card_keyword":"开始","card_astrology":"天王星","card_elements":"风","card_summarize":"desc"}},{"position":"第2号位牌信息","image":"https://yuanfenju.com/Public/img/taluo/19.jpg","card_info":{"cart_reverse":"正位","card_description":{"base_desc":"c","advice":"d"},"card_name":"太阳 The Sun","card_keyword":"希望","card_astrology":"太阳","card_elements":"火","card_summarize":"desc2"}}]}`
	client := newTestClient(t, "/v1/Zhanbu/taluospreads", http.StatusOK, body)
	resp, err := client.Divination.Taluospreads(context.Background(), TaluospreadsRequest{
		TaluoSpreads:     "2",
		TaluoUserChecked: "1,20",
		Lang:             "zh-cn",
	})
	if err != nil {
		t.Fatalf("taluospreads failed: %v", err)
	}
	if len(resp.Data) != 2 || resp.Data[0].CardInfo.CardName != "愚人 The Fool" || resp.Data[1].CardInfo.CardDescription.Advice != "d" {
		t.Fatalf("unexpected taluospreads response: %#v", resp.Data)
	}
}

func TestDivinationTaluospreads_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/taluospreads", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":[]}`)
	_, err := client.Divination.Taluospreads(context.Background(), TaluospreadsRequest{
		TaluoSpreads:     "2",
		TaluoUserChecked: "1,20",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestDivinationTaluospreads_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/taluospreads", http.StatusBadRequest, `bad request`)
	_, err := client.Divination.Taluospreads(context.Background(), TaluospreadsRequest{
		TaluoSpreads:     "2",
		TaluoUserChecked: "1,20",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 400") {
		t.Fatalf("expected http 400 error, got: %v", err)
	}
}

func TestDivinationYunshi_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"运势类型":"白羊座","今日运势":{"速配星座":"处女座","综合分数":"78","今明运势":"today"},"明日运势":{"速配星座":"双鱼座","今明运势":"tomorrow"},"本周运势":{"速配星座":"巨蟹座","本周运势":"week"},"本月运势":{"速配星座":"摩羯座","本月运势":"month"},"本年运势":{"速配星座":"巨蟹座","本年运势":"year"}}}`
	client := newTestClient(t, "/v1/Zhanbu/yunshi", http.StatusOK, body)
	resp, err := client.Divination.Yunshi(context.Background(), DivinationYunshiRequest{
		Type:        "0",
		TitleYunshi: "3",
		Lang:        "zh-cn",
	})
	if err != nil {
		t.Fatalf("divination yunshi failed: %v", err)
	}
	if resp.Data.FortuneType != "白羊座" || resp.Data.DailyFortune.CompatibleSign != "处女座" || resp.Data.WeeklyFortune.WeeklyFortune != "week" {
		t.Fatalf("unexpected divination yunshi response: %#v", resp.Data)
	}
}

func TestDivinationYunshi_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/yunshi", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Yunshi(context.Background(), DivinationYunshiRequest{
		Type:        "0",
		TitleYunshi: "3",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestDivinationYunshi_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/yunshi", http.StatusServiceUnavailable, `unavailable`)
	_, err := client.Divination.Yunshi(context.Background(), DivinationYunshiRequest{
		Type:        "0",
		TitleYunshi: "3",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 503") {
		t.Fatalf("expected http 503 error, got: %v", err)
	}
}

func TestDivinationShengxiaoyunshi_Success(t *testing.T) {
	body := `{"errcode":0,"errmsg":"ok","data":{"运势类型":"属鼠","今日运势":{"速配生肖":"猪","综合分数":"78","今明运势":"today"},"明日运势":{"速配生肖":"龙","今明运势":"tomorrow"},"本周运势":{"速配生肖":"猪","本周运势":"week"},"本月运势":{"速配生肖":"狗","本月运势":"month"},"本年运势":{"速配生肖":"狗","本年运势":"year"}}}`
	client := newTestClient(t, "/v1/Zhanbu/yunshi", http.StatusOK, body)
	resp, err := client.Divination.Shengxiaoyunshi(context.Background(), ShengxiaoyunshiRequest{
		TitleYunshi: "3",
		Lang:        "zh-cn",
	})
	if err != nil {
		t.Fatalf("shengxiaoyunshi failed: %v", err)
	}
	if resp.Data.FortuneType != "属鼠" || resp.Data.DailyFortune.CompatibleZodiac != "猪" || resp.Data.MonthlyFortune.MonthlyFortune != "month" {
		t.Fatalf("unexpected shengxiaoyunshi response: %#v", resp.Data)
	}
}

func TestDivinationShengxiaoyunshi_APIError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/yunshi", http.StatusOK, `{"errcode":-1,"errmsg":"bad request","notice":"n","data":{}}`)
	_, err := client.Divination.Shengxiaoyunshi(context.Background(), ShengxiaoyunshiRequest{
		TitleYunshi: "3",
	})
	if err == nil {
		t.Fatal("expected API error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %v", err)
	}
}

func TestDivinationShengxiaoyunshi_HTTPError(t *testing.T) {
	client := newTestClient(t, "/v1/Zhanbu/yunshi", http.StatusBadRequest, `bad request`)
	_, err := client.Divination.Shengxiaoyunshi(context.Background(), ShengxiaoyunshiRequest{
		TitleYunshi: "3",
	})
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "http 400") {
		t.Fatalf("expected http 400 error, got: %v", err)
	}
}
