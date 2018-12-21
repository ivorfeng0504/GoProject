package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/time"
	"github.com/devfeel/dotweb"
	"strconv"
	"time"
)

// GetTodayNews 获取今日头条
func GetTodayNews(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsList, err := agent.GetTodayNews()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetClosingNewsTop2 获取盘后预测Top2
func GetClosingNewsTop2(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsList, err := agent.GetClosingNews()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	if len(newsList) > 2 {
		newsList = newsList[:2]
	}
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetClosingNews 获取盘后预测资讯
func GetClosingNews(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsList, err := agent.GetClosingNews()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetNewsInfo 获取指定日期的要闻
func GetNewsInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	dateStr := _http.GetRequestValue(ctx, "Date")
	date, err := _time.ParseTime(dateStr)
	if err != nil {
		date = time.Now()
	}
	newsList, err := agent.GetNewsInfo(&date)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetNewsInfoTopN 获取最新的N条要闻
func GetNewsInfoTopN(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsList, err := agent.GetNewsInfoTopN()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetHotNewsInfo 获取热门资讯
func GetHotNewsInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsList, err := agent.GetHotNewsInfo()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetTopicNewsInfo 获取主题的相关资讯
func GetTopicNewsInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	topicIdStr := _http.GetRequestValue(ctx, "TopicId")
	topicId, err := strconv.Atoi(topicIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "TopicId不正确"
		return ctx.WriteJson(response)
	}
	newsList, err := agent.GetTopicNewsInfo(topicId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetNewsDetail 获取资讯详情
func GetNewsDetail(ctx dotweb.Context) error {
	detail := agent.ExpertNewsInfoDetail{}
	response := contract.NewResonseInfo()
	newsTypeStr := _http.GetRequestValue(ctx, "NewsType")
	newsType, err := strconv.Atoi(newsTypeStr)
	if err != nil {
		newsType = 0
	}
	newsInfoIdStr := _http.GetRequestValue(ctx, "NewsInfoId")
	newsInfoId, err := strconv.ParseInt(newsInfoIdStr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "NewsInfoId不正确"
		return ctx.WriteJson(response)
	}
	newInfoDetail, err := agent.GetNewsInfoDetail(newsInfoId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取资讯详情失败"
		return ctx.WriteJson(response)
	}
	if newInfoDetail == nil {
		response.RetCode = -3
		response.RetMsg = "该资讯不存在"
		return ctx.WriteJson(response)
	}
	//当前资讯详情
	detail.CurrentNews = newInfoDetail
	var newsList []*agent.ExpertNewsInfo
	switch newsType {
	case _const.NewsInfoType_ClosingNews:
		newsList, err = agent.GetClosingNews()
		if err != nil {
			newsList = nil
		}
		break
	case _const.NewsInfoType_NewsInfo:
		dateStr := _http.GetRequestValue(ctx, "Date")
		date, err := _time.ParseTime(dateStr)
		if err != nil {
			date = time.Now()
		}
		newsList, err = agent.GetNewsInfo(&date)
		if err != nil {
			newsList = nil
		}
		break
	case _const.NewsInfoType_HotNewsInfo:
		newsList, err = agent.GetHotNewsInfo()
		if err != nil {
			newsList = nil
		}
		break
	case _const.NewsInfoType_TopicNewsInfo:
		topicIdStr := _http.GetRequestValue(ctx, "TopicId")
		topicId, err := strconv.Atoi(topicIdStr)
		if err != nil {
			newsList = nil
			break
		}
		newsList, err = agent.GetTopicNewsInfo(topicId)
		if err != nil {
			newsList = nil
		}
		break
	case _const.NewsInfoType_NewsInfoTopN:
		newsList, err = agent.GetNewsInfoTopN()
		if err != nil {
			newsList = nil
		}
		break
	}
	currentIndex := -1
	newsListLen := 0
	if newsList != nil {
		newsListLen = len(newsList)
		for index, news := range newsList {
			if news.NewsInformationId == newsInfoId {
				currentIndex = index
			}
		}
	}

	if currentIndex >= 0 {
		if currentIndex == 0 {
			detail.PreviousNews = nil
			if newsListLen > 1 {
				detail.NextNews = newsList[currentIndex+1]
			}
		} else if currentIndex == newsListLen-1 {
			detail.PreviousNews = newsList[currentIndex-1]
			detail.NextNews = nil
		} else {
			detail.PreviousNews = newsList[currentIndex-1]
			detail.NextNews = newsList[currentIndex+1]
		}
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = detail

	go func() {
		//异步添加点击量
		agent.AddClick(_const.ClickType_NewsInfo, newsInfoIdStr)
	}()
	return ctx.WriteJson(response)
}

// NotifyUpdate 通知更新资讯数据库
func NotifyUpdate(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	tableName := ctx.QueryString("tableName")
	encrypt := ctx.QueryString("encrypt")
	if len(tableName) == 0 {
		response.RetCode = -1
		response.RetMsg = "tableName不能为空"
		return ctx.WriteJson(response)
	}
	if len(encrypt) == 0 {
		response.RetCode = -2
		response.RetMsg = "encrypt不能为空"
		return ctx.WriteJson(response)
	}
	if encrypt != "2721f23b-cc2d-4c00-9d68-f2fbc7713075" {
		response.RetCode = -3
		response.RetMsg = "授权失败，禁止访问"
		return ctx.WriteJson(response)
	}
	err := agent.NotifyUpdate()
	if err != nil {
		response.RetCode = -4
		response.RetMsg = "通知失败 " + err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}
