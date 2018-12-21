package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"github.com/devfeel/dottask"
)

// Task_InitNextGuessMarketActivity 初始化下一期猜数字活动
func Task_InitNextGuessMarketActivity(context *task.TaskContext) error {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_InitNextGuessMarketActivity】 初始化下一期猜数字活动")
	srv := service.NewGuessMarketActivityService()
	issueNumber, err := srv.InitGuessMarketActivity(true)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextGuessMarketActivity】 初始化下一期猜数字活动 执行异常 期号为【%s】", issueNumber)
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_InitNextGuessMarketActivity】 初始化下一期猜数字活动 执行完毕 期号为【%s】", issueNumber)
	}
	return err
}

// Task_InitCurrentGuessMarketActivity 初始化本期猜数字活动（用于第一次初始化或者初始化下期失败后重试）
func Task_InitCurrentGuessMarketActivity(context *task.TaskContext) error {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_InitCurrentGuessMarketActivity】 初始化本期猜数字活动")
	srv := service.NewGuessMarketActivityService()
	issueNumber, err := srv.InitGuessMarketActivity(false)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitCurrentGuessMarketActivity】 初始化本期猜数字活动 执行异常 期号为【%s】", issueNumber)
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_InitCurrentGuessMarketActivity】 初始化本期猜数字活动 执行完毕 期号为【%s】", issueNumber)
	}
	return err
}

// Task_GrantAward 猜数字活动发放奖励
func Task_GrantAward(context *task.TaskContext) error {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_GrantAward】 猜数字活动发放奖励")
	srv := service.NewGuessMarketActivityService()
	issueNumber, err := srv.GrantAward()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_GrantAward】 猜数字活动发放奖励 执行异常 期号为【%s】", issueNumber)
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_GrantAward】 猜数字活动发放奖励 执行完毕 期号为【%s】", issueNumber)
	}
	return err
}
