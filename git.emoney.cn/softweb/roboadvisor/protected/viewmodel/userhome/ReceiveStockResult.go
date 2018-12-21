package viewmodel

import userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"

type ReceiveStockResult struct {
	//连登天数
	LoginCountSerial int
	//股票列表
	StockList []*userhome_model.StockInfo
	//研报地址
	ReportUrl string
	//是否已经领取股票
	IsReceived bool
}
