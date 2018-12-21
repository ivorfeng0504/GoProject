package home

import (
	"github.com/devfeel/dotweb"
	"emoney.cn/fundchannel/contract"
	"emoney.cn/fundchannel/protected/service"
	"time"
	"emoney.cn/fundchannel/global"
)

// Index 基金频道-首页
func Index(ctx dotweb.Context) error {
	//ctx.ViewData().Set("key", "value")
	//return ctx.View("index.html")
	return contract.RenderHtml(ctx, "index.html")
}

// Live 基金频道-直播+
func Live(ctx dotweb.Context) error {
	global.InnerLogger.DebugFormat("SSO-Page:%s", "直播+")
	return contract.RenderHtml(ctx, "live.html")
}

// Strategy 基金频道-益策略
func Strategy(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "strategy.html")
}

// GetEncryptMobile 获取加密手机号
func GetEncryptMobile(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	service := service.NewStrategyService()
	userId := ctx.QueryString("userId")
	jsonData, err := service.GetEncryptMobile(userId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = jsonData
	return ctx.WriteJson(response)
}

// QueryUserRiskInfo 查询用户等级
func QueryUserRiskInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	service := service.NewStrategyService()

	mobileNumber := ctx.QueryString("mobileNumber")

	jsonData, err := service.QueryUserRiskInfo(mobileNumber)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = jsonData
	return ctx.WriteJson(response)
}

// GetStrategyList 首页：获取基金配置列表信息
func GetStrategyList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	service := service.NewStrategyService()
	jsonData, err := service.GetStrategyList()
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = jsonData
	return ctx.WriteJson(response)
}

// GetStrategyInfoByCode 益策略：根据code + date 获取策略信息
func GetStrategyInfoByCode(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	service := service.NewStrategyService()
	code := ctx.QueryString("code")

	jsonData, err := service.GetStrategyInfoByCode(code)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = jsonData
	return ctx.WriteJson(response)
}

// GetFundTimeLineByCode 益策略：根据code 获取实时估值数据
func GetFundTimeLineByCode(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	service := service.NewStrategyService()
	code := ctx.QueryString("code")

	if len(code) == 0 {
		response.RetCode = -1
		response.RetMsg = "策略代码为空"
		return ctx.WriteJson(response)
	}

	jsonData, err := service.GetFundTimeLineByCode(code)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = jsonData
	return ctx.WriteJson(response)
}

// GetYqqContent 首页：获取直播室最新内容
func GetYqqContent(ctx dotweb.Context) error {
	service := service.NewYqqService()
	response := contract.NewResonseInfo()

	arr, err := service.GetYqqLiveLatestContent()
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = arr
	return ctx.WriteJson(response)
}

// GetLiveHeadInfo 直播+：获取直播室信息（今日话题+播主信息）
func GetLiveHeadInfo(ctx dotweb.Context) error {
	service := service.NewYqqService()
	response := contract.NewResonseInfo()

	roomInfo, err := service.GetLiveRoomInfo()
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = roomInfo
	return ctx.WriteJson(response)
}

// GetYqqLiveAllContent 直播+：获取全部直播内容
func GetYqqLiveAllContent(ctx dotweb.Context) error {
	service := service.NewYqqService()
	response := contract.NewResonseInfo()
	date := ctx.QueryString("date")
	if len(date) == 0 {
		date = time.Now().Format("2006-01-02")
	}
	msgId := ctx.QueryString("msgId")

	json, err := service.GetYqqLiveAllContent(date, msgId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = json
	return ctx.WriteJson(response)
}

// GetYqqAllQuestion 直播+：获取全部问答
func GetYqqAllQuestion(ctx dotweb.Context) error {
	service := service.NewYqqService()
	response := contract.NewResonseInfo()
	msgId := ctx.QueryString("msgId")

	json, err := service.GetYqqAllQuestion(msgId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = json
	return ctx.WriteJson(response)
}

func GetYqqMyQuestion(ctx dotweb.Context) error {
	service := service.NewYqqService()
	response, _ := service.GetYqqMyQuestion(contract.GetSSOQuery(ctx))
	return ctx.WriteJson(response)
}

// GetUserInfoByUID 获取用户信息
func GetUserInfoByUID(ctx dotweb.Context) error {
	service := service.NewStrategyService()
	response := contract.NewResonseInfo()
	userInfo := contract.UserInfo(ctx)

	if userInfo == nil {
		response.RetCode = -1
		response.RetMsg = "没有获取到用户登录信息"
		return ctx.WriteJson(response)
	}

	json, err := service.GetUserInfoByUID(userInfo.Uid)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = json
	return ctx.WriteJson(response)
}
