package scmapi

//订单列表查询参数
type QueryOrderProdListByParamsRequest struct {
	//产品线ID
	ProdLineId string `json:"ProdLineid"`
	//产品ID
	ProdId string `json:"PRODID"`
	//活动代码
	ActivityCode string `json:"ACTIVITY_CODE"`
	//EM账号
	EmCard string `json:"EmCard"`
	//密文电话号码
	MidPwd string `json:"MIDPWD"`
	//开始时间
	OrdAddTimeStart string `json:"ORDADD_TIME_Start"`
	//结束时间
	OrdAddTimeEnd string `json:"ORDADD_TIME_End"`
	//开始时间-备货成交时间
	StockUpDateStart string `json:"StockUpDate_Start"`
	//结束时间-备货成交时间
	StockUpDateEnd string `json:"StockUpDate_End"`
	//退货状态
	Refund_Sign string `json:"Refund_Sign"`
	//成交状态 A9002
	WitchState string `json:"WITCHSTATE"`
}
