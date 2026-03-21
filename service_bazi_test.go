package yuanfenju

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"
	"time"
)

func TestBaziPaipanDataUnmarshal(t *testing.T) {
	raw := `{
  "errcode": 0,
  "errmsg": "ok",
  "data": {
    "base_info": {
      "zhen": {
        "province": "a",
        "city": "b",
        "jingdu": "c",
        "weidu": "d",
        "shicha": "e"
      },
      "sex": "坤造",
      "name": "测试",
      "gongli": "1990年1月1日12时0分",
      "nongli": "己巳年 十二月 初五日 午时",
      "qiyun": "1年5月22天起运",
      "jiaoyun": "1991年6月18日6时31分58秒",
      "zhengge": "正官格"
    },
    "bazi_info": {
      "kw": "戌亥",
      "tg_cg_god": ["伤官"],
      "bazi": ["己巳"],
      "dz_cg": ["丙|庚|戊"],
      "dz_cg_god": ["比肩|偏财|食神"],
      "day_cs": ["临官"],
      "na_yin": ["大林木"]
    },
    "dayun_info": {
      "big_god": ["劫财"],
      "big": ["丁丑"],
      "big_cs": ["养"],
      "xu_sui": [2],
      "big_start_year": [1991],
      "big_start_year_liu_nian": "",
      "big_end_year": [2000],
      "years_info0": [{"year_char":"庚午"}],
      "years_info1": [{"year_char":"辛未"}],
      "years_info2": [{"year_char":"壬申"}],
      "years_info3": [{"year_char":"癸酉"}],
      "years_info4": [{"year_char":"甲戌"}],
      "years_info5": [{"year_char":"乙亥"}],
      "years_info6": [{"year_char":"丙子"}],
      "years_info7": [{"year_char":"丁丑"}],
      "years_info8": [{"year_char":"戊寅"}],
      "years_info9": [{"year_char":"己卯"}]
    },
    "start_info": {
      "jishen": ["天德贵人"],
      "xz": "摩羯座",
      "sx": "蛇"
    },
    "detail_info": {
      "zhuxing": {"year":"伤官","month":"比肩","day":"日元","hour":"偏印"},
      "sizhu": {
        "year":{"tg":"己","dz":"巳"},
        "month":{"tg":"丙","dz":"子"},
        "day":{"tg":"丙","dz":"寅"},
        "hour":{"tg":"甲","dz":"午"}
      },
      "canggan": {"year":["丙"],"month":["癸"],"day":["甲"],"hour":["丁"]},
      "fuxing": {"year":["比肩"],"month":["正官"],"day":["偏印"],"hour":["劫财"]},
      "xingyun": {"year":"临官","month":"胎","day":"长生","hour":"帝旺"},
      "zizuo": {"year":"帝旺","month":"胎","day":"长生","hour":"死"},
      "kongwang": {"year":"戌亥","month":"申酉","day":"戌亥","hour":"辰巳"},
      "nayin": {"year":"大林木","month":"涧下水","day":"炉中火","hour":"沙中金"},
      "shensha": {"year":"甲","month":"乙","day":"丙","hour":"丁"},
      "dayunshensha": [{"tgdz":"丁丑","shensha":"月德合"}]
    }
  }
}`

	var resp CommonResponse[BaziPaipanData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal bazi data failed: %v", err)
	}

	if resp.Data.BaseInfo.Zhengge != "正官格" {
		t.Fatalf("unexpected zhengge: %s", resp.Data.BaseInfo.Zhengge)
	}
	if resp.Data.BaziInfo.KW != "戌亥" {
		t.Fatalf("unexpected kw: %s", resp.Data.BaziInfo.KW)
	}
	if len(resp.Data.DayunInfo.YearsInfo9) != 1 || resp.Data.DayunInfo.YearsInfo9[0].YearChar != "己卯" {
		t.Fatalf("unexpected years_info9: %#v", resp.Data.DayunInfo.YearsInfo9)
	}
	if resp.Data.DetailInfo.Sizhu.Day.TG != "丙" || resp.Data.DetailInfo.DayunShensha[0].TGDZ != "丁丑" {
		t.Fatalf("unexpected detail_info: %#v", resp.Data.DetailInfo)
	}
}

func TestBaziPaipanRequestValidate(t *testing.T) {
	valid := BaziPaipanRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "1990",
		Month:  "01",
		Day:    "01",
		Hours:  "12",
		Zhen:   "0",
		Lang:   "zh-cn",
		Minute: "00",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid request, got err: %v", err)
	}

	cases := []struct {
		name string
		req  BaziPaipanRequest
	}{
		{
			name: "missing sex",
			req: BaziPaipanRequest{
				Type: "1", Year: "1990", Month: "01", Day: "01", Hours: "12",
			},
		},
		{
			name: "invalid sex",
			req: BaziPaipanRequest{
				Sex: "2", Type: "1", Year: "1990", Month: "01", Day: "01", Hours: "12",
			},
		},
		{
			name: "invalid lang",
			req: BaziPaipanRequest{
				Sex: "1", Type: "1", Year: "1990", Month: "01", Day: "01", Hours: "12", Lang: "fr-fr",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("expected ErrValidation, got: %v", err)
			}
		})
	}
}

func TestBaziJiuxingRequestValidate(t *testing.T) {
	okReq := BaziJiuxingRequest{
		Sex:    "0",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
		Lang:   "zh-cn",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := BaziJiuxingRequest{
		Sex:    "2",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestBaziJiuxingDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"base_info":{"sex":"坤造","name":"张三","gongli":"a","nongli":"b","qiyun":"c","jiaoyun":"d"},"jiuxing":{"风水命":"四绿木","九星":"文曲星","论命":"x","特性":"y","机会":"z","忠告":"w"}}}`
	var resp CommonResponse[BaziJiuxingData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal jiuxing data failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.Jiuxing.Jiuxing != "文曲星" {
		t.Fatalf("unexpected jiuxing data: %#v", resp.Data)
	}
}

func TestBaziHehunRequestValidate(t *testing.T) {
	okReq := BaziHehunRequest{
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
		Lang:         "zh-cn",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := BaziHehunRequest{
		MaleType:     "2",
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
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestBaziHehunDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"male":{"name":"男方","bazi":["甲子"]},"female":{"name":"女方","bazi":["乙丑"]},"minggong":{"male_fengshui":"东四命","female_fengshui":"东四命","score":"30","male_minggong":"震","female_minggong":"震","description":"d","detail_description":"dd"},"nianqitongzhi":{"score":"20","male_nian_zhi":"辰","male_nian_zhi_desc":"木","female_nian_zhi":"辰","female_nian_zhi_desc":"木","description":"d","detail_description":"dd"},"yueling":{"score":"5","male_yue_zhi":"亥","female_yue_zhi":"亥","description":"d","detail_description":"dd"},"rigan":{"score":"25","male_yue_zhi":"丁","female_yue_zhi":"丁","description":"d","detail_description":"dd"},"tiangan":{"score":"5","male_yue_zhi":"丁","female_yue_zhi":"丁","description":"d","detail_description":"dd"},"zinv":{"nannv":"男","score":"15","description":"d","detail_description":"dd"},"all_score":100,"male_sx":"龙","female_sx":"龙","male_xz":"天蝎座","female_xz":"天蝎座"}}`
	var resp CommonResponse[BaziHehunData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal hehun data failed: %v", err)
	}
	if resp.Data.AllScore != 100 || resp.Data.Male.Name != "男方" || resp.Data.Zinv.Nannv != "男" {
		t.Fatalf("unexpected hehun data: %#v", resp.Data)
	}
}

func TestBaziHepanRequestValidate(t *testing.T) {
	okReq := BaziHepanRequest{
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
		Lang:         "zh-cn",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := BaziHepanRequest{
		MaleType:     "1",
		MaleYear:     "1988",
		MaleMonth:    "11",
		MaleDay:      "8",
		MaleHours:    "12",
		MaleMinute:   "20",
		FemaleType:   "2",
		FemaleYear:   "1988",
		FemaleMonth:  "11",
		FemaleDay:    "8",
		FemaleHours:  "12",
		FemaleMinute: "20",
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestBaziHepanDataUnmarshal(t *testing.T) {
	raw := `{"errcode":0,"errmsg":"ok","data":{"male":{"name":"甲方","bazi":["甲子"]},"female":{"name":"乙方","bazi":["乙丑"]},"minggong":{"score":"10"},"nianqitongzhi":{"score":"20"},"yueling":{"score":"5"},"rigan":{"score":"25"},"tiangan":{"score":"5"},"jiankang":{"score":"10","description":"d","detail_description":"dd"},"all_score":100,"male_sx":"龙","female_sx":"龙","male_xz":"天蝎座","female_xz":"天蝎座"}}`
	var resp CommonResponse[BaziHepanData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal hepan data failed: %v", err)
	}
	if resp.Data.AllScore != 100 || resp.Data.Male.Name != "甲方" || resp.Data.Jiankang.Score != "10" {
		t.Fatalf("unexpected hepan data: %#v", resp.Data)
	}
}

func TestBaziCesuanRequestValidate(t *testing.T) {
	okReq := BaziCesuanRequest{
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
		Lang:      "zh-cn",
		Factor:    "1",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := BaziCesuanRequest{
		Sex:       "0",
		Type:      "1",
		Year:      "1988",
		Month:     "11",
		Day:       "8",
		Hours:     "12",
		Minute:    "20",
		Zhen:      "3",
		Longitude: "181",
		Latitude:  "39.92",
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestBaziCesuanDataUnmarshal(t *testing.T) {
	raw := `{
		"errcode":0,
		"errmsg":"ok",
		"data":{
			"base_info":{"sex":"坤造","name":"张三","gongli":"a","nongli":"b","qiyun":"c","jiaoyun":"d","zhengge":"伤官格","wuxing_xiji":"喜水"},
			"bazi_info":{"kw":"戌亥","tg_cg_god":["伤","杀"],"bazi":"戊辰 癸亥 丁卯 丙午","na_yin":"大林木"},
			"chenggu":{"year_weight":"1.2","month_weight":"1.8","day_weight":"1.6","hour_weight":"1.0","total_weight":"5.6","description":"desc"},
			"wuxing":{"detail_desc":"先天","simple_desc":"木","simple_description":"后天","detail_description":"详细"},
			"yinyuan":{"sanshishu_yinyuan":"姻缘文案"},
			"caiyun":{"sanshishu_caiyun":{"simple_desc":"先苦后甜","detail_desc":"财运详批"}},
			"sizhu":{"rizhu":"日柱论命"},
			"mingyun":{"sanshishu_mingyun":"命运批示"},
			"sx":"龙",
			"xz":"天蝎座",
			"xiyongshen":{
				"qiangruo":"八字偏弱",
				"xiyongshen":"金，水",
				"jishen":"水",
				"xiyongshen_desc":"说明",
				"jin_number":0,
				"mu_number":1,
				"shui_number":2,
				"huo_number":3,
				"tu_number":2,
				"tonglei":"金水",
				"yilei":"木火土",
				"rizhu_tiangan":"水",
				"zidang":0,
				"yidang":9,
				"zidang_percent":"0%",
				"yidang_percent":"100%",
				"jin_score":75,
				"mu_score":16,
				"shui_score":118,
				"huo_score":36,
				"tu_score":139,
				"jin_score_percent":"19.53%",
				"mu_score_percent":"4.17%",
				"shui_score_percent":"30.73%",
				"huo_score_percent":"9.38%",
				"tu_score_percent":"36.2%",
				"yinyang":"阴阳平衡"
			}
		}
	}`
	var resp CommonResponse[BaziCesuanData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal cesuan data failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.Caiyun.SanshishuCaiyun.SimpleDesc != "先苦后甜" || resp.Data.Xiyongshen.ShuiScore != 118 {
		t.Fatalf("unexpected cesuan data: %#v", resp.Data)
	}
}

func TestBaziJingpanRequestValidate(t *testing.T) {
	okReq := BaziJingpanRequest{
		Sex:        "0",
		Type:       "1",
		Year:       "1994",
		Month:      "4",
		Day:        "30",
		Hours:      "10",
		Minute:     "0",
		LoadMode:   "2",
		DayunIndex: "2",
		Zhen:       "3",
		Longitude:  "116.46",
		Latitude:   "39.92",
		Lang:       "zh-cn",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := BaziJingpanRequest{
		Sex:    "0",
		Type:   "1",
		Year:   "1994",
		Month:  "4",
		Day:    "30",
		Hours:  "10",
		Minute: "0",
		IsZip:  "3",
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestBaziJingpanDataUnmarshal(t *testing.T) {
	raw := `{
		"errcode":0,
		"errmsg":"ok",
		"data":{
			"base_info":{
				"sex":"乾造",
				"name":"张三",
				"gongli":"1994-04-30 10:00:00",
				"nongli":"一九九四年三月二十日 巳时",
				"zhengge":"伤官格",
				"minggua":{"minggua_name":"乾","minggua_fangwei":"西四命"},
				"xiyongshen":{"shui_score":118}
			},
			"detail_info":{
				"dayun_info":[{"dayun_index":1,"dayun_start_year":1996,"dayun_shensha":"禄神","liunian_info":[{"liunian_index":0,"liunian_year":1996}]}],
				"sizhu_info":{"year":{"tg_god":"偏印","tg":"甲","dz":"戌"},"month":{"tg_god":"食神"},"day":{"tg_god":"日主"},"hour":{"tg_god":"正官"}},
				"taishen_info":{"taiyuan":{"ganzhi":"己未"}}
			}
		}
	}`
	var resp CommonResponse[BaziJingpanData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal jingpan data failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.DetailInfo.DayunInfo[0].DayunIndex != 1 || resp.Data.DetailInfo.SizhuInfo.Year.TG != "甲" {
		t.Fatalf("unexpected jingpan data: %#v", resp.Data)
	}
}

func TestBaziJingsuanRequestValidate(t *testing.T) {
	okReq := BaziJingsuanRequest{
		Sex:       "1",
		Type:      "1",
		Year:      "1988",
		Month:     "11",
		Day:       "8",
		Hours:     "12",
		Minute:    "20",
		Zhen:      "3",
		Longitude: "116.46",
		Latitude:  "39.92",
		Lang:      "zh-cn",
		Factor:    "1",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := BaziJingsuanRequest{
		Sex:    "1",
		Type:   "1",
		Year:   "1988",
		Month:  "11",
		Day:    "8",
		Hours:  "12",
		Minute: "20",
		Lang:   "fr-fr",
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestBaziJingsuanDataUnmarshal(t *testing.T) {
	raw := `{
		"errcode":0,
		"errmsg":"ok",
		"data":{
			"base_info":{"name":"张三","sex":"乾造","gongli":"a","nongli":"b"},
			"detail_info":{
				"dayun_info":[{"dayun_index":1,"dayun_shensha":"禄神","dayun_indication":{"shiye":"事业"}}],
				"sizhu_info":{
					"year":{"tg_god":"伤官","tg":"戊","dz":"辰"},
					"month":{"tg_god":"七杀"},
					"day":{"tg_god":"日主"},
					"hour":{"tg_god":"偏印"},
					"sizhu_indication":{
						"chenggu":{"total_weight":"3.8"},
						"wuxing":{"simple_desc":"木"},
						"yinyuan":{"sanshishu_yinyuan":"夫妻和合"},
						"caiyun":{"sanshishu_caiyun":{"simple_desc":"知足常乐"}},
						"xingge":{"rizhu":"癸亥日柱"},
						"mingyun":{"sanshishu_mingyun":"命运批示"}
					}
				}
			}
		}
	}`
	var resp CommonResponse[BaziJingsuanData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal jingsuan data failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.DetailInfo.SizhuInfo.SizhuIndication.Caiyun.SanshishuCaiyun.SimpleDesc != "知足常乐" || resp.Data.DetailInfo.DayunInfo[0].DayunIndication.Shiye != "事业" {
		t.Fatalf("unexpected jingsuan data: %#v", resp.Data)
	}
}

func TestBaziWeilaiRequestValidate(t *testing.T) {
	currentYear := time.Now().Year()
	okReq := BaziWeilaiRequest{
		Sex:        "1",
		Type:       "1",
		Year:       "1988",
		Month:      "11",
		Day:        "8",
		Hours:      "12",
		Minute:     "20",
		YunshiYear: strconv.Itoa(currentYear),
		Zhen:       "3",
		Longitude:  "116.46",
		Latitude:   "39.92",
		Lang:       "zh-cn",
	}
	if err := okReq.Validate(); err != nil {
		t.Fatalf("expected valid request, got: %v", err)
	}

	badReq := BaziWeilaiRequest{
		Sex:        "1",
		Type:       "1",
		Year:       "1988",
		Month:      "11",
		Day:        "8",
		Hours:      "12",
		Minute:     "20",
		YunshiYear: strconv.Itoa(currentYear - 1),
	}
	err := badReq.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestBaziWeilaiDataUnmarshal(t *testing.T) {
	raw := `{
		"errcode":0,
		"errmsg":"ok",
		"data":{
			"base_info":{"name":"张三","sex":"乾造","gongli":"a","nongli":"b"},
			"detail_info":{
				"sizhu_info":{
					"year":{"tg_god":"伤官","tg":"戊","dz":"辰"},
					"month":{"tg_god":"七杀"},
					"day":{"tg_god":"日主"},
					"hour":{"tg_god":"偏印"}
				},
				"yunshi_year_info":{
					"yunshi_year":{"year":2026,"tg_god":"伤官","tg":"甲","dz":"辰","indication":{"shiye":"事业预测"}}
				},
				"yunshi_month_info":[
					{"month":"1月","tg_god":"正官","tg":"乙","dz":"丑","indication":{"caiyun":"月财运"},"yunshi_day_info":[{"day":"1号","tg_god":"日主","tg":"丙","dz":"申","indication":{"yunshi":"日运势"}}]}
				]
			}
		}
	}`
	var resp CommonResponse[BaziWeilaiData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal weilai data failed: %v", err)
	}
	if resp.Data.BaseInfo.Name != "张三" || resp.Data.DetailInfo.YunshiYearInfo.YunshiYear.Year != 2026 || resp.Data.DetailInfo.YunshiMonthInfo[0].YunshiDayInfo[0].Day != "1号" {
		t.Fatalf("unexpected weilai data: %#v", resp.Data)
	}
}
