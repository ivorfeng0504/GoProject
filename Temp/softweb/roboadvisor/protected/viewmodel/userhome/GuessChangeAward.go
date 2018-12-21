package viewmodel

import (
	"github.com/devfeel/mapper"
)

//猜涨跌奖品
type GuessChangeAward struct {
	//活动周期
	//6.18-6.22
	ActivityCycle string
	//猜中奖品名称
	AwardName string
	//开奖时间
	PublishTime mapper.JSONTime
	//研报
	ReportUrl string
}
