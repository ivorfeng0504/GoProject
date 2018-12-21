package model

import "github.com/devfeel/mapper"

type UserInActivity struct {
	UserInActivityId int64
	ActivityId       int64
	UserInfoId       int
	CreateTime       mapper.JSONTime
}
