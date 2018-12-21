package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"github.com/devfeel/dottask"
	"time"
)

// Task_InitNextGuessChangeActivity 初始化下一期猜涨跌活动
func Task_InitNextGuessChangeActivity(context *task.TaskContext) error {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_InitNextGuessChangeActivity】 初始化下一期猜涨跌活动")
	srv := service.NewGuessChangeActivityService()
	issueNumber, err := srv.InitNextGuessChangeActivity()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_InitNextGuessChangeActivity】 初始化下一期猜涨跌活动 执行异常 期号为【%s】", issueNumber)
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_InitNextGuessChangeActivity】 初始化下一期猜涨跌活动 执行完毕 期号为【%s】", issueNumber)
	}
	return err
}

// Task_PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品
func Task_PublishCurrentGuessChangeActivity(context *task.TaskContext) error {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_PublishCurrentGuessChangeActivity】 公布当期猜涨跌结果，并且发放奖品")
	srv := service.NewGuessChangeActivityService()
	issueNumber, err := srv.PublishCurrentGuessChangeActivity(config.CurrentConfig.Activity_GuessChange, time.Now())
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_PublishCurrentGuessChangeActivity】 公布当期猜涨跌结果，并且发放奖品 执行异常 期号为【%s】", issueNumber)
	} else {
		global.TaskLogger.DebugFormat("【Task Finished】 【Task_PublishCurrentGuessChangeActivity】 公布当期猜涨跌结果，并且发放奖品 执行完毕 期号为【%s】", issueNumber)
	}
	return err
}
