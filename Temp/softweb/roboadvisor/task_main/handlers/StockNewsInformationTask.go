package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/myoptional"
	"github.com/devfeel/dottask"
)

// Task_SyncStockNewsInformation 同步最新的资讯信息到关系表及相关缓存
func Task_SyncStockNewsInformation(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_SyncStockNewsInformation】 同步最新的资讯信息到关系表及相关缓存")
	newsSrv := myoptional.NewStockNewsInformationService()
	err = newsSrv.SyncStockNewsInformation(_const.StockNewsInformationPerTop)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_SyncStockNewsInformation】 同步最新的资讯信息到关系表及相关缓存 执行异常 ")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_SyncStockNewsInformation】 同步最新的资讯信息到关系表及相关缓存 执行完毕")
	}
	return err
}
