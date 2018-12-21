package mine

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"time"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"strconv"
	"net/url"
	"git.emoney.cn/softweb/roboadvisor/protected/service/rsa"
	"git.emoney.cn/softweb/roboadvisor/protected/service/user"
	"math/rand"
	mobile2 "git.emoney.cn/softweb/roboadvisor/protected/service/mobile"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	//service2 "git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/captcha"
	"git.emoney.cn/softweb/roboadvisor/protected/model/user"
	"git.emoney.cn/softweb/roboadvisor/web_front_userhome/validate"
)

const(
	SendValidateCode_uniqueTag="userhome.mobile.SendValidateCode"
)

// Index 首页
func MyProfile(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "my_profile.html")
}
// Login_QQ QQ登录页
func Login_QQ(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "login_qq.html")
}
// Login_WeChat 微信登录页
func Login_WeChat(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "login_wechat.html")
}
// Bind_QQ QQ绑定确认页
func Bind_QQ(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "bind_qq.html")
}
// Bind_WeChat 微信绑定确认页
func Bind_WeChat(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "bind_wechat.html")
}

// GetValidateCode 获取验证码
func GetValidateCode(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	//user := contract.UserInfo(ctx)

	mobile := ctx.QueryString("Mobile")
	if len(mobile) == 0 {
		response.RetCode = -1
		response.RetMsg = "手机号码不能为空"
		return ctx.WriteJson(response)
	}

	mobile, err := rsa.DecryptRSA(mobile)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "手机号码解密失败"
		return ctx.WriteJson(response)
	}

	//生成随机6位验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	content := "您的验证码是" + vcode + "，请在页面中提交验证码完成验证。"
	expire := time.Now().Add(1800000000000) //纳秒（30分钟）
	nextTime := time.Now().Add(60000000000) //纳秒（1分钟）

	messageService := mobile2.NewMessageService()
	retMsg, err := messageService.SendValidateCode(mobile, content, vcode, "3", "30|999999999|186", expire, nextTime, SendValidateCode_uniqueTag, "100007")

	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if retMsg == "0" {
		response.RetCode = 0
		response.RetMsg = "发送验证码成功"
	} else {
		response.RetCode = -1
		response.RetMsg = retMsg
	}
	return ctx.WriteJson(response)
}

// Reg_TyUser 注册智投体验版用户
func Reg_TyAccount(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	mobile := ctx.QueryString("Mobile")
	if len(mobile) == 0 {
		response.RetCode = -1
		response.RetMsg = "手机号码不能为空"
		return ctx.WriteJson(response)
	}

	code := ctx.QueryString("code")
	if len(code) == 0 {
		response.RetCode = -1
		response.RetMsg = "验证码不能为空"
		return ctx.WriteJson(response)
	}

	password := ctx.QueryString("password")
	if len(password) == 0 {
		response.RetCode = -1
		response.RetMsg = "登录密码不能为空"
		return ctx.WriteJson(response)
	}

	confirmpassword := ctx.QueryString("confirmpassword")
	if len(confirmpassword) == 0 {
		response.RetCode = -1
		response.RetMsg = "确认密码不能为空"
		return ctx.WriteJson(response)
	}

	mobile, err := rsa.DecryptRSA(mobile)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "手机号码解密失败"
		return ctx.WriteJson(response)
	}
	password, err = rsa.DecryptRSA(password)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "密码解密失败"
		return ctx.WriteJson(response)
	}
	confirmpassword, err = rsa.DecryptRSA(confirmpassword)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "确认密码解密失败"
		return ctx.WriteJson(response)
	}

	if password != confirmpassword {
		response.RetCode = -1
		response.RetMsg = "两次密码输入不一致"
		return ctx.WriteJson(response)
	}

	//验证码是否正确
	messageService := mobile2.NewMessageService()
	retMsg, err := messageService.CheckValidateCode(mobile, code, SendValidateCode_uniqueTag)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if retMsg != "0" {
		//验证码错误
		response.RetCode = -1
		response.RetMsg = retMsg
		return ctx.WriteJson(response)
	}

	//调用注册接口
	userService := user.NewUserService()
	tyuser, _ := userService.Reg_MobileService(mobile, "", loginuser.ClientVersion, loginuser.SID, loginuser.TID)

	if tyuser != nil {
		response.RetCode = int(tyuser.RetCode)
		response.RetMsg = tyuser.RetMsg

		//手机号已注册
		if response.RetCode == -1 {
			return ctx.WriteJson(response)
		}
	}

	return ctx.WriteJson(response)
}

// 获取手机密码并发送短信
func GetPwdByMobile(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	mobile := ctx.QueryString("Mobile")
	if len(mobile) == 0 {
		response.RetCode = -1
		response.RetMsg = "手机号码不能为空"
		return ctx.WriteJson(response)
	}

	mobile, err := rsa.DecryptRSA(mobile)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "手机号码解密失败"
		return ctx.WriteJson(response)
	}

	//获取手机号对应的密码
	userService := user.NewTyUserService()
	//retCode, retMsg, password, err := userService.QKTaste_RegistMobile_Web(mobile, "", loginuser.SID, loginuser.TID)
	retCode, password, err := userService.GetMobilePwd_emoney(mobile, "", loginuser.SID, loginuser.TID)

	if err != nil || retCode !=0 {
		response.RetCode = -1
		response.RetMsg = "获取密码失败"
		return ctx.WriteJson(response)
	}

	vcode := password
	content := "您的验证码是" + vcode + "，请在页面中提交验证码完成验证。"
	expire := time.Now().Add(1800000000000) //纳秒（30分钟）
	nextTime := time.Now().Add(60000000000) //纳秒（1分钟）

	messageService := mobile2.NewMessageService()
	_, err = messageService.SendValidateCode(mobile, content, vcode, "3", "30|999999999|186", expire, nextTime, SendValidateCode_uniqueTag, "100007")

	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if retCode == 0 || retCode == 2 {
		response.RetCode = 0
		response.RetMsg = "发送验证码成功"
	} else {
		response.RetCode = -1
		response.RetMsg = "获取验证码失败"
	}
	return ctx.WriteJson(response)
}

// ModifyPassword 修改密码
func ModifyPassword(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	mobile := ctx.QueryString("Mobile")
	if len(mobile) == 0 {
		response.RetCode = -1
		response.RetMsg = "手机号码不能为空"
		return ctx.WriteJson(response)
	}

	code := ctx.QueryString("code")
	if len(code) == 0 {
		response.RetCode = -1
		response.RetMsg = "验证码不能为空"
		return ctx.WriteJson(response)
	}

	password := ctx.QueryString("password")
	if len(password) == 0 {
		response.RetCode = -1
		response.RetMsg = "登录密码不能为空"
		return ctx.WriteJson(response)
	}

	confirmpassword := ctx.QueryString("confirmpassword")
	if len(confirmpassword) == 0 {
		response.RetCode = -1
		response.RetMsg = "确认密码不能为空"
		return ctx.WriteJson(response)
	}
	password, err := rsa.DecryptRSA(password)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "密码解密失败"
		return ctx.WriteJson(response)
	}
	confirmpassword, err = rsa.DecryptRSA(confirmpassword)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "确认密码解密失败"
		return ctx.WriteJson(response)
	}

	if password != confirmpassword {
		response.RetCode = -1
		response.RetMsg = "两次密码输入不一致"
		return ctx.WriteJson(response)
	}
	//验证码是否正确
	messageService := mobile2.NewMessageService()
	retMsg, err := messageService.CheckValidateCode(mobile, code, SendValidateCode_uniqueTag)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if retMsg != "0" {
		//验证码错误
		response.RetCode = -1
		response.RetMsg = retMsg
		return ctx.WriteJson(response)
	}

	//调用修改密码接口
	userService := user.NewUserService()
	retcode, retmsg, err := userService.Stock_ChangePasswd_ResetService(loginuser.UserName, password)
	if retcode == 0 || retcode == 2 || retcode == 3 || retcode == 4 {
		response.RetCode = 0
		response.RetMsg = "修改密码成功"
	} else {
		response.RetCode = -1
		response.RetMsg = retmsg
	}

	return ctx.WriteJson(response)
}

// ModifyHeadportrait 修改头像
func ModifyHeadportrait(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	headimg := ctx.QueryString("headimg")
	if len(headimg) == 0 {
		response.RetCode = -1
		response.RetMsg = "头像不能为空"
		return ctx.WriteJson(response)
	}

	userService := service.NewUserInfoService()
	retCode,err:=userService.ModifyHeadportrait(loginuser.UID,headimg)

	if err != nil {
		response.RetCode = -1
		response.RetMsg = "系统异常，请稍后再试"
		return ctx.WriteJson(response)
	}
	if retCode > 0 {
		response.RetCode = 0
		response.RetMsg = "成功"
		return ctx.WriteJson(response)
	} else {
		response.RetCode = -1
		response.RetMsg = "头像修改失败"
		return ctx.WriteJson(response)
	}
}

// ModifyNickName 修改昵称
func ModifyNickName(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	nickname := ctx.QueryString("nickname")
	nickname, err := url.QueryUnescape(nickname)

	if len(nickname) == 0 {
		response.RetCode = -1
		response.RetMsg = "昵称不能为空"
		return ctx.WriteJson(response)
	}
	userService := service.NewUserInfoService()
	//isExist,err:=userService.IsExistNickname(nickname)
	//if err != nil {
	//	response.RetCode = -1
	//	response.RetMsg = "系统异常，请稍后再试"
	//	return ctx.WriteJson(response)
	//}
	//if isExist{
	//	response.RetCode = -1
	//	response.RetMsg = "昵称已存在"
	//	return ctx.WriteJson(response)
	//}

	retCode, maskNickname, err := userService.ModifyNickName(loginuser.UID, nickname)

	if err != nil {
		response.RetCode = -1
		response.RetMsg = "系统异常，请稍后再试"
		return ctx.WriteJson(response)
	}
	if retCode > 0 {
		response.RetCode = 0
		response.RetMsg = "成功"
		response.Message = maskNickname
		return ctx.WriteJson(response)
	} else {
		response.RetCode = -1
		response.RetMsg = "头像修改失败"
		return ctx.WriteJson(response)
	}
}

// BindAccountPhone 账号绑定手机号
func BindAccountPhone(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	mobile := ctx.QueryString("Mobile")
	if len(mobile) == 0 {
		response.RetCode = -1
		response.RetMsg = "手机号码不能为空"
		return ctx.WriteJson(response)
	}

	pwd := ctx.QueryString("pwd")
	if len(pwd) == 0 {
		response.RetCode = -1
		response.RetMsg = "密码不能为空"
		return ctx.WriteJson(response)
	}

	code := ctx.QueryString("code")
	if len(code) == 0 {
		response.RetCode = -1
		response.RetMsg = "验证码不能为空"
		return ctx.WriteJson(response)
	}

	captchaId:=ctx.QueryString("captchaId")
	if len(captchaId) == 0 {
		response.RetCode = -1
		response.RetMsg = "captchaId不能为空"
		return ctx.WriteJson(response)
	}

	mobile, err := rsa.DecryptRSA(mobile)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "手机号码解密失败"
		return ctx.WriteJson(response)
	}

	pwd, err = rsa.DecryptRSA(pwd)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "密码解密失败"
		return ctx.WriteJson(response)
	}

	//验证码是否正确
	if !captcha.VerifyString(captchaId, code) {
		//验证码错误
		response.RetCode = -1
		response.RetMsg = "验证码有误，请重新输入"
		return ctx.WriteJson(response)
	}

	//调用绑定手机号接口
	//code：实为手机号对应的密码
	userService := user.NewUserService()
	boundAccountAdd_response, err := userService.BoundGroupAddLogin(loginuser.GID,loginuser.UserName, mobile, pwd)

	if err != nil || boundAccountAdd_response == nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if boundAccountAdd_response.RetCode == 0 {
		response.RetCode = 0
		response.RetMsg = "绑定QQ账号成功"
		response.Message = boundAccountAdd_response
	} else {
		response.RetCode = boundAccountAdd_response.RetCode
		response.RetMsg = boundAccountAdd_response.RetMsg
		response.Message = boundAccountAdd_response
	}

	//if boundAccountAdd_response.RetCode == 0 {
	//	mobilemask := mobile[:3] + "****" + mobile[7:]
	//	mobilex, _ := mobile2.EncryptMobileHex(mobile)
	//	userHomeService := service2.NewUserInfoService()
	//	userHomeService.ModifyMobile(loginuser.UID, mobilemask, mobilex)
	//
	//	return ctx.WriteJson(boundAccountAdd_response)
	//}
	return ctx.WriteJson(boundAccountAdd_response)
}

// BindAccountEM 绑定EM号
func BindAccountEM(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)
	emcard := ctx.QueryString("emcard")
	if len(emcard) == 0 {
		response.RetCode = -1
		response.RetMsg = "em账号不能为空"
		return ctx.WriteJson(response)
	}

	pwd := ctx.QueryString("pwd")
	if len(pwd) == 0 {
		response.RetCode = -1
		response.RetMsg = "密码不能为空"
		return ctx.WriteJson(response)
	}

	code := ctx.QueryString("code")
	if len(code) == 0 {
		response.RetCode = -1
		response.RetMsg = "验证码不能为空"
		return ctx.WriteJson(response)
	}

	captchaId:=ctx.QueryString("captchaId")
	if len(captchaId) == 0 {
		response.RetCode = -1
		response.RetMsg = "captchaId不能为空"
		return ctx.WriteJson(response)
	}

	emcard, err := rsa.DecryptRSA(emcard)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "em号解密失败"
		return ctx.WriteJson(response)
	}
	pwd, err = rsa.DecryptRSA(pwd)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "密码解密失败"
		return ctx.WriteJson(response)
	}

	//图形验证码是否正确
	if !captcha.VerifyString(captchaId, code) {
		//验证码错误
		response.RetCode = -1
		response.RetMsg = "验证码有误，请重新输入"
		return ctx.WriteJson(response)
	}

	//调用绑定手机号接口
	//code：实为手机号对应的密码
	userService := user.NewUserService()
	boundAccountAdd_response, err := userService.BoundGroupAddLogin(loginuser.GID,loginuser.UserName, emcard, pwd)

	if err != nil || boundAccountAdd_response == nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if boundAccountAdd_response.RetCode == 0 {
		response.RetCode = 0
		response.RetMsg = "绑定QQ账号成功"
		response.Message = boundAccountAdd_response
	} else {
		response.RetCode = boundAccountAdd_response.RetCode
		response.RetMsg = boundAccountAdd_response.RetMsg
		response.Message = boundAccountAdd_response
	}

	return ctx.WriteJson(boundAccountAdd_response)
}

// BindAccountWeChat 绑定微信号
func BindAccountWeChat(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	uniqueID := ctx.QueryString("uniqueID")
	if len(uniqueID) == 0 {
		response.RetCode = -1
		response.RetMsg = "uniqueID不能为空"
		return ctx.WriteJson(response)
	}
	nickname := ctx.QueryString("nickname")
	if len(nickname) == 0 {
		response.RetCode = -1
		response.RetMsg = "nickname不能为空"
		return ctx.WriteJson(response)
	}
	nickname, _ = url.PathUnescape(nickname)

	uniqueID = "wx_" + uniqueID
	userService := user.NewUserService()
	boundAccountAdd_response, err := userService.BoundGroupAddLogin(loginuser.GID,loginuser.UserName, uniqueID, nickname)
	if err != nil || boundAccountAdd_response == nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if boundAccountAdd_response.RetCode == 0 {
		response.RetCode = 0
		response.RetMsg = "绑定微信账号成功"
		response.Message = boundAccountAdd_response
	} else {
		response.RetCode = boundAccountAdd_response.RetCode
		response.RetMsg = boundAccountAdd_response.RetMsg
		response.Message = boundAccountAdd_response
	}

	return ctx.WriteJson(response)
}

// BindAccountQQ 绑定QQ号
func BindAccountQQ(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	openID := ctx.QueryString("openID")
	if len(openID) == 0 {
		response.RetCode = -1
		response.RetMsg = "openID不能为空"
		return ctx.WriteJson(response)
	}
	nickname := ctx.QueryString("nickname")
	if len(nickname) == 0 {
		response.RetCode = -1
		response.RetMsg = "nickname不能为空"
		return ctx.WriteJson(response)
	}
	nickname, _ = url.PathUnescape(nickname)

	userService := user.NewUserService()
	boundAccountAdd_response, err := userService.BoundGroupAddLogin(loginuser.GID,loginuser.UserName, openID, nickname)
	if err != nil || boundAccountAdd_response == nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if boundAccountAdd_response.RetCode == 0 {
		response.RetCode = 0
		response.RetMsg = "绑定QQ账号成功"
		response.Message = boundAccountAdd_response
	} else {
		response.RetCode = boundAccountAdd_response.RetCode
		response.RetMsg = boundAccountAdd_response.RetMsg
		response.Message = boundAccountAdd_response
	}

	return ctx.WriteJson(response)
}

// 移除绑定
func RemoveBind(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	cidStr := ctx.QueryString("cid")
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "cid不正确"
		return ctx.WriteJson(response)
	}
	username:=loginuser.UserName

	userService := user.NewUserService()

	_, retCode, err := userService.BoundGroupRmvLogin(loginuser.GID, username, cid)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if retCode >= 0 {
		response.RetCode = 0
		response.RetMsg = "成功解除绑定"
		return ctx.WriteJson(response)
	} else {
		response.RetCode = -1
		response.RetMsg = "解除绑定失败"
		return ctx.WriteJson(response)
	}
}

// GetProfile 获取个人资料
func GetProfile(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)

	userInfoService := service.NewUserInfoService()
	userByredis, _ := userInfoService.GetUserInfoByUID(loginuser.UID)

	userService := user.NewUserService()
	//获取用户绑定账号列表,返回到客户端
	boundAccount_Response, err := userService.BoundGroupQryLogin(loginuser.GID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	profileInfo, err := userService.GetAccountProfile(loginuser.GID, loginuser.UserName)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if profileInfo != nil {
		userByredis.Province = profileInfo.ProvinceName
		userByredis.City = profileInfo.ProvinceName1
		userByredis.ServiceAgentName = profileInfo.ServiceAgentName
	}

	userByredis.PID = loginuser.PID
	userByredis.UserName = loginuser.UserName
	//登录账号是手机号，则获取加密后的账号返回
	if validate.IsMobile(loginuser.UserName) {
		for _, v := range boundAccount_Response {
			if v.AccountType == 1 && v.AccountName != "" {
				userByredis.UserName = v.EncryptMobile
				break
			}
		}
	}

	userByredis.BindAccountList = boundAccount_Response
	response.RetCode = 0
	response.RetMsg = "获取个人资料成功"
	response.Message = userByredis
	return ctx.WriteJson(response)
}

// 获取已绑定的手机号
func GetBindMobileByUID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	loginuser := contract_userhome.UserHomeUserInfo(ctx)
	userService := user.NewUserService()
	//获取用户绑定账号列表,返回到客户端
	boundAccount_Response, err := userService.BoundGroupQryLogin(loginuser.GID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	boundAccount_mobile := new(model.BoundAccount_Response)
	for _, v := range boundAccount_Response {
		if v.AccountType == 1 && v.AccountName != "" {
			boundAccount_mobile.AccountName = v.AccountName
			boundAccount_mobile.EncryptMobile = v.EncryptMobile
			boundAccount_mobile.AccountType = v.AccountType
			boundAccount_mobile.CustomerID = v.CustomerID
			break
		}
	}

	response.RetCode = 0
	response.RetMsg = "获取已绑定手机号成功"
	response.Message = boundAccount_mobile

	return ctx.WriteJson(response)
}


