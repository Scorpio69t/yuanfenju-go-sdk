package yuanfenju

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

var (
	baziAllowedSex  = []string{"0", "1"}
	baziAllowedType = []string{"1", "2"}
	baziAllowedZhen = []string{"0", "1", "3"}
	baziAllowedLang = []string{"zh-cn", "zh-tw", "en-us"}

	jiuxingAllowedType = []string{"0", "1"}

	hehunAllowedType = []string{"0", "1"}
	hehunAllowedLang = []string{"zh-cn", "en-us"}

	cesuanAllowedType   = []string{"0", "1"}
	cesuanAllowedSect   = []string{"1", "2"}
	cesuanAllowedZhen   = []string{"1", "2", "3"}
	cesuanAllowedFactor = []string{"0", "1"}

	jingpanAllowedType      = []string{"0", "1"}
	jingpanAllowedSect      = []string{"1", "2"}
	jingpanAllowedLoadMode  = []string{"1", "2"}
	jingpanAllowedDayun     = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	jingpanAllowedIsZip     = []string{"1", "2"}
	jingpanAllowedZhen      = []string{"1", "2", "3"}
	jingpanAllowedLang      = []string{"zh-cn", "zh-tw"}
	jingsuanAllowedLang     = []string{"zh-cn", "en-us", "zh-tw"}
	weilaiAllowedComputeDay = []string{"1", "2"}
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

type BaziHepanRequest struct {
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

type BaziCesuanRequest struct {
	Name      string
	Sex       string // 0 男，1 女
	Type      string // 0 农历，1 公历
	Year      string
	Month     string
	Day       string
	Hours     string
	Minute    string
	Sect      string // 1 晚子时日柱算明天，2 晚子时日柱算当天
	Zhen      string // 1 中国真太阳时，2 不使用真太阳时，3 全球真太阳时
	Province  string
	City      string
	Longitude string
	Latitude  string
	Timezone  string
	Lang      string // zh-cn / zh-tw / en-us
	Factor    string // 0 不调整，1 调整
}

type BaziJingpanRequest struct {
	Name       string
	Sex        string // 0 男，1 女
	Type       string // 0 农历，1 公历
	Year       string
	Month      string
	Day        string
	Hours      string
	Minute     string
	Sect       string // 1 晚子时日柱算明天，2 晚子时日柱算当天
	LoadMode   string // 1 全量，2 按需
	DayunIndex string // 1~9
	IsZip      string // 1 压缩，2 不压缩
	Zhen       string // 1 中国真太阳时，2 不使用真太阳时，3 全球真太阳时
	Province   string
	City       string
	Longitude  string
	Latitude   string
	Timezone   string
	Lang       string // zh-cn / zh-tw
}

type BaziJingsuanRequest struct {
	Name      string
	Sex       string // 0 男，1 女
	Type      string // 0 农历，1 公历
	Year      string
	Month     string
	Day       string
	Hours     string
	Minute    string
	Sect      string // 1 晚子时日柱算明天，2 晚子时日柱算当天
	Zhen      string // 1 中国真太阳时，2 不使用真太阳时，3 全球真太阳时
	Province  string
	City      string
	Longitude string
	Latitude  string
	Timezone  string
	Lang      string // zh-cn / en-us / zh-tw
	Factor    string // 0 不调整，1 调整
}

type BaziWeilaiRequest struct {
	Name         string
	Sex          string // 0 男，1 女
	Type         string // 0 农历，1 公历
	Year         string
	Month        string
	Day          string
	Hours        string
	Minute       string
	YunshiYear   string // 预测未来年（>= 当前年份）
	ComputeDaily string // 1 计算每日，2 不计算
	Sect         string // 1 晚子时日柱算明天，2 晚子时日柱算当天
	Zhen         string // 1 中国真太阳时，2 不使用真太阳时，3 全球真太阳时
	Province     string
	City         string
	Longitude    string
	Latitude     string
	Timezone     string
	Lang         string // zh-cn / en-us / zh-tw
}

func (r BaziHepanRequest) toValues() url.Values {
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

func (r BaziHepanRequest) Validate() error {
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

func (r BaziCesuanRequest) toValues() url.Values {
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
	if r.Sect != "" {
		v.Set("sect", r.Sect)
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
	if r.Longitude != "" {
		v.Set("longitude", r.Longitude)
	}
	if r.Latitude != "" {
		v.Set("latitude", r.Latitude)
	}
	if r.Timezone != "" {
		v.Set("timezone", r.Timezone)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	if r.Factor != "" {
		v.Set("factor", r.Factor)
	}
	return v
}

func (r BaziCesuanRequest) Validate() error {
	if r.Sex == "" {
		return newRequiredFieldError("sex")
	}
	if !inSet(r.Sex, baziAllowedSex) {
		return newEnumFieldError("sex", r.Sex, baziAllowedSex)
	}

	if r.Type == "" {
		return newRequiredFieldError("type")
	}
	if !inSet(r.Type, cesuanAllowedType) {
		return newEnumFieldError("type", r.Type, cesuanAllowedType)
	}

	required := []struct {
		field string
		value string
	}{
		{"year", r.Year},
		{"month", r.Month},
		{"day", r.Day},
		{"hours", r.Hours},
		{"minute", r.Minute},
	}
	for _, x := range required {
		if x.value == "" {
			return newRequiredFieldError(x.field)
		}
	}

	if r.Sect != "" && !inSet(r.Sect, cesuanAllowedSect) {
		return newEnumFieldError("sect", r.Sect, cesuanAllowedSect)
	}
	if r.Zhen != "" && !inSet(r.Zhen, cesuanAllowedZhen) {
		return newEnumFieldError("zhen", r.Zhen, cesuanAllowedZhen)
	}

	if r.Zhen == "1" {
		if r.Province == "" {
			return newRequiredFieldError("province")
		}
		if r.City == "" {
			return newRequiredFieldError("city")
		}
	}
	if r.Zhen == "3" {
		if r.Longitude == "" {
			return newRequiredFieldError("longitude")
		}
		lng, err := strconv.ParseFloat(r.Longitude, 64)
		if err != nil {
			return &ValidationError{Field: "longitude", Message: "must be a valid float"}
		}
		if lng < -180 || lng > 180 {
			return &ValidationError{Field: "longitude", Message: "must be in range [-180, 180]"}
		}

		if r.Latitude == "" {
			return newRequiredFieldError("latitude")
		}
		lat, err := strconv.ParseFloat(r.Latitude, 64)
		if err != nil {
			return &ValidationError{Field: "latitude", Message: "must be a valid float"}
		}
		if lat < -90 || lat > 90 {
			return &ValidationError{Field: "latitude", Message: "must be in range [-90, 90]"}
		}
	}

	if r.Lang != "" && !inSet(r.Lang, baziAllowedLang) {
		return newEnumFieldError("lang", r.Lang, baziAllowedLang)
	}
	if r.Factor != "" && !inSet(r.Factor, cesuanAllowedFactor) {
		return newEnumFieldError("factor", r.Factor, cesuanAllowedFactor)
	}
	return nil
}

func (r BaziJingpanRequest) toValues() url.Values {
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
	if r.Sect != "" {
		v.Set("sect", r.Sect)
	}
	if r.LoadMode != "" {
		v.Set("load_mode", r.LoadMode)
	}
	if r.DayunIndex != "" {
		v.Set("dayun_index", r.DayunIndex)
	}
	if r.IsZip != "" {
		v.Set("iszip", r.IsZip)
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
	if r.Longitude != "" {
		v.Set("longitude", r.Longitude)
	}
	if r.Latitude != "" {
		v.Set("latitude", r.Latitude)
	}
	if r.Timezone != "" {
		v.Set("timezone", r.Timezone)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r BaziJingpanRequest) Validate() error {
	if r.Sex == "" {
		return newRequiredFieldError("sex")
	}
	if !inSet(r.Sex, baziAllowedSex) {
		return newEnumFieldError("sex", r.Sex, baziAllowedSex)
	}
	if r.Type == "" {
		return newRequiredFieldError("type")
	}
	if !inSet(r.Type, jingpanAllowedType) {
		return newEnumFieldError("type", r.Type, jingpanAllowedType)
	}

	required := []struct {
		field string
		value string
	}{
		{"year", r.Year},
		{"month", r.Month},
		{"day", r.Day},
		{"hours", r.Hours},
		{"minute", r.Minute},
	}
	for _, x := range required {
		if x.value == "" {
			return newRequiredFieldError(x.field)
		}
	}

	if r.Sect != "" && !inSet(r.Sect, jingpanAllowedSect) {
		return newEnumFieldError("sect", r.Sect, jingpanAllowedSect)
	}
	if r.LoadMode != "" && !inSet(r.LoadMode, jingpanAllowedLoadMode) {
		return newEnumFieldError("load_mode", r.LoadMode, jingpanAllowedLoadMode)
	}
	if r.DayunIndex != "" && !inSet(r.DayunIndex, jingpanAllowedDayun) {
		return newEnumFieldError("dayun_index", r.DayunIndex, jingpanAllowedDayun)
	}
	if r.IsZip != "" && !inSet(r.IsZip, jingpanAllowedIsZip) {
		return newEnumFieldError("iszip", r.IsZip, jingpanAllowedIsZip)
	}
	if r.Zhen != "" && !inSet(r.Zhen, jingpanAllowedZhen) {
		return newEnumFieldError("zhen", r.Zhen, jingpanAllowedZhen)
	}
	if r.Zhen == "1" {
		if r.Province == "" {
			return newRequiredFieldError("province")
		}
		if r.City == "" {
			return newRequiredFieldError("city")
		}
	}
	if r.Zhen == "3" {
		if r.Longitude == "" {
			return newRequiredFieldError("longitude")
		}
		lng, err := strconv.ParseFloat(r.Longitude, 64)
		if err != nil {
			return &ValidationError{Field: "longitude", Message: "must be a valid float"}
		}
		if lng < -180 || lng > 180 {
			return &ValidationError{Field: "longitude", Message: "must be in range [-180, 180]"}
		}

		if r.Latitude == "" {
			return newRequiredFieldError("latitude")
		}
		lat, err := strconv.ParseFloat(r.Latitude, 64)
		if err != nil {
			return &ValidationError{Field: "latitude", Message: "must be a valid float"}
		}
		if lat < -90 || lat > 90 {
			return &ValidationError{Field: "latitude", Message: "must be in range [-90, 90]"}
		}
	}
	if r.Lang != "" && !inSet(r.Lang, jingpanAllowedLang) {
		return newEnumFieldError("lang", r.Lang, jingpanAllowedLang)
	}
	return nil
}

func (r BaziJingsuanRequest) toValues() url.Values {
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
	if r.Sect != "" {
		v.Set("sect", r.Sect)
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
	if r.Longitude != "" {
		v.Set("longitude", r.Longitude)
	}
	if r.Latitude != "" {
		v.Set("latitude", r.Latitude)
	}
	if r.Timezone != "" {
		v.Set("timezone", r.Timezone)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	if r.Factor != "" {
		v.Set("factor", r.Factor)
	}
	return v
}

func (r BaziJingsuanRequest) Validate() error {
	if r.Sex == "" {
		return newRequiredFieldError("sex")
	}
	if !inSet(r.Sex, baziAllowedSex) {
		return newEnumFieldError("sex", r.Sex, baziAllowedSex)
	}
	if r.Type == "" {
		return newRequiredFieldError("type")
	}
	if !inSet(r.Type, cesuanAllowedType) {
		return newEnumFieldError("type", r.Type, cesuanAllowedType)
	}

	required := []struct {
		field string
		value string
	}{
		{"year", r.Year},
		{"month", r.Month},
		{"day", r.Day},
		{"hours", r.Hours},
		{"minute", r.Minute},
	}
	for _, x := range required {
		if x.value == "" {
			return newRequiredFieldError(x.field)
		}
	}

	if r.Sect != "" && !inSet(r.Sect, cesuanAllowedSect) {
		return newEnumFieldError("sect", r.Sect, cesuanAllowedSect)
	}
	if r.Zhen != "" && !inSet(r.Zhen, cesuanAllowedZhen) {
		return newEnumFieldError("zhen", r.Zhen, cesuanAllowedZhen)
	}
	if r.Zhen == "1" {
		if r.Province == "" {
			return newRequiredFieldError("province")
		}
		if r.City == "" {
			return newRequiredFieldError("city")
		}
	}
	if r.Zhen == "3" {
		if r.Longitude == "" {
			return newRequiredFieldError("longitude")
		}
		lng, err := strconv.ParseFloat(r.Longitude, 64)
		if err != nil {
			return &ValidationError{Field: "longitude", Message: "must be a valid float"}
		}
		if lng < -180 || lng > 180 {
			return &ValidationError{Field: "longitude", Message: "must be in range [-180, 180]"}
		}

		if r.Latitude == "" {
			return newRequiredFieldError("latitude")
		}
		lat, err := strconv.ParseFloat(r.Latitude, 64)
		if err != nil {
			return &ValidationError{Field: "latitude", Message: "must be a valid float"}
		}
		if lat < -90 || lat > 90 {
			return &ValidationError{Field: "latitude", Message: "must be in range [-90, 90]"}
		}
	}
	if r.Lang != "" && !inSet(r.Lang, jingsuanAllowedLang) {
		return newEnumFieldError("lang", r.Lang, jingsuanAllowedLang)
	}
	if r.Factor != "" && !inSet(r.Factor, cesuanAllowedFactor) {
		return newEnumFieldError("factor", r.Factor, cesuanAllowedFactor)
	}
	return nil
}

func (r BaziWeilaiRequest) toValues() url.Values {
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
	if r.YunshiYear != "" {
		v.Set("yunshi_year", r.YunshiYear)
	}
	if r.ComputeDaily != "" {
		v.Set("compute_daily", r.ComputeDaily)
	}
	if r.Sect != "" {
		v.Set("sect", r.Sect)
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
	if r.Longitude != "" {
		v.Set("longitude", r.Longitude)
	}
	if r.Latitude != "" {
		v.Set("latitude", r.Latitude)
	}
	if r.Timezone != "" {
		v.Set("timezone", r.Timezone)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r BaziWeilaiRequest) Validate() error {
	if r.Sex == "" {
		return newRequiredFieldError("sex")
	}
	if !inSet(r.Sex, baziAllowedSex) {
		return newEnumFieldError("sex", r.Sex, baziAllowedSex)
	}
	if r.Type == "" {
		return newRequiredFieldError("type")
	}
	if !inSet(r.Type, cesuanAllowedType) {
		return newEnumFieldError("type", r.Type, cesuanAllowedType)
	}

	required := []struct {
		field string
		value string
	}{
		{"year", r.Year},
		{"month", r.Month},
		{"day", r.Day},
		{"hours", r.Hours},
		{"minute", r.Minute},
		{"yunshi_year", r.YunshiYear},
	}
	for _, x := range required {
		if x.value == "" {
			return newRequiredFieldError(x.field)
		}
	}

	yunshiYear, err := strconv.Atoi(r.YunshiYear)
	if err != nil {
		return &ValidationError{Field: "yunshi_year", Message: "must be a valid integer year"}
	}
	if yunshiYear < time.Now().Year() {
		return &ValidationError{Field: "yunshi_year", Message: "must be greater than or equal to current year"}
	}

	if r.ComputeDaily != "" && !inSet(r.ComputeDaily, weilaiAllowedComputeDay) {
		return newEnumFieldError("compute_daily", r.ComputeDaily, weilaiAllowedComputeDay)
	}
	if r.Sect != "" && !inSet(r.Sect, cesuanAllowedSect) {
		return newEnumFieldError("sect", r.Sect, cesuanAllowedSect)
	}
	if r.Zhen != "" && !inSet(r.Zhen, cesuanAllowedZhen) {
		return newEnumFieldError("zhen", r.Zhen, cesuanAllowedZhen)
	}
	if r.Zhen == "1" {
		if r.Province == "" {
			return newRequiredFieldError("province")
		}
		if r.City == "" {
			return newRequiredFieldError("city")
		}
	}
	if r.Zhen == "3" {
		if r.Longitude == "" {
			return newRequiredFieldError("longitude")
		}
		lng, err := strconv.ParseFloat(r.Longitude, 64)
		if err != nil {
			return &ValidationError{Field: "longitude", Message: "must be a valid float"}
		}
		if lng < -180 || lng > 180 {
			return &ValidationError{Field: "longitude", Message: "must be in range [-180, 180]"}
		}

		if r.Latitude == "" {
			return newRequiredFieldError("latitude")
		}
		lat, err := strconv.ParseFloat(r.Latitude, 64)
		if err != nil {
			return &ValidationError{Field: "latitude", Message: "must be a valid float"}
		}
		if lat < -90 || lat > 90 {
			return &ValidationError{Field: "latitude", Message: "must be in range [-90, 90]"}
		}
	}
	if r.Lang != "" && !inSet(r.Lang, jingsuanAllowedLang) {
		return newEnumFieldError("lang", r.Lang, jingsuanAllowedLang)
	}
	return nil
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

type BaziHepanData struct {
	Male          BaziHehunPersonInfo   `json:"male"`
	Female        BaziHehunPersonInfo   `json:"female"`
	Minggong      BaziHehunMinggongInfo `json:"minggong"`
	Nianqitongzhi BaziHehunNianqiInfo   `json:"nianqitongzhi"`
	Yueling       BaziHehunYuelingInfo  `json:"yueling"`
	Rigan         BaziHehunRiganInfo    `json:"rigan"`
	Tiangan       BaziHehunTianganInfo  `json:"tiangan"`
	Jiankang      BaziHepanJiankangInfo `json:"jiankang"`
	AllScore      int                   `json:"all_score"`
	MaleSX        string                `json:"male_sx"`
	FemaleSX      string                `json:"female_sx"`
	MaleXZ        string                `json:"male_xz"`
	FemaleXZ      string                `json:"female_xz"`
}

type BaziHepanJiankangInfo struct {
	Score             string `json:"score"`
	Description       string `json:"description"`
	DetailDescription string `json:"detail_description"`
}

type BaziCesuanData struct {
	BaseInfo   BaziCesuanBaseInfo   `json:"base_info"`
	BaziInfo   BaziCesuanBaziInfo   `json:"bazi_info"`
	Chenggu    BaziCesuanChenggu    `json:"chenggu"`
	Wuxing     BaziCesuanWuxing     `json:"wuxing"`
	Yinyuan    BaziCesuanYinyuan    `json:"yinyuan"`
	Caiyun     BaziCesuanCaiyun     `json:"caiyun"`
	Sizhu      BaziCesuanSizhu      `json:"sizhu"`
	Mingyun    BaziCesuanMingyun    `json:"mingyun"`
	SX         string               `json:"sx"`
	XZ         string               `json:"xz"`
	Xiyongshen BaziCesuanXiyongshen `json:"xiyongshen"`
}

type BaziCesuanBaseInfo struct {
	Zhen       *BaziZhenInfo `json:"zhen,omitempty"`
	Sex        string        `json:"sex"`
	Name       string        `json:"name"`
	Gongli     string        `json:"gongli"`
	Nongli     string        `json:"nongli"`
	Qiyun      string        `json:"qiyun"`
	Jiaoyun    string        `json:"jiaoyun"`
	Zhengge    string        `json:"zhengge"`
	WuxingXiji string        `json:"wuxing_xiji"`
}

type BaziCesuanBaziInfo struct {
	KW      string   `json:"kw"`
	TGCGGod []string `json:"tg_cg_god"`
	Bazi    string   `json:"bazi"`
	NaYin   string   `json:"na_yin"`
}

type BaziCesuanChenggu struct {
	YearWeight  string `json:"year_weight"`
	MonthWeight string `json:"month_weight"`
	DayWeight   string `json:"day_weight"`
	HourWeight  string `json:"hour_weight"`
	TotalWeight string `json:"total_weight"`
	Description string `json:"description"`
}

type BaziCesuanWuxing struct {
	DetailDesc        string `json:"detail_desc"`
	SimpleDesc        string `json:"simple_desc"`
	SimpleDescription string `json:"simple_description"`
	DetailDescription string `json:"detail_description"`
}

type BaziCesuanYinyuan struct {
	SanshishuYinyuan string `json:"sanshishu_yinyuan"`
}

type BaziCesuanCaiyun struct {
	SanshishuCaiyun BaziCesuanSimpleDetail `json:"sanshishu_caiyun"`
}

type BaziCesuanSimpleDetail struct {
	SimpleDesc string `json:"simple_desc"`
	DetailDesc string `json:"detail_desc"`
}

type BaziCesuanSizhu struct {
	Rizhu string `json:"rizhu"`
}

type BaziCesuanMingyun struct {
	SanshishuMingyun string `json:"sanshishu_mingyun"`
}

type BaziCesuanXiyongshen struct {
	Qiangruo         string `json:"qiangruo"`
	Xiyongshen       string `json:"xiyongshen"`
	Jishen           string `json:"jishen"`
	XiyongshenDesc   string `json:"xiyongshen_desc"`
	JinNumber        int    `json:"jin_number"`
	MuNumber         int    `json:"mu_number"`
	ShuiNumber       int    `json:"shui_number"`
	HuoNumber        int    `json:"huo_number"`
	TuNumber         int    `json:"tu_number"`
	Tonglei          string `json:"tonglei"`
	Yilei            string `json:"yilei"`
	RizhuTiangan     string `json:"rizhu_tiangan"`
	Zidang           int    `json:"zidang"`
	Yidang           int    `json:"yidang"`
	ZidangPercent    string `json:"zidang_percent"`
	YidangPercent    string `json:"yidang_percent"`
	JinScore         int    `json:"jin_score"`
	MuScore          int    `json:"mu_score"`
	ShuiScore        int    `json:"shui_score"`
	HuoScore         int    `json:"huo_score"`
	TuScore          int    `json:"tu_score"`
	JinScorePercent  string `json:"jin_score_percent"`
	MuScorePercent   string `json:"mu_score_percent"`
	ShuiScorePercent string `json:"shui_score_percent"`
	HuoScorePercent  string `json:"huo_score_percent"`
	TuScorePercent   string `json:"tu_score_percent"`
	Yinyang          string `json:"yinyang"`
}

type BaziJingpanData struct {
	BaseInfo   BaziJingpanBaseInfo   `json:"base_info"`
	DetailInfo BaziJingpanDetailInfo `json:"detail_info"`
}

type BaziJingsuanData struct {
	BaseInfo   BaziJingpanBaseInfo    `json:"base_info"`
	DetailInfo BaziJingsuanDetailInfo `json:"detail_info"`
}

type BaziWeilaiData struct {
	BaseInfo   BaziJingpanBaseInfo  `json:"base_info"`
	DetailInfo BaziWeilaiDetailInfo `json:"detail_info"`
}

type BaziJingpanBaseInfo struct {
	Zhen          *BaziZhenInfo        `json:"zhen,omitempty"`
	Sex           string               `json:"sex"`
	Name          string               `json:"name"`
	Gongli        string               `json:"gongli"`
	Nongli        string               `json:"nongli"`
	Qiyun         string               `json:"qiyun"`
	Jiaoyun       string               `json:"jiaoyun"`
	Taiyuan       string               `json:"taiyuan"`
	TaiyuanNaYin  string               `json:"taiyuan_nayin"`
	Taixi         string               `json:"taixi"`
	TaixiNaYin    string               `json:"taixi_nayin"`
	Minggong      string               `json:"minggong"`
	MinggongNaYin string               `json:"minggong_nayin"`
	Shengong      string               `json:"shengong"`
	ShengongNaYin string               `json:"shengong_nayin"`
	Shengxiao     string               `json:"shengxiao"`
	Xingzuo       string               `json:"xingzuo"`
	Siling        string               `json:"siling"`
	JiaoyunMang   string               `json:"jiaoyun_mang"`
	Xingxiu       string               `json:"xingxiu"`
	Minggua       BaziJingpanMinggua   `json:"minggua"`
	WuxingWangdu  string               `json:"wuxing_wangdu"`
	WuxingXiji    string               `json:"wuxing_xiji"`
	TianganLiuyi  string               `json:"tiangan_liuyi"`
	DizhiLiuyi    string               `json:"dizhi_liuyi"`
	Zhengge       string               `json:"zhengge"`
	Xiyongshen    BaziCesuanXiyongshen `json:"xiyongshen"`
}

type BaziJingpanMinggua struct {
	MingguaName    string `json:"minggua_name"`
	MingguaFangwei string `json:"minggua_fangwei"`
}

type BaziJingpanDetailInfo struct {
	DayunInfo     []BaziJingpanDayunInfo   `json:"dayun_info"`
	DayunBornInfo BaziJingpanDayunBornInfo `json:"dayun_born_info"`
	SizhuInfo     BaziJingpanSizhuInfo     `json:"sizhu_info"`
	TaishenInfo   BaziJingpanTaishenInfo   `json:"taishen_info"`
}

type BaziJingsuanDetailInfo struct {
	DayunInfo []BaziJingpanDayunInfo `json:"dayun_info"`
	SizhuInfo BaziJingsuanSizhuInfo  `json:"sizhu_info"`
}

type BaziWeilaiDetailInfo struct {
	SizhuInfo       BaziWeilaiSizhuInfo       `json:"sizhu_info"`
	YunshiYearInfo  BaziWeilaiYearInfoWrapper `json:"yunshi_year_info"`
	YunshiMonthInfo []BaziWeilaiMonthInfo     `json:"yunshi_month_info"`
}

type BaziJingpanDayunBornInfo []BaziJingpanDayunInfo

func (d *BaziJingpanDayunBornInfo) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*d = nil
		return nil
	}

	var arr []BaziJingpanDayunInfo
	if err := json.Unmarshal(data, &arr); err == nil {
		*d = arr
		return nil
	}

	var single BaziJingpanDayunInfo
	if err := json.Unmarshal(data, &single); err == nil {
		*d = []BaziJingpanDayunInfo{single}
		return nil
	}

	return &ValidationError{Field: "dayun_born_info", Message: "unexpected json type"}
}

type BaziJingpanDayunInfo struct {
	DayunIndex       int                      `json:"dayun_index"`
	DayunStartYear   int                      `json:"dayun_start_year"`
	DayunStartAge    int                      `json:"dayun_start_age"`
	DayunEndYear     int                      `json:"dayun_end_year"`
	DayunEndAge      int                      `json:"dayun_end_age"`
	DayunStartGanzhi string                   `json:"dayun_start_ganzhi"`
	DayunXun         string                   `json:"dayun_xun"`
	DayunKongwang    string                   `json:"dayun_kongwang"`
	DayunShishen     string                   `json:"dayun_shishen"`
	DayunChangsheng  string                   `json:"dayun_changsheng"`
	DayunYearTG      string                   `json:"dayun_year_tg"`
	DayunYearDZ      string                   `json:"dayun_year_dz"`
	DayunYearCG      []string                 `json:"dayun_year_cg"`
	DayunYearTGGod   string                   `json:"dayun_year_tg_god"`
	DayunYearDZGod   []string                 `json:"dayun_year_dz_god"`
	DayunDZStarCS    string                   `json:"dayun_dz_star_cs"`
	DayunDZSelfCS    string                   `json:"dayun_dz_self_cs"`
	DayunNaYin       string                   `json:"dayun_nayin"`
	DayunWuxing      string                   `json:"dayun_wuxing"`
	DayunTGLiuyi     string                   `json:"dayun_tg_liuyi"`
	DayunDZLiuyi     string                   `json:"dayun_dz_liuyi"`
	LiunianInfo      []BaziJingpanLiunianInfo `json:"liunian_info"`
	XiaoyunInfo      []BaziJingpanXiaoyunInfo `json:"xiaoyun_info"`
	DayunShensha     string                   `json:"dayun_shensha"`
	DayunIndication  BaziYunshiIndication     `json:"dayun_indication"`
}

type BaziJingpanLiunianInfo struct {
	LiunianIndex      int                     `json:"liunian_index"`
	LiunianYear       int                     `json:"liunian_year"`
	LiunianAge        int                     `json:"liunian_age"`
	LiunianGanzhi     string                  `json:"liunian_ganzhi"`
	LiunianXun        string                  `json:"liunian_xun"`
	LiunianKongwang   string                  `json:"liunian_kongwang"`
	LiunianShishen    string                  `json:"liunian_shishen"`
	LiunianChangsheng string                  `json:"liunian_changsheng"`
	LiunianJieqi      []BaziJingpanJieqiInfo  `json:"liunian_jieqi"`
	LiunianYearTG     string                  `json:"liunian_year_tg"`
	LiunianYearDZ     string                  `json:"liunian_year_dz"`
	LiunianYearCG     []string                `json:"liunian_year_cg"`
	LiunianTGGod      string                  `json:"liunian_tg_god"`
	LiunianDZGod      []string                `json:"liunian_dz_god"`
	LiunianDZStarCS   string                  `json:"liunian_dz_star_cs"`
	LiunianDZSelfCS   string                  `json:"liunian_dz_self_cs"`
	LiunianNaYin      string                  `json:"liunian_nayin"`
	LiunianWuxing     string                  `json:"liunian_wuxing"`
	LiunianTGLiuyi    string                  `json:"liunian_tg_liuyi"`
	LiunianDZLiuyi    string                  `json:"liunian_dz_liuyi"`
	LiunianShensha    string                  `json:"liunian_shensha"`
	LiuyueInfo        []BaziJingpanLiuyueInfo `json:"liuyue_info"`
	LiunianIndication BaziYunshiIndication    `json:"liunian_indication"`
}

type BaziJingpanJieqiInfo struct {
	JieqiName        string `json:"jieqi_name"`
	JieqiSolarYMDHMS string `json:"jieqi_solar_ymdhms"`
	JieqiLunarIsLeap string `json:"jieqi_lunar_ymdhms_isleap"`
	JieqiLunarYMDHMS string `json:"jieqi_lunar_ymdhms"`
}

type BaziJingpanLiuyueInfo struct {
	LiuyueIndex      int      `json:"liuyue_index"`
	LiuyueMonth      string   `json:"liuyue_month"`
	LiuyueGanzhi     string   `json:"liuyue_ganzhi"`
	LiuyueXun        string   `json:"liuyue_xun"`
	LiuyueKongwang   string   `json:"liuyue_kongwang"`
	LiuyueShishen    string   `json:"liuyue_shishen"`
	LiuyueChangsheng string   `json:"liuyue_changsheng"`
	LiuyueMonthTG    string   `json:"liuyue_month_tg"`
	LiuyueMonthDZ    string   `json:"liuyue_month_dz"`
	LiuyueMonthCG    []string `json:"liuyue_month_cg"`
	LiuyueMonthTGGod string   `json:"liuyue_month_tg_god"`
	LiuyueMonthDZGod []string `json:"liuyue_month_dz_god"`
	LiuyueDZStarCS   string   `json:"liuyue_dz_star_cs"`
	LiuyueDZSelfCS   string   `json:"liuyue_dz_self_cs"`
	LiuyueNaYin      string   `json:"liuyue_nayin"`
	LiuyueWuxing     string   `json:"liuyue_wuxing"`
	LiuyueTGLiuyi    string   `json:"liuyue_tg_liuyi"`
	LiuyueDZLiuyi    string   `json:"liuyue_dz_liuyi"`
	LiuyueShensha    string   `json:"liuyue_shensha"`
}

type BaziJingpanXiaoyunInfo struct {
	XiaoyunIndex      int      `json:"xiaoyun_index"`
	XiaoyunYear       int      `json:"xiaoyun_year"`
	XiaoyunAge        int      `json:"xiaoyun_age"`
	XiaoyunGanzhi     string   `json:"xiaoyun_ganzhi"`
	XiaoyunXun        string   `json:"xiaoyun_xun"`
	XiaoyunKongwang   string   `json:"xiaoyun_kongwang"`
	XiaoyunShishen    string   `json:"xiaoyun_shishen"`
	XiaoyunChangsheng string   `json:"xiaoyun_changsheng"`
	XiaoyunYearTG     string   `json:"xiaoyun_year_tg"`
	XiaoyunYearDZ     string   `json:"xiaoyun_year_dz"`
	XiaoyunYearCG     []string `json:"xiaoyun_year_cg"`
	XiaoyunTGGod      string   `json:"xiaoyun_tg_god"`
	XiaoyunDZGod      []string `json:"xiaoyun_dz_god"`
	XiaoyunDZStarCS   string   `json:"xiaoyun_dz_star_cs"`
	XiaoyunDZSelfCS   string   `json:"xiaoyun_dz_self_cs"`
	XiaoyunNaYin      string   `json:"xiaoyun_nayin"`
	XiaoyunWuxing     string   `json:"xiaoyun_wuxing"`
	XiaoyunTGLiuyi    string   `json:"xiaoyun_tg_liuyi"`
	XiaoyunDZLiuyi    string   `json:"xiaoyun_dz_liuyi"`
}

type BaziJingpanSizhuInfo struct {
	Year  BaziJingpanPillarInfo `json:"year"`
	Month BaziJingpanPillarInfo `json:"month"`
	Day   BaziJingpanPillarInfo `json:"day"`
	Hour  BaziJingpanPillarInfo `json:"hour"`
}

type BaziJingsuanSizhuInfo struct {
	Year            BaziJingpanPillarInfo       `json:"year"`
	Month           BaziJingpanPillarInfo       `json:"month"`
	Day             BaziJingpanPillarInfo       `json:"day"`
	Hour            BaziJingpanPillarInfo       `json:"hour"`
	SizhuIndication BaziJingsuanSizhuIndication `json:"sizhu_indication"`
}

type BaziWeilaiSizhuInfo struct {
	Year  BaziJingpanPillarInfo `json:"year"`
	Month BaziJingpanPillarInfo `json:"month"`
	Day   BaziJingpanPillarInfo `json:"day"`
	Hour  BaziJingpanPillarInfo `json:"hour"`
}

type BaziJingpanPillarInfo struct {
	TGGod    string   `json:"tg_god"`
	TG       string   `json:"tg"`
	DZ       string   `json:"dz"`
	CG       []string `json:"cg"`
	DZGod    []string `json:"dz_god"`
	DZStarCS string   `json:"dz_star_cs"`
	DZSelfCS string   `json:"dz_self_cs"`
	KW       string   `json:"kw"`
	NY       string   `json:"ny"`
	WX       string   `json:"wx"`
	Xun      string   `json:"xun"`
	Shensha  string   `json:"shensha"`
}

type BaziJingpanTaishenInfo struct {
	Taiyuan  BaziJingpanTaishenItem `json:"taiyuan"`
	Taixi    BaziJingpanTaishenItem `json:"taixi"`
	Minggong BaziJingpanTaishenItem `json:"minggong"`
	Shengong BaziJingpanTaishenItem `json:"shengong"`
}

type BaziJingpanTaishenItem struct {
	Ganzhi     string   `json:"ganzhi"`
	Xun        string   `json:"xun"`
	Kongwang   string   `json:"kongwang"`
	Changsheng string   `json:"changsheng"`
	Gan        string   `json:"gan"`
	Zhi        string   `json:"zhi"`
	DZCG       []string `json:"dz_cg"`
	TGGod      string   `json:"tg_god"`
	DZGod      []string `json:"dz_god"`
	DZStarCS   string   `json:"dz_star_cs"`
	DZSelfCS   string   `json:"dz_self_cs"`
	NaYin      string   `json:"nayin"`
	Wuxing     string   `json:"wuxing"`
	TGLiuyi    string   `json:"tg_liuyi"`
	DZLiuyi    string   `json:"dz_liuyi"`
	Shensha    string   `json:"shensha"`
}

type BaziJingsuanSizhuIndication struct {
	Chenggu BaziCesuanChenggu  `json:"chenggu"`
	Wuxing  BaziCesuanWuxing   `json:"wuxing"`
	Yinyuan BaziCesuanYinyuan  `json:"yinyuan"`
	Caiyun  BaziCesuanCaiyun   `json:"caiyun"`
	Xingge  BaziJingsuanXingge `json:"xingge"`
	Mingyun BaziCesuanMingyun  `json:"mingyun"`
}

type BaziJingsuanXingge struct {
	Rizhu string `json:"rizhu"`
}

type BaziYunshiIndication struct {
	Shiye    string `json:"shiye"`
	Xueye    string `json:"xueye"`
	Yunshi   string `json:"yunshi"`
	Yinyuan  string `json:"yinyuan"`
	Caiyun   string `json:"caiyun"`
	Jiankang string `json:"jiankang"`
}

type BaziWeilaiYearInfoWrapper struct {
	YunshiYear BaziWeilaiYearInfo `json:"yunshi_year"`
}

type BaziWeilaiYearInfo struct {
	Year       int                  `json:"year"`
	TGGod      string               `json:"tg_god"`
	TG         string               `json:"tg"`
	DZ         string               `json:"dz"`
	CG         []string             `json:"cg"`
	DZGod      []string             `json:"dz_god"`
	DZStarCS   string               `json:"dz_star_cs"`
	DZSelfCS   string               `json:"dz_self_cs"`
	KW         string               `json:"kw"`
	NY         string               `json:"ny"`
	WX         string               `json:"wx"`
	Xun        string               `json:"xun"`
	Shensha    string               `json:"shensha"`
	Indication BaziYunshiIndication `json:"indication"`
}

type BaziWeilaiMonthInfo struct {
	Month         string               `json:"month"`
	TGGod         string               `json:"tg_god"`
	TG            string               `json:"tg"`
	DZ            string               `json:"dz"`
	CG            []string             `json:"cg"`
	DZGod         []string             `json:"dz_god"`
	DZStarCS      string               `json:"dz_star_cs"`
	DZSelfCS      string               `json:"dz_self_cs"`
	KW            string               `json:"kw"`
	NY            string               `json:"ny"`
	WX            string               `json:"wx"`
	Xun           string               `json:"xun"`
	Shensha       string               `json:"shensha"`
	Indication    BaziYunshiIndication `json:"indication"`
	YunshiDayInfo []BaziWeilaiDayInfo  `json:"yunshi_day_info"`
}

type BaziWeilaiDayInfo struct {
	Day        string               `json:"day"`
	TGGod      string               `json:"tg_god"`
	TG         string               `json:"tg"`
	DZ         string               `json:"dz"`
	CG         []string             `json:"cg"`
	DZGod      []string             `json:"dz_god"`
	DZStarCS   string               `json:"dz_star_cs"`
	DZSelfCS   string               `json:"dz_self_cs"`
	KW         string               `json:"kw"`
	NY         string               `json:"ny"`
	WX         string               `json:"wx"`
	Xun        string               `json:"xun"`
	Shensha    string               `json:"shensha"`
	Indication BaziYunshiIndication `json:"indication"`
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

func (s *BaziService) Hepan(ctx context.Context, req BaziHepanRequest) (*CommonResponse[BaziHepanData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[BaziHepanData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/hepan", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *BaziService) Cesuan(ctx context.Context, req BaziCesuanRequest) (*CommonResponse[BaziCesuanData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[BaziCesuanData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/cesuan", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *BaziService) Jingpan(ctx context.Context, req BaziJingpanRequest) (*CommonResponse[BaziJingpanData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[BaziJingpanData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/jingpan", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *BaziService) Jingsuan(ctx context.Context, req BaziJingsuanRequest) (*CommonResponse[BaziJingsuanData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[BaziJingsuanData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/jingsuan", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *BaziService) Weilai(ctx context.Context, req BaziWeilaiRequest) (*CommonResponse[BaziWeilaiData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[BaziWeilaiData]{}
	if err := s.client.doForm(ctx, "/v1/Bazi/weilai", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}
