package strategyservice

import "github.com/devfeel/mapper"

type NewsInfo struct {

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

	//系列课程所属的课程ID
	//SeriesNewsIDs string

	//策略权限ID（全部：0   1,2,3其他权限）
	//StrategyID string

	//策略权限ID（前端使用）
	//StrategyInfoID int

	//显示位置、是否置顶（json格式：[{"ShowPosition":"1","IsTop":"1"},{"ShowPosition":"2","IsTop":"1"}]）
	ShowPosition string

	//是否置顶
	IsTop bool

	//是否删除
	//IsDeleted bool

	//创建用户
	//CreateUser string

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//LastModifyTime_fomat string

	//所属版本
	//Pid string

	//专家资讯ID
	//ExpertStrategyID int

	//点击浏览量
	ClickNum int64

	//来源
	From string

	//客户端策略ID
	StrategyID string

	//资讯类型（1：重要提示 2解盘文章）
	ToolTips string
}
