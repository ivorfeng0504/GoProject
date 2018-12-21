package model

import "github.com/devfeel/mapper"

//猜涨跌活动表
type GuessChangeActivity struct {

	//主键Id
	GuessChangeActivityId int64

	//期号 20180522
	IssueNumber string

	//开始时间
	BeginTime mapper.JSONTime

	//结束时间
	EndTime mapper.JSONTime

	//创建时间
	CreateTime mapper.JSONTime

	//开奖结果 1涨 -1跌 0无变化
	Result int

	//当前竞猜结果是否已经公布
	IsPublish bool

	//当前竞猜结果公布时间
	PublishTime mapper.JSONTime

	//活动周期
	//6.18-6.22
	ActivityCycle string
}
