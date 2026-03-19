package yuanfenju

import (
	"context"
	"encoding/json"
	"net/url"
)

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

type MeiriData struct {
	Number      int             `json:"number"`
	GuaMing     string          `json:"guaming"`
	Description json.RawMessage `json:"description"`
}

func (s *DivinationService) Meiri(ctx context.Context, req MeiriRequest) (*CommonResponse[MeiriData], error) {
	resp := &CommonResponse[MeiriData]{}
	if err := s.client.doForm(ctx, "/v1/Zhanbu/meiri", req.toValues(), resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}
