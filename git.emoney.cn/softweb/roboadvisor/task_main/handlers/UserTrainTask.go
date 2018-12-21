package handlers

import (
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/train"
	"github.com/devfeel/dottask"
	"git.emoney.cn/softweb/roboadvisor/protected/service/resapi"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	trainmodel "git.emoney.cn/softweb/roboadvisor/protected/model/train"
)

// Task_SyncStrategyStockPool 定时更新智盈培训课程列表
func Task_RefreshTrainList(context *task.TaskContext) (err error) {
	global.TaskLogger.DebugFormat("【Start Task】 【Task_RefreshTrainList】 定时更新培训课程列表")

	var syncZyPIDList []*trainmodel.PIDInfo
	retPID, err := resapi.GetZYPIDConfig()
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshTrainList】 定时更新培训课程列表获取智盈产品ID 执行异常")
		return nil
	}
	err = _json.Unmarshal(retPID, &syncZyPIDList)
	if err != nil {
		global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshTrainList】 定时更新培训课程列表反序列化PID列表 执行异常 PID=%s", retPID)
		return nil
	}
	//syncZyPIDList := []string{"888020000", "888010000", "888010400", "888012000", "888012001", "888012400", "888020000", "888020400", "888030000", "888030400"}
	for _, v := range syncZyPIDList {
		service := train.NewTrainService()
		err := service.RefreshTrainListToRedis(v.PID)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "【Task Finished】 【Task_RefreshTrainList】 定时更新培训课程列表 执行异常  pid=%s", v.PID)
		} else {
			global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshTrainList】 定时更新培训课程列表 执行完毕 pid=%s", v.PID)
		}
	}
	global.TaskLogger.DebugFormat("【Task Finished】 【Task_RefreshTrainList】 定时更新培训课程列表 执行完毕")
	return nil
}