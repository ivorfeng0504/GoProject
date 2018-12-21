package click

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/click"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
)

// AddClick 记录点击量
func AddClick(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.ClickData{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if len(requestData.Identity) == 0 {
		response.RetCode = -3
		response.RetMsg = "Identity不能为空"
		return ctx.WriteJson(response)
	}
	if len(requestData.ClickType) == 0 {
		response.RetCode = -4
		response.RetMsg = "ClickType不能为空"
		return ctx.WriteJson(response)
	}
	count, err := click.AddClick(requestData.ClickType, requestData.Identity)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = count
	}
	return ctx.WriteJson(response)
}

// QueryClick 查询点击量
func QueryClick(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.ClickData{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if len(requestData.Identitys) == 0 {
		response.RetCode = -3
		response.RetMsg = "Identitys不能为空"
		return ctx.WriteJson(response)
	}
	if len(requestData.ClickType) == 0 {
		response.RetCode = -4
		response.RetMsg = "ClickType不能为空"
		return ctx.WriteJson(response)
	}
	resultData := make(map[string]int64)
	result := agent.ClickResponse{
		ClickType: requestData.ClickType,
		Result:    resultData,
	}
	for _, identity := range requestData.Identitys {
		count, err := click.GetClick(requestData.ClickType, identity)
		if err != nil {
			count = 0
		}
		resultData[identity] = count
	}
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = result
	}
	return ctx.WriteJson(response)
}

// HandleClick 手动处理点击量
func HandleClick(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.ClickData{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if len(requestData.Identity) == 0 {
		response.RetCode = -3
		response.RetMsg = "Identity不能为空"
		return ctx.WriteJson(response)
	}
	if len(requestData.ClickType) == 0 {
		response.RetCode = -4
		response.RetMsg = "ClickType不能为空"
		return ctx.WriteJson(response)
	}
	if requestData.ClickNum == 0 {
		response.RetCode = -4
		response.RetMsg = "ClickNum不能为空"
		return ctx.WriteJson(response)
	}
	count, err := click.HandleClick(requestData.ClickType, requestData.Identity, requestData.ClickNum)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = count
	}
	return ctx.WriteJson(response)
}
