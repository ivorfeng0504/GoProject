package model

import "github.com/devfeel/mapper"

//领数字活动
type GuessMarketActivity struct {
	//主键Id
	GuessMarketActivityId int64

	//开奖结果（四位数字，0000-9999）
	Result string

	//期号 20180522
	IssueNumber string

	//开始时间（当前时间在BeginTime与EndTime之间才可以领取数字）
	BeginTime mapper.JSONTime

	//结束时间（当前时间在BeginTime与EndTime之间才可以领取数字）
	EndTime mapper.JSONTime

	//创建时间
	CreateTime mapper.JSONTime

	//开奖时间
	PublishTime mapper.JSONTime

	//是够已经公布结果
	IsPublish bool
}
