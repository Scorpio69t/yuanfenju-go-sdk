package yuanfenju

import (
	"context"
	"net/url"
	"strconv"
)

var divinationAllowedLang = []string{"zh-cn", "en-us"}

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
