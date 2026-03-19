package yuanfenju

import (
	"context"
	"encoding/json"
	"net/url"
)

type FreeService struct {
	client *Client
}

type QueryMerchantData struct {
	MerchantType               string `json:"merchant_type"`
	MerchantEmail              string `json:"merchant_email"`
	MerchantNickname           string `json:"merchant_nickname"`
	MerchantRegisterTime       string `json:"merchant_register_time"`
	MerchantExpireTime         string `json:"merchant_expire_time"`
	MerchantRemainingCallTimes string `json:"merchant_remaining_call_times"`

	// Deprecated: keep legacy aliases for backward compatibility.
	ExpireTime string `json:"-"`
	CanUseNum  string `json:"-"`
}

func (d *QueryMerchantData) UnmarshalJSON(data []byte) error {
	type alias QueryMerchantData
	type payload struct {
		alias
		LegacyExpireTime string `json:"expire_time"`
		LegacyCanUseNum  string `json:"can_use_num"`
	}

	apply := func(p payload) {
		*d = QueryMerchantData(p.alias)

		if d.MerchantExpireTime == "" {
			d.MerchantExpireTime = p.LegacyExpireTime
		}
		if d.MerchantRemainingCallTimes == "" {
			d.MerchantRemainingCallTimes = p.LegacyCanUseNum
		}

		d.ExpireTime = d.MerchantExpireTime
		d.CanUseNum = d.MerchantRemainingCallTimes
	}

	var obj payload
	if err := json.Unmarshal(data, &obj); err == nil {
		apply(obj)
		return nil
	}

	var arr []payload
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		*d = QueryMerchantData{}
		return nil
	}
	apply(arr[0])
	return nil
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
