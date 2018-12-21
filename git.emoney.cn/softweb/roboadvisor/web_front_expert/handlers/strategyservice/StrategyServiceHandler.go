package strategyservice

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/expert"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"github.com/devfeel/dotweb"
	"strconv"
	"strings"
	"time"
)

func StrategyServiceIndex(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "celueLive.html")
}

func StrategyServiceReadingList(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "celueReading.html")
}

func StrategyServiceVideoList(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "celueVideo.html")
}

func StrategyServiceTrainingList(ctx dotweb.Context) error {
	return contract_expert.RenderExpertHtml(ctx, "celueTraining.html")
}

// GetIndexStrategyNewsV1 获取首页最新的一条策略资讯与策略视频
func GetIndexStrategyNewsV1(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyIdStr := _http.GetRequestValue(ctx, "ClientStrategyId")
	clientStrategyId, err := strconv.Atoi(clientStrategyIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略Id不正确"
		return ctx.WriteJson(response)
	}
	colID, err := agent.GetColumnIdByClientStrategyId(clientStrategyId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	result, err := agent.GetIndexStrategyNewsV1(colID)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	return ctx.WriteJson(response)
}

// GetIndexStrategyNewsV2 获取首页最新的一条策略资讯与策略视频
func GetIndexStrategyNewsV2(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyGroupIdStr := _http.GetRequestValue(ctx, "ClientStrategyGroupId")
	clientStrategyGroupId, err := strconv.Atoi(clientStrategyGroupIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略组Id不正确"
		return ctx.WriteJson(response)
	}

	result, err := agent.GetIndexStrategyNewsV2(clientStrategyGroupId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	return ctx.WriteJson(response)
}

// GetIndexStrategyNewsV3 获取首页最新的一条策略资讯与策略视频
func GetIndexStrategyNewsV3(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyIdStr := _http.GetRequestValue(ctx, "ClientStrategyId")
	clientStrategyId, err := strconv.Atoi(clientStrategyIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略Id不正确"
		return ctx.WriteJson(response)
	}
	clientStrategyGroupIdStr := _http.GetRequestValue(ctx, "ClientStrategyGroupId")
	clientStrategyGroupId, err := strconv.Atoi(clientStrategyGroupIdStr)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "策略组Id不正确"
		return ctx.WriteJson(response)
	}

	result, err := agent.GetIndexStrategyNewsV3(clientStrategyId, clientStrategyGroupId)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	return ctx.WriteJson(response)
}

// GetClientStrategyList 获取客户端所有策略信息
func GetClientStrategyList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	result, err := agent.GetClientStrategyInfoList(false)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	return ctx.WriteJson(response)
}

// GetStrategyNewsList 获取策略指定的资讯或视频 IsVideo=true则筛选出策略视频，否则删选策略资讯
func GetStrategyNewsList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	isVideoStr := _http.GetRequestValue(ctx, "IsVideo")
	isVideo, _ := strconv.ParseBool(isVideoStr)
	newsType := _const.NewsType_News //默认文章类型

	clientStrategyIdStr := _http.GetRequestValue(ctx, "ClientStrategyId")
	clientStrategyId, err := strconv.Atoi(clientStrategyIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略Id不正确"
		return ctx.WriteJson(response)
	}
	colID, err := agent.GetColumnIdByClientStrategyId(clientStrategyId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if isVideo {
		newsType = _const.NewsType_Video //视频
	}

	newsList, err := agent.GetStrategyNewsList(colID, newsType, 1, 100)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetStrategyNewsList 获取策略组指定的资讯或视频 IsVideo=true则筛选出策略组视频，否则删选策略组资讯
func GetStrategyGroupNewsList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	isVideoStr := _http.GetRequestValue(ctx, "IsVideo")
	isVideo, _ := strconv.ParseBool(isVideoStr)
	newsType := _const.NewsType_News //默认文章类型

	clientStrategyIdStr := _http.GetRequestValue(ctx, "StrategyGroupId")
	clientStrategyId, err := strconv.Atoi(clientStrategyIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略Id不正确"
		return ctx.WriteJson(response)
	}
	colID, err := agent.GetColumnIdByClientStrategyId(clientStrategyId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if isVideo {
		newsType = _const.NewsType_Video //视频
	}

	newsList, err := agent.GetStrategyNewsList(colID, newsType, 1, 100)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetStrategyNewsDetailList 批量获取指定的策略资讯或视频的详情
func GetStrategyNewsDetailList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsInfoIdListStr := _http.GetRequestValue(ctx, "NewsInfoIdList")
	if len(newsInfoIdListStr) == 0 {
		response.RetCode = -1
		response.RetMsg = "资讯Id不能为空"
		return ctx.WriteJson(response)
	}

	newsInfoIdList := strings.Split(newsInfoIdListStr, ",")
	if len(newsInfoIdList) == 0 {
		response.RetCode = -2
		response.RetMsg = "资讯Id不能为空"
		return ctx.WriteJson(response)
	}
	newsList, err := agent.GetStrategyNewsDetailList(newsInfoIdList)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetHotStrategyNews 获取热门策略文章
func GetHotStrategyNews(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyIdListStr := _http.GetRequestValue(ctx, "ClientStrategyIdList")
	if len(clientStrategyIdListStr) == 0 {
		response.RetCode = -1
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}

	clientStrategyIdList := strings.Split(clientStrategyIdListStr, ",")
	if len(clientStrategyIdList) == 0 {
		response.RetCode = -2
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}

	newsList, err := agent.GetHotStrategyNewsInfo(clientStrategyIdList, _const.NewsType_News)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetHotStrategyVideoNews 获取热门策略视频
func GetHotStrategyVideoNews(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyIdListStr := _http.GetRequestValue(ctx, "ClientStrategyIdList")
	if len(clientStrategyIdListStr) == 0 {
		response.RetCode = -1
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}

	clientStrategyIdList := strings.Split(clientStrategyIdListStr, ",")
	if len(clientStrategyIdList) == 0 {
		response.RetCode = -2
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}
	newsList, err := agent.GetHotStrategyNewsInfo(clientStrategyIdList, _const.NewsType_Video)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetStrategyRelationNewsTopN 获取策略相关最新资讯或视频
func GetStrategyRelationNewsTopN(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	isVideoStr := _http.GetRequestValue(ctx, "IsVideo")
	isVideo, _ := strconv.ParseBool(isVideoStr)
	newsType := _const.NewsType_News //默认文章类型

	clientStrategyIdStr := _http.GetRequestValue(ctx, "ClientStrategyId")
	clientStrategyId, err := strconv.Atoi(clientStrategyIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略Id不正确"
		return ctx.WriteJson(response)
	}
	colID, err := agent.GetColumnIdByClientStrategyId(clientStrategyId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if isVideo {
		newsType = _const.NewsType_Video //视频
	}

	newsList, err := agent.GetStrategyNewsList(colID, newsType, 1, 15)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetNewstStrategyNews 获取最新策略资讯
func GetNewstStrategyNews(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyIdListStr := _http.GetRequestValue(ctx, "ClientStrategyIdList")
	if len(clientStrategyIdListStr) == 0 {
		response.RetCode = -1
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}

	clientStrategyIdList := strings.Split(clientStrategyIdListStr, ",")
	if len(clientStrategyIdList) == 0 {
		response.RetCode = -2
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}
	newsList, err := agent.GetNewstStrategyNews(clientStrategyIdList, _const.NewsType_News)

	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetNewstStrategyVideoNews 获取最新策略视频
func GetNewstStrategyVideoNews(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyIdListStr := _http.GetRequestValue(ctx, "ClientStrategyIdList")
	if len(clientStrategyIdListStr) == 0 {
		response.RetCode = -1
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}

	clientStrategyIdList := strings.Split(clientStrategyIdListStr, ",")
	if len(clientStrategyIdList) == 0 {
		response.RetCode = -2
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}

	newsList, err := agent.GetNewstStrategyNews(clientStrategyIdList, _const.NewsType_Video)

	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetLiveRoomByStrategyId 根据策略Id获取直播间Id信息
func GetLiveRoomByStrategyId(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyId := _http.GetRequestValue(ctx, "ClientStrategyId")
	if len(clientStrategyId) == 0 {
		response.RetCode = -1
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}
	roomId, err := agent.GetStrategyLiveRoomId(clientStrategyId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = roomId
	return ctx.WriteJson(response)
}

// GetStrategyTagNewsList 获取策略指定的标签资讯 （专家资讯栏目）
func GetStrategyTagNewsList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}

	ClientStrategyIDStr := ctx.QueryString("ClientStrategyId")
	ClientStrategyID, err := strconv.Atoi(ClientStrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "ClientStrategyID不正确"
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

	TagIDStr := ctx.QueryString("TagID")
	TagID, err := strconv.Atoi(TagIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "TagID不正确"
		return ctx.WriteJson(response)
	}

	service := expertnews.NewExpertNews_TagService()
	strategyID, err := service.GetStrategyIDByClientStrategyID(ClientStrategyID)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	newsList, totalcount, err := service.GetStrategyTagNewsListPage(ColumnID, strategyID, TagID, int64(currpage), int64(pageSize))
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	response.TotalCount = totalcount
	return ctx.WriteJson(response)
}

// GetStrategyNewsList_MultiMedia 获取策略指定的实战培训（多媒体类型资讯）
func GetStrategyNewsList_MultiMedia(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsType := _const.NewsType_MultiMedia //多媒体资讯类型

	clientStrategyIdStr := _http.GetRequestValue(ctx, "ClientStrategyId")
	dateStr := _http.GetRequestValue(ctx, "Date")
	var date *time.Time
	dateTmp, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = nil
		err = nil
	} else {
		date = &dateTmp
	}
	clientStrategyId, err := strconv.Atoi(clientStrategyIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略Id不正确"
		return ctx.WriteJson(response)
	}
	colID, err := agent.GetColumnIdByClientStrategyId(clientStrategyId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	newsList, err := agent.GetStrategyNewsList_MultiMedia(colID, newsType, date)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

// GetStrategyNewsDetailList_MultiMedia 批量获取指定的策略多媒体资讯详情
func GetStrategyNewsDetailList_MultiMedia(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsInfoIdListStr := _http.GetRequestValue(ctx, "NewsInfoIdList")
	if len(newsInfoIdListStr) == 0 {
		response.RetCode = -1
		response.RetMsg = "资讯Id不能为空"
		return ctx.WriteJson(response)
	}

	newsInfoIdList := strings.Split(newsInfoIdListStr, ",")
	if len(newsInfoIdList) == 0 {
		response.RetCode = -2
		response.RetMsg = "资讯Id不能为空"
		return ctx.WriteJson(response)
	}
	newsList, err := agent.GetStrategyNewsDetailList_MultiMedia(newsInfoIdList)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

// GetHotStrategyMultiMediaNews 获取热门策略培训
func GetHotStrategyMultiMediaNews(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	clientStrategyIdListStr := _http.GetRequestValue(ctx, "ClientStrategyIdList")
	if len(clientStrategyIdListStr) == 0 {
		response.RetCode = -1
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}

	clientStrategyIdList := strings.Split(clientStrategyIdListStr, ",")
	if len(clientStrategyIdList) == 0 {
		response.RetCode = -2
		response.RetMsg = "策略Id不能为空"
		return ctx.WriteJson(response)
	}
	newsList, err := agent.GetHotStrategyNewsInfo(clientStrategyIdList, _const.NewsType_MultiMedia)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}

// GetNewsList_ImportantTips 获取重要提示
func GetNewsList_ImportantTips(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}

	ClientStrategyIDStr := ctx.QueryString("ClientStrategyId")
	ClientStrategyID, err := strconv.Atoi(ClientStrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "ClientStrategyID不正确"
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

	TagIDStr := ctx.QueryString("TagID")
	TagID, err := strconv.Atoi(TagIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "TagID不正确"
		return ctx.WriteJson(response)
	}

	newsList, totalcount, err := agent.GetNewsList_ImportantTips(ColumnID, ClientStrategyID, TagID, int64(currpage), int64(pageSize))

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	response.TotalCount = totalcount
	return ctx.WriteJson(response)
}

// 获取最近一个月的实战培训综合课表
func GetStrategyNewsList_MultiMedia1Month(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	dateStr := _http.GetRequestValue(ctx, "Date")
	var date *time.Time
	dateTmp, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = nil
		err = nil
	} else {
		date = &dateTmp
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
		response.RetCode = -2
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	newsList, totalCount, err := agent.GetStrategyNewsList_MultiMedia1Month(int64(currpage), int64(pageSize),date)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = "获取策略实战综合列表失败"
		response.Message = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "获取策略实战综合列表成功"
	response.Message = newsList
	response.TotalCount = totalCount
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

func GetTagList_SZPX(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	var tagList []*TagInfo
	for i := 2; i <= 4; i++ {
		tagInfo := new(TagInfo)
		tagname := ""
		tagInfo.TagID = i
		if i == 2 {
			tagname = "解读分析"
		}
		if i == 3 {
			tagname = "机会甄选"
		}
		if i == 4 {
			tagname = "回顾教育"
		}
		tagInfo.TagName = tagname

		tagList = append(tagList, tagInfo)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = tagList
	return ctx.WriteJson(response)
}

type TagInfo struct {
	TagID   int
	TagName string
}
