package model

import "github.com/devfeel/mapper"

//双响炮活动
type ReceiveStockActivity struct {
	//主键Id
	ReceiveStockActivityId int64

	//创建时间
	CreateTime mapper.JSONTime

	//期号20060102 股票推荐的目标日期
	IssueNumber string

	//股票列表 JSON格式
	StockList string

	//研报地址
	ReportUrl string
}
