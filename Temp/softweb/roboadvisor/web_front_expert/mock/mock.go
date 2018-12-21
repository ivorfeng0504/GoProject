package mock

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"time"
)

// AppMock create app Mock
func AppMock() dotweb.Mock {
	m := dotweb.NewStandardMock()
	m.RegisterString("/mock", "mock data")

	//获取实战培训的模拟数据
	sznews := struct {
		Title        string;
		Summary      string;
		AudioUrl     string;
		VideoUrl     string;
		ClassUrl     string;
		ClassSummary string;
	}{Title: "Just go it", AudioUrl: "Listen Url", VideoUrl: "See me  See you", ClassSummary: "excellent", ClassUrl: "Roma Street No1"}

	m.RegisterJSON("/mockjson", sznews)

	var multList []*expertnews.ExpertNews_MultiMedia_List
	for i := 1; i <= 10; i++ {
		multInfo := new(expertnews.ExpertNews_MultiMedia_List)
		multInfo.ID = i
		multInfo.NewsType = 3
		multInfo.Summary = "目前操作策略上还是不变，这位置不追涨。只要趋势顶底指标中期和长期线一直在高位，就不能在阳线的时候去操作。60分钟级别的调整，总会来的，把握好操作节奏。回档支撑还是看2730和2770两个位置。"
		multInfo.Title = "红盘普涨迎国庆 再派红包犹未尽"
		multInfo.CoverImg = "http://static.emoney.cn/webupload/JRPTSoftManage/2018/9/anser_aeeCEBHd-fFHD-EaAI-aFDc-AaGIabIJAFGD.jpg"
		multInfo.AudioURL = "http://emoney.gensee.com/webcast/site/vod/play-a2b093ee6fb946d3bfd3fd9de81d038e"
		multInfo.TagInfo="[{\"TagID\":\"4\",\"TagName\":\"回顾教育\"}]"
		if i == 2 {
			multInfo.VideoPlayURL = "http://emoney.gensee.com/webcast/site/vod/play-a2b093ee6fb946d3bfd3fd9de81d038e"
			multInfo.TagInfo="[{\"TagID\":\"2\",\"TagName\":\"解读分析\"}]"
		}
		if i == 3 {
			multInfo.LiveURL = "http://emoney.gensee.com/webcast/site/vod/play-a2b093ee6fb946d3bfd3fd9de81d038e"
			multInfo.Live_StartTime, _ = time.Parse("2006-01-02 15:04:05", "2018-09-29 15:00:00")
			multInfo.Live_EndTime, _ = time.Parse("2006-01-02 15:04:05", "2018-09-29 19:00:00")
			multInfo.TagInfo="[{\"TagID\":\"3\",\"TagName\":\"机会甄选\"}]"
		}
		multInfo.CreateTime = time.Now()
		multInfo.LastModifyTime = time.Now()

		multList = append(multList, multInfo)
	}
	m.RegisterJSON("/strategy/GetStrategyNewsList_Multiple", multList)

	var tagList []*TagInfo
	for i := 2; i <= 4; i++ {
		tagInfo := new(TagInfo)
		tagname := ""
		tagInfo.TagID = i
		if i == 2 {
			tagname = "解读分析"
		}
		if i == 3 {
			tagname = "机会甄选"
		}
		if i == 4 {
			tagname = "回顾教育"
		}
		tagInfo.TagName = tagname

		tagList = append(tagList, tagInfo)
	}

	m.RegisterJSON("/strategy/GetTagList", tagList)
	return m
}

type TagInfo struct{
	TagID int
	TagName string
}

