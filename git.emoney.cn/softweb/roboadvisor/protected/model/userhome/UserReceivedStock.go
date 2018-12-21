package model

import "github.com/devfeel/mapper"

//用户领取的双响炮股票
type UserReceivedStock struct {
	//主键Id
	UserReceivedStockId int64

	//双响炮记录Id
	ReceiveStockActivityId int64

	//创建时间
	CreateTime mapper.JSONTime

	//股票池生成时间
	StockCreateTime mapper.JSONTime

	//期号20060102 股票推荐的目标日期
	IssueNumber string

	//股票列表
	StockList string

	//UserInfoId
	UserInfoId int

	//UID
	UID int64

	//研报地址
	ReportUrl string
}
