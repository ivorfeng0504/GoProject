package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/global"
	"github.com/devfeel/dottask"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
)

// Task_SyncExpertIndex 定时刷新专家资讯首页数据（今日头条、策略看盘、专家直播、主题、盘后预测）入redis
func Task_SyncExpertIndexData(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_SyncExpertIndexData】 定时刷新专家资讯首页数据")

	syncPageDataService := expertnews.NewSyncPageDataService()
	err = syncPageDataService.SyncIndexData()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_SyncExpertIndexData】 定时刷新专家资讯首页数据 执行异常")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_SyncExpertIndexData】 定时刷新专家资讯首页数据 执行完毕")
	}

	return nil
}

//Task_SyncExpertYqqData 定时刷新专家资讯 策略页面 益圈圈 数据（直播间分组，最热直播间，vip直播间，最新直播，直播统计数据）入redis
func Task_SyncExpertYqqData(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_SyncExpertYqqData】 定时刷新专家资讯-策略益圈圈数据")

	syncPageDataService := expertnews.ConExpertNews_YqqService()
	_,err = syncPageDataService.GetExpertYqqEntireData_Organize()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_SyncExpertYqqData】 定时刷新专家资讯-策略益圈圈数据 执行异常")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_SyncExpertYqqData】 定时刷新专家资讯-策略益圈圈数据 执行完毕")
	}

	return nil
}