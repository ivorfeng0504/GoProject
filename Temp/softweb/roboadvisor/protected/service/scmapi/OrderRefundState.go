package scmapi

type OrderRefundStateResponse struct {
	//状态码：0：成功  !=0：失败
	StatusCode int
	//消息
	Message         string
	OrderRefundList []*OrderRefundState
}

//物流订单退款状态
type OrderRefundState struct {
	//订单号
	OrderId string

	//退货申请ID
	RetApplyId string

	//退款来源应用ID
	FromAppId string

	//退款方式（1：银行退款，2：原路退款）
	ReturnMode int

	//退款方式
	ReturnModeName string

	// 退款类型（支付宝，微信，招商银行...）
	RefundType string

	//退款账号（银行退款为银行账号）
	RefundCode string

	//实退金额
	RealbackPrice float32

	//申请时间
	CreateTime string

	//退款状态（待退款:0,退款中:1,退款成功:2,退款失败:3,退款驳回:4）
	RefundState int
}

//GetRefundStateDesc 获取退款状态描述
func GetRefundStateDesc(state int) string {
	stateDescDict := make(map[int]string)
	stateDescDict[0] = "待退款"
	stateDescDict[1] = "退款中"
	stateDescDict[2] = "退款成功"
	stateDescDict[3] = "退款失败"
	stateDescDict[4] = "退款驳回"
	return stateDescDict[state]
}
