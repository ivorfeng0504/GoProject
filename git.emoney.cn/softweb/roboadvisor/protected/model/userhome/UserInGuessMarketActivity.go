package model

import "github.com/devfeel/mapper"

// UserInGuessMarketActivity 领数字活动
type UserInGuessMarketActivity struct {
	// 主键Id
	UserInGuessMarketActivityId int64
	// UID
	UID int64
	// 用户主键Id
	UserInfoId int

	//领取时间
	CreateTime mapper.JSONTime

	//期号 20180522
	IssueNumber string

	//领取数字（四位数字，0000-9999）
	Result string

	//昵称
	NickName string

	//是否猜中
	IsGuess bool

	//开奖时间
	PublishTime mapper.JSONTime

	//猜中奖品Id
	AwardId int64

	//猜中奖品名称
	AwardName string

	//本期猜数字活动Id
	GuessMarketActivityId int64
}
