package model

import "github.com/devfeel/mapper"

// Activity 活动表
type Activity struct {
	ActivityId     int64
	Title          string
	Summary        string
	BeginTime      mapper.JSONTime
	EndTime        mapper.JSONTime
	ImageUrl       string
	ActivityUrl    string
	CreateUser     string
	CreateTime     mapper.JSONTime
	LastModifyTime mapper.JSONTime
	TargetUserType string
	IsDeleted      bool
	IsEnabled      bool

	//打开方式（0 当前页打开 1 浏览器新窗口外链）
	OpenMode int

	//是否需要绑定手机号 0 否 1 是
	NeedBind int

	//是否需要附加SSO 0 否 1 是
	NeedSSO int
}
