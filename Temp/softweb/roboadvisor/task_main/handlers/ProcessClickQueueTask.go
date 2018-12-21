package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/click"
	expertnews_srv "git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/service/panwai"
	"github.com/devfeel/dottask"
	"strconv"
)

// Task_ProcessClickQueue 处理点击请求 更新数据库
func Task_ProcessClickQueue(context *task.TaskContext) (err error) {
	global.TaskLogger.InfoFormat("【Start Task】 【Task_ProcessClickQueue】 处理点击请求 更新数据库")
	err = click.ProcessQueue()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_ProcessClickQueue】 处理点击请求 更新数据库 执行异常 ")
	} else {
		global.TaskLogger.InfoFormat("【Task Finished】 【Task_ProcessClickQueue】 处理点击请求 更新数据库 执行完毕")
	}
	return err
}

func init() {
	//资讯点击落库
	click.Register(_const.ClickType_NewsInfo, func(identity string, clickNum int64) (err error) {
		newsSrv := expertnews_srv.NewNewsInformationService()
		newsId, err := strconv.ParseInt(identity, 10, 64)
		if err != nil {
			global.TaskLogger.DebugFormat("资讯点击落库异常 无效的资讯Id clickType=%s identity=%s clickNum=%d", _const.ClickType_NewsInfo, identity, clickNum)
			return err
		}
		err = newsSrv.UpdateClickNum(newsId, clickNum)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "资讯点击落库异常 更新异常 clickType=%s newsId=%d clickNum=%d", _const.ClickType_NewsInfo, newsId, clickNum)
		}
		return err
	})

	//股票相关资讯点击落库
	click.Register(_const.ClickType_StockNews, func(identity string, clickNum int64) (err error) {
		newsSrv := expertnews_srv.NewNewsInformationService()
		newsId, err := strconv.ParseInt(identity, 10, 64)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "股票相关资讯点击落库异常 无效的资讯Id clickType=%s identity=%s clickNum=%d", _const.ClickType_StockNews, identity, clickNum)
			return err
		}
		err = newsSrv.UpdateStockNewsClickNum(newsId, clickNum)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "股票相关资讯点击落库异常 更新异常 clickType=%s newsId=%d clickNum=%d", _const.ClickType_StockNews, newsId, clickNum)
		}
		return err
	})

	//策略资讯点击落库
	click.Register(_const.ClickType_StrategyNewsInfo, func(identity string, clickNum int64) (err error) {
		newsSrv := service.NewNewsInfoService()
		newsId, err := strconv.ParseInt(identity, 10, 64)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "专家资讯策略点击落库异常 无效的策略Id clickType=%s identity=%s clickNum=%d", _const.ClickType_StrategyNewsInfo, identity, clickNum)
			return err
		}
		err = newsSrv.UpdateClickNum(newsId, clickNum)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "专家资讯策略点击落库异常 更新异常 clickType=%s newsId=%d clickNum=%d", _const.ClickType_StrategyNewsInfo, newsId, clickNum)
		}
		return err
	})

	//策略点击落库
	click.Register(_const.ClickType_StrategyInfo, func(identity string, clickNum int64) (err error) {
		newsSrv := expertnews_srv.NewExpertNews_StrategyService()
		newsId, err := strconv.ParseInt(identity, 10, 64)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "策略资讯点击落库异常 无效的资讯Id clickType=%s identity=%s clickNum=%d", _const.ClickType_StrategyInfo, identity, clickNum)
			return err
		}
		err = newsSrv.UpdateClickNum(newsId, clickNum)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "策略资讯点击落库异常 更新异常 clickType=%s newsId=%d clickNum=%d", _const.ClickType_StrategyInfo, newsId, clickNum)
		}
		return err
	})

	//策略资讯-视频播放点击落库
	click.Register(_const.ClickType_StrategyVideoNewsInfo, func(identity string, playNum int64) (err error) {
		newsSrv := service.NewNewsInfoService()
		newsId, err := strconv.ParseInt(identity, 10, 64)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "专家资讯视频播放点击落库异常 无效的策略Id clickType=%s identity=%s playNum=%d", _const.ClickType_StrategyNewsInfo, identity, playNum)
			return err
		}
		err = newsSrv.UpdateVideoPlayNum(newsId, playNum)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "专家资讯视频播放点击落库异常 更新异常 clickType=%s newsId=%d playNum=%d", _const.ClickType_StrategyNewsInfo, newsId, playNum)
		}
		return err
	})

	//策略点播量落库
	click.Register(_const.ClickType_StrategyVideoInfo, func(identity string, playNum int64) (err error) {
		newsSrv := expertnews_srv.NewExpertNews_StrategyService()
		newsId, err := strconv.ParseInt(identity, 10, 64)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "策略资讯点击落库异常 无效的资讯Id clickType=%s identity=%s clickNum=%d", _const.ClickType_StrategyInfo, identity, playNum)
			return err
		}
		err = newsSrv.UpdateVideoPlayNum(newsId, playNum)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "策略资讯点击落库异常 更新异常 clickType=%s newsId=%d clickNum=%d", _const.ClickType_StrategyInfo, newsId, playNum)
		}
		return err
	})
}
