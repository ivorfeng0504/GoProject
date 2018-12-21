package scmapi

type ReturnbackAndRefundRequest struct {
	//需要退款订单ID
	OrderId string
	// 调用方
	// a001：物流 a002：官网 a003：智投
	//参见 RefundAppId.go
	AppId string
	//退货原因
	Reason string
	// 退货产品列表
	//例:
	//[’xxxxxxxxxxx’,’xxxxxxxxx’]
	ReturnbackProd []string
	//退货流程类型（1：常规流程，2：快速流程）
	ReturnFlow int
	//退款方式（1：银行退款，2：原路退款）
	ReturnMode int
	//退款申请银行信息
	OrderRetBankInfo OrderRetBankInfoDo
	//是否整单退 1表示整单退
	IsAllReturn int
}

type OrderRetBankInfoDo struct {
	//退款账号
	ReProCode string `json:"RE_PRO_CODE"`
	//收款人姓名
	ReProName string `json:"RE_PRO_NAME"`
	//开户行省份
	ReBankProvince string `json:"RE_BANK_PROVINCE"`
	//开户行城市
	ReBankCity string `json:"RE_BANK_CITY"`
	//开户行网点
	ReBankArea string `json:"RE_BANK_AREA"`
	//开户行代码
	ReBankCode string `json:"RE_BANK_CODE"`
	//开户行
	ReBankName string `json:"RE_BANK_NAME"`
}
