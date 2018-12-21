package server

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/captcha"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/error"
	"github.com/devfeel/cache"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/session"
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
	appConfig := config.MustInitConfig(configPath + "/dotweb.conf")
	global.DotApp = dotweb.ClassicWithConf(appConfig)

	if RunEnv == RunEnv_Production {
		global.DotApp.SetProductionMode()
	} else {
		global.DotApp.SetDevelopmentMode()
	}

	//设置超时监控
	global.DotApp.UseTimeoutHook(global.AppLogTimeoutHookHandler, time.Second*_const.UseTimeoutHookSecond)

	//global.DotApp.UseRequestLog()
	//global.DotApp.Use(cors.Middleware(cors.NewConfig().UseDefault()))

	//设置404页面
	global.DotApp.SetNotFoundHandle(func(ctx dotweb.Context) {
		ctx.Response().Header().Set(dotweb.HeaderContentType, dotweb.CharsetUTF8)
		_error.NotFound(ctx)
	})

	//异常处理
	global.DotApp.SetExceptionHandle(func(ctx dotweb.Context, err error) {
		ctx.Response().Header().Set(dotweb.HeaderContentType, dotweb.CharsetUTF8)
		global.InnerLogger.ErrorFormat(err, "服务器发生未处理的错误！")
		if global.DotApp.IsDevelopmentMode() {
			stack := string(debug.Stack())
			ctx.WriteStringC(http.StatusInternalServerError, fmt.Sprintln(err)+stack)
		} else {
			_error.ServerError(ctx)
		}
	})

	//设置路由
	InitRoute(global.DotApp.HttpServer)
	//启用Session
	global.DotApp.HttpServer.SetEnabledSession(true)
	//使用Redis来保存Session
	//session过期时间 24小时
	var sessionExpire int64 = 60 * 60 * 24
	cookieName := session.DefaultSessionCookieName + "_userhome"
	storeConfig := _cache.GetSessionRedisConfig(protected.SessionRedisConfig, sessionExpire, "userhome", cookieName)
	global.DotApp.HttpServer.SetSessionConfig(storeConfig)

	//指定验证码的StoreRedis
	captchaRedis := protected.CaptchaStoreRedisConfig
	captcha.SetCustomStore(captcha.NewRedisStore(captchaRedis, 3*time.Minute))

	//使用服务器运行时缓存
	global.DotApp.SetCache(cache.NewRuntimeCache())

	//设置模板路径
	//global.DotApp.HttpServer.Renderer().SetTemplatePath("./views")

	global.InnerLogger.Debug("dotweb.StartServer => " + strconv.Itoa(appConfig.Server.Port))
	err := global.DotApp.Start()
	if err != nil {
		panic(err)
	}
	return err
}
