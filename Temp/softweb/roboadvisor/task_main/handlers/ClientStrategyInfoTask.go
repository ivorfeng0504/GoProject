package handlers

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/panwai"
	strategyservice_srv "git.emoney.cn/softweb/roboadvisor/protected/service/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dottask"
)

// Task_UpdateHotArticleListByClickNum 每日定时更新热门文章资讯
func Task_UpdateHotArticleListByClickNum(context *task.TaskContext) error {
	global.TaskLogger.InfoFormat("【Start Task】 【Task_UpdateHotArticleListByClickNum】 每日定时更新热门文章资讯")

	clientStrategySrv := strategyservice_srv.NewClientStrategyInfoService()
	strategyList, err := clientStrategySrv.GetClientStrategyInfoListDB(true)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task doing】 【Task_UpdateHotArticleListByClickNum】 GetClientStrategyInfoList获取策略信息失败")
		return err
	}

	for _, v := range strategyList {
		var colid = v.ColumnInfoId
		newsSrv := service.NewNewsInfoService()
		newsList, err := newsSrv.GetHotArticleList(colid)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task doing】 【Task_UpdateHotArticleListByClickNum】 热门文章资讯处理被动更新队列 执行异常  ClientStrategyId=%s ClientStrategyName=%s", v.ClientStrategyId, v.ClientStrategyName)
			continue
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_UpdateHotArticleListByClickNum】 单条策略热门文章任务执行完毕 ClientStrategyId=%d newsList=%s ", v.ClientStrategyId, _json.GetJsonString(newsList))
		}
	}

	global.TaskLogger.InfoFormat("【Task Finished】 【Task_UpdateHotArticleListByClickNum】 热门文章资讯处理被动更新队列 执行完毕 已处理的策略列表=【%s】", _json.GetJsonString(strategyList))
	return err
}

// Task_UpdateHotVideoListByClickNum 每日定时更新热门视频资讯
func Task_UpdateHotVideoListByClickNum(context *task.TaskContext) error {
	global.TaskLogger.InfoFormat("【Start Task】 【Task_UpdateHotVideoListByClickNum】 每日定时更新热门视频资讯")

	clientStrategySrv := strategyservice_srv.NewClientStrategyInfoService()
	strategyList, err := clientStrategySrv.GetClientStrategyInfoListDB(true)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task doing】 【Task_UpdateHotVideoListByClickNum】 GetClientStrategyInfoList获取策略信息失败")
		return err
	}

	for _, v := range strategyList {
		var colid = v.ColumnInfoId
		newsSrv := service.NewNewsInfoService()
		newsList, err := newsSrv.GetHotVideoList(colid)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task doing】 【Task_UpdateHotVideoListByClickNum】 热门视频资讯处理被动更新队列 执行异常 ClientStrategyId=%s ClientStrategyName=%s", v.ClientStrategyId, v.ClientStrategyName)
			continue
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_UpdateHotVideoListByClickNum】 单条策略热门视频任务执行完毕 ClientStrategyId=%d newsList=%s ", v.ClientStrategyId, _json.GetJsonString(newsList))
		}
	}

	global.TaskLogger.InfoFormat("【Task Finished】 【Task_UpdateHotArticleListByClickNum】 热门文章资讯处理被动更新队列 执行完毕 已处理的策略列表=【%s】", _json.GetJsonString(strategyList))
	return err
}

// Task_SyncClientStrategyInfo 定时同步策略信息
func Task_SyncClientStrategyInfo(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_SyncClientStrategyInfo】 定时同步策略信息")
	srv := strategyservice_srv.NewClientStrategyInfoService()
	err = srv.SyncClientStrategyInfo()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_SyncClientStrategyInfo】 定时同步策略信息 执行异常 ")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_SyncClientStrategyInfo】 定时同步策略信息 执行完毕")
	}
	return err
}

// Task_RefreshAllIndexStrategyNewsCache 刷新所有组策略下的最新一条视频与资讯
func Task_RefreshAllIndexStrategyNewsCache(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_RefreshAllIndexStrategyNewsCache】 刷新所有组策略下的最新一条视频与资讯")
	srv := strategyservice_srv.NewClientStrategyInfoService()
	tree, err := srv.GetClientStrategyInfoTreeDB()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshAllIndexStrategyNewsCache】 刷新所有组策略下的最新一条视频与资讯 执行异常->GetClientStrategyInfoTreeDB ")
		return err
	}
	if tree == nil || len(tree) == 0 {
		err = errors.New("未获取到策略树")
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshAllIndexStrategyNewsCache】 刷新所有组策略下的最新一条视频与资讯 执行异常->GetClientStrategyInfoTreeDB ")
		return err
	}
	columnSrv := strategyservice_srv.NewColumnInfoService()
	for _, item := range tree {
		groupId := item.ClientStrategyId
		if groupId <= 0 {
			continue
		}
		var columnList []int
		columnList = append(columnList, item.ColumnInfoId)
		for _, sub := range item.Children {
			columnList = append(columnList, sub.ColumnInfoId)
		}
		result, err := columnSrv.RefreshIndexStrategyNewsCache(groupId, columnList)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshAllIndexStrategyNewsCache】 刷新所有组策略下的最新一条视频与资讯 RefreshIndexStrategyNewsCache 执行异常 groupId=%d columnList=%s", groupId, _json.GetJsonString(columnList))
			continue
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshAllIndexStrategyNewsCache】 子任务执行完毕 groupId=%d columnList=%s result=%s", groupId, _json.GetJsonString(columnList), _json.GetJsonString(result))
		}
	}
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshAllIndexStrategyNewsCache】 刷新所有组策略下的最新一条视频与资讯 执行异常 ")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshAllIndexStrategyNewsCache】 刷新所有组策略下的最新一条视频与资讯 执行完毕")
	}
	return err
}

// Task_UpdateHotMultiMediaListByClickNum 每日定时更新热门培训资讯
func Task_UpdateHotMultiMediaListByClickNum(context *task.TaskContext) error {
	global.TaskLogger.InfoFormat("【Start Task】 【Task_UpdateHotMultiMediaListByClickNum】 每日定时更新热门培训资讯")

	clientStrategySrv := strategyservice_srv.NewClientStrategyInfoService()
	strategyList, err := clientStrategySrv.GetClientStrategyInfoListDB(true)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task doing】 【Task_UpdateHotMultiMediaListByClickNum】 GetClientStrategyInfoList获取策略信息失败")
		return err
	}

	for _, v := range strategyList {
		var colid = v.ColumnInfoId
		newsSrv := service.NewNewsInfoService()
		newsList, err := newsSrv.GetHotMultiMediaList(colid)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task doing】 【Task_UpdateHotMultiMediaListByClickNum】 热门培训资讯处理被动更新队列 执行异常 ClientStrategyId=%s ClientStrategyName=%s", v.ClientStrategyId, v.ClientStrategyName)
			continue
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_UpdateHotMultiMediaListByClickNum】 单条策略热门培训任务执行完毕 ClientStrategyId=%d newsList=%s ", v.ClientStrategyId, _json.GetJsonString(newsList))
		}
	}

	global.TaskLogger.InfoFormat("【Task Finished】 【Task_UpdateHotMultiMediaListByClickNum】 热门培训资讯处理被动更新队列 执行完毕 已处理的策略列表=【%s】", _json.GetJsonString(strategyList))
	return err
}
