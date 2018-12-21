package server

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/web_api_stocktalk/handlers/stocknews"
	"git.emoney.cn/softweb/roboadvisor/web_api_stocktalk/handlers/stocktalk"
	"github.com/devfeel/dotweb"
)

func InitRoute(server *dotweb.HttpServer) {
	server.SetVirtualPath(config.CurrentConfig.ServerVirtualPath)
	server.Router().ServerFile("/static/*filepath", config.CurrentConfig.ResourcePath+`static`)

	//微股吧
	g := server.Group("/api/stocktalk")
	g.POST("/insertstocktalk", stocktalk.InsertStockTalk)
	g.POST("/getstocktalklistbydate", stocktalk.GetStockTalkListByDate)
	g.POST("/getstocktalklistbypage", stocktalk.GetStockTalkListByPage)
	g.POST("/sendstocktalkmsg", stocktalk.SendStockTalkMsg)
	g.POST("/getstocktalkbystockcode", stocktalk.GetStockTalkByStockCode)
	g.GET("/getstocktalkbystockcode", stocktalk.GetStockTalkByStockCode)
	g.GET("/getstockcodelist", stocktalk.GetStockCodeList)
	g.POST("/getstockcodelist", stocktalk.GetStockCodeList)
	//手动同步微股吧基本面数据
	g.GET("/syncstocktalkbasic", stocktalk.SyncStockTalkForStockBasicInfoByStockList)

	//股票相关资讯
	g = server.Group("/api/stocknews")
	g.POST("/getstocknewsinfo", stocknews.GetStockNewsInfo)
	g.POST("/getstocknewsinfobyhset", stocknews.GetStockNewsInfoByHSet)
}
