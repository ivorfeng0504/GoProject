package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	"github.com/devfeel/mapper"
)
//通用-专家资讯
type ExpertNews_StrategyNewsInfo struct{
	//资讯
	NewsInfo *model.NewsInfo

	//策略
	StrategyInfo *ExpertNews_StrategyInfo
}


//首页热门大咖接口返回model-精简
type ExpertNews_StrategyNewsInfo_index struct {
	NewsInfo *NewsInfo_index
	StrategyInfo *ExpertNews_StrategyInfo_index
}

type NewsInfo_index struct{
	//唯一标识
	ID int

	//标题
	Title string

	//简介
	Summary string

	//封面图
	CoverImg string

	//策略权限ID（全部：0   1,2,3其他权限）
	StrategyID string

	//策略权限ID（前端使用）
	StrategyInfoID int

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//专家资讯ID
	ExpertStrategyID int

	//点击浏览量
	ClickNum int64

	//视频播放量
	VideoPlayNum int64
}

type ExpertNews_StrategyInfo_index struct {
	//主键 自增长
	ID int

	//策略名称
	StrategyName string

	//策略简介
	StrategySummary string

	//策略图
	StrategyImg string

	//投顾编号
	StrategyTGNo string

	//圈圈直播类型
	LiveType int

	//圈圈直播ID
	LiveID int

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//点击浏览量
	ClickNum int64

	//策略是否有文章
	HasArticle bool

	//策略是否有视频
	HasVideo bool
}


//专家资讯详情接口返回model-精简
type ExpertNews_StrategyNewsInfo_detail struct {
	NewsInfo *NewsInfo_detail
	StrategyInfo *ExpertNews_StrategyInfo_detail
}
type NewsInfo_detail struct{
	//唯一标识
	ID int

	//文章类型（0：文章 1：单课程 2：系列课程）
	NewsType int

	//标题
	Title string

	//内容
	NewsContent string

	//简介
	Summary string

	//封面图
	CoverImg string

	//课程播放地址
	VideoUrl string

	//策略权限ID（全部：0   1,2,3其他权限）
	StrategyID string

	//策略权限ID（前端使用）
	StrategyInfoID int

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//专家资讯ID
	ExpertStrategyID int

	//点击浏览量
	ClickNum int64

	//视频播放量
	VideoPlayNum int64
}

type ExpertNews_StrategyInfo_detail struct {
	//主键 自增长
	ID int

	//策略名称
	StrategyName string

	//策略简介
	StrategySummary string

	//策略图
	StrategyImg string

	//投顾编号
	StrategyTGNo string

	//圈圈直播类型
	LiveType int

	//圈圈直播ID
	LiveID int

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//点击浏览量
	ClickNum int64

	//策略是否有文章
	HasArticle bool

	//策略是否有视频
	HasVideo bool
}
