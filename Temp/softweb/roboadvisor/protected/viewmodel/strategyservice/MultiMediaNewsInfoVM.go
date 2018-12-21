package strategyservice

import (
	"github.com/devfeel/mapper"
	"time"
)

type ExpertNews_MultiMedia_List struct {
	//唯一标识
	ID int

	//文章类型（0：文章 1：单课程 2：系列课程 3:多媒体课程）
	NewsType int

	//标题
	Title string

	//简介
	Summary string

	//封面图
	CoverImg string

	//所属标签
	TagInfo string

	//浏览量
	ClickNum int64

	//音频播放地址
	AudioURL string

	//录播播放地址
	VideoPlayURL string

	//直播视频播放地址
	LiveURL string

	//直播视频录播地址
	LiveVideoURL string

	//直播课程开始时间
	Live_StartTime mapper.JSONTime

	//直播课程结束时间
	Live_EndTime mapper.JSONTime

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//显示位置、是否置顶（json格式：[{"ShowPosition":"1","IsTop":"1"},{"ShowPosition":"2","IsTop":"1"}]）
	ShowPosition string

	//来源
	From string
}

type ExpertNews_MultiMedia_Details struct {
	//唯一标识
	ID int

	//文章类型（0：文章 1：单课程 2：系列课程 3:多媒体课程）
	NewsType int

	//标题
	Title string

	//内容
	NewsContent string

	//简介
	Summary string

	//封面图
	CoverImg string

	//所属标签
	TagInfo string

	//浏览量
	ClickNum int64

	//音频播放地址
	AudioURL string

	//录播播放地址
	VideoPlayURL string

	//直播视频播放地址
	LiveURL string

	//直播视频录播地址
	LiveVideoURL string

	//直播课程开始时间
	Live_StartTime mapper.JSONTime

	//直播课程结束时间
	Live_EndTime mapper.JSONTime

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//显示位置、是否置顶（json格式：[{"ShowPosition":"1","IsTop":"1"},{"ShowPosition":"2","IsTop":"1"}]）
	ShowPosition string

	//来源
	From string
}

type ExpertNews_MultiMedia_TrainList struct {
	//唯一标识
	ID int

	//文章类型（0：文章 1：单课程 2：系列课程 3:多媒体课程）
	NewsType int

	//标题
	Title string

	//简介
	Summary string

	//封面图
	CoverImg string

	//所属标签
	TagInfo string

	//浏览量
	ClickNum int64

	//音频播放地址
	AudioURL string

	//录播播放地址
	VideoPlayURL string

	//直播视频播放地址
	LiveURL string

	//直播视频录播地址
	LiveVideoURL string

	//直播课程开始时间
	Live_StartTime time.Time

	//直播课程结束时间
	Live_EndTime time.Time

	//创建时间
	CreateTime time.Time

	//最后更新时间
	LastModifyTime time.Time

	//显示位置、是否置顶（json格式：[{"ShowPosition":"1","IsTop":"1"},{"ShowPosition":"2","IsTop":"1"}]）
	ShowPosition string

	//来源
	From string
}