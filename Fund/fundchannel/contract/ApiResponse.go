package contract

import (
	"emoney.cn/fundchannel/util/json"
	"emoney.cn/fundchannel/util/crypto"
)

//Api接口响应结果
type ApiResponse struct {
	ResponseInfo

	//客户端是否已经缓存,并且与服务器的数据一致
	ClientCached bool

	//响应数据Message的唯一Hash值
	MessageHash string
}

func NewApiResponse() *ApiResponse {
	return &ApiResponse{}
}

func GetResponseMessageHash(responseMessage interface{}) string {
	hash := _crypto.MD5(_json.GetJsonString(responseMessage))
	return hash
}

//SupportClientCache 处理ApiResponse，如果客户端启用缓存，并且客户端的数据的Hash与服务端一致，则不响应给客户端具体的Message
func SupportClientCache(request *ApiRequest, response *ApiResponse) {
	if request == nil || response == nil || request.AllowClientCache == false {
		return
	}
	messageHash := GetResponseMessageHash(response.Message)
	response.MessageHash = messageHash
	if request.MessageHash == messageHash && messageHash != "" {
		response.ClientCached = true
		response.Message = nil
	}
}
