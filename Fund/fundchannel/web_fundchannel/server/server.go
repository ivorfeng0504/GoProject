package server

import (
	fundconfig "emoney.cn/fundchannel/config"
	"emoney.cn/fundchannel/global"
	"emoney.cn/fundchannel/protected"
	"emoney.cn/fundchannel/util/cache"
	_ "fmt"
	"github.com/devfeel/cache"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/session"
	"github.com/devfeel/middleware/cors"
	"strconv"
)

func StartServer(configPath string) error {

	appConfig := config.MustInitConfig(configPath + "/dotweb.conf")
	global.DotApp = dotweb.ClassicWithConf(appConfig)

	global.DotApp.SetDevelopmentMode()
	global.DotApp.UseRequestLog()
	global.DotApp.Use(cors.Middleware(cors.NewConfig().UseDefault()))

	//启用Session
	global.DotApp.HttpServer.SetEnabledSession(true)
	//使用Redis来保存Session
	//session过期时间 24小时
	var sessionExpire int64 = 60 * 60 * 24
	cookieName := session.DefaultSessionCookieName + "_fundchan"
	storeConfig := _cache.GetSessionRedisConfig(protected.SessionRedisConfig, sessionExpire, "fundchan", cookieName)
	global.DotApp.HttpServer.SetSessionConfig(storeConfig)

	//使用服务器运行时缓存
	global.DotApp.SetCache(cache.NewRuntimeCache())
	//设置路由
	InitRoute(global.DotApp.HttpServer)

	//设置模板路径
	global.DotApp.HttpServer.Renderer().SetTemplatePath(fundconfig.CurrentConfig.ResourcePath + "./views")

	global.InnerLogger.Debug("dotweb.StartServer => " + strconv.Itoa(appConfig.Server.Port))
	err := global.DotApp.Start()
	if err != nil {
		panic(err)
	}
	return err
}
