package expertnews

import (
	"github.com/devfeel/mapper"
	"time"
)

//type FocusStrategyInfo struct {
//	ExpertNewsStrategy *ExpertNews_StrategyInfo
//	LatestLive *ReciveLive
//	LatestNews *model.NewsInfo
//}

type FocusStrategyInfo struct {
	//策略ID
	StrategyID int
	//策略名称
	StrategyName string
	//策略图
	StrategyImg string
	//资讯类型（1：直播 2：资讯）
	NewsType int
	//资讯ID
	NewsID	int
	//直播ID
	LiveId int
	//标题
	Title string
	//简介
	Summary string
	//发布时间
	SendTime time.Time
}

type ReciveLive struct {
	LId int
	LiveName string
	LiveContent string
	LiveImg string
	SendTime mapper.JSONTime
}