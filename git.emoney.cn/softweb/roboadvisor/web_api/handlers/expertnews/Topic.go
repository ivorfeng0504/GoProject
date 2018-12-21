package expertnews

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
)

// GetTopicList 获取主题列表
func GetTopicList(ctx dotweb.Context) error {
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

	expertTopicService := expertnews.NewExpertNews_TopicService()
	topiclist, totalCount, err := expertTopicService.GetTopicListPage(currpage, pageSize)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = topiclist
	response.TotalCount = totalCount

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetTopicInfoByID 获取主题详情
func GetTopicInfoByID(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	TopicIDStr := ctx.QueryString("TopicID")
	TopicID, err := strconv.Atoi(TopicIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "TopicID不正确"
		return ctx.WriteJson(response)
	}

	expertTopicService := expertnews.NewExpertNews_TopicService()
	topicInfo, err := expertTopicService.GetTopicInfoByID(TopicID)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = topicInfo

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
