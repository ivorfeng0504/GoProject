package server

import (
	"github.com/devfeel/dotweb"

	"emoney.cn/fundchannel/config"
	"emoney.cn/fundchannel/web_fundchannel/handlers/home"
	"emoney.cn/fundchannel/web_fundchannel/handlers/error"
	"emoney.cn/fundchannel/middleware"
)

func InitRoute(server *dotweb.HttpServer) {
	server.SetVirtualPath(config.CurrentConfig.ServerVirtualPath)

	server.Router().ServerFile("/static/*filepath", config.CurrentConfig.ResourcePath+`static`)

	g := server.Group("/fundchannel/home")
	g.GET("/index", home.Index).Use(middleware.NewSSOMiddleware())
	g.GET("/strategy", home.Strategy).Use(middleware.NewSSOMiddleware())
	g.GET("/live", home.Live).Use(middleware.NewSSOMiddleware())

	g.GET("/QueryUserRiskInfo", home.QueryUserRiskInfo)
	g.GET("/GetEncryptMobile", home.GetEncryptMobile)

	g.GET("/GetStrategyList", home.GetStrategyList)
	g.GET("/GetStrategyInfoByCode", home.GetStrategyInfoByCode)
	g.GET("/GetFundTimeLineByCode", home.GetFundTimeLineByCode)

	g.GET("/GetYqqContent", home.GetYqqContent)
	g.GET("/GetLiveHeadInfo", home.GetLiveHeadInfo)
	g.GET("/GetYqqLiveAllContent", home.GetYqqLiveAllContent)
	g.GET("/GetYqqAllQuestion", home.GetYqqAllQuestion)
	g.GET("/GetYqqMyQuestion", home.GetYqqMyQuestion)

	g.GET("/GetUserInfoByUID", home.GetUserInfoByUID)

	//错误页面
	g = server.Group("/error")
	//访问地址不存在
	g.GET("/404", _error.NotFound)
	//服务器出错
	g.GET("/500", _error.ServerError)
	//身份认证出错
	g.GET("/401", _error.NotAuth)
}
