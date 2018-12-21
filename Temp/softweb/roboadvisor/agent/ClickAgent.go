package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// HandleClick 记录点击量
// clickType 点击服务名称 详见const/ClickType.go
// identity为资讯等的唯一标识
func HandleClick(clickType string, identity string,clicknum int) (count int64, err error) {
	req := contract.NewApiRequest()
	requestData := ClickData{
		ClickType: clickType,
		Identity:  identity,
		ClickNum:clicknum,
	}
	req.RequestData = requestData
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/click/handleclick", req)
	if err != nil {
		return count, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return count, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return count, nil
	}
	countF := response.Message.(float64)
	count = int64(countF)
	return count, nil
}

// AddClick 记录点击量
// clickType 点击服务名称 详见const/ClickType.go
// identity为资讯等的唯一标识
func AddClick(clickType string, identity string) (count int64, err error) {
	req := contract.NewApiRequest()
	requestData := ClickData{
		ClickType: clickType,
		Identity:  identity,
	}
	req.RequestData = requestData
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/click/addclick", req)
	if err != nil {
		return count, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return count, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return count, nil
	}
	countF := response.Message.(float64)
	count = int64(countF)
	return count, nil
}

// QueryClick 查询点击量
// clickType 点击服务名称 详见const/ClickType.go
// identitys为资讯等的唯一标识集合（一个或多个）
func QueryClick(clickType string, identitys ...string) (result *ClickResponse, err error) {
	req := contract.NewApiRequest()
	requestData := ClickData{
		ClickType: clickType,
		Identitys: identitys,
	}
	req.RequestData = requestData
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/click/queryclick", req)
	if err != nil {
		return nil, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return nil, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, err
}

//点击请求数据
type ClickData struct {
	//资讯等唯一标识符
	Identity string
	//资讯等唯一标识符集合
	Identitys []string
	//点击类型 详情见const/ClickType.go
	ClickType string
	//点击量
    ClickNum int

}

//点击查询响应数据
type ClickResponse struct {
	//点击类型 详情见const/ClickType.go
	ClickType string
	//点击查询结果 key为identity value为点击数
	Result map[string]int64
}
