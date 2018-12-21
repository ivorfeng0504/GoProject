package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetActivityAwardListByActivityId 获取活动的奖品信息
func GetActivityAwardListByActivityId(activityId int64, ignoreExpired bool) (awardList []*userhome_model.ActivityAward, err error) {
	req := contract.NewApiRequest()
	awardActivityReq := ActivityAwardRequest{
		ActivityId:    activityId,
		IgnoreExpired: ignoreExpired,
	}
	req.RequestData = awardActivityReq
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/activity/getactivityawardlist", req)
	if err != nil {
		return awardList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return awardList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &awardList)
	return awardList, nil
}

// GetGuessMarketAward 获取猜数字活动的奖品信息
func GetGuessMarketAward() (awardList []*userhome_model.ActivityAward, err error) {
	activityId := config.CurrentConfig.Activity_GuessNumber
	return GetActivityAwardListByActivityId(activityId, true)
}

// GetGuessChangeAward 获取猜涨跌奖品信息
func GetGuessChangeAward() (awardList []*userhome_model.ActivityAward, err error) {
	activityId := config.CurrentConfig.Activity_GuessChange
	return GetActivityAwardListByActivityId(activityId, true)
}

type ActivityAwardRequest struct {
	ActivityId    int64
	IgnoreExpired bool
}
