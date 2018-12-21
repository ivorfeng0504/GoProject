package viewmodel

import (
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
)

//用户领取的双响炮股票
type UserReceivedStockViewModel struct {

	//创建时间
	CreateTime string

	//期号20060102 股票推荐的目标日期
	IssueNumber string

	//股票列表
	StockList []*userhome_model.StockInfo

	//研报地址
	ReportUrl string
}
