package yuanfenju

import (
	"context"
	"net/url"
)

var (
	baziAllowedSex  = []string{"0", "1"}
	baziAllowedType = []string{"1", "2"}
	baziAllowedZhen = []string{"0", "1", "3"}
	baziAllowedLang = []string{"zh-cn", "zh-tw", "en-us"}

	jiuxingAllowedType = []string{"0", "1"}

	hehunAllowedType = []string{"0", "1"}
	hehunAllowedLang = []string{"zh-cn", "en-us"}
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

type BaziJiuxingRequest struct {
	Name   string
	Sex    string // 0 男，1 女
	Type   string // 0 农历，1 公历
	Year   string
	Month  string
	Day    string
	Hours  string
	Minute string
	Lang   string // zh-cn / zh-tw / en-us
}

func (r BaziJiuxingRequest) toValues() url.Values {
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
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r BaziJiuxingRequest) Validate() error {
	if r.Sex == "" {
		return newRequiredFieldError("sex")
	}
	if !inSet(r.Sex, baziAllowedSex) {
		return newEnumFieldError("sex", r.Sex, baziAllowedSex)
	}

	if r.Type == "" {
		return newRequiredFieldError("type")
	}
	if !inSet(r.Type, jiuxingAllowedType) {
		return newEnumFieldError("type", r.Type, jiuxingAllowedType)
	}

	if r.Year == "" {
		return newRequiredFieldError("year")
	}
	if r.Month == "" {
		return newRequiredFieldError("month")
	}
	if r.Day == "" {
		return newRequiredFieldError("day")
	}
	if r.Hours == "" {
		return newRequiredFieldError("hours")
	}
	if r.Minute == "" {
		return newRequiredFieldError("minute")
	}

	if r.Lang != "" && !inSet(r.Lang, baziAllowedLang) {
		return newEnumFieldError("lang", r.Lang, baziAllowedLang)
	}
	return nil
}

type BaziHehunRequest struct {
	MaleName   string
	MaleType   string // 0 农历，1 公历
	MaleYear   string
	MaleMonth  string
	MaleDay    string
	MaleHours  string
	MaleMinute string

	FemaleName   string
	FemaleType   string // 0 农历，1 公历
	FemaleYear   string
	FemaleMonth  string
	FemaleDay    string
	FemaleHours  string
	FemaleMinute string

	Lang string // zh-cn / en-us
}

func (r BaziHehunRequest) toValues() url.Values {
	v := url.Values{}
	if r.MaleName != "" {
		v.Set("male_name", r.MaleName)
	}
	if r.MaleType != "" {
		v.Set("male_type", r.MaleType)
	}
	if r.MaleYear != "" {
		v.Set("male_year", r.MaleYear)
	}
	if r.MaleMonth != "" {
		v.Set("male_month", r.MaleMonth)
	}
	if r.MaleDay != "" {
		v.Set("male_day", r.MaleDay)
	}
	if r.MaleHours != "" {
		v.Set("male_hours", r.MaleHours)
	}
	if r.MaleMinute != "" {
		v.Set("male_minute", r.MaleMinute)
	}

	if r.FemaleName != "" {
		v.Set("female_name", r.FemaleName)
	}
	if r.FemaleType != "" {
		v.Set("female_type", r.FemaleType)
	}
	if r.FemaleYear != "" {
		v.Set("female_year", r.FemaleYear)
	}
	if r.FemaleMonth != "" {
		v.Set("female_month", r.FemaleMonth)
	}
	if r.FemaleDay != "" {
		v.Set("female_day", r.FemaleDay)
	}
	if r.FemaleHours != "" {
		v.Set("female_hours", r.FemaleHours)
	}
	if r.FemaleMinute != "" {
		v.Set("female_minute", r.FemaleMinute)
	}

	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r BaziHehunRequest) Validate() error {
	if r.MaleType == "" {
		return newRequiredFieldError("male_type")
	}
	if !inSet(r.MaleType, hehunAllowedType) {
		return newEnumFieldError("male_type", r.MaleType, hehunAllowedType)
	}
	if r.FemaleType == "" {
		return newRequiredFieldError("female_type")
	}
	if !inSet(r.FemaleType, hehunAllowedType) {
		return newEnumFieldError("female_type", r.FemaleType, hehunAllowedType)
	}

	required := []struct {
		field string
		value string
	}{
		{"male_year", r.MaleYear},
		{"male_month", r.MaleMonth},
		{"male_day", r.MaleDay},
		{"male_hours", r.MaleHours},
		{"male_minute", r.MaleMinute},
		{"female_year", r.FemaleYear},
		{"female_month", r.FemaleMonth},
		{"female_day", r.FemaleDay},
		{"female_hours", r.FemaleHours},
		{"female_minute", r.FemaleMinute},
	}
	for _, x := range required {
		if x.value == "" {
			return newRequiredFieldError(x.field)
		}
	}

	if r.Lang != "" && !inSet(r.Lang, hehunAllowedLang) {
		return newEnumFieldError("lang", r.Lang, hehunAllowedLang)
	}
	return nil
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

func (r BaziPaipanRequest) Validate() error {
	if r.Sex == "" {
		return newRequiredFieldError("sex")
	}
	if !inSet(r.Sex, baziAllowedSex) {
		return newEnumFieldError("sex", r.Sex, baziAllowedSex)
	}

	if r.Type == "" {
		return newRequiredFieldError("type")
	}
	if !inSet(r.Type, baziAllowedType) {
		return newEnumFieldError("type", r.Type, baziAllowedType)
	}

	if r.Year == "" {
		return newRequiredFieldError("year")
	}
	if r.Month == "" {
		return newRequiredFieldError("month")
	}
	if r.Day == "" {
		return newRequiredFieldError("day")
	}
	if r.Hours == "" {
		return newRequiredFieldError("hours")
	}

	if r.Zhen != "" && !inSet(r.Zhen, baziAllowedZhen) {
		return newEnumFieldError("zhen", r.Zhen, baziAllowedZhen)
	}
	if r.Lang != "" && !inSet(r.Lang, baziAllowedLang) {
		return newEnumFieldError("lang", r.Lang, baziAllowedLang)
	}
	return nil
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

type BaziJiuxingData struct {
	BaseInfo BaziJiuxingBaseInfo `json:"base_info"`
	Jiuxing  BaziJiuxingInfo     `json:"jiuxing"`
}

type BaziJiuxingBaseInfo struct {
	Sex     string `json:"sex"`
	Name    string `json:"name"`
	Gongli  string `json:"gongli"`
	Nongli  string `json:"nongli"`
	Qiyun   string `json:"qiyun"`
	Jiaoyun string `json:"jiaoyun"`
}

type BaziJiuxingInfo struct {
	FengshuiMing string `json:"风水命"`
	Jiuxing      string `json:"九星"`
	Lunming      string `json:"论命"`
	Texing       string `json:"特性"`
	Jihui        string `json:"机会"`
	Zhonggao     string `json:"忠告"`
}

type BaziHehunData struct {
	Male          BaziHehunPersonInfo   `json:"male"`
	Female        BaziHehunPersonInfo   `json:"female"`
	Minggong      BaziHehunMinggongInfo `json:"minggong"`
	Nianqitongzhi BaziHehunNianqiInfo   `json:"nianqitongzhi"`
	Yueling       BaziHehunYuelingInfo  `json:"yueling"`
	Rigan         BaziHehunRiganInfo    `json:"rigan"`
	Tiangan       BaziHehunTianganInfo  `json:"tiangan"`
	Zinv          BaziHehunZinvInfo     `json:"zinv"`
	AllScore      int                   `json:"all_score"`
	MaleSX        string                `json:"male_sx"`
	FemaleSX      string                `json:"female_sx"`
	MaleXZ        string                `json:"male_xz"`
	FemaleXZ      string                `json:"female_xz"`
}

type BaziHehunPersonInfo struct {
	Bazi    []string `json:"bazi"`
	GLYear  string   `json:"gl_year"`
	GLMonth string   `json:"gl_month"`
	GLDay   string   `json:"gl_day"`
	GLHours string   `json:"gl_hours"`
	NLYear  string   `json:"nl_year"`
	NLMonth string   `json:"nl_month"`
	NLDay   string   `json:"nl_day"`
	NLHours string   `json:"nl_hours"`
	Sex     string   `json:"sex"`
	Name    string   `json:"name"`
	Gongli  string   `json:"gongli"`
	Nongli  string   `json:"nongli"`
	KW      string   `json:"kw"`
	TGCGGod []string `json:"tg_cg_god"`
	DZCG    []string `json:"dz_cg"`
	DZCGGod []string `json:"dz_cg_god"`
	DayCS   []string `json:"day_cs"`
	NaYin   []string `json:"na_yin"`
}

type BaziHehunMinggongInfo struct {
	MaleFengshui      string `json:"male_fengshui"`
	FemaleFengshui    string `json:"female_fengshui"`
	Score             string `json:"score"`
	MaleMinggong      string `json:"male_minggong"`
	FemaleMinggong    string `json:"female_minggong"`
	Description       string `json:"description"`
	DetailDescription string `json:"detail_description"`
}

type BaziHehunNianqiInfo struct {
	Score             string `json:"score"`
	MaleNianZhi       string `json:"male_nian_zhi"`
	MaleNianZhiDesc   string `json:"male_nian_zhi_desc"`
	FemaleNianZhi     string `json:"female_nian_zhi"`
	FemaleNianZhiDesc string `json:"female_nian_zhi_desc"`
	Description       string `json:"description"`
	DetailDescription string `json:"detail_description"`
}

type BaziHehunYuelingInfo struct {
	Score             string `json:"score"`
	MaleYueZhi        string `json:"male_yue_zhi"`
	FemaleYueZhi      string `json:"female_yue_zhi"`
	Description       string `json:"description"`
	DetailDescription string `json:"detail_description"`
}

type BaziHehunRiganInfo struct {
	Score             string `json:"score"`
	MaleYueZhi        string `json:"male_yue_zhi"`
	FemaleYueZhi      string `json:"female_yue_zhi"`
	Description       string `json:"description"`
	DetailDescription string `json:"detail_description"`
}

type BaziHehunTianganInfo struct {
	Score             string `json:"score"`
	MaleYueZhi        string `json:"male_yue_zhi"`
	FemaleYueZhi      string `json:"female_yue_zhi"`
	Description       string `json:"description"`
	DetailDescription string `json:"detail_description"`
}

type BaziHehunZinvInfo struct {
	Nannv             string `json:"nannv"`
	Score             string `json:"score"`
	Description       string `json:"description"`
	DetailDescription string `json:"detail_description"`
}

func (s *BaziService) Paipan(ctx context.Context, req BaziPaipanRequest) (*CommonResponse[BaziPaipanData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[BaziPaipanData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/paipan", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *BaziService) Jiuxing(ctx context.Context, req BaziJiuxingRequest) (*CommonResponse[BaziJiuxingData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[BaziJiuxingData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/jiuxing", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *BaziService) Hehun(ctx context.Context, req BaziHehunRequest) (*CommonResponse[BaziHehunData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[BaziHehunData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/hehun", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}
