package contract

import "time"

type ResponseInfo struct {

	//响应代码，0表示正常，小于0为服务出错代码，大于0为业务代码
	RetCode int

	// 响应代码对应的提示消息
	RetMsg string

	//服务返回的数据
	Message interface{}

	//总数量
	TotalCount int

	//系统时间
	SystemTime time.Time
}

func NewResonseInfo() *ResponseInfo {
	return &ResponseInfo{}
}

func CreateResponse(retCode int, retMsg string, message interface{},totalCount int) *ResponseInfo {
	return &ResponseInfo{
		RetCode:    retCode,
		RetMsg:     retMsg,
		Message:    message,
		TotalCount: totalCount,
	}
}
