package index

import (
	"encoding/json"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/expert"
	model "git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/model/yqq"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	expertnews_srv "git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers"
	"github.com/devfeel/dotweb"
	"sort"
	"strings"
	"time"

	"github.com/devfeel/dotlog/util/json"
	"strconv"
)

//获取runtimcache
var rtcache = handlers.NewWebRuntimeCache("index_", 60)

// Home 首页
func Home(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "home.html")
}

func Default(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "default.html")
}

// 视频
func Live_video(ctx dotweb.Context) error {
	loginuser := contract_expert.ExpertUserInfo(ctx)
	if loginuser != nil {
		pid := loginuser.PID

		//付费版、体验期内跳转zhuti.html
		//付费过期、体验过期、游客跳转zhuti_ad.html
		//FreePID:="888010400,888012400,888030000,888030400,888020400"
		PayPID := "888010000,888020000,888012000,888012001"
		if strings.Contains(PayPID, strconv.Itoa(pid)) {
			return contract_expert.RenderExpertHtml(ctx, "live_video.html")
		} else {
			return contract_expert.RenderExpertHtml(ctx, "video_ad.html")
		}
	} else {
		return contract_expert.RenderExpertHtml(ctx, "video_ad.html")
	}
}

// 要闻-详情页
func Yaowen_article(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "yaowen_article.html")
}

// 要闻
func Yaowen_home(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "yaowen_home.html")
}

// 要闻
func Yaowen_home_static_byapi(ctx dotweb.Context) error {
	pageData, _ := agent.GetYaoWenHomePageData()
	return contract_expert.RenderExpertHtmlWithPageData(ctx, "yaowen_home_static.html", pageData)
}

// 要闻
func Yaowen_home_static(ctx dotweb.Context) error {
	cachekey := "GetYaoWenHomePageDataByRedis"
	var pageData *agent.YaoWenHomePageData
	obj, exist := rtcache.GetCache(cachekey)
	if exist == false {
		newsSrv := expertnews_srv.NewNewsInformationService()
		pageData, _ = newsSrv.GetYaoWenHomePageData()
		if pageData != nil {
			rtcache.SetCache(cachekey, pageData)
		}
	} else {
		pageData = obj.(*agent.YaoWenHomePageData)
	}
	return contract_expert.RenderExpertHtmlWithPageData(ctx, "yaowen_home_static.html", pageData)
}

// 预测
func Yuce_article(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "yuce_article.html")
}

// 策略-详情页
func Zhuanjiacelue_article(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "zhuanjiacelue_article.html")
}

// 策略-首页
func Zhuanjiacelue_home(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "zhuanjiacelue_home.html")
}

func CelueArticle(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "celue_article.html")
}

func CelueHome(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "celue_home.html")
}

func CelueHomeStatic(ctx dotweb.Context) error {
	cachekey := "Expert_EntireYqq"
	obj, exist := rtcache.GetCache(cachekey)

	if exist == false {
		YqqLiveService := expertnews.ConExpertNews_YqqService()
		entire_obj, err := YqqLiveService.GetExpertYqqEntireData_RedisCache()
		if err == nil {
			obj = struct {
				RetCode string
				RetMsg  string
				Data    interface{}
			}{
				RetCode: "0", RetMsg: "OK", Data: entire_obj}
			rtcache.SetCache(cachekey, obj)
		}
	}
	return contract_expert.RenderExpertHtmlWithPageData(ctx, "celue_home_static.html", obj)
}

// 策略-列表页
func Zhuanjiacelue_list(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "zhuanjiacelue_list.html")
}

// 主题-首页
func Zhuti(ctx dotweb.Context) error {
	loginuser := contract_expert.ExpertUserInfo(ctx)
	if loginuser != nil {
		pid := loginuser.PID

		//付费版、体验期内跳转zhuti.html
		//付费过期、体验过期、游客跳转zhuti_ad.html
		//FreePID:="888010400,888012400,888030000,888030400,888020400"
		PayPID := "888010000,888020000,888012000,888012001"
		if strings.Contains(PayPID, strconv.Itoa(pid)) {
			return contract_expert.RenderExpertHtml(ctx, "zhuti.html")
		} else {
			return contract_expert.RenderExpertHtml(ctx, "zhuti_ad.html")
		}
	} else {
		return contract_expert.RenderExpertHtml(ctx, "zhuti_ad.html")
	}
}

// 主题-首页-伪静态
func Zhuti_Static(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	currpage := 1
	pageSize := 30
	cachekey := "GetTopicList:" + strconv.Itoa(currpage) + ":" + strconv.Itoa(pageSize)
	cacheKey_total := "GetTopicList_totalCount:" + strconv.Itoa(currpage) + ":" + strconv.Itoa(pageSize)
	obj, exist := rtcache.GetCache(cachekey)
	obj_totalcount, exist := rtcache.GetCache(cacheKey_total)
	if !exist {
		expertTopicService := expertnews.NewExpertNews_TopicService()
		topiclist, totalCount, err := expertTopicService.GetTopicListPage(currpage, pageSize)
		if err != nil {
			response.RetCode = -1
			response.RetMsg = err.Error()
		} else {
			rtcache.SetCache(cachekey, topiclist)
			rtcache.SetCache(cacheKey_total, totalCount)
			obj = topiclist
			obj_totalcount = totalCount
		}

	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = obj
	response.TotalCount = obj_totalcount.(int)
	return contract_expert.RenderExpertHtmlWithPageData(ctx, "zhuti_static.html", response)
}

// 主题-详情页
func Zhuti_article(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "zhuti_article.html")
}

// 益圈圈主页
func Yqq_home(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "yqq_home.html")
}

// 益圈圈主页
func ShareArticle(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "shareArticle.html")
}

func HomeView(ctx dotweb.Context) error {
	//今日头条
	newsList, err := agent.GetTodayNews()
	if err != nil {
		ctx.ViewData().Set("jrtt", nil)
	} else {
		ctx.ViewData().Set("jrtt", newsList)
	}

	//大咖看盘
	expertNewsService := expertnews.NewExpertNews_StrategyService()
	strategyNewslist, err := expertNewsService.GetStrategyNewsList_Top6()
	if err != nil {
		ctx.ViewData().Set("dkkp", nil)
	} else {
		var newslist []*model.ExpertNews_StrategyNewsInfo
		for i, _ := range strategyNewslist {
			var news = strategyNewslist[i]
			var jsontime_str = news.NewsInfo.CreateTime.String()
			stringtime := strings.Replace(jsontime_str, "T", " ", -1)

			var newsinfo = news
			lastmodifytime, _ := time.Parse("2006-01-02 15:04:05", stringtime)
			newsinfo.NewsInfo.LastModifyTime_fomat = lastmodifytime.Format("01-02  15:04")

			summary := newsinfo.NewsInfo.Summary
			rs := []rune(summary)
			if len(rs) > 39 {
				summary = string(rs[:39]) + "..."
			}
			newsinfo.NewsInfo.Summary = summary
			newslist = append(newslist, newsinfo)
		}
		ctx.ViewData().Set("dkkp", newslist)
	}

	//主题
	expertTopicService := expertnews.NewExpertNews_TopicService()
	topiclist, _, err := expertTopicService.GetTopicListPage(1, 3)
	if err != nil {
		ctx.ViewData().Set("topic", nil)
	} else {
		var ret_topicList []*model.ExpertNews_Topic
		for i, _ := range topiclist {
			topicinfo := topiclist[i]
			stocklist := topicinfo.RelatedStockList
			topicsummary := topicinfo.TopicSummary
			rs := []rune(topicsummary)
			if len(rs) > 24 {
				topicsummary = string(rs[:24]) + "..."
			}

			if len(stocklist) > 2 {
				sort.Sort(RelationStocks(stocklist))
				topicinfo.RelatedStockList = stocklist[:2]
			}
			topicinfo.TopicSummary = topicsummary
			ret_topicList = append(ret_topicList, topicinfo)
		}
		ctx.ViewData().Set("topic", topiclist)
	}

	//专家直播
	YqqLiveService := expertnews.ConExpertNews_YqqService()
	RoomListStr, err := YqqLiveService.GetLiveRoomList()
	var YqqLive yqq.YqqRet
	err = json.Unmarshal([]byte(RoomListStr), &YqqLive)
	if err != nil {
		ctx.ViewData().Set("live", nil)
	} else {
		yqqroomdata := YqqLive.Data
		if len(yqqroomdata) > 4 {
			sort.Sort(YqqRoomDatas(yqqroomdata))
			yqqroomdata = yqqroomdata[:4]
		}
		ctx.ViewData().Set("live", yqqroomdata)
	}

	//盘后预测
	ClosingNews, err := agent.GetClosingNews()
	var phyc_newslist []*agent.ExpertNewsInfo
	for i, _ := range ClosingNews {
		var news = ClosingNews[i]
		var jsontime_str = news.PublishTime.String()
		stringtime := strings.Replace(jsontime_str, "T", " ", -1)

		var newsinfo = news
		PublishTime, _ := time.Parse("2006-01-02 15:04:05", stringtime)
		newsinfo.PublistTime_format = PublishTime.Format("01-02 15:04")

		articleSummary := news.ArticleSummary
		rs := []rune(articleSummary)
		if len(rs) > 56 {
			articleSummary = string(rs[:56]) + "..."
		}
		newsinfo.ArticleSummary = articleSummary
		phyc_newslist = append(phyc_newslist, newsinfo)
	}
	ctx.ViewData().Set("phyc", phyc_newslist)

	ctx.ViewData().Set("ServerVirtualPath", config.CurrentConfig.ServerVirtualPath)
	ctx.ViewData().Set("StaticServerHost", config.CurrentConfig.StaticServerHost)
	ctx.ViewData().Set("ResourceVersion", config.CurrentConfig.ResourceVersion)

	err = ctx.View("home.html")
	return err
}

// 专家资讯首页-输出数据
func HomeView_New(ctx dotweb.Context) error {
	ctx.ViewData().Set("jrtt", nil)
	ctx.ViewData().Set("dkkp", nil)
	ctx.ViewData().Set("topic", nil)
	ctx.ViewData().Set("live", nil)
	ctx.ViewData().Set("phyc", nil)
	ctx.ViewData().Set("dkkp0", nil)

	service := expertnews.NewSyncPageDataService()

	cachekey := "GetSyncIndexData"
	indexCacheData, exist := rtcache.GetCache(cachekey)
	if !exist {
		SyncIndexData, err := service.GetSyncIndexData()
		if err != nil {
			return ctx.View("home.html")
		}
		rtcache.SetCache(cachekey, SyncIndexData)
		indexCacheData = SyncIndexData
	}
	indexData, isOk := indexCacheData.(*model.IndexDataInfo)
	if !isOk {
		SyncIndexData, err := service.GetSyncIndexData()
		if err != nil {
			return ctx.View("home.html")
		}
		rtcache.SetCache(cachekey, SyncIndexData)
		indexCacheData = SyncIndexData
		indexData, isOk = indexCacheData.(*model.IndexDataInfo)
		if !isOk {
			return ctx.View("home.html")
		}
	}

	//今日头条
	todaynewsList := indexData.JrttList
	var retnewslist []*agent.ExpertNewsInfo
	for i, _ := range todaynewsList {
		var news = todaynewsList[i]
		var jsontime_str = news.PublishTime.String()
		stringtime := strings.Replace(jsontime_str, "T", " ", -1)

		var newsinfo = news
		lastmodifytime, _ := time.Parse("2006-01-02 15:04:05", stringtime)
		newsinfo.PublistTime_format = lastmodifytime.Format("15:04")

		retnewslist = append(retnewslist, news)
	}
	ctx.ViewData().Set("jrtt", retnewslist)

	//大咖看盘
	strategyNewslist := indexData.ClkpList
	var newslist []*model.ExpertNews_StrategyNewsInfo
	for i, _ := range strategyNewslist {
		var news = strategyNewslist[i]
		var jsontime_str = news.NewsInfo.CreateTime.String()
		stringtime := strings.Replace(jsontime_str, "T", " ", -1)

		var newsinfo = news
		lastmodifytime, _ := time.Parse("2006-01-02 15:04:05", stringtime)
		newsinfo.NewsInfo.LastModifyTime_fomat = lastmodifytime.Format("01-02  15:04")

		summary := newsinfo.NewsInfo.Summary
		rs := []rune(summary)
		if len(rs) > 39 {
			summary = string(rs[:39]) + "..."
		}
		newsinfo.NewsInfo.Summary = summary
		newslist = append(newslist, newsinfo)
	}
	ctx.ViewData().Set("dkkp", newslist)
	ctx.ViewData().Set("dkkp0", newslist[0])

	//主题
	topiclist := indexData.TopicList
	var ret_topicList []*model.ExpertNews_Topic
	for i, _ := range topiclist {
		topicinfo := topiclist[i]
		stocklist := topicinfo.RelatedStockList
		topicsummary := topicinfo.TopicSummary
		rs := []rune(topicsummary)
		if len(rs) > 24 {
			topicsummary = string(rs[:24]) + "..."
		}

		if len(stocklist) > 2 {
			sort.Sort(RelationStocks(stocklist))
			topicinfo.RelatedStockList = stocklist[:2]
		}
		topicinfo.TopicSummary = topicsummary
		ret_topicList = append(ret_topicList, topicinfo)
	}
	ctx.ViewData().Set("topic", topiclist)

	//专家直播
	yqqroomdata := indexData.ZjzbList
	ctx.ViewData().Set("live", yqqroomdata)

	//盘后预测
	ClosingNews := indexData.PhycList
	var phyc_newslist []*agent.ExpertNewsInfo
	for i, _ := range ClosingNews {
		var news = ClosingNews[i]
		var jsontime_str = news.PublishTime.String()
		stringtime := strings.Replace(jsontime_str, "T", " ", -1)

		var newsinfo = news
		PublishTime, _ := time.Parse("2006-01-02 15:04:05", stringtime)
		newsinfo.PublistTime_format = PublishTime.Format("01-02 15:04")

		articleSummary := news.ArticleSummary
		rs := []rune(articleSummary)
		if len(rs) > 56 {
			articleSummary = string(rs[:56]) + "..."
		}
		newsinfo.ArticleSummary = articleSummary
		phyc_newslist = append(phyc_newslist, newsinfo)
	}
	ctx.ViewData().Set("phyc", phyc_newslist)

	ctx.ViewData().Set("ServerVirtualPath", config.CurrentConfig.ServerVirtualPath)
	ctx.ViewData().Set("StaticServerHost", config.CurrentConfig.StaticServerHost)
	ctx.ViewData().Set("ResourceVersion", config.CurrentConfig.ResourceVersion)
	ctx.ViewData().Set("ColumnID", config.CurrentConfig.ColumnID)

	SSOURL := ctx.Request().RawQuery()
	ctx.ViewData().Set("SSOURL", SSOURL)

	return ctx.View("home.html")
}

//专家直播条件排序
type YqqRoomDatas []yqq.YqqRoom

//Len()
func (s YqqRoomDatas) Len() int {
	return len(s)
}
func (s YqqRoomDatas) Less(i, j int) bool {
	return s[i].FansNum > s[j].FansNum
}
func (s YqqRoomDatas) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}

//主题-个股排序
type RelationStocks []*model.Topic_Stock

func (s RelationStocks) Len() int {
	return len(s)
}
func (s RelationStocks) Less(i, j int) bool {
	return s[i].StockSortIndex > s[j].StockSortIndex
}
func (s RelationStocks) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}

// 专家资讯首页接口合并(今日头条、盘后预测)
func GetExpertNewsIndex(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()

	//今日头条
	todayNewsList, err := agent.GetTodayNews()
	if err != nil {
		todayNewsList = nil
	}

	//大咖看盘
	//expertNewsService := expertnews.NewExpertNews_StrategyService()
	//strategyNewslist, err := expertNewsService.GetStrategyNewsList_Top6()
	//var strategyNewslist_index []*expertnews2.ExpertNews_StrategyNewsInfo_index
	//
	//if err != nil {
	//	strategyNewslist_index = nil
	//} else {
	//	jsonstr, _ := _json.Marshal(strategyNewslist)
	//	err = _json.Unmarshal(jsonstr, &strategyNewslist_index)
	//	if err != nil {
	//		strategyNewslist_index = nil
	//	}
	//}

	//专家直播
	//zjzbList, err := gethottop4room()
	//if err != nil {
	//	zjzbList = nil
	//}

	//盘后预测
	closingNewsList, err := agent.GetClosingNews()
	if err != nil {
		closingNewsList = nil
	} else {
		if len(closingNewsList) > 2 {
			closingNewsList = closingNewsList[:2]
		}
	}

	//接口合并返回
	indexObj := struct {
		JrttList []*agent.ExpertNewsInfo
		PhycList []*agent.ExpertNewsInfo
	}{
		JrttList: todayNewsList,
		PhycList: closingNewsList}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = indexObj
	return ctx.WriteJson(response)
}

// 专家资讯首页数据接口整合-json输出数据
func GetIndexData(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()

	service := expertnews.NewSyncPageDataService()

	cachekey := "GetSyncIndexData"
	indexCacheData, exist := rtcache.GetCache(cachekey)
	if !exist {
		SyncIndexData, err := service.GetSyncIndexData()
		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		rtcache.SetCache(cachekey, SyncIndexData)
		indexCacheData = SyncIndexData
	}
	indexData, isOk := indexCacheData.(*model.IndexDataInfo)
	if !isOk {
		SyncIndexData, err := service.GetSyncIndexData()
		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		rtcache.SetCache(cachekey, SyncIndexData)
		indexCacheData = SyncIndexData
		indexData, isOk = indexCacheData.(*model.IndexDataInfo)
		if !isOk {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
	}

	if indexData != nil {
		//策略看盘字段优化
		jsonstr, _ := _json.Marshal(indexData.ClkpList)
		var strategyNewslist_index []*model.ExpertNews_StrategyNewsInfo_index
		err := _json.Unmarshal(jsonstr, &strategyNewslist_index)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		//主题字段优化
		topiclist := indexData.TopicList
		var ret_topicList []*model.ExpertNews_Topic
		for i, _ := range topiclist {
			topicinfo := topiclist[i]
			if topicinfo != nil {
				//stocklist := topicinfo.RelatedStockList
				//topicsummary := topicinfo.TopicSummary
				//rs := []rune(topicsummary)
				//if len(rs) > 24 {
				//	topicsummary = string(rs[:24]) + "..."
				//}
				//topicinfo.TopicSummary = topicsummary
				//if len(stocklist) > 2 {
				//	sort.Sort(RelationStocks(stocklist))
				//	topicinfo.RelatedStockList = stocklist[:2]
				//}
				topicinfo.RelatedBKInfo = ""
				topicinfo.RelatedStockInfo = ""
				ret_topicList = append(ret_topicList, topicinfo)
			}
		}

		//接口合并返回
		indexObj := struct {
			JrttList  []*agent.ExpertNewsInfo
			PhycList  []*agent.ExpertNewsInfo
			ClkpList  []*model.ExpertNews_StrategyNewsInfo_index
			TopicList []*model.ExpertNews_Topic
			ZjzbList  []yqq.YqqRoom
		}{
			JrttList:  indexData.JrttList,
			PhycList:  indexData.PhycList,
			ClkpList:  strategyNewslist_index,
			TopicList: ret_topicList,
			ZjzbList:  indexData.ZjzbList}

		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = indexObj
		return ctx.WriteJson(response)
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = nil
		return ctx.WriteJson(response)
	}
}

//专家策略视频
func ZhuanjiacelueVideo(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "zhuanjiacelue_video.html")
}

//专家策略文章
func ZhuanjiacelueArticle(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "zhuanjiacelue_article.html")
}
