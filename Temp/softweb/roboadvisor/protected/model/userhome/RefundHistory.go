package model

import "github.com/devfeel/mapper"

//用户中心退款记录表
type RefundHistory struct {
	//主键Id
	RefundHistoryId int64
	//账号
	Account string
	//加密手机号
	Mobile string
	//订单Id
	OrderId string
	//退款原因
	Reason string
	//退款方式（1：银行退款，2：原路退款）
	RefundMode int
	//退款提交数据 退款接口的报文数据原文
	SubmitData string
	//创建时间
	CreateTime mapper.JSONTime
}
