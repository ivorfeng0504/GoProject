package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"strconv"
)

// ReceiveLatestLiveToStrategy 我的自选-接收直播消息
func ReceiveLatestLiveToStrategy(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	request := contract.NewApiRequest()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	requestData := new(expertnews.ReciveLive)
	err = _json.Unmarshal(_json.GetJsonString(request.RequestData), &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if requestData != nil {
		strategyService := expertnews2.NewExpertNews_StrategyService()
		err = strategyService.ReceiveUpdateMsgToStrategy(1, nil, requestData)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = "直播消息接收失败"
			return ctx.WriteJson(response)
		}
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// ReceiveLatestLiveToStrategy 我的自选-接收策略相关资讯消息
func ReceiveLatestNewsToStrategy(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	request := contract.NewApiRequest()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	requestData := new(model.NewsInfo)
	err = _json.Unmarshal(_json.GetJsonString(request.RequestData), &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	strategyService := expertnews2.NewExpertNews_StrategyService()
	err = strategyService.ReceiveUpdateMsgToStrategy(2, requestData, nil)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "资讯消息接收失败"
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// GetFocusStrategyList 获取已关注策略列表
func GetFocusStrategyList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()

	UIDStr := ctx.QueryString("UID")
	UID, err := strconv.ParseInt(UIDStr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "UID不正确"
		return ctx.WriteJson(response)
	}

	strategyService := expertnews2.NewExpertNews_FocusStrategyService()

	focusStrategyList, hasFocus, err := strategyService.GetFocusStrategyByUID(UID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取已关注策略列表失败"
		response.Message = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = strconv.FormatBool(hasFocus)
	response.Message = focusStrategyList
	return ctx.WriteJson(response)
}

// GetStrategyInfoBySID 获取策略详情
func GetStrategyInfoBySID(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "StrategyID不正确"
		return ctx.WriteJson(response)
	}

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyInfo, err := expertNewsService.GetStrategyInfoByStrategyID(StrategyID)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyInfo

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStrategyNewsList 根据策略ID获取资讯信息
func GetStrategyNewsList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()

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

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
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

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}



// GetStrategyNewsList_Top6 首页-热门观点
func GetStrategyNewsList_Top6(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyNewslist, err := expertNewsService.GetStrategyNewsList_Top6()

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	jsonstr, _ := _json.Marshal(strategyNewslist)
	var strategyNewslist_index []*expertnews.ExpertNews_StrategyNewsInfo_index
	err = _json.Unmarshal(jsonstr, &strategyNewslist_index)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.TotalCount = 6
	response.Message = strategyNewslist_index
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStrategyNewsListByPage 分页获取策略资讯
func GetStrategyNewsListByPage(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
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

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyNewslist, totalnum, err := expertNewsService.GetStrategyNewsList_RmgdByPage(currpage, pageSize)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	jsonstr, _ := _json.Marshal(strategyNewslist)
	var strategyNewslist_detail []*expertnews.ExpertNews_StrategyNewsInfo_detail
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

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}


// GetStrategyNewsInfoByNewsID 获取策略资讯详情
func GetStrategyNewsInfoByNewsID(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	NewsIDStr := ctx.QueryString("NewsID")
	NewsID, err := strconv.Atoi(NewsIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "NewsID不正确"
		return ctx.WriteJson(response)
	}
	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyNewsInfo, err := expertNewsService.GetStrategyNewsInfoByNewsID(NewsID)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyNewsInfo

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetNewsListByStrategyIDAndType 策略下的资讯列表
func GetNewsListByStrategyIDAndType(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
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

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyNewslist, totalCount, err := expertNewsService.GetNewsListByStrategyIDAndTypePage(ColumnID, StrategyID, NewsType, int64(currpage), int64(pageSize))

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	jsonstr, _ := _json.Marshal(strategyNewslist)
	var strategyNewslist_index []*expertnews.ExpertNews_StrategyNewsInfo_index
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

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetNewsDetailsByStrategyIDAndType 策略下的资讯详情
func GetNewsDetailsByStrategyIDAndType(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
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

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyNewslist, totalCount, err := expertNewsService.GetNewsListByStrategyIDAndTypePage(ColumnID, StrategyID, NewsType, int64(currpage), int64(pageSize))

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	jsonstr, _ := _json.Marshal(strategyNewslist)
	var strategyNewslist_details []*expertnews.ExpertNews_StrategyNewsInfo_detail
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

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

