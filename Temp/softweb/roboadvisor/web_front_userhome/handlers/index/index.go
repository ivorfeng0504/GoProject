package index

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/panwai"
	user2 "git.emoney.cn/softweb/roboadvisor/protected/service/user"
	service2 "git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"github.com/devfeel/dotweb"
	"strconv"
)

// Index 首页
func Index(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "index.html")
}

func Rsatest(ctx dotweb.Context) error {
	//_, err := mystocksyn.CopyAndWrite("1000664674", "2011647269", "2011647269")
	//loginuser := contract_userhome.UserHomeUserInfo(ctx)
	userService := user2.NewUserService()
	gid, err := userService.CidFindGid(1018618276)
	//softdataService := user2.NewSoftDataSyncService()
	//retCode, err := softdataService.SoftUserDataSync(2011647269, 1000478450)
	response := contract.NewResonseInfo()
	if err != nil {
		response.RetCode = -1
		response.Message = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = gid
	return ctx.WriteJson(response)
}

// GetUserInfo 获取登录用户信息
func GetUserInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)
	emcard := ""
	enddate := ""

	//头像昵称更新
	userInfoService := service2.NewUserInfoService()
	userByredis, _ := userInfoService.GetUserInfoByUID(loginuser.UID)

	loginuser.Headportrait = userByredis.Headportrait
	loginuser.NickName = userByredis.NickName
	userService := user2.NewUserService()

	//查询已绑定EM号，如未绑定em号则不显示到期日
	boundAccount_Response, _ := userService.BoundGroupQryLogin(loginuser.GID)
	for _, v := range boundAccount_Response {
		if v.AccountType == 0 && v.AccountName != "" {
			emcard = v.AccountName
			break
		}
	}

	enddate, _ = userService.GetEndDate(emcard)
	loginuser.EndDate = enddate
	loginuser.UserName = ""
	//用户等级
	loginuser.UserLevel, _ = agent.GetUserLevel(*loginuser)
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = loginuser
	return ctx.WriteJson(response)
}

// GetTopNewsInfoByColumnID 根据栏目id获取置顶资讯
func GetTopNewsInfoByColumnID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}

	PID := ctx.QueryString("PID")
	if PID == "" {
		response.RetCode = -1
		response.RetMsg = "PID不能为空"
		return ctx.WriteJson(response)
	}

	newsService := service.NewNewsInfoService()
	newslist, err := newsService.GetTopNewsInfoByColumnID_userHome(ColumnID, PID)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newslist
	return ctx.WriteJson(response)
}

// GetNewsListByColumnID 根据栏目id获取资讯信息
func GetNewsListByColumnID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	ColumnIDStr := ctx.QueryString("ColumnID")
	ColumnID, err := strconv.Atoi(ColumnIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "栏目编号不正确"
		return ctx.WriteJson(response)
	}

	PID := ctx.QueryString("PID")
	if PID == "" {
		response.RetCode = -1
		response.RetMsg = "PID不能为空"
		return ctx.WriteJson(response)
	}

	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	newsService := service.NewNewsInfoService()
	newslist, totalCount, err := newsService.GetNewsListByColumnIDPage_userhome(ColumnID, PID, currpage, pageSize)

	//无可用页码
	if (currpage - 1) > totalCount/pageSize {
		response.RetCode = -1
		response.RetMsg = "无可用页码"
		return ctx.WriteJson(response)
	}

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newslist
	response.TotalCount = totalCount
	return ctx.WriteJson(response)
}

// GetNewsInfoByID 查看单条资讯
func GetNewsInfoByID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	NewsIdStr := ctx.QueryString("NewsId")
	NewsId, err := strconv.Atoi(NewsIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "资讯编号不正确"
		return ctx.WriteJson(response)
	}

	newsService := service.NewNewsInfoService()
	newsInfo, err := newsService.GetNewsInfoByID(NewsId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsInfo
	return ctx.WriteJson(response)
}
