package strategyservice

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	expertnews_model "git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/service/panwai"
	"git.emoney.cn/softweb/roboadvisor/protected/service/resapi"
	strategyservice_srv "git.emoney.cn/softweb/roboadvisor/protected/service/strategyservice"
	strategyservice_vm "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"net/http"
	"sort"
)

// GetClientStrategyInfoList 获取所有可用的策略及栏目信息
func GetClientStrategyInfoList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	includeParent := request.RequestData.(bool)
	srv := strategyservice_srv.NewClientStrategyInfoService()
	result, err := srv.GetClientStrategyInfoList(includeParent)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetClientStrategyInfoListMock 获取所有可用的策略及栏目信息-测试
func GetClientStrategyInfoListMock(ctx dotweb.Context) error {
	mockJson := `{"RetCode":0,"RetMsg":"SUCCESS","Message":[{"ClientStrategyInfoId":0,"ColumnInfoId":126,"ClientStrategyId":70009,"ClientStrategyName":"突破压力","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":true,"ColumnName":"","Children":null,"ParentId":0,"ParentName":""},{"ClientStrategyInfoId":0,"ColumnInfoId":109,"ClientStrategyId":100003,"ClientStrategyName":"阶段新高","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70009,"ParentName":"突破压力"},{"ClientStrategyInfoId":0,"ColumnInfoId":120,"ClientStrategyId":100014,"ClientStrategyName":"资金博弈（金叉）","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70009,"ParentName":"突破压力"},{"ClientStrategyInfoId":0,"ColumnInfoId":121,"ClientStrategyId":100015,"ClientStrategyName":"龙腾四海（再强）","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70009,"ParentName":"突破压力"},{"ClientStrategyInfoId":0,"ColumnInfoId":123,"ClientStrategyId":100017,"ClientStrategyName":"按部就班（走强）","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70009,"ParentName":"突破压力"},{"ClientStrategyInfoId":0,"ColumnInfoId":127,"ClientStrategyId":70010,"ClientStrategyName":"基本面·策略","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":true,"ColumnName":"","Children":null,"ParentId":0,"ParentName":""},{"ClientStrategyInfoId":0,"ColumnInfoId":88,"ClientStrategyId":80019,"ClientStrategyName":"高成长","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70010,"ParentName":"基本面·策略"},{"ClientStrategyInfoId":0,"ColumnInfoId":89,"ClientStrategyId":80020,"ClientStrategyName":"高分红","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70010,"ParentName":"基本面·策略"},{"ClientStrategyInfoId":0,"ColumnInfoId":90,"ClientStrategyId":80021,"ClientStrategyName":"高盈利","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70010,"ParentName":"基本面·策略"},{"ClientStrategyInfoId":0,"ColumnInfoId":91,"ClientStrategyId":80022,"ClientStrategyName":"增持回购","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70010,"ParentName":"基本面·策略"},{"ClientStrategyInfoId":0,"ColumnInfoId":146,"ClientStrategyId":70012,"ClientStrategyName":"抄底三剑客","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":true,"ColumnName":"","Children":null,"ParentId":0,"ParentName":""},{"ClientStrategyInfoId":0,"ColumnInfoId":147,"ClientStrategyId":100018,"ClientStrategyName":"冰谷+潜龙","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70012,"ParentName":"抄底三剑客"},{"ClientStrategyInfoId":0,"ColumnInfoId":148,"ClientStrategyId":100019,"ClientStrategyName":"资金潜龙","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70012,"ParentName":"抄底三剑客"},{"ClientStrategyInfoId":0,"ColumnInfoId":149,"ClientStrategyId":100020,"ClientStrategyName":"冰谷火焰","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70012,"ParentName":"抄底三剑客"},{"ClientStrategyInfoId":0,"ColumnInfoId":150,"ClientStrategyId":100021,"ClientStrategyName":"锅底右侧","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70012,"ParentName":"抄底三剑客"},{"ClientStrategyInfoId":0,"ColumnInfoId":151,"ClientStrategyId":100022,"ClientStrategyName":"冰谷+锅底","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70012,"ParentName":"抄底三剑客"},{"ClientStrategyInfoId":0,"ColumnInfoId":152,"ClientStrategyId":100023,"ClientStrategyName":"潜龙+锅底","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70012,"ParentName":"抄底三剑客"},{"ClientStrategyInfoId":0,"ColumnInfoId":153,"ClientStrategyId":100024,"ClientStrategyName":"三剑客","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70012,"ParentName":"抄底三剑客"},{"ClientStrategyInfoId":0,"ColumnInfoId":154,"ClientStrategyId":70013,"ClientStrategyName":"长波突进","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":true,"ColumnName":"","Children":null,"ParentId":0,"ParentName":""},{"ClientStrategyInfoId":0,"ColumnInfoId":155,"ClientStrategyId":100025,"ClientStrategyName":"长波变速","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70013,"ParentName":"长波突进"},{"ClientStrategyInfoId":0,"ColumnInfoId":156,"ClientStrategyId":100026,"ClientStrategyName":"冲量变速","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70013,"ParentName":"长波突进"},{"ClientStrategyInfoId":0,"ColumnInfoId":157,"ClientStrategyId":100027,"ClientStrategyName":"龙腾长波","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70013,"ParentName":"长波突进"},{"ClientStrategyInfoId":0,"ColumnInfoId":124,"ClientStrategyId":70007,"ClientStrategyName":"拐点形态","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":true,"ColumnName":"","Children":null,"ParentId":0,"ParentName":""},{"ClientStrategyInfoId":0,"ColumnInfoId":111,"ClientStrategyId":100005,"ClientStrategyName":"底部量变","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70007,"ParentName":"拐点形态"},{"ClientStrategyInfoId":0,"ColumnInfoId":112,"ClientStrategyId":100006,"ClientStrategyName":"深跌回弹","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70007,"ParentName":"拐点形态"},{"ClientStrategyInfoId":0,"ColumnInfoId":119,"ClientStrategyId":100013,"ClientStrategyName":"趋势顶底（回转）","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70007,"ParentName":"拐点形态"},{"ClientStrategyInfoId":0,"ColumnInfoId":165,"ClientStrategyId":100028,"ClientStrategyName":"审时夺势","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70007,"ParentName":"拐点形态"},{"ClientStrategyInfoId":0,"ColumnInfoId":125,"ClientStrategyId":70008,"ClientStrategyName":"中继趋势","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":true,"ColumnName":"","Children":null,"ParentId":0,"ParentName":""},{"ClientStrategyInfoId":0,"ColumnInfoId":117,"ClientStrategyId":100011,"ClientStrategyName":"大单比率（推升）","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70008,"ParentName":"中继趋势"},{"ClientStrategyInfoId":0,"ColumnInfoId":118,"ClientStrategyId":100012,"ClientStrategyName":"超级资金（吸筹）","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70008,"ParentName":"中继趋势"},{"ClientStrategyInfoId":0,"ColumnInfoId":122,"ClientStrategyId":100016,"ClientStrategyName":"资金流变（造势）","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70008,"ParentName":"中继趋势"},{"ClientStrategyInfoId":0,"ColumnInfoId":166,"ClientStrategyId":100007,"ClientStrategyName":"小步上扬","IsDeleted":false,"CreateTime":"0001-01-01T00:00:00","IsTop":false,"ColumnName":"","Children":null,"ParentId":70008,"ParentName":"中继趋势"}],"TotalCount":0,"SystemTime":"0001-01-01T00:00:00Z","ClientCached":false,"MessageHash":""}`
	return ctx.WriteJsonBlobC(http.StatusOK, []byte(mockJson))
}

// GetStrategyLiveRoomDict 根据策略与直播间关系的字典表
func GetStrategyLiveRoomDict(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	result, err := resapi.GetStrategyLiveRoomConfigDict()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStrategyNewsList 根据栏目id和资讯类型获取资讯列表
func GetStrategyNewsList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StrategyNewsRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := strategyservice_srv.NewColumnInfoService()
	newsList, err := srv.GetStrategyNewsList(requestData.ColumnId, requestData.NewsType, requestData.PageIndex, requestData.PageSize)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStrategyNewsListByStrategyId 根据策略Id和资讯类型获取资讯列表
func GetStrategyNewsListByStrategyId(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	var requestData agent.StrategyNewsRequest
	err := agent.Bind(ctx, request, &requestData)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	clientStrategySrv := strategyservice_srv.NewClientStrategyInfoService()
	strategyDict, err := clientStrategySrv.GetClientStrategyInfoDict()
	if err != nil || len(strategyDict) == 0 {
		response.RetCode = -3
		response.RetMsg = "获取策略映射关系异常"
		return ctx.WriteJson(response)
	}

	requestData.ColumnId = strategyDict[requestData.ClientStrategyId]
	if requestData.ColumnId <= 0 {
		response.RetCode = -4
		response.RetMsg = "策略对应的栏目不存在"
		return ctx.WriteJson(response)
	}
	if requestData.PageIndex <= 0 {
		requestData.PageIndex = 1
	}
	if requestData.PageSize <= 0 {
		requestData.PageSize = 100
	}
	srv := strategyservice_srv.NewColumnInfoService()
	var newsListResult []*model.NewsInfo
	newsList, err := srv.GetStrategyNewsList(requestData.ColumnId, requestData.NewsType, requestData.PageIndex, requestData.PageSize)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if newsList != nil {
		newsListResult = append(newsListResult, newsList...)
	}
	//是否查询子集策略的新闻
	if requestData.ContainChild {
		childList, err := clientStrategySrv.GetClientStrategyInfoListByParentId(requestData.ClientStrategyId)
		if err != nil {
			response.RetCode = -6
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		//查询所有子集策略的新闻或视频
		if childList != nil && len(childList) > 0 {
			for _, childClientStrategyInfo := range childList {
				if childClientStrategyInfo == nil || childClientStrategyInfo.ClientStrategyId == requestData.ClientStrategyId || requestData.ClientStrategyId <= 0 {
					continue
				}
				childNewsList, err := srv.GetStrategyNewsList(childClientStrategyInfo.ColumnInfoId, requestData.NewsType, requestData.PageIndex, requestData.PageSize)
				if err != nil {
					response.RetCode = -7
					response.RetMsg = err.Error()
					return ctx.WriteJson(response)
				}
				if childNewsList != nil {
					newsListResult = append(newsListResult, childNewsList...)
				}
			}

		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsListResult
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStrategyNewsDetailList 根据栏目id和资讯类型批量获取资讯列表
func GetStrategyNewsDetailList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData []string
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := strategyservice_srv.NewColumnInfoService()
	newsList, err := srv.GetStrategyNewsDetailList(requestData)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetHotStrategyNewsInfo 获取热门视频或者资讯
func GetHotStrategyNewsInfo(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StrategyNewsInfoRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var newsList []*model.NewsInfo
	newsService := service.NewNewsInfoService()
	if requestData.NewsType == _const.NewsType_News {
		newsList, err = newsService.GetHotArticleListFromRedis(requestData.ClientStrategyIdList)
	}
	if requestData.NewsType == _const.NewsType_Video {
		newsList, err = newsService.GetHotVideoListFromRedis(requestData.ClientStrategyIdList)
	}
	if requestData.NewsType == _const.NewsType_MultiMedia {
		newsList, err = newsService.GetHotMultiMediaListFromRedis(requestData.ClientStrategyIdList)
	}
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	newsInfoList := strategyservice_vm.SortedByClickNumNewsInfoList(newsList)
	if newsList != nil {
		//排序并且取前N条
		sort.Sort(newsInfoList)
		if len(newsInfoList) > requestData.Count {
			newsInfoList = newsInfoList[:requestData.Count]
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsInfoList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetNewstStrategyNews 获取最新策略资讯或视频
func GetNewstStrategyNews(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StrategyNewsInfoRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := strategyservice_srv.NewColumnInfoService()
	newsList, err := srv.GetNewstStrategyNews(requestData.ClientStrategyIdList, requestData.NewsType)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	newsInfoList := strategyservice_vm.NewsInfoList(newsList)
	if newsList != nil {
		//排序并且取前N条
		sort.Sort(newsInfoList)
		if len(newsInfoList) > requestData.Count {
			newsInfoList = newsInfoList[:requestData.Count]
		}
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsInfoList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetIndexStrategyNews 获取首页最新的一条策略资讯与策略视频
func GetIndexStrategyNews(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StrategyNewsInfoRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := strategyservice_srv.NewColumnInfoService()
	result, err := srv.GetIndexStrategyNewsCache(requestData.ClientStrategyGroupId)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStrategyNewsListByStrategyId_MultiMedia 根据策略Id和资讯类型获取多媒体资讯列表
func GetStrategyNewsList_MultiMedia(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StrategyNewsRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if requestData.ColumnId <= 0 {
		response.RetCode = -4
		response.RetMsg = "策略对应的栏目不存在"
		return ctx.WriteJson(response)
	}
	srv := strategyservice_srv.NewColumnInfoService()
	var newsList []*strategyservice_vm.ExpertNews_MultiMedia_List
	if requestData.FilterDate == nil {
		newsList, err = srv.GetStrategyNewsList_MultiMedia(requestData.ColumnId, requestData.NewsType)
	} else {
		//根据日期筛选
		newsList, err = srv.GetStrategyNewsList_MultiMedia_ByDate(requestData.ColumnId, requestData.NewsType, requestData.FilterDate)
	}

	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStrategyNewsDetailList_MultiMedia 根据栏目id和资讯类型批量获取资讯列表
func GetStrategyNewsDetailList_MultiMedia(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData []string
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := strategyservice_srv.NewColumnInfoService()
	newsList, err := srv.GetStrategyNewsDetailList_MultiMedia(requestData)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetNewsList_ImportantTips 获取重要提示
func GetNewsList_ImportantTips(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StrategyNewsRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	service := expertnews.NewExpertNews_TagService()
	strategyID, err := service.GetStrategyIDByClientStrategyID(requestData.ClientStrategyId)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	var retResult []*model.NewsInfo
	if requestData.PageIndex == 1 {
		var tipsNewsList []*model.NewsInfo
		topTops_newsList, err := service.GetTopNewsList_ImportantTips(requestData.ClientStrategyId)
		if err != nil {
			response.RetCode = -3
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		tipsNewsList = append(tipsNewsList, topTops_newsList...)

		Tips_newsList, _, err := service.GetNewsListPage_ImportantTips(requestData.ClientStrategyId, 1, 4)
		if err != nil {
			response.RetCode = -3
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		tipsNewsList = append(tipsNewsList, Tips_newsList...)

		for _, v := range tipsNewsList {
			v.ToolTips = "1" //重要提示
			retResult = append(retResult, v)
		}
	}

	newsList, totalcount, err := service.GetStrategyTagNewsListPage(requestData.ColumnId, strategyID, requestData.TagID, requestData.PageIndex, requestData.PageSize)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	retResult = append(retResult, newsList...)

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = retResult
	response.TotalCount = totalcount
	return ctx.WriteJson(response)
}

// 获取最近一个月的实战培训综合课表
func GetStrategyNewsList_MultiMedia1Month(ctx dotweb.Context) error {

	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StrategyNewsRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	srv := strategyservice_srv.NewColumnInfoService()
	var newsList []*expertnews_model.ExpertNews_MultiMedia_List
	totalcount := 0
	if requestData.FilterDate == nil {
		newsList, totalcount, err = srv.GetStrategyNewsListByPage_MultiMedia1Month(requestData.PageIndex, requestData.PageSize)
	} else {
		//根据日期筛选
		newsList, totalcount, err = srv.GetStrategyNewsListByPage_MultiMedia_ByDate_Page(requestData.PageIndex, requestData.PageSize, requestData.FilterDate)
	}

	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	response.TotalCount = totalcount
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)

}
