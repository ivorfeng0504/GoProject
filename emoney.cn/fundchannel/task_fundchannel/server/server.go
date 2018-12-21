package server

import (
	"emoney.cn/fundchannel/global"
	"emoney.cn/fundchannel/task_fundchannel/handlers"
	"fmt"
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
	service.RegisterHandler("Task_Print", handlers.Task_Print)
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
