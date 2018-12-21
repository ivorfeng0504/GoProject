package model

import "github.com/devfeel/mapper"

type Award struct {
	AwardId        int64
	AwardName      string
	Summary        string
	IntroduceVideo string
	//奖品图片
	AwardImg string
	AwardType      int
	AvailableDay   int
	JRPTFunc       int
	LinkUrl        string
	QQ             string
	IsDeleted      bool
	CreateUser     string
	CreateTime     mapper.JSONTime
	LastModifyTime mapper.JSONTime
}
