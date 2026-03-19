package yuanfenju

import (
	"context"
	"net/url"
)

type FreeService struct {
	client *Client
}

type QueryMerchantData struct {
	MerchantType string `json:"merchant_type"`
	ExpireTime   string `json:"expire_time"`
	CanUseNum    string `json:"can_use_num"`
}

func (s *FreeService) QueryMerchant(ctx context.Context) (*CommonResponse[QueryMerchantData], error) {
	resp := &CommonResponse[QueryMerchantData]{}
	if err := s.client.doForm(ctx, "/v1/Free/querymerchant", nil, resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

type QueryTimesData struct {
	CallTimes         string `json:"call_times"`
	ExpireTime        int64  `json:"expire_time"`
	ExpireTimeMessage string `json:"expire_time_message"`
}

func (s *FreeService) QueryTimes(ctx context.Context) (*CommonResponse[QueryTimesData], error) {
	resp := &CommonResponse[QueryTimesData]{}
	if err := s.client.doForm(ctx, "/v1/Free/querytimes", nil, resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}

func (s *FreeService) QueryTimesWithParams(ctx context.Context, params url.Values) (*CommonResponse[QueryTimesData], error) {
	resp := &CommonResponse[QueryTimesData]{}
	if err := s.client.doForm(ctx, "/v1/Free/querytimes", params, resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, &APIError{Code: resp.ErrCode, Message: resp.ErrMsg, Notice: resp.Notice}
	}
	return resp, nil
}
