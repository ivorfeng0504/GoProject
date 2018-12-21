package mine

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/mapper"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userproduct"
	"github.com/devfeel/dotweb"
	"strconv"
)

// Index 我的
func Index(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "mine.html")
}

// GetUserAwardList 获取用户奖品列表
func GetUserAwardList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	userAwardStateStr := ctx.PostFormValue("state")
	userAwardState, err := strconv.Atoi(userAwardStateStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "奖品状态不正确"
		return ctx.WriteJson(response)
	}
	userInfoId := contract_userhome.UserHomeUserInfo(ctx).ID
	userAwardList, err := agent.GetUserAwardList(userInfoId, userAwardState)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取用户奖品异常"
		global.InnerLogger.ErrorFormat(err, "获取用户奖品异常 state=%d", userAwardState)
		return ctx.WriteJson(response)
	}

	vmUserAwardList := mapper.MapperUserAward(userAwardList)

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = vmUserAwardList
	return ctx.WriteJson(response)
}

// GetUserProductList 获取用户产品列表
func GetUserProductList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	uid := contract_userhome.UserHomeUserInfo(ctx).UID
	if uid <= 0 {
		response.RetCode = -1
		response.RetMsg = "未获取到用户uid"
		return ctx.WriteJson(response)
	}
	userProductList, err := userproduct.GetUserProductList(uid)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取用户产品信息异常"
		global.InnerLogger.ErrorFormat(err, "获取用户产品信息异常 uid=%s", uid)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = userProductList
	return ctx.WriteJson(response)
}
