package panwai

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/panwai"
	"strconv"
)

// Aticle 资讯详情页
func Article(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "article.html")
}

// PWGF 盘外功夫
func PWGF(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "pwgf.html")
}

// 策略详情-课堂
func Lesson(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "lesson.html")
}

// 策略详情-课堂
func SeriesLesson(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "serieslesson.html")
}

// GetBannerInfo 获取banner信息
func GetBannerInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	bannerService := service.NewBannerInfoService()
	bannerlist, err := bannerService.GetBannerInfoList()
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = bannerlist
	return ctx.WriteJson(response)
}

func GetTopNewsInfoByColumnID(ctx dotweb.Context) error {
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
		response.RetMsg = "策略编号不正确"
		return ctx.WriteJson(response)
	}

	newsService := service.NewNewsInfoService()
	newslist, err := newsService.GetTopNewsInfoByColumnID(ColumnID, StrategyID)

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
// GetNewsListByColumnID 根据栏目id获取资讯信息
func GetNewsListByColumnID(ctx dotweb.Context) error {
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
		response.RetMsg = "策略编号不正确"
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

	newsService := service.NewNewsInfoService()
	//newslist, err := newsService.GetNewsListByColumnID(ColumnID,StrategyID)
	newslist, totalCount, err := newsService.GetNewsListByColumnIDPage(ColumnID, StrategyID, currpage, pageSize)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newslist
	response.TotalCount = totalCount
	return ctx.WriteJson(response)
}

func GetSeriesNewsListByNewsID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()

	NewsIdStr := ctx.QueryString("NewsId")
	NewsId, err := strconv.Atoi(NewsIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "资讯编号不正确"
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

	newsService := service.NewNewsInfoService()
	//newslist, err := newsService.GetSeriesNewsListByNewsID(NewsId)
	newslist, totalCount, err := newsService.GetSeriesNewsListByNewsIDPage(NewsId, currpage, pageSize)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newslist
	response.TotalCount = totalCount
	return ctx.WriteJson(response)
}


// GetNewsInfoByID 查看单条资讯
func GetNewsInfoByID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
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





