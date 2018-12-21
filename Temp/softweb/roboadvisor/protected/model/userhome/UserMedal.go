package model

import "github.com/devfeel/mapper"

// UserMedal 用户勋章
type UserMedal struct {
	//主键Id
	UserMedalId int64

	//用户Id
	UserInfoId int

	//用户UID
	UID int64

	//勋章类型 详见MedalType.go
	MedalType int

	//勋章名称
	MedalName string

	//勋章等级
	MedalLevel int

	//创建时间
	CreateTime mapper.JSONTime

	//最后修改时间
	LastModifyTime mapper.JSONTime
}
