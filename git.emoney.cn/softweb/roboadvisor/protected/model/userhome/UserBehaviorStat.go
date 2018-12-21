package model

import "github.com/devfeel/mapper"

type UserBehaviorStat struct {
	UserBehaviorStatId int64
	StatFrom           int
	UID                int64
	ThirdIdentity      string
	CreateTime         mapper.JSONTime
}
