package viewmodel

import "github.com/devfeel/mapper"

//猜涨跌历史竞猜结果
type GuessChangeResult struct {
	//期号 20180522
	IssueNumber string
	//猜测结果 1涨 -1跌 0无变化
	Result int
	//看涨  看跌
	ResultDesc string
	//参与时间
	CreateTime mapper.JSONTime
	//是否猜中
	IsGuessed bool
	//当前竞猜结果是否已经公布
	IsPublish bool
	//参与状态描述 未参与 未开始 未猜中 猜中 未开奖
	StateDesc string
}
