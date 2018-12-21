package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dataapi"
	"git.emoney.cn/softweb/roboadvisor/protected/service/strategyapi"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"github.com/devfeel/dottask"
	"time"
)

// Task_InitNextToDayReport 初始化今日推荐研报【龙头猜想】
func Task_InitNextToDayReport(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_InitNextToDayReport】 初始化下一个交易日的推荐研报【龙头猜想】")
	iniDate := time.Now().Format("20060102")
	activityAwardSrv := service.NewAwardInActivityService()
	awardList, err := activityAwardSrv.GetActivityAwardListByActivityIdDB(config.CurrentConfig.Activity_ReceiveStock, true)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextToDayReport】 执行异常 获取奖品信息异常")
		return err
	}
	if awardList == nil || len(awardList) == 0 {
		global.TaskLogger.WarnFormat("【Task Finished】 【Task_InitNextToDayReport】 执行异常 未获取到奖品信息")
		return err
	}
	award := awardList[0]
	if len(award.LinkUrl) == 0 {
		global.TaskLogger.WarnFormat("【Task Finished】 【Task_InitNextToDayReport】 执行异常 未获取到研报地址")
		return err
	}
	srv := service.NewReceiveStockActivityService()
	stockPool, err := srv.GetNewstStockPool()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextToDayReport】 执行异常 获取最近一期活动异常")
		return err
	}
	if stockPool != nil && stockPool.IssueNumber == iniDate && award.LinkUrl == stockPool.ReportUrl {
		global.TaskLogger.DebugFormat("当前活动已经初始化 并且奖品未发生变化 无需再次初始化  iniDate=%s", iniDate)
		return nil
	}

	//如果当天活动不存在 则新增 否则更新研报
	if stockPool == nil || stockPool.IssueNumber != iniDate {
		_, err = srv.InsertStockPoolByReportUrl(award.LinkUrl, iniDate)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextToDayReport】 初始化下一个交易日的推荐研报【龙头猜想】 执行异常 iniDate=%s ReportUrl=%s", iniDate, stockPool.ReportUrl)
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_InitNextToDayReport】 初始化下一个交易日的推荐研报【龙头猜想】 执行完毕 iniDate=%s ReportUrl=%s", iniDate, stockPool.ReportUrl)
		}
	} else {
		err = srv.UpdateStockPoolByReportUrl(stockPool.ReceiveStockActivityId, award.LinkUrl)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextToDayReport】 初始化下一个交易日的推荐研报【龙头猜想】 执行异常 ReceiveStockActivityId=%d iniDate=%s ReportUrl=%s", stockPool.ReceiveStockActivityId, iniDate, stockPool.ReportUrl)
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_InitNextToDayReport】 初始化下一个交易日的推荐研报【龙头猜想】 执行完毕 ReceiveStockActivityId=%d  iniDate=%s ReportUrl=%s", stockPool.ReceiveStockActivityId, iniDate, stockPool.ReportUrl)
		}
	}
	return err
}

// Task_InitNextTradeDayStockList 初始化下一个交易日的推荐股票【双响炮】-弃用
func Task_InitNextTradeDayStockList(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_InitNextTradeDayStockList】 初始化下一个交易日的推荐股票【双响炮】")
	nextTradeDay, err := dataapi.GetNextTradeDay(time.Now())
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextTradeDayStockList】 执行异常 获取下一个交易日失败")
		return err
	}
	iniDate := nextTradeDay.Format("20060102")
	srv := service.NewReceiveStockActivityService()
	stockPool, err := srv.GetNewstStockPool()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextTradeDayStockList】 执行异常 获取最近一期股池异常")
		return err
	}
	if stockPool != nil && stockPool.IssueNumber == iniDate {
		global.TaskLogger.DebugFormat("当前股池已经初始化 无需再次初始化  iniDate=%s", iniDate)
		return nil
	}

	//获取推荐股票
	stockList, err := strategyapi.GetShuangXiangPaoPool(5)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextTradeDayStockList】 执行异常 获取股池失败")
		return err
	}
	_, err = srv.InsertStockPool(stockList, iniDate)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextTradeDayStockList】 初始化下一个交易日的推荐股票【双响炮】 执行异常 iniDate=%s", iniDate)
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_InitNextTradeDayStockList】 初始化下一个交易日的推荐股票【双响炮】 执行完毕 iniDate=%s", iniDate)
	}
	return err
}
