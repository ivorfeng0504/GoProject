package activity

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/mapper"
	"github.com/devfeel/dotweb"
	"strconv"
)

// Index 活动首页
func Index(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "activity.html")
}

// GetActivityList 获取活动列表
func GetActivityList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	activityStateStr := ctx.PostFormValue("state")
	activityState, err := strconv.Atoi(activityStateStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "活动状态不正确"
		return ctx.WriteJson(response)
	}
	activityList, err := agent.GetActivityList(activityState)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取活动异常"
		global.InnerLogger.ErrorFormat(err, "获取活动异常 state=%d", activityState)
		return ctx.WriteJson(response)
	}
	activityList = agent.FilterEnable(activityList)
	userInfo := contract_userhome.UserHomeUserInfo(ctx)
	userInfoId := userInfo.ID
	activityList = agent.FilterPID(activityList, userInfo.PID)
	var userActivityMap map[int64]bool
	//往期活动才查询用户是否已经参与
	if activityState == _const.ActivityState_Finish {
		userActivityMap, err = agent.GetUserInActivityMap(userInfoId)
		if err != nil {
			response.RetCode = -3
			response.RetMsg = "获取用户活动参与信息出错"
			global.InnerLogger.ErrorFormat(err, "获取用户活动参与信息出错 state=%d userInfoId=", activityState, userInfoId)
			return ctx.WriteJson(response)
		}
	} else {
		userActivityMap = nil
	}

	vmActivityList := mapper.MapperActivity(activityList, userActivityMap)
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = vmActivityList
	return ctx.WriteJson(response)
}

// GetSerialLoginRule 获取当前用户的连登奖励
func GetSerialLoginRule(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	rule, err := agent.GetUserSerialLoginRule(*user)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取连登奖励异常"
		global.InnerLogger.ErrorFormat(err, "获取连登奖励异常 uid=%d", user.UID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = rule
	return ctx.WriteJson(response)
}

// GetUserMedal 获取用户勋章
func GetUserMedal(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	userMedal, err := agent.GetUserMedal(*user)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取用户勋章异常"
		global.InnerLogger.ErrorFormat(err, "获取用户勋章异常 userInfoId=%d", user.ID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = userMedal
	return ctx.WriteJson(response)
}
