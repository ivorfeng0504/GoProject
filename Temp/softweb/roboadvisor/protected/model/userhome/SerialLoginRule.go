package model

import "github.com/devfeel/mapper"

type SerialLoginRule struct {
	SerialLoginRuleId int64
	LoginDay          int
	AwardId           int64
	AwardName         string
	TargetUserType    int
	PID               string
	UserPrompt        string
	ActivityUrl       string
	CreateUser        string
	CreateTime        mapper.JSONTime
	LastModifyTime    mapper.JSONTime
	IsDeleted         bool
}
