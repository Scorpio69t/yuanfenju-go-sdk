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

func TestZhiwenRequestValidate(t *testing.T) {
	okReq := ZhiwenRequest{
		Muzhi:     "0",
		Shizhi:    "1",
		Zhongzhi:  "0",
		Wumingzhi: "1",
		Xiaozhi:   "0",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := ZhiwenRequest{
		Muzhi:     "2",
		Shizhi:    "1",
		Zhongzhi:  "0",
		Wumingzhi: "1",
		Xiaozhi:   "0",
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestZhiwenDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"muzhi":"箩纹","shizhi":"箩纹","zhongzhi":"箩纹","wumingzhi":"箩纹","xiaozhi":"箩纹","description":{"分析":"a","诗曰":"b","性格":"c","婚姻":"d","职业":"e","健康":"f","运势":"g"}}}`
	var resp CommonResponse[ZhiwenData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal zhiwen data failed: %v", err)
	}
	if resp.Data.Muzhi != "箩纹" || resp.Data.Description.Fenxi != "a" || resp.Data.Description.Yunshi != "g" {
		t.Fatalf("unexpected zhiwen data: %#v", resp.Data)
	}
}

func TestYaoguaRequestValidate(t *testing.T) {
	if err := (YaoguaRequest{Lang: "zh-tw"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	err := (YaoguaRequest{Lang: "fr-fr"}).Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestYaoguaDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"id":21,"common_desc1":"火雷噬嗑","common_desc2":"象曰","common_desc3":"解卦","shiye":"事业","jingshang":"经商","qiuming":"求名","waichu":"外出","hunlian":"婚恋","juece":"决策","image":"https://yuanfenju.com/Public/img/zhouyi64gua/21.jpg"}}`
	var resp CommonResponse[YaoguaData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal yaogua data failed: %v", err)
	}
	if resp.Data.ID != 21 || resp.Data.CommonDesc1 != "火雷噬嗑" || resp.Data.Image == "" {
		t.Fatalf("unexpected yaogua data: %#v", resp.Data)
	}
}

func TestTaluojieduRequestValidate(t *testing.T) {
	okReq := TaluojieduRequest{SpreadID: "2", TopicID: "1", Lang: "zh-cn"}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := TaluojieduRequest{SpreadID: "18", TopicID: "1"}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestTaluojieduDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"cards":[{"positions_index":1,"positions_name":"选项A","positions_desc":"desc","orientation_code":1,"orientation_text":"正位","card_no":10,"card_name":"隐者","card_keywords":"内省","card_astrology":"处女座","card_element":"土","card_description":"牌面","card_interpretation":{"general":"g","topic":"t","advice":"a"},"image_id":10,"image_url":"https://img"}],"overall_interpretation":{"summary_message":"summary","oracle_message":"oracle"},"environment":{"calculation_time":"2025-12-18 14:58:09","time_ganzhi":"乙未","time_element":"土局"}}}`
	var resp CommonResponse[TaluojieduData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal taluojiedu data failed: %v", err)
	}
	if len(resp.Data.Cards) != 1 || resp.Data.Cards[0].CardName != "隐者" || resp.Data.OverallInterpretation.OracleMessage != "oracle" {
		t.Fatalf("unexpected taluojiedu data: %#v", resp.Data)
	}
}

func TestTaluoxipaiRequestValidate(t *testing.T) {
	if err := (TaluoxipaiRequest{TaluoSpreads: "3"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	err := (TaluoxipaiRequest{TaluoSpreads: "10"}).Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestTaluoxipaiDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"0":21,"1":10,"2":14,"image":"https://yuanfenju.com/Public/img/taluo/back.jpg"}}`
	var resp CommonResponse[TaluoxipaiData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal taluoxipai data failed: %v", err)
	}
	if len(resp.Data.CardNos) != 3 || resp.Data.CardNos[0] != 21 || resp.Data.CardNos[2] != 14 {
		t.Fatalf("unexpected card_nos: %#v", resp.Data.CardNos)
	}
	if resp.Data.Image == "" {
		t.Fatalf("unexpected image: %s", resp.Data.Image)
	}
}

func TestTaluozhanbuRequestValidate(t *testing.T) {
	okReq := TaluozhanbuRequest{TaluoInverse: "0", Lang: "zh-cn"}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := TaluozhanbuRequest{TaluoInverse: "2", Lang: "zh-cn"}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestTaluozhanbuDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"牌名":"月亮 The Moon","关键字":"不安","星相":"双鱼座","四要素":"水","牌面描述":"desc","正位含义":{"基本含义":"a","建议":"b"},"逆位含义":{"基本含义":"c","建议":"d"},"含义":{"基本含义":"e","建议":"f"},"正逆":"正位","id":19,"image":"https://yuanfenju.com/Public/img/taluo/18.jpg"}}`
	var resp CommonResponse[TaluozhanbuData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal taluozhanbu data failed: %v", err)
	}
	if resp.Data.CardName != "月亮 The Moon" || resp.Data.Orientation != "正位" || resp.Data.Meaning.Advice != "f" {
		t.Fatalf("unexpected taluozhanbu data: %#v", resp.Data)
	}
}

func TestTaluospreadsRequestValidate(t *testing.T) {
	okReq := TaluospreadsRequest{
		TaluoSpreads:     "3",
		TaluoUserChecked: "1,20,9",
		Lang:             "zh-cn",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := TaluospreadsRequest{
		TaluoSpreads:     "3",
		TaluoUserChecked: "1,20",
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestTaluospreadsDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":[{"position":"第1号位牌信息","image":"https://yuanfenju.com/Public/img/taluo/0.jpg","card_info":{"cart_reverse":"逆位","card_description":{"base_desc":"a","advice":"b"},"card_name":"愚人 The Fool","card_keyword":"开始","card_astrology":"天王星","card_elements":"风","card_summarize":"desc"}}]}`
	var resp CommonResponse[[]TaluospreadsData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal taluospreads data failed: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].CardInfo.CardName != "愚人 The Fool" || resp.Data[0].CardInfo.CardDescription.Advice != "b" {
		t.Fatalf("unexpected taluospreads data: %#v", resp.Data)
	}
}

func TestDivinationYunshiRequestValidate(t *testing.T) {
	okReq := DivinationYunshiRequest{Type: "0", TitleYunshi: "3", Lang: "zh-cn", ParameterStyle: "chinese"}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := DivinationYunshiRequest{Type: "2", TitleYunshi: "3"}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestShengxiaoyunshiRequestValidate(t *testing.T) {
	okReq := ShengxiaoyunshiRequest{TitleYunshi: "3", Lang: "en-us", ParameterStyle: "english"}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := ShengxiaoyunshiRequest{TitleYunshi: "12"}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestDivinationYunshiDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"运势类型":"白羊座","今日运势":{"速配星座":"处女座","提防星座":"狮子座","幸运颜色":"黄色","幸运数字":"73","幸运宝石":"草绿石","综合分数":"78","爱情分数":"65","事业分数":"85","心情分数":"73","交际分数":"81","财富分数":"76","健康分数":"87","今明运势":"a","爱情运势":"b","事业运势":"c","财富运势":"d","健康运势":"e"},"明日运势":{"速配星座":"双鱼座","提防星座":"摩羯座","今明运势":"x"},"本周运势":{"速配星座":"巨蟹座","本周运势":"w"},"本月运势":{"速配星座":"摩羯座","本月运势":"m"},"本年运势":{"速配星座":"巨蟹座","本年运势":"y"}}}`
	var resp CommonResponse[DivinationYunshiData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal divination yunshi data failed: %v", err)
	}
	if resp.Data.FortuneType != "白羊座" || resp.Data.DailyFortune.LuckyColor != "黄色" || resp.Data.WeeklyFortune.WeeklyFortune != "w" {
		t.Fatalf("unexpected yunshi data: %#v", resp.Data)
	}
}
