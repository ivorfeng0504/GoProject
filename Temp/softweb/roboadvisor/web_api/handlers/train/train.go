package train

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"strconv"
	"time"
	"git.emoney.cn/softweb/roboadvisor/protected/service/train"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// 根据标签获取课程列表
func GetTrainListByTag(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	trainTagStr := ctx.QueryString("trainTag")
	trainTag, err := strconv.Atoi(trainTagStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "培训标签不正确"
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

	pid := ctx.QueryString("pid")
	if pid == "" {
		response.RetCode = -1
		response.RetMsg = "pid不能为空"
		return ctx.WriteJson(response)
	}

	service := train.NewTrainService()
	trainList, totalcount, err := service.GetTrainListByTag(pid, trainTag, currpage, pageSize)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取培训列表失败"
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = trainList
	response.TotalCount = totalcount

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// 根据地区获取课程列表
func GetTrainListByDateAndArea(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()

	area := ctx.QueryString("area")
	if area == "" {
		response.RetCode = -1
		response.RetMsg = "地区不能为空"
		return ctx.WriteJson(response)
	}
	req_date := ctx.QueryString("date")
	if req_date == "" {
		response.RetCode = -1
		response.RetMsg = "日期不正确"
		return ctx.WriteJson(response)
	}
	date, err := time.Parse("2006-01-02 15:04:05", req_date)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "日期格式不正确"
		return ctx.WriteJson(response)
	}
	pid := ctx.QueryString("pid")
	if pid == "" {
		response.RetCode = -1
		response.RetMsg = "pid不能为空"
		return ctx.WriteJson(response)
	}

	if area == "默认" {
		area = "全国"
	}

	service := train.NewTrainService()
	trainList, err := service.GetTrainListByDateAndArea(pid, date, area)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取培训列表失败"
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = trainList

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}



// 当日最新（支持日期筛选）
func GetTrainListByDate(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.UserTrainRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	service := train.NewTrainService()
	trainList, err := service.GetTrainListByNow(requestData.PID, requestData.FilterDate)

	if err != nil {
		response.RetCode = -3
		response.RetMsg = "获取培训列表失败"
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = trainList
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)

}

// 解盘晨会和指标学习（支持日期和标签筛选）
func GetTrainListByDateAndTag(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.UserTrainRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	service := train.NewTrainService()
	trainList, totalCount, err := service.GetTrainListByDateAndPidAndTag(requestData.PID, requestData.FilterDate, requestData.TrainTag, requestData.TrainType, requestData.PageSize, requestData.CurrentPage)

	if err != nil {
		response.RetCode = -3
		response.RetMsg = "获取培训列表失败"
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = trainList
	response.TotalCount = totalCount
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)

}
