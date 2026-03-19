package yuanfenju

import (
	"context"
	"encoding/json"
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

// BaziPaipanData 使用结构体建模关键字段；复杂扩展字段用 json.RawMessage 承接，避免 map 类型。
type BaziPaipanData struct {
	BaseInfo     BaziBaseInfo    `json:"base_info"`
	BaziInfo     json.RawMessage `json:"bazi_info"`
	FiveElements json.RawMessage `json:"five_elements"`
	Luck         json.RawMessage `json:"luck"`
	Shensha      json.RawMessage `json:"shensha"`
	Nayin        json.RawMessage `json:"nayin"`
	Extra        json.RawMessage `json:"extra"`
}

type BaziBaseInfo struct {
	Zhen    *BaziZhenInfo `json:"zhen,omitempty"`
	Sex     string        `json:"sex"`
	Name    string        `json:"name"`
	Gongli  string        `json:"gongli"`
	Nongli  string        `json:"nongli"`
	Qiyun   string        `json:"qiyun"`
	Jiaoyun string        `json:"jiaoyun"`
}

type BaziZhenInfo struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Jingdu   string `json:"jingdu"`
	Weidu    string `json:"weidu"`
	Shicha   string `json:"shicha"`
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
