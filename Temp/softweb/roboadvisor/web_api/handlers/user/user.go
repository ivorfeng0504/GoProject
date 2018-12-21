package user

import (
	"github.com/devfeel/dotweb"
	//"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/user"
	"git.emoney.cn/softweb/roboadvisor/protected/model/user"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/contract"
)

func Reg_Mobile(ctx dotweb.Context) error {
	userService := user.NewUserService()
	mobile := ctx.QueryString("mobile")
	hardwareInfo := ctx.QueryString("hardwareInfo")
	strclientVersion := ctx.QueryString("clientVersion")
	sidstr := ctx.QueryString("sid")
	tidstr := ctx.QueryString("tid")

	tyuser := new(model.TyRegUser_Response)
	if len(mobile) == 0 {
		tyuser.RetCode = -1
		tyuser.RetMsg = "手机不能为空"
		return ctx.WriteJson(tyuser)
	}
	sid, err := strconv.Atoi(sidstr)
	if err != nil {
		tyuser.RetCode = -1
		tyuser.RetMsg = "sid不正确"
		return ctx.WriteJson(tyuser)
	}
	tid, err := strconv.Atoi(tidstr)
	if err != nil {
		tyuser.RetCode = -1
		tyuser.RetMsg = "tid不正确"
		return ctx.WriteJson(tyuser)
	}
	clientVersion, err := strconv.Atoi(strclientVersion)
	if err != nil {
		tyuser.RetCode = -1
		tyuser.RetMsg = "clientVersion不正确"
		return ctx.WriteJson(tyuser)
	}
	tyuser, _ = userService.Reg_MobileService(mobile,hardwareInfo,clientVersion, sid, tid)

	return ctx.WriteJson(tyuser)
}

func Reg_QQorWeChat(ctx dotweb.Context) error {
	userService := user.NewUserService()
	openid := ctx.QueryString("openid")
	regtypestr := ctx.QueryString("regtype")
	sidstr := ctx.QueryString("sid")
	tidstr := ctx.QueryString("tid")

	tyuser := new(model.TyRegUser_Response)
	if len(openid) == 0 {
		tyuser.RetCode = -1
		tyuser.RetMsg = "openid不能为空"
		return ctx.WriteJson(tyuser)
	}
	regtype, err := strconv.Atoi(regtypestr)
	if err != nil {
		tyuser.RetCode = -1
		tyuser.RetMsg = "regtype不正确"
		return ctx.WriteJson(tyuser)
	}
	sid, err := strconv.Atoi(sidstr)
	if err != nil {
		tyuser.RetCode = -1
		tyuser.RetMsg = "sid不正确"
		return ctx.WriteJson(tyuser)
	}
	tid, err := strconv.Atoi(tidstr)
	if err != nil {
		tyuser.RetCode = -1
		tyuser.RetMsg = "tid不正确"
		return ctx.WriteJson(tyuser)
	}
	tyuser, _ = userService.Reg_QQorWeChatService(openid, regtype, sid, tid)
	return ctx.WriteJson(tyuser)
}

func ResetPassWord(ctx dotweb.Context) error {
	userService := user.NewUserService()
	response := contract.NewResonseInfo()
	username := ctx.QueryString("username")
	password := ctx.QueryString("pwd")

	if len(username) == 0 {
		response.RetCode = -1
		response.RetMsg = "username不能为空"
		return ctx.WriteJson(response)
	}
	if len(password) == 0 {
		response.RetCode = -1
		response.RetMsg = "password不能为空"
		return ctx.WriteJson(response)
	}

	retcode, retmsg, _ := userService.Stock_ChangePasswd_ResetService(username, password)

	if retcode == 1 || retcode == 2 || retcode == 3 || retcode == 4 {
		response.RetCode = 0
		response.RetMsg = "重置密码成功"
	} else {
		response.RetCode = -1
		response.RetMsg = retmsg
	}
	return ctx.WriteJson(response)
}

// BoundGroupQryLogin 查询已绑定账号列表
func BoundGroupQryLogin(ctx dotweb.Context) error {
	userService := user.NewUserService()
	response := contract.NewResonseInfo()

	uidstr := ctx.QueryString("uid")
	uid, err := strconv.ParseInt(uidstr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "uid不正确"
		return ctx.WriteJson(response)
	}

	boundAccount, err := userService.BoundGroupQryLogin(uid)

	if err != nil {
		response.RetCode = -1
		response.RetMsg = "查询已绑定账号列表失败"
	} else {
		response.RetCode = 0
		response.RetMsg = "查询已绑定账号列表成功"
		response.Message = boundAccount
	}
	return ctx.WriteJson(response)
}

// BoundGroupAddLogin 绑定
func BoundGroupAddLogin(ctx dotweb.Context) error {
	userService := user.NewUserService()
	response := contract.NewResonseInfo()

	curUID := ctx.QueryString("curUID")
	curUserName := ctx.QueryString("curUserName")
	addUserName := ctx.QueryString("addUserName")
	addPassword := ctx.QueryString("addPassword")

	if len(curUID) == 0 {
		response.RetCode = -1
		response.RetMsg = "curUID不能为空"
		return ctx.WriteJson(response)
	}
	if len(curUserName) == 0 {
		response.RetCode = -1
		response.RetMsg = "curUserName不能为空"
		return ctx.WriteJson(response)
	}
	if len(addUserName) == 0 {
		response.RetCode = -1
		response.RetMsg = "addUserName不能为空"
		return ctx.WriteJson(response)
	}
	if len(addPassword) == 0 {
		response.RetCode = -1
		response.RetMsg = "addPassword不能为空"
		return ctx.WriteJson(response)
	}
	curUIDint, err := strconv.ParseInt(curUID, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "curUID格式不正确"
		return ctx.WriteJson(response)
	}

	boundAccuntAdd, err := userService.BoundGroupAddLogin(curUIDint, curUserName, addUserName, addPassword)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "绑定接口调用失败"
	} else {
		response.RetCode = 0
		response.RetMsg = "绑定接口调用成功"
		response.Message = boundAccuntAdd
	}
	return ctx.WriteJson(response)
}

// BoundGroupRmvLogin 解绑
func BoundGroupRmvLogin(ctx dotweb.Context) error {
	userService := user.NewUserService()
	response := contract.NewResonseInfo()

	curUserName := ctx.QueryString("curUserName")

	if len(curUserName) == 0 {
		response.RetCode = -1
		response.RetMsg = "curUserName不能为空"
		return ctx.WriteJson(response)
	}
	curUIDstr := ctx.QueryString("curUID")
	curUID, err := strconv.ParseInt(curUIDstr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "curUID不正确"
		return ctx.WriteJson(response)
	}

	rmvCIDstr := ctx.QueryString("rmvCID")
	rmvCID, err := strconv.ParseInt(rmvCIDstr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "rmvCID不正确"
		return ctx.WriteJson(response)
	}

	rmvNewPid, retCode, err := userService.BoundGroupRmvLogin(curUID,curUserName, rmvCID)

	if err != nil {
		response.RetCode = -1
		response.RetMsg = "解绑接口调用失败"
	} else {
		response.RetCode = retCode
		response.RetMsg = "解绑接口调用成功"
		response.Message = rmvNewPid
	}

	return ctx.WriteJson(response)
}

// GetLoginIDByName 获取账号CID DID等
func GetLoginIDByName(ctx dotweb.Context) error {
	userService := user.NewUserService()
	response := contract.NewResonseInfo()

	userName := ctx.QueryString("userName")
	userPasswd := ctx.QueryString("userPasswd")
	createLoginstr := ctx.QueryString("createLogin")

	if len(userName) == 0 {
		response.RetCode = -1
		response.RetMsg = "userName不能为空"
		return ctx.WriteJson(response)
	}

	createLogin, err := strconv.Atoi(createLoginstr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "createLogin不正确"
		return ctx.WriteJson(response)
	}

	if createLogin!=0 && len(userPasswd) == 0 {
		response.RetCode = -1
		response.RetMsg = "userPasswd不能为空"
		return ctx.WriteJson(response)
	}

	accountLoginIDInfo, err := userService.GetLoginIDByName(userName, userPasswd, createLogin)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取账号登录ID信息调用失败"
	} else {
		response.RetCode = 0
		response.RetMsg = "获取账号登录ID信息调用成功"
		response.Message = accountLoginIDInfo
	}
	return ctx.WriteJson(response)
}

