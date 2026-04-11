package yuanfenju

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

var toolsAllowedLang = []string{"zh-cn", "en-us"}

type ToolsService struct {
	client *Client
}

type ToolsQqRequest struct {
	Qq   string
	Lang string // zh-cn / en-us
}

func (r ToolsQqRequest) toValues() url.Values {
	v := url.Values{}
	if r.Qq != "" {
		v.Set("qq", r.Qq)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r ToolsQqRequest) Validate() error {
	if r.Qq == "" {
		return newRequiredFieldError("qq")
	}
	if r.Lang != "" && !inSet(r.Lang, toolsAllowedLang) {
		return newEnumFieldError("lang", r.Lang, toolsAllowedLang)
	}
	return nil
}

type ToolsShoujiRequest struct {
	Shouji string
	Lang   string // zh-cn / en-us
}

func (r ToolsShoujiRequest) toValues() url.Values {
	v := url.Values{}
	if r.Shouji != "" {
		v.Set("shouji", r.Shouji)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r ToolsShoujiRequest) Validate() error {
	if r.Shouji == "" {
		return newRequiredFieldError("shouji")
	}
	if r.Lang != "" && !inSet(r.Lang, toolsAllowedLang) {
		return newEnumFieldError("lang", r.Lang, toolsAllowedLang)
	}
	return nil
}

type ToolsShuziRequest struct {
	Shuzi string
	Lang  string // zh-cn / en-us
}

func (r ToolsShuziRequest) toValues() url.Values {
	v := url.Values{}
	if r.Shuzi != "" {
		v.Set("shuzi", r.Shuzi)
	}
	if r.Lang != "" {
		v.Set("lang", r.Lang)
	}
	return v
}

func (r ToolsShuziRequest) Validate() error {
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
	if r.Lang != "" && !inSet(r.Lang, toolsAllowedLang) {
		return newEnumFieldError("lang", r.Lang, toolsAllowedLang)
	}
	return nil
}

type ToolsJixiongData struct {
	Desc      string `json:"desc"`
	Xiongji   string `json:"xiongji"`
	Desc1     string `json:"desc1"`
	Score     string `json:"score"`
	InputData string `json:"data"`
	Shuli     int    `json:"shuli"`
}

func (d *ToolsJixiongData) UnmarshalJSON(data []byte) error {
	type rawPayload struct {
		Desc    json.RawMessage `json:"desc"`
		Xiongji json.RawMessage `json:"xiongji"`
		Desc1   json.RawMessage `json:"desc1"`
		Score   json.RawMessage `json:"score"`
		Data    json.RawMessage `json:"data"`
		Shuli   json.RawMessage `json:"shuli"`
	}

	var raw rawPayload
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	desc, err := parseJSONStringOrNumber(raw.Desc)
	if err != nil {
		return err
	}
	xiongji, err := parseJSONStringOrNumber(raw.Xiongji)
	if err != nil {
		return err
	}
	desc1, err := parseJSONStringOrNumber(raw.Desc1)
	if err != nil {
		return err
	}
	score, err := parseJSONStringOrNumber(raw.Score)
	if err != nil {
		return err
	}
	inputData, err := parseJSONStringOrNumber(raw.Data)
	if err != nil {
		return err
	}
	shuli, err := parseJSONInt(raw.Shuli)
	if err != nil {
		return err
	}

	*d = ToolsJixiongData{
		Desc:      desc,
		Xiongji:   xiongji,
		Desc1:     desc1,
		Score:     score,
		InputData: inputData,
		Shuli:     shuli,
	}
	return nil
}

func (s *ToolsService) Qq(ctx context.Context, req ToolsQqRequest) (*CommonResponse[ToolsJixiongData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[ToolsJixiongData]{}
	if err := s.client.doForm(ctx, "/v1/Jixiong/qq", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *ToolsService) Shouji(ctx context.Context, req ToolsShoujiRequest) (*CommonResponse[ToolsJixiongData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[ToolsJixiongData]{}
	if err := s.client.doForm(ctx, "/v1/Jixiong/shouji", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *ToolsService) Shuzi(ctx context.Context, req ToolsShuziRequest) (*CommonResponse[ToolsJixiongData], error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &CommonResponse[ToolsJixiongData]{}
	if err := s.client.doForm(ctx, "/v1/Jixiong/shuzi", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}
