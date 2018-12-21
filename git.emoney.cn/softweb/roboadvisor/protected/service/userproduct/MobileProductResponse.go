package userproduct

type MobileProductResponse struct {
	Protocol string              `json:"protocol"`
	Result   MobileApiResult     `json:"result"`
	Detail   []MobileProductInfo `json:"detail"`
}

type MobileApiResult struct {
	Code        int   `json:"code"`
	UpdateTime  int64 `json:"updateTime"`
	MinInterval int   `json:"minInterval"`
}

type MobileProductInfo struct {
	AuthName   string `json:"authName"`
	CreateTime int64  `json:"createTime"`
	ExpiryTime int64  `json:"expiryTime"`
}
