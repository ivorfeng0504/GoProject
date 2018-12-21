package handlers

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dataapi"
	expertnews_srv "git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dottask"
	"github.com/garyburd/redigo/redis"
	"time"
)

// 每日定时更新当天的资讯信息 是否在运行
//如果在运行则延后30秒钟执行
var Task_UpdateNewsInformationByDateIsRun = false

// Task_UpdateNewsInformationByDate 每日定时更新当天的资讯信息
func Task_UpdateNewsInformationByDate(context *task.TaskContext) error {
	Task_UpdateNewsInformationByDateIsRun = true
	global.TaskLogger.DebugFormat("【Start Task】 【Task_UpdateNewsInformationByDate】 每日定时更新当天的咨询信息")
	now := time.Now()
	newsList, err := dataapi.UpdateNewsInformationByDate(&now)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_UpdateNewsInformationByDate】 每日定时更新当天的咨询信息 执行异常 date=【%s】", now.Format("2006-01-02"))
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_UpdateNewsInformationByDate】 每日定时更新当天的咨询信息 执行完毕 date=【%s】 newsList=【%s】", now.Format("2006-01-02"), _json.GetJsonString(newsList))
	}
	Task_UpdateNewsInformationByDateIsRun = false
	return err
}

// Task_ProcessUpdateQueue 处理被动更新队列
func Task_ProcessUpdateQueue(context *task.TaskContext) (err error) {
	if Task_UpdateNewsInformationByDateIsRun {
		//延后30秒钟执行 防止并发更新冲突
		time.Sleep(time.Second * 30)
		global.TaskLogger.InfoFormat("【Start Task】 【Task_ProcessUpdateQueue】 处理被动更新队列 Task_UpdateNewsInformationByDateIsRun 延后处理队列")
	}
	global.TaskLogger.InfoFormat("【Start Task】 【Task_ProcessUpdateQueue】 处理被动更新队列")
	newsSrv := expertnews_srv.NewNewsInformationService()
	tableName, err := newsSrv.PopUpdateQueue()
	if err == redis.ErrNil {
		err = nil
		tableName = ""
	}
	if err == nil && tableName == fmt.Sprintf("[[%s]]", _const.NewsInformationTable) {
		//清除队列  合并所有已存在的更新请求
		newsSrv.ClearUpdateQueue()
		//更新数据库
		newsList, err := dataapi.UpdateNewsInformation()
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_ProcessUpdateQueue】 处理被动更新队列 数据库更新 执行异常 tableName=【%s】", tableName)
			return err
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_ProcessUpdateQueue】 处理被动更新队列 数据库更新完毕 tableName=【%s】 newsList=【%s】", tableName, _json.GetJsonString(newsList))
		}
		//刷新缓存
		err = newsSrv.RefreshCache()
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_ProcessUpdateQueue】 处理被动更新队列 刷新缓存 执行异常 tableName=【%s】", tableName)
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_ProcessUpdateQueue】 处理被动更新队列 更新缓存完毕 tableName=【%s】 newsList=【%s】", tableName, _json.GetJsonString(newsList))
		}
	}
	return err
}

// Task_UpdateNewsInformationTemplate 更新资讯模板
func Task_UpdateNewsInformationTemplate(context *task.TaskContext) error {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_UpdateNewsInformationTemplate】 更新资讯模板")
	templates, err := dataapi.GetNewsTemplate()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_UpdateNewsInformationTemplate】 更新资讯模板 执行异常")
		return err
	}
	if templates == nil || len(templates) == 0 {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_UpdateNewsInformationTemplate】 更新资讯模板 资讯模板为空")
		return err
	}
	newsSrv := expertnews_srv.NewNewsInformationService()
	//更新资讯模板缓存
	err = newsSrv.SetNewsTemplateCache(templates)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_UpdateNewsInformationTemplate】 更新资讯模板 执行异常 templates=%s", _json.GetJsonString(templates))
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_UpdateNewsInformationTemplate】 更新资讯模板 执行完毕 templates=%s", _json.GetJsonString(templates))
	}
	return err
}

// Task_RefreshTodayNewsCache 刷新获取今日头条缓存
func Task_RefreshTodayNewsCache(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_RefreshTodayNewsCache】 刷新获取今日头条缓存")
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.RefreshTodayNewsCacheV2()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshTodayNewsCache】 刷新获取今日头条缓存 执行异常 newsList=【%s】", _json.GetJsonString(newsList))
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshTodayNewsCache】 刷新获取今日头条缓存 执行完毕  newsList=【%s】", _json.GetJsonString(newsList))
	}
	return err
}

// Task_RefreshClosingNewsCache 刷新盘后预测资讯缓存
func Task_RefreshClosingNewsCache(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_RefreshClosingNewsCache】 刷新盘后预测资讯缓存")
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.RefreshClosingNewsCache()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshClosingNewsCache】 刷新盘后预测资讯缓存 执行异常 newsList=【%s】", _json.GetJsonString(newsList))
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshClosingNewsCache】 刷新盘后预测资讯缓存 执行完毕  newsList=【%s】", _json.GetJsonString(newsList))
	}
	return err
}

// Task_RefreshNewsInfoCache 刷新近几日要闻缓存-近3天
func Task_RefreshNewsInfoCache(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_RefreshNewsInfoCache】 刷新近几日要闻缓存-近3天")
	newsSrv := expertnews_srv.NewNewsInformationService()
	dateList := []time.Time{time.Now(), time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, -2)}
	for _, date := range dateList {
		newsList, err := newsSrv.RefreshNewsInfoCache(date)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshNewsInfoCache】 刷新近几日要闻缓存-近3天 执行异常 date=%s newsList=【%s】", date.Format("2006-01-02"), _json.GetJsonString(newsList))
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshNewsInfoCache】 刷新近几日要闻缓存-近3天 执行完毕  date=%s newsList=【%s】", date.Format("2006-01-02"), _json.GetJsonString(newsList))
		}
	}
	return err
}

// Task_RefreshNewsInfoTopNCache 刷新最新的N条要闻缓存
func Task_RefreshNewsInfoTopNCache(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_RefreshNewsInfoTopNCache】 刷新最新的N条要闻缓存")
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.RefreshNewsInfoTopNCache(expertnews_srv.NewsInfoTopN)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshNewsInfoTopNCache】 刷新最新的N条要闻缓存 执行异常 top=%d newsList=【%s】", expertnews_srv.NewsInfoTopN, _json.GetJsonString(newsList))
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshNewsInfoTopNCache】 刷新最新的N条要闻缓存 执行完毕  top=%d newsList=【%s】", expertnews_srv.NewsInfoTopN, _json.GetJsonString(newsList))
	}
	return err
}

// Task_RefreshHotNewsInfoCache 刷新热门资讯缓存
func Task_RefreshHotNewsInfoCache(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_RefreshHotNewsInfoCache】 刷新热门资讯缓存")
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.RefreshHotNewsInfoCache()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshHotNewsInfoCache】 刷新热门资讯缓存 执行异常 newsList=【%s】", _json.GetJsonString(newsList))
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshHotNewsInfoCache】 刷新热门资讯缓存 执行完毕  newsList=【%s】", _json.GetJsonString(newsList))
	}
	return err
}

// Task_RefreshTopicNewsInfoCache 刷新所有主题的相关资讯
func Task_RefreshTopicNewsInfoCache(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_RefreshTopicNewsInfoCache】 刷新所有主题的相关资讯")
	newsSrv := expertnews_srv.NewNewsInformationService()
	topicSrv := expertnews_srv.NewExpertNews_TopicService()
	topicList, err := topicSrv.GetTopicList()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshTopicNewsInfoCache】 刷新所有主题的相关资讯 获取所有主题异常")
		return err
	}
	if topicList == nil || len(topicList) == 0 {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshTopicNewsInfoCache】 刷新所有主题的相关资讯 执行完毕 没有需要刷新的主题资讯")
		return nil
	}
	for _, topic := range topicList {
		newsList, err := newsSrv.RefreshTopicNewsInfoCache(topic)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshTopicNewsInfoCache】 刷新所有主题的相关资讯 执行异常 topicId=%d newsList=【%s】", topic.ID, _json.GetJsonString(newsList))
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshTopicNewsInfoCache】 刷新所有主题的相关资讯 执行完毕 topicId=%d newsList=【%s】", topic.ID, _json.GetJsonString(newsList))
		}
	}
	return err
}

// Task_SyncBlockNewsInformation 同步板块资讯信息
func Task_SyncBlockNewsInformation(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_SyncBlockNewsInformation】 同步板块资讯信息")
	newsSrv := expertnews_srv.NewNewsInformationService()
	err = newsSrv.SyncBlockNewsInformation()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_SyncBlockNewsInformation】 同步板块资讯信息 执行异常 ")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_SyncBlockNewsInformation】 同步板块资讯信息 执行完毕 ")
	}
	return err
}
