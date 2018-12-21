package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	expertnews_srv "git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/mapper"
	"time"
)

// NotifyUpdateNews 通知更新新闻资讯
func NotifyUpdateNews(ctx dotweb.Context) error {
	response := contract.NewApiResponse()
	tableName := ctx.QueryString("tableName")
	if len(tableName) == 0 {
		response.RetCode = -1
		response.RetMsg = "表名不能为空"
		return ctx.WriteJson(response)
	}
	newsSrv := expertnews_srv.NewNewsInformationService()
	err := newsSrv.AddUpdateQueue(tableName)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
	}
	return ctx.WriteJson(response)
}

// GetTodayNews 获取今日头条
func GetTodayNews(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	var result []*agent.ExpertNewsInfo
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.GetTodayNewsCacheV2()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(newsList, &result)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = "Mapper出错"
			return ctx.WriteJson(response)
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetClosingNews 获取盘后预测资讯
func GetClosingNews(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	var result []*agent.ExpertNewsInfo
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.GetClosingNewsCache()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(newsList, &result)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = "Mapper出错"
			return ctx.WriteJson(response)
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetNewsInfo 获取指定日期的要闻
func GetNewsInfo(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var result []*agent.ExpertNewsInfo
	newsSrv := expertnews_srv.NewNewsInformationService()
	date := time.Now()
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &date)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "传递的时间不正确"
		return ctx.WriteJson(response)
	}
	newsList, err := newsSrv.GetNewsInfoCache(date)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	//如果查询当天的要闻 并且返回结果为空 则返回缓存中最新的要闻
	if newsList == nil || len(newsList) == 0 {
		if time.Now().Format("2006-01-02") == date.Format("2006-01-02") {
			newsList, err = newsSrv.GetNewstNewsInfoCache()
			if err != nil {
				response.RetCode = -4
				response.RetMsg = err.Error()
				return ctx.WriteJson(response)
			}
		}
	}
	if newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(newsList, &result)
		if err != nil {
			response.RetCode = -5
			response.RetMsg = "Mapper出错"
			return ctx.WriteJson(response)
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetNewsInfoTopN 获取最新的N条要闻
func GetNewsInfoTopN(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()

	var result []*agent.ExpertNewsInfo
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.GetNewsInfoTopNCache()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(newsList, &result)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = "Mapper出错"
			return ctx.WriteJson(response)
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetYaoWenHomePageData 获取热读页面数据
func GetYaoWenHomePageData(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	newsSrv := expertnews_srv.NewNewsInformationService()
	pageData, err := newsSrv.GetYaoWenHomePageData()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = pageData
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetHotNewsInfo 获取热门资讯
func GetHotNewsInfo(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	var result []*agent.ExpertNewsInfo
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.GetHotNewsInfoCache()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(newsList, &result)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = "Mapper出错"
			return ctx.WriteJson(response)
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetTopicNewsInfo 获取主题的相关资讯
func GetTopicNewsInfo(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	topicIdF, success := request.RequestData.(float64)
	if !success {
		response.RetCode = -2
		response.RetMsg = "TopicId不正确"
		return ctx.WriteJson(response)
	}
	topicId := int(topicIdF)
	if topicId <= 0 {
		response.RetCode = -3
		response.RetMsg = "TopicId不正确"
		return ctx.WriteJson(response)
	}
	var result []*agent.ExpertNewsInfo
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.GetTopicNewsInfoCache(topicId)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(newsList, &result)
		if err != nil {
			response.RetCode = -5
			response.RetMsg = "Mapper出错"
			return ctx.WriteJson(response)
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetNewsInfoDetail 获取指定的资讯信息
func GetNewsInfoDetail(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	newsInfoIdF, success := request.RequestData.(float64)
	if !success {
		response.RetCode = -2
		response.RetMsg = "资讯Id不正确"
		return ctx.WriteJson(response)
	}
	newsInfoId := int64(newsInfoIdF)
	if newsInfoId <= 0 {
		response.RetCode = -3
		response.RetMsg = "资讯Id不正确"
		return ctx.WriteJson(response)
	}
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsInfo, err := newsSrv.GetNewsInfoById(newsInfoId)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsInfo
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetBlockNewsInfomation 获取板块相关资讯
func GetBlockNewsInfomation(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.BlockNewsInfomationRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	newsSrv := expertnews_srv.NewNewsInformationService()
	newsList, err := newsSrv.GetBlockNewsInfomation(requestData.BlockCode, requestData.Top)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	var result []*agent.ExpertNewsInfo
	if newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(newsList, &result)
		if err != nil {
			response.RetCode = -4
			response.RetMsg = "Mapper出错"
			return ctx.WriteJson(response)
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
