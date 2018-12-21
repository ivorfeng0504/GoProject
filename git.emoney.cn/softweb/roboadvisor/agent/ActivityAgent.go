package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"strconv"
	"strings"
)

// GetActivityList 根据活动状态获取活动列表
func GetActivityList(state int) (activityList []*userhome_model.Activity, err error) {
	req := contract.NewApiRequest()
	req.RequestData = state
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/activity/getactivitylist", req)
	if err != nil {
		return activityList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return activityList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &activityList)
	return activityList, nil
}

// FilterEnable 过滤未启用的活动
func FilterEnable(src []*userhome_model.Activity) (activityList []*userhome_model.Activity) {
	if src == nil {
		return nil
	}
	for _, activity := range src {
		if activity.IsEnabled {
			activityList = append(activityList, activity)
		}
	}
	return activityList
}

// FilterPIDType 过滤活动权限
// Deprecated: 弃用 请改用FilterPID
func FilterPIDType(src []*userhome_model.Activity, pidType int) (activityList []*userhome_model.Activity) {
	if src == nil {
		return nil
	}
	for _, activity := range src {
		//如果是0 则不过滤任何版本
		if strings.Contains(activity.TargetUserType, "0") || strings.Contains(activity.TargetUserType, strconv.Itoa(pidType)) {
			activityList = append(activityList, activity)
		}
	}
	return activityList
}

// FilterPID 过滤活动权限，通过PID
func FilterPID(src []*userhome_model.Activity, pid int) (activityList []*userhome_model.Activity) {
	if src == nil {
		return nil
	}
	for _, activity := range src {
		//如果是0 则不过滤任何版本
		if len(activity.TargetUserType) > 0 {
			pidList := strings.Split(activity.TargetUserType, ",")
			for _, pidValue := range pidList {
				if len(pidValue) > 0 && (pidValue == strconv.Itoa(pid) || pidValue == "0") {
					activityList = append(activityList, activity)
				}
			}
		}
	}
	return activityList
}
