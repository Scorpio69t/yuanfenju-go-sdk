package yuanfenju

import (
	"encoding/json"
	"errors"
	"testing"
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
