package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetUserAwardList 获取用户的奖品列表
func GetUserAwardList(userInfoId int, state int) (userAwardList []*userhome_model.UserAward, err error) {
	req := contract.NewApiRequest()
	requestInfo := UserAwardRequest{
		UserInfoId: userInfoId,
		State:      state,
	}
	req.RequestData = requestInfo
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/activity/getuserawardlist", req)
	if err != nil {
		return userAwardList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return userAwardList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &userAwardList)
	return userAwardList, nil
}

type UserAwardRequest struct {
	UserInfoId int
	State      int
}
