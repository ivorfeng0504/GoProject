package panwai

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/panwai"
	"strconv"
)

// GetNewsByColIDAndStrategyID 盘外功夫
func GetNewsByColIDAndStrategyID(ctx dotweb.Context) error {
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
		response.RetMsg = "策略编号不正确"
		return ctx.WriteJson(response)
	}

	newsService := service.NewNewsInfoService()
	newslist, err := newsService.GetNewsListByColumnID(ColumnID,StrategyID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newslist
	return ctx.WriteJson(response)
}

// GetSeriesNewsListByNewsID 系列课程
func GetSeriesNewsListByNewsID(ctx dotweb.Context) error {
	response := contract.NewApiResponse()

	NewsIdStr := ctx.QueryString("NewsId")
	NewsId, err := strconv.Atoi(NewsIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "资讯编号不正确"
		return ctx.WriteJson(response)
	}

	newsService := service.NewNewsInfoService()
	newslist, err := newsService.GetSeriesNewsListByNewsID(NewsId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newslist
	return ctx.WriteJson(response)
}


// GetNewsInfoByID 查看单条资讯
func GetNewsInfoByID(ctx dotweb.Context) error {
	response := contract.NewApiResponse()
	NewsIdStr := ctx.QueryString("NewsId")
	NewsId, err := strconv.Atoi(NewsIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "资讯编号不正确"
		return ctx.WriteJson(response)
	}

	newsService := service.NewNewsInfoService()
	newsInfo, err := newsService.GetNewsInfoByID(NewsId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsInfo
	return ctx.WriteJson(response)
}