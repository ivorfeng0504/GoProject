package expertnews

import "github.com/devfeel/mapper"

type ExpertNews_StrategyInfo struct {
	//主键 自增长
	ID int

	//策略名称
	StrategyName string

	//策略简介
	StrategySummary string

	//策略标签
	StrategyTags string

	//策略图
	StrategyImg string

	//投顾编号
	StrategyTGNo string

	//圈圈直播类型
	LiveType int

	//圈圈直播ID
	LiveID int

	//是否删除
	IsDeleted bool

	//创建人
	CreateUser string

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//标签列表
	TagsList []StrategyTags

	//点击浏览量
	ClickNum int64

	//点播播放量
	VideoPlayNum int64

	//关注策略的粉丝数
	FansNum int64

	//策略是否有文章
	HasArticle bool

	//策略是否有视频
	HasVideo bool
}

type StrategyTags struct {
	//标签
	Tags string
}
