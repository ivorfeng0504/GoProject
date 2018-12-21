package contract

import (
	"emoney.cn/fundchannel/util/json"
	"emoney.cn/fundchannel/util/crypto"
)

//Api接口请求数据
type ApiRequest struct {

	//是否允许客户端缓存
	AllowClientCache bool

	//响应数据Message的唯一Hash值
	MessageHash string

	//客户端业务请求数据
	RequestData interface{}
}

func NewApiRequest() *ApiRequest {
	return &ApiRequest{}
}

// GetRequestHash 构建请求的Hash键值
func GetRequestHash(url string, requestData interface{}) string {
	requestKey := "API_Request_" + url + "|" + _json.GetJsonString(requestData)
	hash := _crypto.MD5(requestKey)
	return hash
}