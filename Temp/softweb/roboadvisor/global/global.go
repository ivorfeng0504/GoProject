package global

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/core"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/dottask"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/mapper"
)

//全局map
var GlobalItemMap *core.CMap
var DotApp *dotweb.DotWeb
var InnerLogger dotlog.Logger
var TaskLogger dotlog.Logger
var DotTask *task.TaskService

func Init(configPath string) error {
	GlobalItemMap = core.NewCMap()
	err := dotlog.StartLogService(configPath + "/dotlog.conf")
	if err != nil {
		return errors.New("log service start error => " + err.Error())
	}
	InnerLogger = dotlog.GetLogger(_const.LoggerName_Inner)
	TaskLogger = dotlog.GetLogger(_const.LoggerName_Task)
	//设置JSONTime序列化的格式
	mapper.SetTimeJSONFormat("2006-01-02T15:04:05")
	return nil
}

// AppLogTimeoutHookHandler 记录超时日志
func AppLogTimeoutHookHandler(ctx dotweb.Context) {
	realDration := ctx.Items().GetTimeDuration(dotweb.ItemKeyHandleDuration)
	logs := fmt.Sprintf("【AppLogTimeoutHookHandler Timeout】 req %v, cost %v", ctx.Request().Url(), realDration.Seconds())
	InnerLogger.Trace(logs)
}
