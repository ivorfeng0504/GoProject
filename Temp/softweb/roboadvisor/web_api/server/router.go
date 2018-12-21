package server

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/click"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/expertnews"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/live"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/market"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/myoptional"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/panwai"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/strategy"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/train"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/user"
	"git.emoney.cn/softweb/roboadvisor/web_api/handlers/userhome"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/topic"
	"github.com/devfeel/dotweb"
)

func InitRoute(server *dotweb.HttpServer) {
	server.SetVirtualPath(config.CurrentConfig.ServerVirtualPath)
	server.Router().ServerFile("/static/*filepath", config.CurrentConfig.ResourcePath+`static`)
	g := server.Group("/api/live")

	//调用示例
	//$.post("http://127.0.0.1:8083/api/live/istradetime",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/istradetime", live.IsTradeTime)
	g.POST("/addquestion", live.AddQuestion)
	g.POST("/livequestion", live.GetLiveQuestionAnswerList)
	g.POST("/livecontent", live.GetLiveContent)
	g.POST("/hasnewmessage", live.HasNewMessage)
	g.POST("/addroompermit", live.AddLiveUserInRoom)
	g.POST("/removeroompermit", live.RemoveLiveUserInRoom)
	g.POST("/getuserroomlist", live.GetUserRoomList)

	//用户操作
	g.POST("/getuserbyuserid", live.GetUserById)
	g.POST("/getuserbyaccount", live.GetUserByAccount)
	g.POST("/getuserbyuid", live.GetUserByUID)
	g.POST("/adduser", live.AddUser)

	//websocket服务回调地址  授权校验
	g.GET("/checktoken", live.CheckToken)
	//long poll服务回调地址  获取消息
	g.GET("/hasnewmessage", live.GetHasNewMessage)

	g = server.Group("/api/panwai")
	//调用示例
	//$.get("http://127.0.0.1:8083/api/panwai/newslist",{ColumnID:1,StrategyID:0},function(data){console.log(JSON.stringify(data))})
	g.GET("/newslist", panwai.GetNewsByColIDAndStrategyID)
	g.GET("/newsinfo", panwai.GetNewsInfoByID)
	g.GET("/seriesnewslist", panwai.GetSeriesNewsListByNewsID)

	g = server.Group("/api/strategy")
	g.GET("/getstrategybyid", strategy.GetStrategyInfo)
	g.GET("/getstrategylist", strategy.GetStrategyList)

	//客户端行情接口
	g = server.Group("/api/market")
	g.GET("/getatmosphere", market.GetAtmosphere)
	g.GET("/getstrategypool", market.GetStrategyPool)           //策略池 CL_id对应不同操作
	g.GET("/bigorderstrategyhero", market.BigOrderStrategyHero) //英雄榜 CL_id对应不同操作
	g.GET("/unusualgroup", market.UnusualGroup)                 //行情指数标签
	g.GET("/gettodaystrategy", market.GetTodayStrategy)         //今日策略
	g.GET("/getstrategylist", market.GetStrategyList)           //全部策略
	g.GET("/activegoodsprimary", market.ActiveGoodsPrimary)     //策略预选池
	g.GET("/commonpool", market.CommonPool)                     //通用池
	g.GET("/strategyoutpool", market.StrategyOutCommonPool)
	g.GET("/ontradingpool", market.OnTradingPool)
	g.GET("/aftertradingpool", market.AfterTradingPool) //短线宝类通用股池 6370

	//活动相关接口
	g = server.Group("/api/activity")
	g.POST("/getactivitylist", userhome.GetActivityList)
	g.POST("/getuserinactivitymap", userhome.GetUserInActivityMap)
	g.POST("/getuserserialaward", userhome.GetUserSerialLoginRule)
	g.POST("/getuserlevel", userhome.GetUserLevel)
	g.POST("/getuserawardlist", userhome.GetUserAwardList)
	g.POST("/getguessmarkethistoryinfo", userhome.GetGuessMarketHistoryInfo)
	g.POST("/getactivityawardlist", userhome.GetActivityAwardListByActivityId)
	g.POST("/getguessmarketnumber", userhome.GetGuessMarketNumber)
	g.POST("/getcurrentguessinfo", userhome.GetCurrentGuessInfo)
	g.POST("/getusermedal", userhome.GetUserMedal)
	g.POST("/insertuserstock", userhome.InsertOrUpdateUserStock)
	g.POST("/getuserstocktoday", userhome.GetUserStockToday)
	g.POST("/getuserstockhistory", userhome.GetUserStockHistory)
	g.POST("/getnewststockpool", userhome.GetNewstStockPool)
	//猜涨跌相关
	g.POST("/currentguesschange", userhome.GetCurrentGuessChange)
	g.POST("/guesschangesubmit", userhome.GuessChangeSubmit)
	g.POST("/myguesschangeinfocurrentweek", userhome.GetMyGuessChangeInfoCurrentWeek)
	g.POST("/myguessawardlist", userhome.GetMyGuessageAwardList)
	g.POST("/myguesschangeinfonewst", userhome.GetMyGuessChangeInfoNewst)

	//刷新种子用户信息
	g.POST("/refreshseeduserlist", userhome.RefreshSeedUserInfoList)
	g.GET("/refreshseeduserlist", userhome.RefreshSeedUserInfoList)

	//专家资讯相关接口
	g = server.Group("/api/expertnews")
	g.GET("/notifyupdate", expertnews.NotifyUpdateNews)
	g.POST("/notifyupdate", expertnews.NotifyUpdateNews)
	g.POST("/gettodaynews", expertnews.GetTodayNews)
	g.POST("/getclosingnews", expertnews.GetClosingNews)
	g.POST("/getnews", expertnews.GetNewsInfo)
	g.POST("/getnewstopn", expertnews.GetNewsInfoTopN)
	g.POST("/gethotnews", expertnews.GetHotNewsInfo)
	g.POST("/getyaowenpagedata", expertnews.GetYaoWenHomePageData)
	g.POST("/gettopicnews", expertnews.GetTopicNewsInfo)
	g.POST("/getnewsdetail", expertnews.GetNewsInfoDetail)
	g.POST("/receivelatestlivetostrategy", expertnews.ReceiveLatestLiveToStrategy)
	g.POST("/receivelatestnewstostrategy", expertnews.ReceiveLatestNewsToStrategy)
	g.POST("/getblocknewsinfomation", expertnews.GetBlockNewsInfomation)
	g.GET("/getstrategynewslist", expertnews.GetStrategyNewsList)
	g.GET("/getstrategyinfobysid", expertnews.GetStrategyInfoBySID)
	g.GET("/getstrategynewslist_top6", expertnews.GetStrategyNewsList_Top6)
	g.GET("/getstrategynewslistbypage", expertnews.GetStrategyNewsListByPage)
	g.GET("/getstrategynewsinfobynewsid", expertnews.GetStrategyNewsInfoByNewsID)

	g.GET("/getstrategynewslistbypage_app", expertnews.GetStrategyNewsListByPage_app)
	g.GET("/getstrategynewsinfobynewsid_app", expertnews.GetStrategyNewsInfoByNewsID_app)
	g.GET("/getstrategyinfobysid_app", expertnews.GetStrategyInfoBySID_app)
	g.GET("/getstrategynewslistbysid_app", expertnews.GetStrategyNewsListBySID_app)
	g.GET("/getstrategylist_app", expertnews.GetStrategyList_app)

	g.GET("/getnewslistbystrategyidandtype", expertnews.GetNewsListByStrategyIDAndType)
	g.GET("/getnewsdetailsbystrategyidandtype", expertnews.GetNewsDetailsByStrategyIDAndType)

	//专家资讯-主题
	g = server.Group("/api/topic")
	g.GET("/gettopiclist", topic.GetTopicList)
	g.GET("/gettopicinfobyid", topic.GetTopicInfoByID)

	//点击量统计相关接口
	g = server.Group("/api/click")
	g.POST("/addclick", click.AddClick)
	g.POST("/queryclick", click.QueryClick)
	g.POST("/handleclick", click.HandleClick)

	//退款相关接口
	g = server.Group("/api/refundorder")
	g.POST("/queryorderlist", userhome.QueryOrderList)
	g.POST("/validateorder", userhome.ValidateOrder)
	g.POST("/refundsubmit", userhome.RefundSubmit)
	g.POST("/getrefundstatus", userhome.GetRefundStatus)

	//我的自选相关接口
	g = server.Group("/api/myoptional")
	g.GET("/getfocusstrategylist", expertnews.GetFocusStrategyList)

	//微股吧
	g = server.Group("/api/stocktalk")
	g.POST("/insertstocktalk", myoptional.InsertStockTalk)
	g.POST("/getstocktalklistbydate", myoptional.GetStockTalkListByDate)
	g.POST("/getstocktalklistbypage", myoptional.GetStockTalkListByPage)
	g.POST("/sendstocktalkmsg", myoptional.SendStockTalkMsg)
	g.POST("/getstocktalkbystockcode", myoptional.GetStockTalkByStockCode)
	g.GET("/getstocktalkbystockcode", myoptional.GetStockTalkByStockCode)
	g.GET("/getstockcodelist", myoptional.GetStockCodeList)
	g.POST("/getstockcodelist", myoptional.GetStockCodeList)
	//手动同步微股吧基本面数据
	g.GET("/syncstocktalkbasic", myoptional.SyncStockTalkForStockBasicInfoByStockList)

	//自选股测评
	g = server.Group("/api/evaluation")
	g.POST("/getresult", myoptional.GetResult)
	g.POST("/submitresult", myoptional.SubmitResult)
	g.POST("/getstrategylist", myoptional.GetStrategyList)

	//股票相关资讯
	g = server.Group("/api/stocknews")
	g.POST("/getstocknewsinfo", myoptional.GetStockNewsInfo)
	g.POST("/getstocknewsinfobyhset", myoptional.GetStockNewsInfoByHSet)

	//策略资讯副窗口
	g = server.Group("/api/strategyservice")
	g.POST("/getstrategylist", strategyservice.GetClientStrategyInfoList)
	g.POST("/getstrategylistmock", strategyservice.GetClientStrategyInfoListMock)
	g.POST("/getstrategyliveroomdict", strategyservice.GetStrategyLiveRoomDict)
	g.POST("/getstrategynewslist", strategyservice.GetStrategyNewsList)
	g.POST("/getstrategynewslistbystrategyid", strategyservice.GetStrategyNewsListByStrategyId)
	g.POST("/getstrategynewsdetaillist", strategyservice.GetStrategyNewsDetailList)
	g.POST("/gethotstrategynewsinfo", strategyservice.GetHotStrategyNewsInfo)
	g.POST("/getnewststrategynews", strategyservice.GetNewstStrategyNews)
	g.POST("/getindexstrategynews", strategyservice.GetIndexStrategyNews)
	g.POST("/getstrategynewslist_multimedia", strategyservice.GetStrategyNewsList_MultiMedia)
	g.POST("/getstrategynewsdetaillist_multimedia", strategyservice.GetStrategyNewsDetailList_MultiMedia)
	g.POST("/getnewslist_importanttips", strategyservice.GetNewsList_ImportantTips)
	g.POST("/getstrategynewslist_multimedia1month",strategyservice.GetStrategyNewsList_MultiMedia1Month)

	//账号体系相关接口
	g = server.Group("/api/user")
	g.GET("/RegMobile", user.Reg_Mobile)
	g.GET("/RegQQorWechat", user.Reg_QQorWeChat)
	g.GET("/ResetPwd", user.ResetPassWord)
	g.GET("/getloginidbyname", user.GetLoginIDByName)
	g.GET("/boundgroupqrylogin", user.BoundGroupQryLogin)
	g.GET("/boundgroupaddlogin", user.BoundGroupAddLogin)
	g.GET("/boundgrouprmvlogin", user.BoundGroupRmvLogin)

	g = server.Group("/api/train")
	g.GET("/gettrainlistbydateandarea", train.GetTrainListByDateAndArea)
	g.GET("/gettrainlistbytag", train.GetTrainListByTag)
	g.POST("/gettrainlistbydate", train.GetTrainListByDate)
	g.POST("/gettrainlistbydateandtag", train.GetTrainListByDateAndTag)
}
