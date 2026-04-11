package yuanfenju

import (
	"context"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

var divinationAllowedLang = []string{"zh-cn", "en-us"}
var yaoguaAllowedLang = []string{"zh-cn", "en-us", "zh-tw"}
var divinationYunshiAllowedLang = []string{"zh-cn", "zh-tw", "en-us"}
var divinationYunshiAllowedType = []string{"0", "1"}
var divinationYunshiAllowedParameterStyle = []string{"chinese", "english"}
var taluojieduAllowedLang = []string{"zh-cn", "en-us", "zh-tw"}
var taluoSpreadsAllowedRange = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
var taluoInverseAllowedValues = []string{"0", "1"}

type DivinationService struct {
	client *Client
}

type MeiriRequest struct {
	Lang string // zh-cn / en-us
}

func (r MeiriRequest) toValues() url.Values {
	v := url.Values{}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r MeiriRequest) Validate() error {
	if r.Lang != "" && !inSet(r.Lang, divinationAllowedLang) {
		return newEnumFieldError("lang", r.Lang, divinationAllowedLang)
	}
	return nil
}

type MeiriData struct {
	Number      int              `json:"number"`
	GuaMing     string           `json:"guaming"`
	Description MeiriDescription `json:"description"`
}

type MeiriDescription struct {
	GuaYue  string `json:"卦曰"`
	JieYue  string `json:"解曰"`
	XiongJi string `json:"凶吉"`
	YunShi  string `json:"运势"`
	CaiFu   string `json:"财富"`
	GanQing string `json:"感情"`
	ShiYe   string `json:"事业"`
	ShenTi  string `json:"身体"`
	ShenGui string `json:"神鬼"`
	XingRen string `json:"行人"`
}

type XiaoliurenRequest struct {
	Shuzi string // 0~999999999999
	Lang  string // zh-cn / en-us
}

func (r XiaoliurenRequest) toValues() url.Values {
	v := url.Values{}
	v.Set("shuzi", r.Shuzi)
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r XiaoliurenRequest) Validate() error {
	if r.Shuzi == "" {
		return newRequiredFieldError("shuzi")
	}
	n, err := strconv.ParseInt(r.Shuzi, 10, 64)
	if err != nil {
		return &ValidationError{Field: "shuzi", Message: "must be a valid integer"}
	}
	if n < 0 || n > 999999999999 {
		return &ValidationError{Field: "shuzi", Message: "must be in range [0, 999999999999]"}
	}
	if r.Lang != "" && !inSet(r.Lang, divinationAllowedLang) {
		return newEnumFieldError("lang", r.Lang, divinationAllowedLang)
	}
	return nil
}

type XiaoliurenData struct {
	Number      int              `json:"number"`
	GuaMing     string           `json:"guaming"`
	Description MeiriDescription `json:"description"`
}

type ZhiwenRequest struct {
	Muzhi     string // 0 箩纹，1 簸箕纹
	Shizhi    string // 0 箩纹，1 簸箕纹
	Zhongzhi  string // 0 箩纹，1 簸箕纹
	Wumingzhi string // 0 箩纹，1 簸箕纹
	Xiaozhi   string // 0 箩纹，1 簸箕纹
}

func (r ZhiwenRequest) toValues() url.Values {
	v := url.Values{}
	if r.Muzhi != "" {
		v.Set("muzhi", r.Muzhi)
	}
	if r.Shizhi != "" {
		v.Set("shizhi", r.Shizhi)
	}
	if r.Zhongzhi != "" {
		v.Set("zhongzhi", r.Zhongzhi)
	}
	if r.Wumingzhi != "" {
		v.Set("wumingzhi", r.Wumingzhi)
	}
	if r.Xiaozhi != "" {
		v.Set("xiaozhi", r.Xiaozhi)
	}
	return v
}

func (r ZhiwenRequest) Validate() error {
	allowed := []string{"0", "1"}
	required := []struct {
		field string
		value string
	}{
		{"muzhi", r.Muzhi},
		{"shizhi", r.Shizhi},
		{"zhongzhi", r.Zhongzhi},
		{"wumingzhi", r.Wumingzhi},
		{"xiaozhi", r.Xiaozhi},
	}
	for _, x := range required {
		if x.value == "" {
			return newRequiredFieldError(x.field)
		}
		if !inSet(x.value, allowed) {
			return newEnumFieldError(x.field, x.value, allowed)
		}
	}
	return nil
}

type ZhiwenData struct {
	Muzhi       string            `json:"muzhi"`
	Shizhi      string            `json:"shizhi"`
	Zhongzhi    string            `json:"zhongzhi"`
	Wumingzhi   string            `json:"wumingzhi"`
	Xiaozhi     string            `json:"xiaozhi"`
	Description ZhiwenDescription `json:"description"`
}

type ZhiwenDescription struct {
	Fenxi    string `json:"分析"`
	Shiyue   string `json:"诗曰"`
	Xingge   string `json:"性格"`
	Hunyin   string `json:"婚姻"`
	Zhiye    string `json:"职业"`
	Jiankang string `json:"健康"`
	Yunshi   string `json:"运势"`
}

type YaoguaRequest struct {
	Lang string // zh-cn / en-us / zh-tw
}

func (r YaoguaRequest) toValues() url.Values {
	v := url.Values{}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r YaoguaRequest) Validate() error {
	if r.Lang != "" && !inSet(r.Lang, yaoguaAllowedLang) {
		return newEnumFieldError("lang", r.Lang, yaoguaAllowedLang)
	}
	return nil
}

type YaoguaData struct {
	ID          int    `json:"id"`
	CommonDesc1 string `json:"common_desc1"`
	CommonDesc2 string `json:"common_desc2"`
	CommonDesc3 string `json:"common_desc3"`
	Shiye       string `json:"shiye"`
	Jingshang   string `json:"jingshang"`
	Qiuming     string `json:"qiuming"`
	Waichu      string `json:"waichu"`
	Hunlian     string `json:"hunlian"`
	Juece       string `json:"juece"`
	Image       string `json:"image"`
}

type TaluojieduRequest struct {
	SpreadID string // 1~17
	TopicID  string // 1~5
	Lang     string // zh-cn / en-us / zh-tw
}

func (r TaluojieduRequest) toValues() url.Values {
	v := url.Values{}
	if r.SpreadID != "" {
		v.Set("spread_id", r.SpreadID)
	}
	if r.TopicID != "" {
		v.Set("topic_id", r.TopicID)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r TaluojieduRequest) Validate() error {
	if r.SpreadID == "" {
		return newRequiredFieldError("spread_id")
	}
	spreadID, err := strconv.Atoi(r.SpreadID)
	if err != nil {
		return &ValidationError{Field: "spread_id", Message: "must be a valid integer"}
	}
	if spreadID < 1 || spreadID > 17 {
		return &ValidationError{Field: "spread_id", Message: "must be in range [1, 17]"}
	}

	if r.TopicID == "" {
		return newRequiredFieldError("topic_id")
	}
	topicID, err := strconv.Atoi(r.TopicID)
	if err != nil {
		return &ValidationError{Field: "topic_id", Message: "must be a valid integer"}
	}
	if topicID < 1 || topicID > 5 {
		return &ValidationError{Field: "topic_id", Message: "must be in range [1, 5]"}
	}

	if r.Lang != "" && !inSet(r.Lang, taluojieduAllowedLang) {
		return newEnumFieldError("lang", r.Lang, taluojieduAllowedLang)
	}
	return nil
}

type TaluojieduData struct {
	Cards                 []TaluojieduCard                `json:"cards"`
	OverallInterpretation TaluojieduOverallInterpretation `json:"overall_interpretation"`
	Environment           TaluojieduEnvironment           `json:"environment"`
}

type TaluojieduCard struct {
	PositionsIndex     int                          `json:"positions_index"`
	PositionsName      string                       `json:"positions_name"`
	PositionsDesc      string                       `json:"positions_desc"`
	OrientationCode    int                          `json:"orientation_code"`
	OrientationText    string                       `json:"orientation_text"`
	CardNo             int                          `json:"card_no"`
	CardName           string                       `json:"card_name"`
	CardKeywords       string                       `json:"card_keywords"`
	CardAstrology      string                       `json:"card_astrology"`
	CardElement        string                       `json:"card_element"`
	CardDescription    string                       `json:"card_description"`
	CardInterpretation TaluojieduCardInterpretation `json:"card_interpretation"`
	ImageID            int                          `json:"image_id"`
	ImageURL           string                       `json:"image_url"`
}

type TaluojieduCardInterpretation struct {
	General string `json:"general"`
	Topic   string `json:"topic"`
	Advice  string `json:"advice"`
}

type TaluojieduOverallInterpretation struct {
	SummaryMessage string `json:"summary_message"`
	OracleMessage  string `json:"oracle_message"`
}

type TaluojieduEnvironment struct {
	CalculationTime string `json:"calculation_time"`
	TimeGanzhi      string `json:"time_ganzhi"`
	TimeElement     string `json:"time_element"`
}

type TaluoxipaiRequest struct {
	TaluoSpreads string // 1~9，默认 1
}

func (r TaluoxipaiRequest) toValues() url.Values {
	v := url.Values{}
	if r.TaluoSpreads != "" {
		v.Set("taluo_spreads", r.TaluoSpreads)
	}
	return v
}

func (r TaluoxipaiRequest) Validate() error {
	if r.TaluoSpreads == "" {
		return nil
	}
	if !inSet(r.TaluoSpreads, taluoSpreadsAllowedRange) {
		return newEnumFieldError("taluo_spreads", r.TaluoSpreads, taluoSpreadsAllowedRange)
	}
	return nil
}

type TaluoxipaiData struct {
	CardNos []int  `json:"-"`
	Image   string `json:"image"`
}

func (d *TaluoxipaiData) UnmarshalJSON(data []byte) error {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		var numbers []int
		if err2 := json.Unmarshal(data, &numbers); err2 == nil {
			*d = TaluoxipaiData{CardNos: numbers}
			return nil
		}

		var array []json.RawMessage
		if err2 := json.Unmarshal(data, &array); err2 == nil && len(array) == 0 {
			*d = TaluoxipaiData{}
			return nil
		}
		return err
	}

	type cardPair struct {
		Position int
		CardNo   int
	}

	pairs := make([]cardPair, 0, len(object))
	for k, v := range object {
		if k == "image" {
			if err := json.Unmarshal(v, &d.Image); err != nil {
				return err
			}
			continue
		}

		position, err := strconv.Atoi(k)
		if err != nil {
			continue
		}

		var cardNo int
		if err := json.Unmarshal(v, &cardNo); err != nil {
			return err
		}
		pairs = append(pairs, cardPair{Position: position, CardNo: cardNo})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Position < pairs[j].Position
	})

	d.CardNos = make([]int, 0, len(pairs))
	for _, p := range pairs {
		d.CardNos = append(d.CardNos, p.CardNo)
	}

	return nil
}

type TaluozhanbuRequest struct {
	TaluoInverse string // 0 正位，1 逆位
	Lang         string // zh-cn / en-us / zh-tw
}

func (r TaluozhanbuRequest) toValues() url.Values {
	v := url.Values{}
	if r.TaluoInverse != "" {
		v.Set("taluo_inverse", r.TaluoInverse)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r TaluozhanbuRequest) Validate() error {
	if r.TaluoInverse != "" && !inSet(r.TaluoInverse, taluoInverseAllowedValues) {
		return newEnumFieldError("taluo_inverse", r.TaluoInverse, taluoInverseAllowedValues)
	}
	if r.Lang != "" && !inSet(r.Lang, taluojieduAllowedLang) {
		return newEnumFieldError("lang", r.Lang, taluojieduAllowedLang)
	}
	return nil
}

type TaluozhanbuData struct {
	CardName        string                  `json:"牌名"`
	CardKeyword     string                  `json:"关键字"`
	CardAstrology   string                  `json:"星相"`
	CardElements    string                  `json:"四要素"`
	CardDescription string                  `json:"牌面描述"`
	UprightMeaning  TaluozhanbuMeaningBlock `json:"正位含义"`
	ReverseMeaning  TaluozhanbuMeaningBlock `json:"逆位含义"`
	Meaning         TaluozhanbuMeaningBlock `json:"含义"`
	Orientation     string                  `json:"正逆"`
	ID              int                     `json:"id"`
	Image           string                  `json:"image"`
}

type TaluozhanbuMeaningBlock struct {
	BaseMeaning         string `json:"基本含义"`
	LoveMarriage        string `json:"恋爱婚姻"`
	WorkStudy           string `json:"工作学业"`
	InterpersonalWealth string `json:"人际财富"`
	HealthLife          string `json:"健康生活"`
	Other               string `json:"其它"`
	Advice              string `json:"建议"`
	LoveMarriageAdvice  string `json:"恋爱婚姻建议"`
	WorkStudyAdvice     string `json:"工作学业建议"`
	InterWealthAdvice   string `json:"人际财富建议"`
	HealthLifeAdvice    string `json:"健康生活建议"`
}

type TaluospreadsRequest struct {
	TaluoSpreads     string // 2~9
	TaluoUserChecked string // 牌号列表，英文逗号分隔
	Lang             string // zh-cn / en-us / zh-tw
}

func (r TaluospreadsRequest) toValues() url.Values {
	v := url.Values{}
	if r.TaluoSpreads != "" {
		v.Set("taluo_spreads", r.TaluoSpreads)
	}
	if r.TaluoUserChecked != "" {
		v.Set("taluo_user_checked", r.TaluoUserChecked)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r TaluospreadsRequest) Validate() error {
	if r.TaluoSpreads == "" {
		return newRequiredFieldError("taluo_spreads")
	}
	spreads, err := strconv.Atoi(r.TaluoSpreads)
	if err != nil {
		return &ValidationError{Field: "taluo_spreads", Message: "must be a valid integer"}
	}
	if spreads < 2 || spreads > 9 {
		return &ValidationError{Field: "taluo_spreads", Message: "must be in range [2, 9]"}
	}

	if r.TaluoUserChecked == "" {
		return newRequiredFieldError("taluo_user_checked")
	}

	parts := strings.Split(r.TaluoUserChecked, ",")
	cardNos := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		cardNos = append(cardNos, p)
	}
	if len(cardNos) != spreads {
		return &ValidationError{
			Field:   "taluo_user_checked",
			Message: "count must match taluo_spreads",
		}
	}

	for _, no := range cardNos {
		cardNo, err := strconv.Atoi(no)
		if err != nil {
			return &ValidationError{
				Field:   "taluo_user_checked",
				Message: "must contain valid integers",
			}
		}
		if cardNo < 0 || cardNo > 21 {
			return &ValidationError{
				Field:   "taluo_user_checked",
				Message: "card number must be in range [0, 21]",
			}
		}
	}

	if r.Lang != "" && !inSet(r.Lang, taluojieduAllowedLang) {
		return newEnumFieldError("lang", r.Lang, taluojieduAllowedLang)
	}
	return nil
}

type TaluospreadsData struct {
	Position string               `json:"position"`
	Image    string               `json:"image"`
	CardInfo TaluospreadsCardInfo `json:"card_info"`
}

type TaluospreadsCardInfo struct {
	CardReverse     string                      `json:"cart_reverse"`
	CardDescription TaluospreadsCardDescription `json:"card_description"`
	CardName        string                      `json:"card_name"`
	CardKeyword     string                      `json:"card_keyword"`
	CardAstrology   string                      `json:"card_astrology"`
	CardElements    string                      `json:"card_elements"`
	CardSummarize   string                      `json:"card_summarize"`
}

type TaluospreadsCardDescription struct {
	BaseDesc           string `json:"base_desc"`
	LoveMarriage       string `json:"love_marriage"`
	WorkStudy          string `json:"work_study"`
	InterWealth        string `json:"inter_wealth"`
	HealthLife         string `json:"health_life"`
	Other              string `json:"other"`
	Advice             string `json:"advice"`
	LoveMarriageAdvice string `json:"love_marriage_advice"`
	WorkStudyAdvice    string `json:"work_study_advice"`
	InterWealthAdvice  string `json:"inter_wealth_advice"`
	HealthLifeAdvice   string `json:"health_life_advice"`
}

type DivinationYunshiRequest struct {
	Type           string // 0 星座，1 生肖
	TitleYunshi    string // 0~11
	Lang           string // zh-cn / zh-tw / en-us
	ParameterStyle string // chinese / english
}

func (r DivinationYunshiRequest) toValues() url.Values {
	v := url.Values{}
	if r.Type != "" {
		v.Set("type", r.Type)
	}
	if r.TitleYunshi != "" {
		v.Set("title_yunshi", r.TitleYunshi)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	if r.ParameterStyle != "" {
		v.Set("parameter_style", r.ParameterStyle)
	}
	return v
}

func (r DivinationYunshiRequest) Validate() error {
	if r.Type == "" {
		return newRequiredFieldError("type")
	}
	if !inSet(r.Type, divinationYunshiAllowedType) {
		return newEnumFieldError("type", r.Type, divinationYunshiAllowedType)
	}

	if r.TitleYunshi == "" {
		return newRequiredFieldError("title_yunshi")
	}
	titleID, err := strconv.Atoi(r.TitleYunshi)
	if err != nil {
		return &ValidationError{Field: "title_yunshi", Message: "must be a valid integer"}
	}
	if titleID < 0 || titleID > 11 {
		return &ValidationError{Field: "title_yunshi", Message: "must be in range [0, 11]"}
	}

	if r.Lang != "" && !inSet(r.Lang, divinationYunshiAllowedLang) {
		return newEnumFieldError("lang", r.Lang, divinationYunshiAllowedLang)
	}
	if r.ParameterStyle != "" && !inSet(r.ParameterStyle, divinationYunshiAllowedParameterStyle) {
		return newEnumFieldError("parameter_style", r.ParameterStyle, divinationYunshiAllowedParameterStyle)
	}
	return nil
}

type ShengxiaoyunshiRequest struct {
	TitleYunshi    string // 0~11
	Lang           string // zh-cn / zh-tw / en-us
	ParameterStyle string // chinese / english
}

func (r ShengxiaoyunshiRequest) toYunshiRequest() DivinationYunshiRequest {
	return DivinationYunshiRequest{
		Type:           "1",
		TitleYunshi:    r.TitleYunshi,
		Lang:           r.Lang,
		ParameterStyle: r.ParameterStyle,
	}
}

func (r ShengxiaoyunshiRequest) Validate() error {
	return r.toYunshiRequest().Validate()
}

type DivinationYunshiData struct {
	FortuneType     string                        `json:"运势类型"`
	DailyFortune    DivinationYunshiDailyFortune  `json:"今日运势"`
	TomorrowFortune DivinationYunshiDailyFortune  `json:"明日运势"`
	WeeklyFortune   DivinationYunshiPeriodFortune `json:"本周运势"`
	MonthlyFortune  DivinationYunshiPeriodFortune `json:"本月运势"`
	YearlyFortune   DivinationYunshiPeriodFortune `json:"本年运势"`
}

type DivinationYunshiDailyFortune struct {
	CompatibleSign       string `json:"速配星座"`
	CautionarySign       string `json:"提防星座"`
	CompatibleZodiac     string `json:"速配生肖"`
	CautionaryZodiac     string `json:"提防生肖"`
	LuckyColor           string `json:"幸运颜色"`
	LuckyNumber          string `json:"幸运数字"`
	LuckyGem             string `json:"幸运宝石"`
	OverallScore         string `json:"综合分数"`
	LoveScore            string `json:"爱情分数"`
	CareerScore          string `json:"事业分数"`
	MoodScore            string `json:"心情分数"`
	SocialScore          string `json:"交际分数"`
	WealthScore          string `json:"财富分数"`
	HealthScore          string `json:"健康分数"`
	TodayTomorrowFortune string `json:"今明运势"`
	LoveFortune          string `json:"爱情运势"`
	CareerFortune        string `json:"事业运势"`
	WealthFortune        string `json:"财富运势"`
	HealthFortune        string `json:"健康运势"`
}

type DivinationYunshiPeriodFortune struct {
	CompatibleSign   string `json:"速配星座"`
	CautionarySign   string `json:"提防星座"`
	CompatibleZodiac string `json:"速配生肖"`
	CautionaryZodiac string `json:"提防生肖"`
	LuckyColor       string `json:"幸运颜色"`
	LuckyNumber      string `json:"幸运数字"`
	LuckyGem         string `json:"幸运宝石"`
	OverallScore     string `json:"综合分数"`
	LoveScore        string `json:"爱情分数"`
	CareerScore      string `json:"事业分数"`
	MoodScore        string `json:"心情分数"`
	SocialScore      string `json:"交际分数"`
	WealthScore      string `json:"财富分数"`
	HealthScore      string `json:"健康分数"`
	WeeklyFortune    string `json:"本周运势"`
	MonthlyFortune   string `json:"本月运势"`
	YearlyFortune    string `json:"本年运势"`
	LoveFortune      string `json:"爱情运势"`
	CareerFortune    string `json:"事业运势"`
	WealthFortune    string `json:"财富运势"`
	HealthFortune    string `json:"健康运势"`
}

func (s *DivinationService) Meiri(ctx context.Context, req MeiriRequest) (*CommonResponse[MeiriData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[MeiriData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/meiri", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Xiaoliuren(ctx context.Context, req XiaoliurenRequest) (*CommonResponse[XiaoliurenData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[XiaoliurenData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/xiaoliuren", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Zhiwen(ctx context.Context, req ZhiwenRequest) (*CommonResponse[ZhiwenData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[ZhiwenData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/zhiwen", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Yaogua(ctx context.Context, req YaoguaRequest) (*CommonResponse[YaoguaData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[YaoguaData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/yaogua", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Taluojiedu(ctx context.Context, req TaluojieduRequest) (*CommonResponse[TaluojieduData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[TaluojieduData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/taluojiedu", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Taluoxipai(ctx context.Context, req TaluoxipaiRequest) (*CommonResponse[TaluoxipaiData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[TaluoxipaiData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/taluoxipai", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Taluozhanbu(ctx context.Context, req TaluozhanbuRequest) (*CommonResponse[TaluozhanbuData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[TaluozhanbuData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/taluozhanbu", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Taluospreads(ctx context.Context, req TaluospreadsRequest) (*CommonResponse[[]TaluospreadsData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[[]TaluospreadsData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/taluospreads", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Yunshi(ctx context.Context, req DivinationYunshiRequest) (*CommonResponse[DivinationYunshiData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[DivinationYunshiData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/yunshi", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *DivinationService) Shengxiaoyunshi(ctx context.Context, req ShengxiaoyunshiRequest) (*CommonResponse[DivinationYunshiData], error) {
	yunshiReq := req.toYunshiRequest()
	if err := yunshiReq.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[DivinationYunshiData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/yunshi", yunshiReq.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}
