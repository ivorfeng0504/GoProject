package handlers

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/global"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dataapi"
	"git.emoney.cn/softweb/roboadvisor/protected/service/myoptional"
	myoptional_srv "git.emoney.cn/softweb/roboadvisor/protected/service/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service/strategyapi"
	strategyservice_srv "git.emoney.cn/softweb/roboadvisor/protected/service/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dottask"
	"strconv"
	"time"
)

// Task_AutoSendStrategyStockTalkTask 策略入选信息自动发送到微股吧Task
func Task_AutoSendStrategyStockTalkTask(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task")
	//判断是否为交易日 非交易日不需要发送消息
	isTradeDay, err := dataapi.IsTradeDayToDay()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task 判断是否为交易日异常 ")
		return err
	}
	if isTradeDay == false {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task 当前不是交易日 不需要发送消息")
		return nil
	}
	srv := strategyservice_srv.NewClientStrategyInfoService()
	strategyInfoList, err := srv.GetClientStrategyInfoList(false)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task 执行异常 GetClientStrategyInfoList失败")
		return err
	}
	if strategyInfoList == nil || len(strategyInfoList) == 0 {
		err = errors.New("策略信息为空")
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task 执行异常 策略信息为空")
		return err
	}
	strategyStockTalkMap := make(map[string]*myoptional_model.StrategyStockTalk)
	nickName := "智盈策略君"
	//目前限制指定的几个父级策略才能发送消息
	//70007	拐点形态
	//70008	中继趋势
	//70009	突破压力
	//70010	基本面·策略
	//70012 抄底三剑客
	//70013 长波突进
	//70014 中继回踩
	enableStrategyMap := map[int]bool{70007: true, 70008: true, 70009: true, 70010: true, 70012: true, 70013: true, 70014: true}

	strategyPoolSrv := myoptional_srv.NewStrategyPoolService()

	//循环所有策略 获取入选的股票集合 构造消息内容
	for _, strategyInfo := range strategyInfoList {
		exist := enableStrategyMap[strategyInfo.ParentId]
		if exist == false {
			continue
		}
		//读取接口的实时数据 如果没有数据则放弃
		stockPool, innErr := strategyapi.GetCommonPoolStockList(-1, strconv.Itoa(strategyInfo.ClientStrategyId), strategyapi.NoCache)
		if innErr != nil && innErr != strategyapi.NilStockPool {
			global.TaskLogger.ErrorFormat(innErr, "【Task Continue】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task 执行异常 获取股池异常 ClientStrategyId=%d", strategyInfo.ClientStrategyId)
			continue
		}
		if innErr == strategyapi.NilStockPool || stockPool == nil || len(stockPool) == 0 {
			global.TaskLogger.DebugFormat("【Task Continue】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task  获取股池为空 自动忽略 ClientStrategyId=%d", strategyInfo.ClientStrategyId)
			continue
		}
		for _, stock := range stockPool {
			strategyStockTalk, exist := strategyStockTalkMap[stock.StockCode]
			if exist == false {
				strategyStockTalk = &myoptional_model.StrategyStockTalk{
					StockCode: stock.StockCode,
					StockName: stock.StockName,
					NickName:  nickName,
				}
				strategyStockTalkMap[stock.StockCode] = strategyStockTalk
			}
			strategyDesc := strategyInfo.ParentName + "-" + strategyInfo.ClientStrategyName
			strategyStockTalk.StrategyDescList = append(strategyStockTalk.StrategyDescList, strategyDesc)

			//将股池记录插入记录表
			strategyPool := &myoptional_model.StrategyPool{
				ClientStrategyId:   strategyInfo.ClientStrategyId,
				ClientStrategyName: strategyInfo.ClientStrategyName,
				ParentId:           strategyInfo.ParentId,
				ParentName:         strategyInfo.ParentName,
				StockCode:          stock.StockCode,
				StockName:          stock.StockName,
				IssueNumber:        time.Now().Format("20060102"),
			}
			_, insertErr := strategyPoolSrv.InsertStrategyPool(strategyPool)
			if insertErr != nil {
				global.TaskLogger.ErrorFormat(insertErr, "【Task_AutoSendStrategyStockTalkTask】 将股池记录插入记录表 异常 strategyPool=%s", _json.GetJsonString(strategyPool))
			}
		}
	}

	global.TaskLogger.DebugFormat("【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task 数据准备完毕 strategyStockTalkMap=【%s】", _json.GetJsonString(strategyStockTalkMap))
	stockTalkSrv := myoptional.NewStockTalkService()
	err = stockTalkSrv.AutoInsertStrategyStockTalk(strategyStockTalkMap)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task 执行异常 ")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task 执行完毕")
		//写入成功后  清理一次旧的重复数据
		err = stockTalkSrv.ProcessRepeatStockTalk()
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task ProcessRepeatStockTalk执行异常 ")
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_AutoSendStrategyStockTalkTask】 策略入选信息自动发送到微股吧Task->处理重复的微股吧数据 执行完毕")
		}
	}
	return err
}

// Task_ProcessStockTalkMsgQueue 处理直播精选消息队列，将队列中的数据插入数据库
func Task_ProcessStockTalkMsgQueue(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_ProcessStockTalkMsgQueue】 处理直播精选消息队列，将队列中的数据插入数据库")
	srv := myoptional.NewStockTalkMsgService()
	err = srv.ProcessStockTalkMsgQueue()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_ProcessStockTalkMsgQueue】 处理直播精选消息队列 执行异常")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_ProcessStockTalkMsgQueue】 处理直播精选消息队列 执行完毕")
	}
	return err
}

// Task_ProcessRepeatStockTalk 处理重复的微股吧数据  逻辑删除
func Task_ProcessRepeatStockTalk(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_ProcessRepeatStockTalk】 处理重复的微股吧数据")
	srv := myoptional.NewStockTalkService()
	err = srv.ProcessRepeatStockTalk()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_ProcessRepeatStockTalk】 处理重复的微股吧数据 执行异常")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_ProcessRepeatStockTalk】 处理重复的微股吧数据 执行完毕")
	}
	return err
}

// Task_SyncStockTalkForStockBasicInfo 同步最新的基本面信息到微股吧
func Task_SyncStockTalkForStockBasicInfo(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_SyncStockTalkForStockBasicInfo】 同步最新的基本面信息到微股吧")
	srv := myoptional.NewStockTalkService()
	err = srv.SyncStockTalkForStockBasicInfo()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_SyncStockTalkForStockBasicInfo】 同步最新的基本面信息到微股吧 执行异常")
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_SyncStockTalkForStockBasicInfo】 同步最新的基本面信息到微股吧 执行完毕")
	}
	return err
}
