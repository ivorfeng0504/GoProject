package scmapi

type ReturnbackAndRefundResponse struct {
	//状态码：0：成功  !=0：失败
	StatusCode int
	//消息
	Message string
	//退货执行结果
	//StatusCode：状态码：0：成功  !=0：失败
	//Message：消息
	ReturnbackValue *ReturnbackAndRefundResponseBase
	//退款执行结果
	//StatusCode：状态码：0：成功  !=0：失败
	//Message：消息
	RefundValue *ReturnbackAndRefundResponseBase
}

type ReturnbackAndRefundResponseBase struct {
	StatusCode int
	Message    string
}
