package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetUserInActivityMap 获取用户参与活动Map
func GetUserInActivityMap(userInfoId int) (userActivityMap map[int64]bool, err error) {
	req := contract.NewApiRequest()
	req.RequestData = userInfoId
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/activity/getuserinactivitymap", req)
	if err != nil {
		return userActivityMap, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return userActivityMap, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &userActivityMap)
	return userActivityMap, nil
}
