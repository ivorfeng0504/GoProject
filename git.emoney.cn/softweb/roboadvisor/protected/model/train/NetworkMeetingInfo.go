package train

import "time"

type NetworkMeetingInfo struct{
	ID int `mapper:"Id"`
	//课程编号
	Mtg_id int `mapper:"mtg_id"`
	//课程名称
	Mtg_name string `mapper:"mtg_name"`
	//讲师名称
	Txtteachar string `mapper:"txtteachar"`
	//课程类型
	Meeting_type string `mapper:"meeting_type"`
	//直播地址
	Gensee_URL string `mapper:"Gensee_URL"`
	//录播地址
	Video_url string `mapper:"video_url"`
	//课程封面
	CoverImg string
	//课程开始时间
	Class_date time.Time `mapper:"class_date"`
	//课程结束时间
	EndDate time.Time `mapper:"endDate"`
	//地区
	Ddlarea string `mapper:"ddlarea"`
	//频道
	Meeting_channel int `mapper:"meeting_channel"`
	//标签
	TrainTag int
	//标签名称
	TrainTagName string
}

type PIDInfo struct {
	PID string
	PIDName string
}


type TrainClientInfo struct {
	//策略id
	ParentId string
	//策略名称
	ParentName string
	//显示名称
	ShowName string
}

type TrainTagInfo struct {
	//策略id
	TrainTagID string
	//策略名称
	TrainTagName string
	//显示名称
	ShowName string
	//分类 jpch:解盘晨会 zbxx：指标学习
	Type string
}