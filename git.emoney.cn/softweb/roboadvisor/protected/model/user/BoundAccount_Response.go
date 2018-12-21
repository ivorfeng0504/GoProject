package model

type BoundAccount_Response struct {
	//账号显示名
	AccountName string
	//CustomerID
	CustomerID int64
	//账号类型  0: em帐号; 1: 手机号; 2: 微信帐号; 3: QQ帐号
	AccountType int64
	//加密手机
	EncryptMobile string

}
