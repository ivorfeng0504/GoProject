package model

import "github.com/devfeel/mapper"

//活动奖品视图
type ActivityAward struct {
	AwardName      string
	Summary        string
	IntroduceVideo string
	//奖品图片
	AwardImg     string
	AwardType    int
	AvailableDay int
	JRPTFunc     int
	LinkUrl      string
	QQ           string

	//主键Id
	AwardInActivityId int64
	//活动唯一标识
	ActivityId int64
	//奖品唯一标识
	AwardId int64
	//开始时间
	BeginTime mapper.JSONTime
	//结束时间
	EndTime mapper.JSONTime
	//创建时间
	CreateTime mapper.JSONTime
	//最近修改时间
	//LastModifyTime mapper.JSONTime
	//创建人
	//CreateUser string
	//是否删除
	//IsDeleted bool
}
