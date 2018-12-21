package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"strconv"
)


// GetStrategyNewsInfoByNewsID 获取策略资讯详情
func GetStrategyNewsInfoByNewsID_app(ctx dotweb.Context) error {
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

// GetStrategyNewsListByPage 获取策略看盘资讯(区分文章和视频）
func GetStrategyNewsListByPage_app(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()

	newsTypeStr := ctx.QueryString("newsType")
	newsType, err := strconv.Atoi(newsTypeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "newsType不正确"
		return ctx.WriteJson(response)
	}

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyNewslist, err := expertNewsService.GetStrategyNewsList_Top1_video()
	if newsType==0{
		//读取资讯类型为文章的列表
		strategyNewslist,_, err = expertNewsService.GetStrategyNewsList_RmgdByPage(1,1000)
	}

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
	response.TotalCount = len(strategyNewslist_index)
	response.Message = strategyNewslist_index

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStrategyInfoBySID_app 获取策略详情
func GetStrategyInfoBySID_app(ctx dotweb.Context) error {
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

// GetStrategyNewsListBySID_app 获取策略下的文章资讯列表
func GetStrategyNewsListBySID_app(ctx dotweb.Context) error {
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

	newsTypeStr := ctx.QueryString("newsType")
	newsType, err := strconv.Atoi(newsTypeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "newsType不正确"
		return ctx.WriteJson(response)
	}

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyNewslist, totalCount, err := expertNewsService.GetNewsListByStrategyIDAndTypePage(ColumnID, StrategyID, newsType,int64(currpage), int64(pageSize))

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

// GetStrategyList_app 获取专家策略列表
func GetStrategyList_app(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()

	expertNewsService := expertnews2.NewExpertNews_StrategyService()
	strategyList,err:= expertNewsService.GetExpertStrategyList()
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyList
	response.TotalCount = len(strategyList)

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
