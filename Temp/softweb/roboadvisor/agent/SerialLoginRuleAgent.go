package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetUserSerialLoginRule 获取用户连登奖励
func GetUserSerialLoginRule(userInfo userhome_model.UserInfo) (rule *userhome_model.SerialLoginRule, err error) {
	req := contract.NewApiRequest()
	req.RequestData = userInfo
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/activity/getuserserialaward", req)
	if err != nil {
		return rule, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return rule, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &rule)
	return rule, nil
}

// GetUserLevel 获取用户等级
func GetUserLevel(userInfo userhome_model.UserInfo) (level int, err error) {
	req := contract.NewApiRequest()
	req.RequestData = userInfo
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/activity/getuserlevel", req)
	if err != nil {
		return level, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return level, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return level, nil
	}
	level = int(response.Message.(float64))
	return level, err
}
