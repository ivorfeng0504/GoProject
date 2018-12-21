package model

import "github.com/devfeel/mapper"

type UserAward struct {
	//主键Id
	UserAwardId int64
	//UserInfo主键Id
	UserInfoId int
	//奖品名称
	AwardName string
	//奖品Id
	AwardId int64
	//奖品图片
	AwardImg string
	//创建时间
	CreateTime mapper.JSONTime
	//介绍视频
	IntroduceVideo string
	//活动Id
	ActivityId int64
	//活动名称
	ActivityName int64
	//奖品发放状态 详情见UserAwardState.go
	State int
	//有效期
	AvailableDay int
	//奖品类型
	AwardType int
}
