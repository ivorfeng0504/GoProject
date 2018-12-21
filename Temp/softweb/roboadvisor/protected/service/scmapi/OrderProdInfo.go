package scmapi

type OrderProdInfo struct {
	//订单号码
	OrderId string `json:"ORDER_ID"`

	//客户名称
	CustomerName string `json:"CUSTOMER_NAME"`

	//密文电话号码
	MidPwd string `json:"MIDPWD"`

	//加*电话号码
	Mid string `json:"MID"`

	//EM账号
	EmCard string `json:"EmCard"`

	//录单时间
	OrdAddTime string `json:"ORDADD_TIME"`

	//备货日期
	StockUpDate string `json:"StockUpDate"`

	//产品线ID
	ProdLineId string `json:"ProdLineid"`

	//产品线名称
	ProdLineName string `json:"ProdLIneName"`

	//产品ID
	ProdId string `json:"PRODID"`

	//产品名称
	ProdName string `json:"PRODNAME"`

	//活动代码
	ActivityCode string `json:"ACTIVITY_CODE"`

	//活动名称
	ActivityName string `json:"ACTIVITY_NAME"`

	//订单金额
	SPrice float32 `json:"SPRICE"`

	//卡号
	OldProd string `json:"OLDPROD"`

	//退货状态 --0 未退货 -1 已退货；1 全部
	RefundSign int `json:"Refund_Sign"`

	//是否支持快速退货退款0：不支持；1：支持
	//IsQuickReturn int `json:"IsQuickReturn"`

	//是否支持快速退货退款0：不支持；1：支持
	IsRefundByCustomer int `json:"IsRefundByCustomer"`
}
