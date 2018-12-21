package activity

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/mapper"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userloginlog"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel/userhome"
	"github.com/devfeel/dotweb"
)

// ReceiveStock 双响炮活动首页
func ReceiveStock(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "receive_stock.html")
}

// GetUserSerialLoginDay 获取用户连登天数，如果已经当前已经领了股票 则直接返回股票信息
func GetUserSerialLoginDay(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	result := viewmodel.ReceiveStockResult{}
	loginLog, err := userloginlog.GetUserLoginCountSerial(user.UID)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取连登天数异常，请稍后重试"
		global.InnerLogger.ErrorFormat(err, "获取连登天数异常 uid=%d", user.UID)
		return ctx.WriteJson(response)
	}
	result.LoginCountSerial = loginLog.LoginCountSerial
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	return ctx.WriteJson(response)
}

// GetStockList 解密今日股票
func GetStockList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	loginLog, err := userloginlog.GetUserLoginCountSerial(user.UID)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取连登天数异常，请稍后重试"
		global.InnerLogger.ErrorFormat(err, "获取连登天数异常 uid=%d", user.UID)
		return ctx.WriteJson(response)
	}
	if loginLog.LoginCountSerial < 3 {
		response.RetCode = -2
		response.RetMsg = "连登天数不足3天，您还没有不能领取股票哦"
		global.InnerLogger.ErrorFormat(err, "连登天数不足3天 uid=%d", user.UID)
		return ctx.WriteJson(response)
	}
	userStock, err := agent.InsertUserStock(user)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = "用户领取股票失败，请稍后重试"
		global.InnerLogger.ErrorFormat(err, "用户领取股票失败 userInfoId=%d", user.ID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = mapper.MapperUserReceivedStock(userStock)
	return ctx.WriteJson(response)
}

// GetUserStockHistory 获取用户领取的股票历史
func GetUserStockHistory(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	userStockList, err := agent.GetUserStockHistory(user.ID)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取用户历史双响炮数据异常，请稍后重试"
		global.InnerLogger.ErrorFormat(err, "获取用户历史双响炮数据异常 userInfoId=%d", user.ID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = mapper.MapperUserReceivedStockList(userStockList)
	return ctx.WriteJson(response)
}
