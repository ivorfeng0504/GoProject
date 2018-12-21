package viewmodel

import "git.emoney.cn/softweb/roboadvisor/agent"

type RefundOrderInfoDetail struct {
	//订单号
	OrderId string
	//产品列表
	ProductList      []*agent.RefundProductInfo
	RefundDetailList []*RefundDetail
}

type RefundDetail struct {
	//订单号
	OrderId string

	//退款名称
	ReturnModeName string

	//退款方式
	RefundTypeDesc string

	//退款金额
	RealbackPrice float32

	//退款状态描述
	RefundStatusDesc string

	//退款时间
	RefundTime string
}
