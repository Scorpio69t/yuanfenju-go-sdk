package yuanfenju

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
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
