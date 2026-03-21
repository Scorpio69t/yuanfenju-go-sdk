package yuanfenju

import (
	"context"
	"net/url"
	"strconv"
)

var divinationAllowedLang = []string{"zh-cn", "en-us"}
var yaoguaAllowedLang = []string{"zh-cn", "en-us", "zh-tw"}

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
