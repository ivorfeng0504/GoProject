package server

import (
	"git.emoney.cn/softweb/roboadvisor/config"

	"git.emoney.cn/softweb/roboadvisor/middleware/expert"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/bk1minute"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/click"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/error"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/expertnews"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/index"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/strategy"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/topic"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers/yqq"
	"github.com/devfeel/dotweb"
)

func InitRoute(server *dotweb.HttpServer) {
	server.SetVirtualPath(config.CurrentConfig.ServerVirtualPath)
	server.Router().ServerFile("/static/*filepath", config.CurrentConfig.ResourcePath+`static`)

	//mock测试路由
	server.Router().GET("/mockjson", nil)

	g := server.Group("/page")
	//2018-08-30 优化首页更改为动态数据页面
	g.GET("/home", index.Home).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/home.html", index.Home).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/home2", index.HomeView_New).Use(middleware.NewExpertSSOMiddleware())

	g.GET("/default", index.Default)
	g.GET("/live_video", index.Live_video)
	g.GET("/livevideo", index.Live_video).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/livevideo.html", index.Live_video).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/yaowen_article", index.Yaowen_article)
	g.GET("/yaowen_home", index.Yaowen_home)
	g.GET("/yaowenhome", index.Yaowen_home).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/yaowenhome.html", index.Yaowen_home).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/yaowenhome2.html", index.Yaowen_home_static_byapi).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/yaowenhome3.html", index.Yaowen_home_static).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/yuce_article", index.Yuce_article)
	g.GET("/zhuanjiacelue_article", index.Zhuanjiacelue_article)
	g.GET("/zhuanjiacelue_home", index.Zhuanjiacelue_home)
	g.GET("/celuearticle", index.CelueArticle).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/celue_home", index.CelueHome)
	g.GET("/celuehome", index.CelueHome).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/celuehome.html", index.CelueHome).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/celuehome2.html", index.CelueHomeStatic).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/zhuanjiaceluelist", index.Zhuanjiacelue_list).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/zhuti", index.Zhuti).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/zhuti.html", index.Zhuti).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/zhuti2.html", index.Zhuti_Static).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/zhuti_article", index.Zhuti_article)
	g.GET("/yqq_home", index.Yqq_home)
	g.GET("/staticHome", index.HomeView)
	g.GET("/shareArticle", index.ShareArticle)
	g.GET("/zhuanjiacelueVideo", index.ZhuanjiacelueVideo).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/zhuanjiacelueArticle", index.ZhuanjiacelueArticle).Use(middleware.NewExpertSSOMiddleware())
	// 策略资讯接口相关
	g = server.Group("/strategy")

	//获取策略资讯列表
	//调用示例
	//$.get("http://127.0.0.1:8085/strategy/GetStrategyNewsList",{ColumnID:1,StrategyID:0,currpage:1,pageSize:20},function(data){console.log(JSON.stringify(data))})
	g.GET("/GetStrategyNewsList", strategy.GetStrategyNewsList)

	//获取策略详情
	//$.get("http://127.0.0.1:8085/strategy/GetStrategyInfoBySID",{StrategyID:1},function(data){console.log(JSON.stringify(data))})
	g.GET("/GetStrategyInfoBySID", strategy.GetStrategyInfoBySID)

	//获取策略资讯详情
	//$.get("http://127.0.0.1:8085/strategy/GetStrategyNewsInfoByNewsID",{NewsID:1},function(data){console.log(JSON.stringify(data))})
	g.GET("/GetStrategyNewsInfoByNewsID", strategy.GetStrategyNewsInfoByNewsID)

	//获取热门观点
	//$.get("http://127.0.0.1:8085/strategy/GetStrategyNewsList_Top6",{},function(data){console.log(JSON.stringify(data))})
	g.GET("/GetStrategyNewsList_Top6", strategy.GetStrategyNewsList_Top6)
	g.GET("/GetStrategyNewsList_Rmgd", strategy.GetStrategyNewsList_Rmgd)

	g.GET("/GetStrategyNewsList_Multiple", nil)
	g.GET("/GetTagList", nil)

	//获取热门文章-根据点击量排序
	//$.get("http://127.0.0.1:8085/strategy/GetHotStrategyNewsList_Top10",{ColumnID:1},function(data){console.log(JSON.stringify(data))})
	g.GET("/GetHotStrategyNewsList_Top10", strategy.GetHotStrategyNewsList_Top10)

	//获取当前资讯的上一篇和下一篇资讯
	//$.get("http://127.0.0.1:8085/strategy/GetNextStrategyNewsListByCurrNews",{NewsID:1,ColumnID:1,StrategyID:0,currpage:1,pageSize:20},function(data){console.log(JSON.stringify(data))})
	g.GET("/GetNextStrategyNewsListByCurrNews", strategy.GetNextStrategyNewsListByCurrNews)

	g.GET("/AddFocusStrategy", strategy.AddFocusStrategy).Use(middleware.NewExpertAuthMiddleware())
	g.GET("/RemoveFocusStrategy", strategy.RemoveFocusStrategy).Use(middleware.NewExpertAuthMiddleware())
	g.GET("/HasFocusStrategy", strategy.HasFocusStrategy).Use(middleware.NewExpertAuthMiddleware())

	g.GET("/getnewslistbystrategyidandtype", strategy.GetNewsListByStrategyIDAndType)
	g.GET("/getnewsdetailsbystrategyidandtype", strategy.GetNewsDetailsByStrategyIDAndType)

	//主题接口相关
	g = server.Group("/topic")

	//获取主题列表（可分页）
	//$.get("http://127.0.0.1:8085/topic/GetTopicList",{},function(data){console.log(JSON.stringify(data))})
	g.GET("/GetTopicList", topic.GetTopicList)

	//获取主题详情
	//$.get("http://127.0.0.1:8085/topic/GetTopicList",{TopicID:1},function(data){console.log(JSON.stringify(data))})
	g.GET("/GetTopicInfoByID", topic.GetTopicInfoByID)

	g.GET("/HotspotHead", topic.GetHotspotById)
	g.GET("/HotspotContent", topic.GetHotspotContentById)

	//访问益圈圈数据API
	g = server.Group("/yqq")
	g.GET("/getliveroomlist", yqq.GetLiveRoomList)
	g.GET("/gethotvipliverooms", yqq.GetHotVipLiveRooms)
	g.GET("/getskilbag", yqq.GetSkilbag)
	g.GET("/getyqqstats", yqq.GetExpertLiveData)
	g.GET("/gettagsliverooms", yqq.GetTagLiveRoomInfo)
	g.GET("/getlatestlive", yqq.GetLatestRoomInfo)

	//专家资讯的策略页面合并请求
	//g.GET("/expert_strategy", yqq.Expert_StrategyChannel)
	g.GET("/getexpertliveindex", yqq.GetExpertLiveIndexRooms)
	g.GET("/gethottop4room", yqq.GetHotTop4LiveRoom)
	g.GET("/getyqqhomedata", yqq.GetYqqHomeData)
	g.GET("/expert_entireyqq", yqq.Expert_EntireYqq)

	//错误页面
	g = server.Group("/error")
	//访问地址不存在
	g.GET("/404", _error.NotFound)
	//服务器出错
	g.GET("/500", _error.ServerError)
	//身份认证出错
	g.GET("/401", _error.NotAuth)

	//点击统计相关
	g = server.Group("/click")
	//记录点击量
	//$.get("http://127.0.0.1:8085/click/addclick",{identity:1,clickType:"news.information"},function(data){console.log(JSON.stringify(data))})
	g.GET("/addclick", click.AddClick)
	//查询点击量
	//$.get("http://127.0.0.1:8085/click/queryclick",{identity:1,clickType:"news.information"},function(data){console.log(JSON.stringify(data))})
	g.GET("/queryclick", click.QueryClick)

	g.GET("/handleclick", click.HandleClick)

	//资讯相关接口
	g = server.Group("/expertnews")
	//获取今日头条
	//$.post("http://127.0.0.1:8085/expertnews/gettodaynews",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/gettodaynews", expertnews.GetTodayNews)
	g.GET("/gettodaynews", expertnews.GetTodayNews)
	//获取盘后预测Top2
	//$.post("http://127.0.0.1:8085/expertnews/getclosingnewstop2",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getclosingnewstop2", expertnews.GetClosingNewsTop2)
	g.GET("/getclosingnewstop2", expertnews.GetClosingNewsTop2)
	//获取盘后预测
	//$.post("http://127.0.0.1:8085/expertnews/getclosingnews",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getclosingnews", expertnews.GetClosingNews)
	g.GET("/getclosingnews", expertnews.GetClosingNews)
	//获取指定日期的要闻
	//$.post("http://127.0.0.1:8085/expertnews/getnews",{Date:"2018-6-12"},function(data){console.log(JSON.stringify(data))})
	g.POST("/getnews", expertnews.GetNewsInfo)
	g.GET("/getnews", expertnews.GetNewsInfo)
	//获取最新的30条要闻
	//$.post("http://127.0.0.1:8085/expertnews/getnewstop30",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getnewstop30", expertnews.GetNewsInfoTopN)
	g.GET("/getnewstop30", expertnews.GetNewsInfoTopN)
	//获取热门资讯
	//$.post("http://127.0.0.1:8085/expertnews/gethotnews",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/gethotnews", expertnews.GetHotNewsInfo)
	g.GET("/gethotnews", expertnews.GetHotNewsInfo)
	//获取主题的相关资讯
	//$.post("http://127.0.0.1:8085/expertnews/gettopicnews",{TopicId:1},function(data){console.log(JSON.stringify(data))})
	g.POST("/gettopicnews", expertnews.GetTopicNewsInfo)
	g.GET("/gettopicnews", expertnews.GetTopicNewsInfo)
	//获取资讯详情 NewsType为资讯来源类型 详见const/NewsInfoType.go  NewsInfoId为资讯Id 当查看要闻资讯时需要传递Date 当查看主题资讯时需要传递TopicId
	//$.post("http://127.0.0.1:8085/expertnews/getnewsdetail",{NewsType:1,NewsInfoId:1,Date:"2018-6-12",TopicId:1},function(data){console.log(JSON.stringify(data))})
	g.POST("/getnewsdetail", expertnews.GetNewsDetail)
	g.GET("/getnewsdetail", expertnews.GetNewsDetail)
	//通知更新资讯数据库
	//示例： http://127.0.0.1:8085/expertnews/notifyupdate?tableName=news.information&encrypt=2721f23b-cc2d-4c00-9d68-f2fbc7713075
	g.GET("/notifyupdate", expertnews.NotifyUpdate)

	//专家资讯首页接口合并(今日头条、盘后预测)
	g.GET("/index", index.GetExpertNewsIndex)
	//专家资讯首页接口合并（今日头条、策略看盘、主题、专家直播、盘后预测）
	g.GET("/indexdata", index.GetIndexData)

	// 策略资讯服务副窗口
	g = server.Group("/strategyservice")
	g.GET("/celuelive", strategyservice.StrategyServiceIndex).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/celuereading", strategyservice.StrategyServiceReadingList).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/celuevideo", strategyservice.StrategyServiceVideoList).Use(middleware.NewExpertSSOMiddleware())
	g.GET("/celuetraining", strategyservice.StrategyServiceTrainingList).Use(middleware.NewExpertSSOMiddleware())
	//获取首页最新的一条策略资讯与策略视频V1  ClientStrategyId 策略Id
	//$.post("http://127.0.0.1:8085/strategyservice/getindexstrategynews",{ClientStrategyId:80013},function(data){console.log(JSON.stringify(data))})
	g.POST("/getindexstrategynews", strategyservice.GetIndexStrategyNewsV1)
	g.GET("/getindexstrategynews", strategyservice.GetIndexStrategyNewsV1)
	//获取首页最新的一条策略资讯与策略视频V2  ClientStrategyGroupId 策略组Id
	//$.post("http://127.0.0.1:8085/strategyservice/getindexstrategynewsv2",{ClientStrategyGroupId:70001},function(data){console.log(JSON.stringify(data))})
	g.POST("/getindexstrategynewsv2", strategyservice.GetIndexStrategyNewsV2)
	g.GET("/getindexstrategynewsv2", strategyservice.GetIndexStrategyNewsV2)
	//获取首页最新的一条策略资讯与策略视频V3 ClientStrategyId 策略Id ClientStrategyGroupId 策略组Id
	//$.post("http://127.0.0.1:8085/strategyservice/getindexstrategynewsv3",{ClientStrategyId:80013,ClientStrategyGroupId:70001},function(data){console.log(JSON.stringify(data))})
	g.POST("/getindexstrategynewsv3", strategyservice.GetIndexStrategyNewsV3)
	g.GET("/getindexstrategynewsv3", strategyservice.GetIndexStrategyNewsV3)
	//获取客户端所有策略信息 ParentId为0则表示顶级菜单
	//$.post("http://127.0.0.1:8085/strategyservice/getclientstrategylist",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getclientstrategylist", strategyservice.GetClientStrategyList)
	g.GET("/getclientstrategylist", strategyservice.GetClientStrategyList)
	//获取策略指定的资讯或视频  ClientStrategyId 策略Id IsVideo=true则筛选出策略视频，否则筛选策略资讯
	//$.post("http://127.0.0.1:8085/strategyservice/getstrategynewslist",{ClientStrategyId:80013,IsVideo:false},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstrategynewslist", strategyservice.GetStrategyNewsList)
	g.GET("/getstrategynewslist", strategyservice.GetStrategyNewsList)
	//批量获取指定的策略资讯或视频的详情 NewsInfoIdList 策略资讯Id集合，逗号分隔
	//$.post("http://127.0.0.1:8085/strategyservice/getstrategynewsdetaillist",{NewsInfoIdList:"1001,1002,1003"},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstrategynewsdetaillist", strategyservice.GetStrategyNewsDetailList)
	g.GET("/getstrategynewsdetaillist", strategyservice.GetStrategyNewsDetailList)
	//获取热门策略文章 ClientStrategyIdList 用户的策略权限列表，逗号分隔
	//$.post("http://127.0.0.1:8085/strategyservice/gethotstrategynews",{ClientStrategyIdList:"80013,80015,80017"},function(data){console.log(JSON.stringify(data))})
	g.POST("/gethotstrategynews", strategyservice.GetHotStrategyNews)
	g.GET("/gethotstrategynews", strategyservice.GetHotStrategyNews)
	//获取策略相关最新资讯或视频 ClientStrategyId 策略Id IsVideo=true则筛选出策略视频，否则筛选策略资讯
	//$.post("http://127.0.0.1:8085/strategyservice/getstrategyrelationnewstopn",{ClientStrategyId:80013,IsVideo:true},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstrategyrelationnewstopn", strategyservice.GetStrategyRelationNewsTopN)
	g.GET("/getstrategyrelationnewstopn", strategyservice.GetStrategyRelationNewsTopN)
	//获取热门策略视频 ClientStrategyIdList 用户的策略权限列表，逗号分隔
	//$.post("http://127.0.0.1:8085/strategyservice/gethotstrategyvideonews",{ClientStrategyIdList:"80013,80015,80017"},function(data){console.log(JSON.stringify(data))})
	g.POST("/gethotstrategyvideonews", strategyservice.GetHotStrategyVideoNews)
	g.GET("/gethotstrategyvideonews", strategyservice.GetHotStrategyVideoNews)
	//获取最新策略资讯 ClientStrategyIdList 用户的策略权限列表，逗号分隔
	//$.post("http://127.0.0.1:8085/strategyservice/getnewststrategynews",{ClientStrategyIdList:"80013,80015,80017"},function(data){console.log(JSON.stringify(data))})
	g.POST("/getnewststrategynews", strategyservice.GetNewstStrategyNews)
	g.GET("/getnewststrategynews", strategyservice.GetNewstStrategyNews)
	//获取最新策略视频 ClientStrategyIdList 用户的策略权限列表，逗号分隔
	//$.post("http://127.0.0.1:8085/strategyservice/getnewststrategyvideonews",{ClientStrategyIdList:"80013,80015,80017"},function(data){console.log(JSON.stringify(data))})
	g.POST("/getnewststrategyvideonews", strategyservice.GetNewstStrategyVideoNews)
	g.GET("/getnewststrategyvideonews", strategyservice.GetNewstStrategyVideoNews)
	//根据策略Id获取直播间Id信息 ClientStrategyId 策略Id
	//$.post("http://127.0.0.1:8085/strategyservice/getliveroom",{ClientStrategyId:80013},function(data){console.log(JSON.stringify(data))})
	g.POST("/getliveroom", strategyservice.GetLiveRoomByStrategyId)
	g.GET("/getliveroom", strategyservice.GetLiveRoomByStrategyId)
	//$.post("http://127.0.0.1:8085/strategyservice/getstrategygroupnewslist",{StrategyGroupId:80013,IsVideo:false},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstrategygroupnewslist", strategyservice.GetStrategyGroupNewsList)
	g.GET("/getstrategygroupnewslist", strategyservice.GetStrategyGroupNewsList)

	//GetStrategyTagNewsList 获取策略指定的标签资讯 （专家资讯栏目）
	g.GET("/getstrategytagnewslist", strategyservice.GetStrategyTagNewsList)

	//获取策略指定的资讯或视频  ClientStrategyId 策略Id IsVideo=true则筛选出策略视频，否则筛选策略资讯
	//$.post("http://127.0.0.1:8085/strategyservice/getstrategynewslist_multimedia",{ClientStrategyId:80013},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstrategynewslist_multimedia", strategyservice.GetStrategyNewsList_MultiMedia)
	g.GET("/getstrategynewslist_multimedia", strategyservice.GetStrategyNewsList_MultiMedia)
	//批量获取指定的策略资讯或视频的详情 NewsInfoIdList 策略资讯Id集合，逗号分隔
	//$.post("http://127.0.0.1:8085/strategyservice/getstrategynewsdetaillist_multimedia",{NewsInfoIdList:"1001,1002,1003"},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstrategynewsdetaillist_multimedia", strategyservice.GetStrategyNewsDetailList_MultiMedia)
	g.GET("/getstrategynewsdetaillist_multimedia", strategyservice.GetStrategyNewsDetailList_MultiMedia)

	//获取热门培训
	g.POST("/gethotstrategymultimedianews", strategyservice.GetHotStrategyMultiMediaNews)
	g.GET("/gethotstrategymultimedianews", strategyservice.GetHotStrategyMultiMediaNews)

	//获取实战培训列表接口
	g.GET("/gettaglist_szpx", strategyservice.GetTagList_SZPX)

	//获取重要提示
	g.POST("/getnewslist_importanttips",strategyservice.GetNewsList_ImportantTips)
	g.GET("/getnewslist_importanttips",strategyservice.GetNewsList_ImportantTips)

	//获取最近一个月实战培训数据列表
	g.POST("/getstrategynewslist_multimedia1month",strategyservice.GetStrategyNewsList_MultiMedia1Month)
	g.GET("/getstrategynewslist_multimedia1month",strategyservice.GetStrategyNewsList_MultiMedia1Month)

	// 板块1分钟
	g = server.Group("/bk1minutes")
	g.GET("/index.html", bk1minute.IndexView)
	g.GET("/index", bk1minute.IndexView)
	g.GET("/bkminute_up.html", bk1minute.BkMinute_UpView)
	g.GET("/bkminute_up", bk1minute.BkMinute_UpView)
	g.GET("/bkminute_down.html", bk1minute.BkMinute_DownView)
	g.GET("/bkminute_down", bk1minute.BkMinute_DownView)
}
