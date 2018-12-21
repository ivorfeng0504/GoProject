package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/strategyapi"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dottask"
)

// Task_SyncStrategyStockPool 定时更新股池信息
func Task_SyncStrategyStockPool(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_SyncStrategyStockPool】 定时更新股池信息")
	// 80018 双响炮 80019 高成长 80021 高盈利
	// 100006 深跌回弹 100013 趋势顶底（回转）100014 资金博弈（金叉） 100011 大单比率（推升）100015 龙腾四海（再强） 100005 底部量变
	syncStrategyStockPoolList := []string{"80018", "80019", "80021", "100006", "100013", "100014", "100011", "100015", "100005"}
	for _, strategyKey := range syncStrategyStockPoolList {
		stockPool, err := strategyapi.GetCommPool(strategyKey, config.CurrentConfig.AppId, strategyapi.ReadCacheAfter)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_SyncStrategyStockPool】 定时更新股池信息 执行异常  strategyKey=%s", strategyKey)
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_SyncStrategyStockPool】 定时更新股池信息 执行完毕 strategyKey=%s stockPool=%s ", strategyKey, _json.GetJsonString(stockPool))
		}
	}
	global.TaskLogger.DebugFormat("【Task Finished】 【Task_SyncStrategyStockPool】 定时更新股池信息 执行完毕")
	return nil
}
