package yuanfenju

import (
	"context"
	"net/url"
)

type BaziService struct {
	client *Client
}

// BaziPaipanRequest 对应 /v1/Bazi/paipan。
// 参数命名与官方文档字段保持一致，便于直接对照接口文档。
type BaziPaipanRequest struct {
	Name     string
	Sex      string // 1 男，0 女
	Type     string // 1 公历，2 农历
	Year     string
	Month    string
	Day      string
	Hours    string
	Minute   string
	Zhen     string // 0 不使用，1 国内真太阳时，3 全球真太阳时
	Province string
	City     string
	IANATime string // zhen=3 时建议传入
	Calendar string
	Lang     string // zh-cn / zh-tw / en-us
}

func (r BaziPaipanRequest) toValues() url.Values {
	v := url.Values{}
	if r.Name != "" {
		v.Set("name", r.Name)
	}
	if r.Sex != "" {
		v.Set("sex", r.Sex)
	}
	if r.Type != "" {
		v.Set("type", r.Type)
	}
	if r.Year != "" {
		v.Set("year", r.Year)
	}
	if r.Month != "" {
		v.Set("month", r.Month)
	}
	if r.Day != "" {
		v.Set("day", r.Day)
	}
	if r.Hours != "" {
		v.Set("hours", r.Hours)
	}
	if r.Minute != "" {
		v.Set("minute", r.Minute)
	}
	if r.Zhen != "" {
		v.Set("zhen", r.Zhen)
	}
	if r.Province != "" {
		v.Set("province", r.Province)
	}
	if r.City != "" {
		v.Set("city", r.City)
	}
	if r.IANATime != "" {
		v.Set("iana_time", r.IANATime)
	}
	if r.Calendar != "" {
		v.Set("calendar", r.Calendar)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

type BaziPaipanData struct {
	BaseInfo   BaziBaseInfo   `json:"base_info"`
	BaziInfo   BaziInfo       `json:"bazi_info"`
	DayunInfo  BaziDayunInfo  `json:"dayun_info"`
	StartInfo  BaziStartInfo  `json:"start_info"`
	DetailInfo BaziDetailInfo `json:"detail_info"`
}

type BaziBaseInfo struct {
	Zhen    *BaziZhenInfo `json:"zhen,omitempty"`
	Sex     string        `json:"sex"`
	Name    string        `json:"name"`
	Gongli  string        `json:"gongli"`
	Nongli  string        `json:"nongli"`
	Qiyun   string        `json:"qiyun"`
	Jiaoyun string        `json:"jiaoyun"`
	Zhengge string        `json:"zhengge"`
}

type BaziZhenInfo struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Jingdu   string `json:"jingdu"`
	Weidu    string `json:"weidu"`
	Shicha   string `json:"shicha"`
}

type BaziInfo struct {
	KW      string   `json:"kw"`
	TGCGGod []string `json:"tg_cg_god"`
	Bazi    []string `json:"bazi"`
	DZCG    []string `json:"dz_cg"`
	DZCGGod []string `json:"dz_cg_god"`
	DayCS   []string `json:"day_cs"`
	NaYin   []string `json:"na_yin"`
}

type BaziDayunInfo struct {
	BigGod              []string           `json:"big_god"`
	Big                 []string           `json:"big"`
	BigCS               []string           `json:"big_cs"`
	XuSui               []int              `json:"xu_sui"`
	BigStartYear        []int              `json:"big_start_year"`
	BigStartYearLiuNian string             `json:"big_start_year_liu_nian"`
	BigEndYear          []int              `json:"big_end_year"`
	YearsInfo0          []BaziYearCharInfo `json:"years_info0"`
	YearsInfo1          []BaziYearCharInfo `json:"years_info1"`
	YearsInfo2          []BaziYearCharInfo `json:"years_info2"`
	YearsInfo3          []BaziYearCharInfo `json:"years_info3"`
	YearsInfo4          []BaziYearCharInfo `json:"years_info4"`
	YearsInfo5          []BaziYearCharInfo `json:"years_info5"`
	YearsInfo6          []BaziYearCharInfo `json:"years_info6"`
	YearsInfo7          []BaziYearCharInfo `json:"years_info7"`
	YearsInfo8          []BaziYearCharInfo `json:"years_info8"`
	YearsInfo9          []BaziYearCharInfo `json:"years_info9"`
}

type BaziYearCharInfo struct {
	YearChar string `json:"year_char"`
}

type BaziStartInfo struct {
	Jishen []string `json:"jishen"`
	XZ     string   `json:"xz"`
	SX     string   `json:"sx"`
}

type BaziDetailInfo struct {
	Zhuxing      BaziFourPillarsText    `json:"zhuxing"`
	Sizhu        BaziFourPillarsGanZhi  `json:"sizhu"`
	Canggan      BaziFourPillarsStrings `json:"canggan"`
	Fuxing       BaziFourPillarsStrings `json:"fuxing"`
	Xingyun      BaziFourPillarsText    `json:"xingyun"`
	Zizuo        BaziFourPillarsText    `json:"zizuo"`
	Kongwang     BaziFourPillarsText    `json:"kongwang"`
	Nayin        BaziFourPillarsText    `json:"nayin"`
	Shensha      BaziFourPillarsText    `json:"shensha"`
	DayunShensha []BaziDayunShenshaInfo `json:"dayunshensha"`
}

type BaziFourPillarsText struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
	Hour  string `json:"hour"`
}

type BaziPillarGanZhi struct {
	TG string `json:"tg"`
	DZ string `json:"dz"`
}

type BaziFourPillarsGanZhi struct {
	Year  BaziPillarGanZhi `json:"year"`
	Month BaziPillarGanZhi `json:"month"`
	Day   BaziPillarGanZhi `json:"day"`
	Hour  BaziPillarGanZhi `json:"hour"`
}

type BaziFourPillarsStrings struct {
	Year  []string `json:"year"`
	Month []string `json:"month"`
	Day   []string `json:"day"`
	Hour  []string `json:"hour"`
}

type BaziDayunShenshaInfo struct {
	TGDZ    string `json:"tgdz"`
	Shensha string `json:"shensha"`
}

func (s *BaziService) Paipan(ctx context.Context, req BaziPaipanRequest) (*CommonResponse[BaziPaipanData], error) {
	resp := &CommonResponse[BaziPaipanData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/paipan", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}
