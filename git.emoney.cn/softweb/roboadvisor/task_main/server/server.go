package server

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/task_main/handlers"
	"github.com/devfeel/dottask"
)

const (
	RunEnv_Flag       = "RunEnv"
	RunEnv_Develop    = "develop"
	RunEnv_Test       = "test"
	RunEnv_Production = "production"
)

var (
	RunEnv string
)

func registerTask(service *task.TaskService) {
	//TODO register task to service
	service.RegisterHandler("Task_GrantAward", handlers.Task_GrantAward)
	service.RegisterHandler("Task_InitCurrentGuessMarketActivity", handlers.Task_InitCurrentGuessMarketActivity)
	service.RegisterHandler("Task_InitNextGuessMarketActivity", handlers.Task_InitNextGuessMarketActivity)
	//service.RegisterHandler("Task_InitNextTradeDayStockList", handlers.Task_InitNextTradeDayStockList)
	service.RegisterHandler("Task_UpdateNewsInformationByDate", handlers.Task_UpdateNewsInformationByDate)
	service.RegisterHandler("Task_ProcessUpdateQueue", handlers.Task_ProcessUpdateQueue)
	service.RegisterHandler("Task_ProcessClickQueue", handlers.Task_ProcessClickQueue)
	service.RegisterHandler("Task_UpdateNewsInformationTemplate", handlers.Task_UpdateNewsInformationTemplate)
	service.RegisterHandler("Task_RefreshTodayNewsCache", handlers.Task_RefreshTodayNewsCache)
	service.RegisterHandler("Task_RefreshClosingNewsCache", handlers.Task_RefreshClosingNewsCache)
	service.RegisterHandler("Task_RefreshNewsInfoCache", handlers.Task_RefreshNewsInfoCache)
	service.RegisterHandler("Task_RefreshHotNewsInfoCache", handlers.Task_RefreshHotNewsInfoCache)
	service.RegisterHandler("Task_RefreshTopicNewsInfoCache", handlers.Task_RefreshTopicNewsInfoCache)
	service.RegisterHandler("Task_UpdateHotNewsInfoByClickNum", handlers.Task_UpdateHotNewsInfoByClickNum)
	service.RegisterHandler("Task_UpdateTopicInfo", handlers.Task_UpdateTopicInfo)
	service.RegisterHandler("Task_RefreshNewsInfoTopNCache", handlers.Task_RefreshNewsInfoTopNCache)
	service.RegisterHandler("Task_InitNextGuessChangeActivity", handlers.Task_InitNextGuessChangeActivity)
	service.RegisterHandler("Task_PublishCurrentGuessChangeActivity", handlers.Task_PublishCurrentGuessChangeActivity)
	service.RegisterHandler("Task_SyncStockNewsInformation", handlers.Task_SyncStockNewsInformation)
	service.RegisterHandler("Task_UpdateHotArticleListByClickNum", handlers.Task_UpdateHotArticleListByClickNum)
	service.RegisterHandler("Task_UpdateHotVideoListByClickNum", handlers.Task_UpdateHotVideoListByClickNum)
	service.RegisterHandler("Task_SyncClientStrategyInfo", handlers.Task_SyncClientStrategyInfo)
	service.RegisterHandler("Task_SyncBlockNewsInformation", handlers.Task_SyncBlockNewsInformation)
	service.RegisterHandler("Task_RefreshAllIndexStrategyNewsCache", handlers.Task_RefreshAllIndexStrategyNewsCache)
	service.RegisterHandler("Task_SyncStrategyStockPool", handlers.Task_SyncStrategyStockPool)
	service.RegisterHandler("Task_InitNextToDayReport", handlers.Task_InitNextToDayReport)
	service.RegisterHandler("Task_SyncExpertIndexData", handlers.Task_SyncExpertIndexData)
	service.RegisterHandler("Task_AutoSendStrategyStockTalkTask", handlers.Task_AutoSendStrategyStockTalkTask)
	service.RegisterHandler("Task_SyncExpertYqqData", handlers.Task_SyncExpertYqqData)
	service.RegisterHandler("Task_ProcessStockTalkMsgQueue", handlers.Task_ProcessStockTalkMsgQueue)
	service.RegisterHandler("Task_RefreshTrainList", handlers.Task_RefreshTrainList)
	service.RegisterHandler("Task_ProcessRepeatStockTalk", handlers.Task_ProcessRepeatStockTalk)
	service.RegisterHandler("Task_UpdateHotMultiMediaListByClickNum", handlers.Task_UpdateHotMultiMediaListByClickNum)
	service.RegisterHandler("Task_SyncStockTalkForStockBasicInfo", handlers.Task_SyncStockTalkForStockBasicInfo)
	service.RegisterHandler("Task_RefreshMultiMediaNewsList_1Month", handlers.Task_RefreshMultiMediaNewsList_1Month)
}

func StartTaskService(configPath string) {
	global.DotTask = task.StartNewService()

	//register all task handler
	registerTask(global.DotTask)

	//load config file
	global.DotTask.LoadConfig(configPath + "/dottask.conf")

	//start all task
	global.DotTask.StartAllTask()

	global.InnerLogger.Debug(fmt.Sprint("StartTaskService", " ", configPath, " ", global.DotTask.PrintAllCronTask()))
}

func StopTaskService() {
	if global.DotTask != nil {
		global.DotTask.StopAllTask()
	}
}
