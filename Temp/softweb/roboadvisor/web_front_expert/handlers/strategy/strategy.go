package strategy

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"strconv"
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/service/panwai"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers"
	"git.emoney.cn/softweb/roboadvisor/contract/expert"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/util/encoding"
)

//获取runtimcache
var rtcache = handlers.NewWebRuntimeCache("strategy_",60)

// GetStrategyNewsList 根据策略ID获取资讯信息
func GetStrategyNewsList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}

	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	//取消内存缓存-会影响阅读量数据展示
	expertNewsService := expertnews.NewExpertNews_StrategyService()
	strategyNewslist, totalCount, err := expertNewsService.GetStrategyNewsListPage(ColumnID, StrategyID, currpage, pageSize)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyNewslist
	response.TotalCount = totalCount
	return ctx.WriteJson(response)
}

// GetStrategyInfoBySID 获取策略详情
func GetStrategyInfoBySID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}

	cachekey := "GetStrategyInfoBySID:"+StrategyIDStr
	obj, exist := rtcache.GetCache(cachekey)
	if !exist {
		expertNewsService := expertnews.NewExpertNews_StrategyService()
		strategyInfo, err := expertNewsService.GetStrategyInfoByStrategyID(StrategyID)

		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		rtcache.SetCache(cachekey,strategyInfo)
		obj = strategyInfo
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = obj
	return ctx.WriteJson(response)
}

// GetStrategyNewsInfoByNewsID 获取策略资讯详情
func GetStrategyNewsInfoByNewsID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	NewsIDStr := ctx.QueryString("NewsID")
	NewsID, err := strconv.Atoi(NewsIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "NewsID不正确"
		return ctx.WriteJson(response)
	}

	cachekey := "GetStrategyNewsInfoByNewsID:"+NewsIDStr
	obj, exist := rtcache.GetCache(cachekey)
	if !exist {
		expertNewsService := expertnews.NewExpertNews_StrategyService()
		strategyNewsInfo, err := expertNewsService.GetStrategyNewsInfoByNewsID(NewsID)

		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		rtcache.SetCache(cachekey,strategyNewsInfo)
		obj = strategyNewsInfo
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = obj
	return ctx.WriteJson(response)
}

// GetStrategyNewsList_Top6 首页-热门观点
func GetStrategyNewsList_Top6(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()

	cachekey := "GetStrategyNewsList_Top6"
	obj, exist := rtcache.GetCache(cachekey)
	if !exist {
		expertNewsService := expertnews.NewExpertNews_StrategyService()
		strategyNewslist, err := expertNewsService.GetStrategyNewsList_Top6()

		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}

		jsonstr, _ := _json.Marshal(strategyNewslist)
		var strategyNewslist_index []*expertnews2.ExpertNews_StrategyNewsInfo_index
		err = _json.Unmarshal(jsonstr, &strategyNewslist_index)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}

		rtcache.SetCache(cachekey, strategyNewslist_index)
		obj = strategyNewslist_index
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.TotalCount = 6
	response.Message = obj
	return ctx.WriteJson(response)
}

// GetStrategyNewsList_Rmgd 首页-热门观点
func GetStrategyNewsList_Rmgd(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	expertNewsService := expertnews.NewExpertNews_StrategyService()
	strategyNewslist, totalnum, err := expertNewsService.GetStrategyNewsList_RmgdByPage(currpage, pageSize)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	jsonstr, _ := _json.Marshal(strategyNewslist)
	var strategyNewslist_detail []*expertnews2.ExpertNews_StrategyNewsInfo_detail
	err = _json.Unmarshal(jsonstr, &strategyNewslist_detail)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.TotalCount = totalnum
	response.Message = strategyNewslist_detail
	return ctx.WriteJson(response)
}

func GetHotStrategyNewsList_Top10(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}
	cachekey := "GetHotStrategyNewsList_Top10:" + ColumnIDStr
	obj, exist := rtcache.GetCache(cachekey)
	if !exist {
		newsService := service.NewNewsInfoService()
		newsList, err := newsService.GetNewsListByClicknumFromRedis(ColumnID)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		rtcache.SetCache(cachekey, newsList)
		obj = newsList
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.TotalCount = 10
	response.Message = obj
	return ctx.WriteJson(response)
}

func GetNextStrategyNewsListByCurrNews(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}

	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	NewsIDStr := ctx.QueryString("NewsID")
	NewsID, err := strconv.Atoi(NewsIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "NewsID不正确"
		return ctx.WriteJson(response)
	}

	cachekey := "GetNextStrategyNewsListByCurrNews:" + ColumnIDStr + ":" + StrategyIDStr + ":" + currpageStr + ":" + pageSizeStr + ":" + NewsIDStr
	obj, exist := rtcache.GetCache(cachekey)
	if !exist {
		expertNewsService := expertnews.NewExpertNews_StrategyService()
		strategyNewslist, err := expertNewsService.GetNextStrategyNewsListByCurrNews(NewsID, ColumnID, StrategyID, currpage, pageSize)

		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		rtcache.SetCache(cachekey, strategyNewslist)
		obj = strategyNewslist
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.TotalCount = 0
	response.Message = obj
	return ctx.WriteJson(response)
}

// 关注某条策略
func AddFocusStrategy(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_expert.ExpertUserInfo(ctx)

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}

	UIDStr := ctx.QueryString("UID")
	UID, err := strconv.ParseInt(UIDStr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "UID不正确"
		return ctx.WriteJson(response)
	}
	LiveIdStr := ctx.QueryString("LiveId")
	LiveId, err := strconv.Atoi(LiveIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "LiveId不正确"
		return ctx.WriteJson(response)
	}

	if loginuser != nil {
		UID = loginuser.UID
	}

	expertNewsService := expertnews.NewExpertNews_FocusStrategyService()
	err = expertNewsService.AddFocusStrategy(UID, StrategyID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	//同步调用益圈圈关注接口
	//cookie获取ssourl
	ssourl, err := ctx.ReadCookieValue("expertnews.focusssourl")
	ssourl = encoding.Base64Encode([]byte(ssourl))
	err = expertNewsService.FocusLive(loginuser.Account, ssourl, 0, LiveId)

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// 取消关注某条策略
func RemoveFocusStrategy(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_expert.ExpertUserInfo(ctx)

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}
	LiveIdStr := ctx.QueryString("LiveId")
	LiveId, err := strconv.Atoi(LiveIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "LiveId不正确"
		return ctx.WriteJson(response)
	}

	UIDStr := ctx.QueryString("UID")
	UID, err := strconv.ParseInt(UIDStr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "UID不正确"
		return ctx.WriteJson(response)
	}
	if loginuser != nil {
		UID = loginuser.UID
	}

	expertNewsService := expertnews.NewExpertNews_FocusStrategyService()
	err = expertNewsService.RemoveFocusStrategy(UID, StrategyID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	//同步调用益圈圈关注接口
	//cookie获取ssourl
	ssourl, err := ctx.ReadCookieValue("expertnews.focusssourl")
	ssourl = encoding.Base64Encode([]byte(ssourl))
	if ssourl != "" {
		err = expertNewsService.FocusLive(loginuser.Account, ssourl, 1, LiveId)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// 是否关注了某条策略
func HasFocusStrategy(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser:=contract_expert.ExpertUserInfo(ctx)

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}

	UIDStr := ctx.QueryString("UID")
	UID, err := strconv.ParseInt(UIDStr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "UID不正确"
		return ctx.WriteJson(response)
	}
	if loginuser != nil {
		UID = loginuser.UID
	}

	expertNewsService := expertnews.NewExpertNews_FocusStrategyService()
	hasFocus, err := expertNewsService.HasFocusStrategy(UID, StrategyID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = hasFocus
	return ctx.WriteJson(response)
}

// 策略下的资讯列表
func GetNewsListByStrategyIDAndType(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}

	NewsTypeStr := ctx.QueryString("NewsType")
	NewsType, err := strconv.Atoi(NewsTypeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "NewsType不正确"
		return ctx.WriteJson(response)
	}

	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	expertNewsService := expertnews.NewExpertNews_StrategyService()
	strategyNewslist, totalCount, err := expertNewsService.GetNewsListByStrategyIDAndTypePage(ColumnID, StrategyID, NewsType, int64(currpage), int64(pageSize))

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	jsonstr, _ := _json.Marshal(strategyNewslist)
	var strategyNewslist_index []*expertnews2.ExpertNews_StrategyNewsInfo_index
	err = _json.Unmarshal(jsonstr, &strategyNewslist_index)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyNewslist_index
	response.TotalCount = totalCount

	return ctx.WriteJson(response)
}

// 策略下的资讯详情
func GetNewsDetailsByStrategyIDAndType(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}

	NewsTypeStr := ctx.QueryString("NewsType")
	NewsType, err := strconv.Atoi(NewsTypeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "NewsType不正确"
		return ctx.WriteJson(response)
	}

	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	expertNewsService := expertnews.NewExpertNews_StrategyService()
	strategyNewslist, totalCount, err := expertNewsService.GetNewsListByStrategyIDAndTypePage(ColumnID, StrategyID, NewsType, int64(currpage), int64(pageSize))

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	jsonstr, _ := _json.Marshal(strategyNewslist)
	var strategyNewslist_details []*expertnews2.ExpertNews_StrategyNewsInfo_detail
	err = _json.Unmarshal(jsonstr, &strategyNewslist_details)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyNewslist_details
	response.TotalCount = totalCount

	return ctx.WriteJson(response)
}

