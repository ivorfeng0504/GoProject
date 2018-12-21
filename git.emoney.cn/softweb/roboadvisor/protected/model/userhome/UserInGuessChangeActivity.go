package model

import "github.com/devfeel/mapper"

//用户参与猜涨跌活动
type UserInGuessChangeActivity struct {

	//主键Id
	UserInGuessChangeActivityId int64

	// UID
	UID int64

	// 用户主键Id
	UserInfoId int

	//昵称
	NickName string

	//参与时间
	CreateTime mapper.JSONTime

	//期号 20180522
	IssueNumber string

	//猜测结果 1涨 -1跌 0无变化
	Result int

	//看涨  看跌
	ResultDesc string

	//是否猜中
	IsGuessed bool

	//当前竞猜结果是否已经公布
	IsPublish bool

	//当前竞猜结果公布时间
	PublishTime mapper.JSONTime

	//猜中奖品Id
	AwardId int64

	//猜中奖品名称
	AwardName string

	//研报
	ReportUrl string

	//活动周期
	//6.18-6.22
	ActivityCycle string
}
