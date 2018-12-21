package server

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	myoptional_middleware "git.emoney.cn/softweb/roboadvisor/middleware/myoptional"
	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/error"
	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/evaluation"
	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/myoptional"
	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/stocknews"
	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/stocktalk"
	"github.com/devfeel/dotweb"

	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/blocknews"
	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/guide"
	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/validatecode"
)

func InitRoute(server *dotweb.HttpServer) {
	server.SetVirtualPath(config.CurrentConfig.ServerVirtualPath)
	server.Router().ServerFile("/static/*filepath", config.CurrentConfig.ResourcePath+`static`)

	g := server.Group("/myoptional")
	// 相关资讯详情
	g.GET("/relatedarticle", myoptional.RelatedArticle).Use(myoptional_middleware.NewMyOptionalSSOMiddleware())
	// 微股吧和策略首页
	g.GET("/index", myoptional.Index).Use(myoptional_middleware.NewMyOptionalSSOMiddleware())
	// 个股微股吧评论
	g.GET("/saysomething", myoptional.SaySomething).Use(myoptional_middleware.NewMyOptionalSSOMiddleware())
	// 微股吧和策略首页-静态页面-测试临时跳转
	g.GET("/indexstatic", myoptional.IndexStatic).Use(myoptional_middleware.NewMyOptionalSSOMiddleware())
	//板块相关资讯
	g.GET("/bknews", myoptional.BKNews)
	g = server.Group("/stocktalk")
	//调用示例-提交评论
	g.POST("/addstocktalk", stocktalk.InsertStockTalk).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	//调用示例-查询最新的评论信息
	g.POST("/getstocktalk", stocktalk.GetStockTalkListByDate).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	g.GET("/getstocktalk", stocktalk.GetStockTalkListByDate).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	//调用示例-分页获取微股吧评论
	//$.post("http://127.0.0.1:8086/stocktalk/getstocktalkbypage",{PageIndex:0,PageSize:10,StockCodeList:"000001,000002"},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstocktalkbypage", stocktalk.GetStockTalkListByPage).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	g.GET("/getstocktalkbypage", stocktalk.GetStockTalkListByPage).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())

	g = server.Group("/evaluation")
	// 测评首页
	g.GET("/index", evaluation.Index).Use(myoptional_middleware.NewMyOptionalSSOMiddleware())
	// 测评弹窗
	g.GET("/gotest", evaluation.GoTest).Use(myoptional_middleware.NewMyOptionalSSOMiddleware())
	g.GET("/discernpic", evaluation.DiscernPic)
	//调用示例-查看当前用户的评测结果，如果Message为null，则未参加测评
	//$.post("http://127.0.0.1:8086/evaluation/getresult",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getresult", evaluation.GetResult).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	//查看当前用户的评测结果 客户端调用接口
	g.GET("/getresult", evaluation.GetResultForClient)
	//调用示例-提交评测
	//$.post("http://127.0.0.1:8086/evaluation/submitresult",{
	//        "InvestTarget": "投资目标A1",
	//        "InvestTargetDesc": "投资目标-描述",
	//        "InvestTargetTip": "投资目标-提示",
	//        "ChooseStockReason": "选股主要考虑",
	//        "ChooseStockReasonDesc": "选股主要考虑-描述",
	//        "ChooseStockReasonTip": "选股主要考虑-提示",
	//        "HoldStockTime": "股票一般拿多久",
	//        "HoldStockTimeDesc": "股票一般拿多久-描述",
	//        "HoldStockTimeTip": "股票一般拿多久-提示",
	//        "BuyStyle": "操作买点",
	//        "BuyStyleDesc": "操作买点-描述",
	//        "BuyStyleTip": "操作买点-提示",
	//},function(data){console.log(JSON.stringify(data))})
	g.POST("/submitresult", evaluation.SubmitResult).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	//调用示例-上传图片，识别股票
	//http://127.0.0.1:8086/evaluation/analyzestock
	g.GET("/upfile", evaluation.Upfile)

	g.POST("/analyzestock", evaluation.AnalyzeStock)
	//调用示例-获取推荐策略及股票列表
	//$.post("http://127.0.0.1:8086/evaluation/getstrategylist",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstrategylist", evaluation.GetStrategyList).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	g.GET("/getstrategylist", evaluation.GetStrategyList).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	//股票相关资讯
	g = server.Group("/stocknews")
	//调用示例-获取股票相关资讯
	//$.post("http://127.0.0.1:8086/stocknews/getstocknews",{StockList:"000001,000002"},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstocknews", stocknews.GetStockNewsInfoByHSet).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())
	g.GET("/getstocknews", stocknews.GetStockNewsInfoByHSet).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())

	g = server.Group("/strategy")
	g.GET("/getfocusstrategylist", myoptional.GetFocusStrategyList).Use(myoptional_middleware.NewMyOptionalAuthMiddleware())

	g = server.Group("/blocknews")
	//获取板块相关资讯
	//$.post("http://127.0.0.1:8086/blocknews/getblocknews",{BlockCode:"1077",PageSize:2},function(data){console.log(JSON.stringify(data))})
	g.POST("/getblocknews", blocknews.GetBlockNewsInfomation)
	g.GET("/getblocknews", blocknews.GetBlockNewsInfomation)

	g = server.Group("/captcha")
	g.GET("/page", validatecode.ShowCaptchaPage)
	g.GET("/image", validatecode.BuffImage)
	g.GET("/reloadimage", validatecode.BuffNewImage)
	g.POST("/verify", validatecode.VerifyCaptcha)
	g.GET("/fetchid", validatecode.GetCaptchaId)

	g = server.Group("/guide")
	// 新手引导 获取推荐股池
	g.GET("/getstocklist", guide.GetRecommendStockList)

	//错误页面
	g = server.Group("/error")
	//访问地址不存在
	g.GET("/404", _error.NotFound)
	//服务器出错
	g.GET("/500", _error.ServerError)
	//身份认证出错
	g.GET("/401", _error.NotAuth)
}
