package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/panwai"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dottask"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"encoding/json"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
)

// Task_UpdateHotNewsInfoByClickNum 每日定时更新热门策略文章
func Task_UpdateHotNewsInfoByClickNum(context *task.TaskContext) error {
	global.TaskLogger.InfoFormat("【Start Task】 【Task_UpdateHotNewsInfoByClickNum】 每日定时更新热门策略文章")

	colid := config.CurrentConfig.ColumnID
	intColID, err := strconv.Atoi(colid)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task doing】 【Task_UpdateHotNewsInfoByClickNum】 处理被动更新队列 栏目id转换int失败 tableName=【%s】", colid)
	} else {
		newsSrv := service.NewNewsInfoService()
		newsList, err := newsSrv.GetNewsListByClicknum(intColID)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_UpdateHotNewsInfoByClickNum】 处理被动更新队列 执行异常 tableName=【%s】", newsList)
		} else {
			global.TaskLogger.InfoFormat("【Task Finished】 【Task_UpdateHotNewsInfoByClickNum】 处理被动更新队列 执行完毕 tableName=【%s】 newsList=【%s】", newsList, _json.GetJsonString(newsList))
		}
	}

	return err
}

// Task_UpdateTopicInfo 每日定时更新主题相关板块个股行情信息
func Task_UpdateTopicInfo(context *task.TaskContext) error {
	global.TaskLogger.InfoFormat("【Start Task】 【Task_UpdateTopicInfo】 每日定时更新主题相关板块个股行情信息")

	topicSrv := expertnews.NewExpertNews_TopicService()
	topicList, err := topicSrv.GetTopicList()

	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task doing】 【Task_UpdateTopicInfo】 获取主题列表失败")
	} else {
		for i, _ := range topicList {
			topicInfo := topicList[i]


			if topicInfo != nil && topicInfo.ID != 0 {
				var topic_bk []*expertnews2.Topic_BK
				var topic_stock []*expertnews2.Topic_Stock
				err = json.Unmarshal([]byte(topicInfo.RelatedBKInfo), &topic_bk)
				err = json.Unmarshal([]byte(topicInfo.RelatedStockInfo), &topic_stock)
				topicInfo.RelatedBKList = topic_bk
				topicInfo.RelatedStockList = topic_stock

				_, err = topicSrv.GetTopicInfoInnerQuotes(topicInfo)
			}
		}
	}

	return err
}

// Task_RefreshMultiMediaNewsList_1Month 每日定时更新最近一个月多媒体课程
func Task_RefreshMultiMediaNewsList_1Month(context *task.TaskContext) error {
	global.TaskLogger.InfoFormat("【Start Task】 【Task_RefreshMultiMediaNewsList_1Month】 每日定时更新最近一个月多媒体课程")

	newsSrv := service.NewNewsInfoService()
	newsList, err := newsSrv.RefreshMultiMediaNewsList_1MonthToRedis()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshMultiMediaNewsList_1Month】 处理被动更新队列 执行异常 ")
	} else {
		global.TaskLogger.InfoFormat("【Task Finished】 【Task_RefreshMultiMediaNewsList_1Month】 处理被动更新队列 执行完毕 newsList=【%s】", _json.GetJsonString(newsList))
	}

	return err
}