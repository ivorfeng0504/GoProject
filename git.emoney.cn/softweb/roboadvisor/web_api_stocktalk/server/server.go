package server

import (
	"fmt"
	_ "fmt"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/config"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
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

func StartServer(configPath string) error {
	//初始化DotServer
	appConfig := config.MustInitConfig(configPath + "/dotweb.conf")
	global.DotApp = dotweb.ClassicWithConf(appConfig)

	if RunEnv == RunEnv_Production {
		global.DotApp.SetProductionMode()
	} else {
		global.DotApp.SetDevelopmentMode()
	}
	//设置超时监控
	global.DotApp.UseTimeoutHook(global.AppLogTimeoutHookHandler, time.Second*_const.UseTimeoutHookSecond)

	//异常处理
	global.DotApp.SetExceptionHandle(func(ctx dotweb.Context, err error) {
		ctx.Response().Header().Set(dotweb.HeaderContentType, dotweb.CharsetUTF8)
		global.InnerLogger.ErrorFormat(err, "服务器发生未处理的错误！")
		if global.DotApp.IsDevelopmentMode() {
			stack := string(debug.Stack())
			ctx.WriteStringC(http.StatusInternalServerError, fmt.Sprintln(err)+stack)
		} else {
			ctx.WriteStringC(http.StatusInternalServerError, "WebApi服务器发生异常")
		}
	})

	//设置路由
	InitRoute(global.DotApp.HttpServer)
	//启用Session
	global.DotApp.HttpServer.SetEnabledSession(true)
	//使用Redis来保存Session
	//session过期时间 24小时
	var sessionExpire int64 = 60 * 60 * 24
	storeConfig := _cache.GetSessionRedisConfig(protected.SessionRedisConfig, sessionExpire, "webapi_stocktalk", "")
	global.DotApp.HttpServer.SetSessionConfig(storeConfig)

	//设置模板路径
	//global.DotApp.HttpServer.Renderer().SetTemplatePath("./views")

	global.InnerLogger.Debug("dotweb.StartServer => " + strconv.Itoa(appConfig.Server.Port))
	err := global.DotApp.Start()
	if err != nil {
		panic(err)
	}
	return err
}
