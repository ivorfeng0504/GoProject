package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetBlockNewsInfomation 获取板块相关资讯
func GetBlockNewsInfomation(blockCode string, pageSize int) (result []*ExpertNewsInfo, err error) {
	if len(blockCode) == 0 {
		return nil, nil
	}
	if pageSize <= 0 {
		pageSize = 15
	} else if pageSize > 50 {
		pageSize = 50
	}
	req := contract.NewApiRequest()
	requestData := BlockNewsInfomationRequest{
		Top:       int64(pageSize),
		BlockCode: blockCode,
	}
	req.RequestData = requestData
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/getblocknewsinfomation", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

type BlockNewsInfomationRequest struct {
	BlockCode string
	//返回的数据总数
	Top int64
}
