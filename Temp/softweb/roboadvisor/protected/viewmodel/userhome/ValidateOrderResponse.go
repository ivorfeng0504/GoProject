package viewmodel

type ValidateOrderResponse struct {
	NickName       string
	ReFundMode     int
	ReFundModeDesc string
	BankInfoList   []*BankInfo
}

type BankInfo struct {
	//银行名称
	BankName string
	//值
	BankValue string
}
