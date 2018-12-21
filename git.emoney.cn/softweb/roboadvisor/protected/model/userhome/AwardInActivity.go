package model

import "github.com/devfeel/mapper"

//活动所属的奖品
type AwardInActivity struct {

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
	LastModifyTime mapper.JSONTime

	//创建人
	CreateUser string

	//是否删除
	IsDeleted bool

	//奖品详情-非数据库字段
	AwardInfo Award
}
