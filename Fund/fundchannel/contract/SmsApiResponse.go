package contract

import (
	"time"
)

type WebApiResponse struct {

	//响应代码，0表示正常，小于0为服务出错代码，大于0为业务代码
	RetCode string

	// 响应代码对应的提示消息
	RetMsg string

	//服务返回的数据
	Message interface{}
}

type SendCodeMessage struct{
	//验证码
	Msg string

	//过期时间
	Expire time.Time

	//下次发送间隔时间
	NextTime time.Time
}