package scmapi

type ValidateReturnbackAndRefundResponse struct {
	//状态码：0：成功  !=0：失败
	StatusCode int

	//消息
	Message string

	//是否支持快速退款(1:是，0:否)
	//IsQuickRefund bool

	//是否支持原路退款(1:是，0:否)
	IsOriginRefund bool

	// 三方退款支付列表
	PaymentDetailList []*ReturnbackAndRefundPaymentDetail

	// 是否支持快速退货退款0：不支持；1：支持
	IsRefundByCustomer bool
}
