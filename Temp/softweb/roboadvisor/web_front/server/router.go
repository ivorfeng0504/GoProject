package server

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/middleware"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/error"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/live"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/panwai"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/strategy"
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/train"
	//middleware2 "git.emoney.cn/softweb/roboadvisor/middleware/train"
	trainMiddleware "git.emoney.cn/softweb/roboadvisor/middleware/train"
)

func InitRoute(server *dotweb.HttpServer) {
	server.SetVirtualPath(config.CurrentConfig.ServerVirtualPath)
	server.Router().ServerFile("/static/*filepath", config.CurrentConfig.ResourcePath+`static`)
	g := server.Group("/live")

	//调用示例-查看直播内容
	//最新全部 $.post("http://127.0.0.1:8082/live/livecontent",{Date:"2018-5-2"},function(data){console.log(JSON.stringify(data))})
	//最新递增 $.post("http://127.0.0.1:8082/live/livecontent",{Date:"2018-5-2",IsIncre:"1"},function(data){console.log(JSON.stringify(data))})
	g.POST("/livecontent", live.GetLiveContent).Use(middleware.NewAuthMiddleware())

	//调用示例-查看直播问答
	//查看所有    $.post("http://127.0.0.1:8082/live/livequestion",{Date:"2018-5-2"},function(data){console.log(JSON.stringify(data))})
	//仅查看自己  $.post("http://127.0.0.1:8082/live/livequestion",{Date:"2018-5-2",IsSelf:"1"},function(data){console.log(JSON.stringify(data))})
	g.POST("/livequestion", live.GetLiveQuestionAnswerList).Use(middleware.NewAuthMiddleware())

	//调用示例-用户提问
	//$.post("http://127.0.0.1:8082/live/addquestion",{AskContent:"老师今天行情怎么样？"},function(data){console.log(JSON.stringify(data))})
	g.POST("/addquestion", live.AddQuestion).Use(middleware.NewAuthMiddleware())

	//调用示例-查看用户直播间权限
	//$.post("http://127.0.0.1:8082/live/getuserroomlist",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getuserroomlist", live.GetRoomList)

	//图文直播首页
	g.GET("/index", live.Index).Use(middleware.NewSSOMiddleware())

	//调用示例
	//$.post("http://127.0.0.1:8082/live/istradetime",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/istradetime", live.IsTradeTime)

	//暂不使用 已迁移到api中
	//调用示例
	//http://127.0.0.1:8082/live/hasnewmessage?querykey=2018-04-04%2008:11:12&groupid=180
	//g.GET("/hasnewmessage", live.HasNewMessage)

	g = server.Group("/panwai")
	//调用示例
	//$.get("http://127.0.0.1:8082/panwai/bannerlist",{},function(data){console.log(JSON.stringify(data))})
	g.GET("/bannerlist", panwai.GetBannerInfo).Use(middleware.NewAuthMiddleware())
	g.GET("/newslist", panwai.GetNewsListByColumnID).Use(middleware.NewAuthMiddleware())
	g.GET("/seriesnewslist", panwai.GetSeriesNewsListByNewsID).Use(middleware.NewAuthMiddleware())

	//$.get("http://127.0.0.1:8082/panwai/newsinfo",{NewsId:1},function(data){console.log(JSON.stringify(data))})
	g.GET("/newsinfo", panwai.GetNewsInfoByID).Use(middleware.NewAuthMiddleware())

	//$.get("http://127.0.0.1:8082/panwai/topnewslist",{ColumnID:1},function(data){console.log(JSON.stringify(data))})
	g.GET("/topnewslist", panwai.GetTopNewsInfoByColumnID).Use(middleware.NewAuthMiddleware())

	g.GET("/article", panwai.Article).Use(middleware.NewSSOMiddleware())
	g.GET("/lesson", panwai.Lesson).Use(middleware.NewSSOMiddleware())
	g.GET("/pwgf", panwai.PWGF).Use(middleware.NewSSOMiddleware())
	g.GET("/serieslesson", panwai.SeriesLesson).Use(middleware.NewSSOMiddleware())

	g = server.Group("/strategy")
	//智投软件访问策略接口
	g.GET("/getstrategy", strategy.GetStrategyInfo)

	//策略访问页面
	g.GET("/strategy", strategy.Strategy)
	g.GET("/detail", strategy.StrategyDetail)

	g = server.Group("/train")
	//g.GET("/index.html", train.TrainIndex).Use(trainMiddleware.NewTrainSSOMiddleware())
	//g.GET("/index", train.TrainIndex).Use(trainMiddleware.NewTrainSSOMiddleware())
	g.GET("/index", train.Index).Use(trainMiddleware.NewTrainSSOMiddleware())
	g.GET("/userguidor.html", train.TrainGuidor).Use(trainMiddleware.NewTrainAuthMiddleware())
	g.GET("/gettrainlistbydateandarea", train.GetTrainListByDateAndArea).Use(trainMiddleware.NewTrainAuthMiddleware())
	g.GET("/gettrainlistbytag", train.GetTrainListByTag).Use(trainMiddleware.NewTrainAuthMiddleware())
	g.GET("/gettrainlistbynew",train.GetTrainListByNew).Use(trainMiddleware.NewTrainAuthMiddleware())

	g.GET("/gettrainlistbydate",train.GetTrainListByDate).Use(trainMiddleware.NewTrainAuthMiddleware())
	g.GET("/gettrainlistbydateandtag",train.GetTrainListByDateAndTag).Use(trainMiddleware.NewTrainAuthMiddleware())
	g.GET("/getstrategynewslist_multimedia1month",train.GetStrategyNewsList_MultiMedia1Month).Use(trainMiddleware.NewTrainAuthMiddleware())
	g.GET("/getstrategynewslist_multimedia",train.GetStrategyNewsList_MultiMedia).Use(trainMiddleware.NewTrainAuthMiddleware())
	g.GET("/gettrainclientinfo",train.GetTrainClientInfo)
	g.GET("/gettraintaginfo",train.GetTrainTagInfo)

	//错误页面
	g = server.Group("/error")
	//访问地址不存在
	g.GET("/404", _error.NotFound)
	//服务器出错
	g.GET("/500", _error.ServerError)
	//身份认证出错
	g.GET("/401", _error.NotAuth)

}
