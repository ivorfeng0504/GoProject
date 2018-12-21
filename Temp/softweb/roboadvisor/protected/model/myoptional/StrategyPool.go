package myoptional

import "github.com/devfeel/mapper"

//策略股池
type StrategyPool struct {
	//主键
	StrategyPoolId int64
	//策略Id
	ClientStrategyId int
	//策略名称
	ClientStrategyName string
	//父级Id
	ParentId int
	//父级名称
	ParentName string
	//股票代码
	StockCode string
	//股票名称
	StockName string
	//期号20060102
	IssueNumber string
	//创建时间
	CreateTime mapper.JSONTime
	//是否删除
	IsDeleted bool
}
