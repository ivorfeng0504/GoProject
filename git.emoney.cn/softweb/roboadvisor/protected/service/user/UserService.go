package user

import (
	"encoding/hex"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/user"
	"git.emoney.cn/softweb/roboadvisor/protected/repository/user"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/mobile"
	"git.emoney.cn/softweb/roboadvisor/protected/service/mystocksyn"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/web_front_userhome/validate"
	"github.com/devfeel/dotlog"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	service.BaseService
	userRepo *user.UserRepository
}

var (
	// shareUserRepo UserRepository
	shareUserRepo *user.UserRepository

	// shareUserLogger 共享的Logger实例
	shareUserLogger dotlog.Logger
)

const (
	userServiceName = "UserServiceLogger"
	TJSZTXkey_pre   = "Emoney.LearnStock.TouJiaoSZTXStatus."
)

func NewUserService() *UserService {
	userService := &UserService{
		userRepo: shareUserRepo,
	}
	userService.RedisCache = _cache.GetRedisCacheProvider(protected.TJSZTXRedisConfig)
	return userService
}

// Reg_MobileService 手机注册
func (service *UserService) Reg_MobileService(mobile string,hardwareInfo string,clientVersion int, sid int, tid int) (*model.TyRegUser_Response, error) {
	var mapRet []map[string]interface{}
	mapRet, err := service.userRepo.Mobile_Reg_For_ZT(mobile, hardwareInfo, clientVersion, sid, tid)

	shareUserLogger.InfoFormat("[手机注册] mobile:%s 注册手机号返回：", mobile, mapRet)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "注册手机号接口失败", mapRet)
		return nil, err
	}

	if len(mapRet) > 0 {
		retCode := mapRet[0]["retCode"].(int64)
		retMsg := mapRet[0]["retMsg"].(string)
		var customerid int64
		password := ""
		if retCode != -2 {
			customerid = mapRet[0]["customerID"].(int64)
			password = mapRet[0]["userPasswd"].(string)
		}

		tyuser := new(model.TyRegUser_Response)
		tyuser.RetCode = retCode
		tyuser.RetMsg = retMsg
		tyuser.CustomerID = customerid
		tyuser.UserPasswd = password

		return tyuser, nil
	}

	return nil, err
}

// Reg_QQorWeChatService 微信、QQ注册
func (service *UserService) Reg_QQorWeChatService(openid string, rettype int, sid int, tid int) (*model.TyRegUser_Response, error) {
	var mapRet []map[string]interface{}
	mapRet, err := service.userRepo.QqAndWechat_Reg_For_ZT(openid, rettype, sid, tid)

	shareUserLogger.InfoFormat("[微信QQ注册] openid:%s Reg_QQorWeChatService注册返回：", openid, mapRet)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "第三方账号注册接口失败", mapRet)
		return nil, err
	}

	if len(mapRet) > 0 {
		retCode := mapRet[0]["retCode"].(int64)
		retMsg := mapRet[0]["retMsg"].(string)
		var customerid int64
		password := ""
		if retCode != -2 {
			customerid = mapRet[0]["customerID"].(int64)
			password = mapRet[0]["userPasswd"].(string)
		}

		tyuser := new(model.TyRegUser_Response)
		tyuser.RetCode = retCode
		tyuser.RetMsg = retMsg
		tyuser.CustomerID = customerid
		tyuser.UserPasswd = password

		return tyuser, nil
	}

	return nil, err

}

// Stock_ChangePasswd_ResetService 重置密码
func (service *UserService) Stock_ChangePasswd_ResetService(username string, newpwd string) (int64, string, error) {
	var mapRet []map[string]interface{}
	mapRet, err := service.userRepo.Stock_ChangePasswd_Reset_For_ZT(username, newpwd)

	shareUserLogger.InfoFormat("[重置密码] username:%s newpwd:%s Stock_ChangePasswd_Reset_result返回：", username, newpwd, mapRet)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "重置密码接口失败", mapRet)
		return -1, "密码修改失败", err
	}

	var retCode int64
	var retMsg string

	if len(mapRet) > 0 {
		if mapRet[0]["retCode"] != nil {
			retCode = mapRet[0]["retCode"].(int64)
		}
		if mapRet[0]["retMsg"] != nil {
			retMsg = mapRet[0]["retMsg"].(string)
		}
	}
	return retCode, retMsg, err
}

// BoundGroupQryLogin 查询已绑定账号列表
func (service *UserService) BoundGroupQryLogin(gid int64) ([]*model.BoundAccount_Response, error) {
	var boundAccountList []*model.BoundAccount_Response
	var mapRet []map[string]interface{}
	mapRet, err := service.userRepo.BoundGroupQryLogin(gid)

	shareUserLogger.InfoFormat("[查询已绑定列表] uid:%d BoundGroupQryLogin返回：%s", gid, _json.GetJsonString(mapRet))
	if err != nil {
		shareUserLogger.ErrorFormat(err, "查询已绑定账号列表接口失败", mapRet)
		return nil, err
	}
	if len(mapRet) > 0 {
		for i, _ := range mapRet {
			boundAccount := new(model.BoundAccount_Response)
			if mapRet[i]["ShowName"] != nil {
				boundAccount.AccountName = mapRet[i]["ShowName"].(string)
			}
			if mapRet[i]["CustomerID"] != nil {
				boundAccount.CustomerID = mapRet[i]["CustomerID"].(int64)
			}
			if mapRet[i]["UserType"] != nil {
				boundAccount.AccountType = mapRet[i]["UserType"].(int64)
			}
			if mapRet[i]["mobile"] != nil {
				var mobilebyte = mapRet[i]["mobile"].([]byte)
				mobileHex := hex.EncodeToString(mobilebyte)
				mobileHex = "0x" + strings.ToUpper(mobileHex)

				boundAccount.EncryptMobile = mobileHex
			}

			boundAccountList = append(boundAccountList, boundAccount)
		}
	}

	return boundAccountList, err
}

// BoundGroupAddLogin 添加绑定
func (service *UserService) BoundGroupAddLogin(curGID int64, curUserName string, addUserName string, addPassword string) (*model.BoundAccountAdd_Response, error) {
	var boundAccountAdd = new(model.BoundAccountAdd_Response)
	var mapRet []map[string]interface{}
	mapRet, err := service.userRepo.BoundGroupAddLogin(curUserName, addUserName, addPassword)

	shareUserLogger.InfoFormat("[添加绑定] curGID:%d curUserName:%s addUserName:%s  addPassword:%s  BoundGroupAddLogin_ResultSet返回：", curGID, curUserName, addUserName, addPassword, mapRet)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "添加绑定接口失败", mapRet)
		return nil, err
	}
	if len(mapRet) > 0 {
		for i, _ := range mapRet {
			if mapRet[i]["addInCID"] != nil {
				boundAccountAdd.AddInCID = mapRet[i]["addInCID"].(int64)
			}
			if mapRet[i]["addInShowName"] != nil {
				boundAccountAdd.AddInShowName = mapRet[i]["addInShowName"].(string)
			}
			if mapRet[i]["addInType"] != nil {
				boundAccountAdd.AddInType = int(mapRet[i]["addInType"].(int64))
			}
			if mapRet[i]["newCurDID"] != nil {
				boundAccountAdd.NewCurDID = mapRet[i]["newCurDID"].(int64)
			}
			if mapRet[i]["addInOldPID"] != nil {
				boundAccountAdd.AddInOldPID = mapRet[i]["addInOldPID"].(int64)
			}
			if mapRet[i]["retMsg"] != nil {
				boundAccountAdd.RetMsg = mapRet[i]["retMsg"].(string)
			}
			if mapRet[i]["retCode"] != nil {
				boundAccountAdd.RetCode = int(mapRet[i]["retCode"].(int64))
			}
		}
	}

	//boundAccountAdd.NewCurDID!=0 && boundAccountAdd.AddInOldPID!=0
	if boundAccountAdd != nil && boundAccountAdd.RetCode == 0 {
		//一个帐号，只有当初始DID=PID，然后添加别的帐号的时候，才会取新的DID，之后就不会再取（未取到则不调用云同步）
		if boundAccountAdd.NewCurDID != 0 {
			//自选股云同步
			retFlag, err := mystocksyn.CopyAndWrite(strconv.FormatInt(curGID, 10), strconv.FormatInt(boundAccountAdd.NewCurDID, 10), strconv.FormatInt(boundAccountAdd.AddInOldPID, 10))
			shareUserLogger.InfoFormat("添加绑定-自选股云同步开始 uidFrom:%d  uidTo:%d  uidOld:%d  返回结果:%s", curGID, boundAccountAdd.NewCurDID, boundAccountAdd.AddInOldPID, retFlag)
			if err != nil {
				shareUserLogger.ErrorFormat(err, "添加绑定-自选股云同步失败 uidFrom:%d  uidTo:%d  uidOld:%d", curGID, boundAccountAdd.NewCurDID, boundAccountAdd.AddInOldPID)
			}
			//客户端指标数据等同步
			softdataService := NewSoftDataSyncService()
			retCode, err := softdataService.SoftUserDataSync(boundAccountAdd.NewCurDID, curGID)
			shareUserLogger.InfoFormat("添加绑定-客户端指标数据 uidFrom:%d  uidTo:%d 返回结果:%d", curGID, boundAccountAdd.NewCurDID, retCode)
			if err != nil {
				shareUserLogger.ErrorFormat(err, "添加绑定-客户端指标数据等同步失败 uidFrom:%d  uidTo:%d", curGID, boundAccountAdd.NewCurDID)
			}
		}

		//绑定手机号变更-通知投教实战体系
		if (!validate.IsMobile(curUserName) && len(curUserName) < 20) || (!validate.IsMobile(addUserName) && len(addUserName) < 20) {
			service.NoticeRedis_TouJiaoSZTX(curUserName)
		}
	}

	return boundAccountAdd, err
}

// BoundGroupRmvLogin 移除绑定
func (service *UserService) BoundGroupRmvLogin(curGID int64, curUserName string, rmvCID int64) (rmvNewPID int64, retCode int, err error) {
	var mapRet []map[string]interface{}
	mapRet, err = service.userRepo.BoundGroupRmvLogin(curUserName, rmvCID)

	shareUserLogger.InfoFormat("[移除绑定] curGID:%d curUserName:%s rmvCID:%d BoundGroupRmvLogin_ResultSet返回：", curGID, curUserName, rmvCID, mapRet)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "移除绑定接口失败", mapRet)
		return -1, -1, err
	}
	if len(mapRet) > 0 {
		for i, _ := range mapRet {
			if mapRet[i]["rmvNewPID"] != nil {
				rmvNewPID = mapRet[i]["rmvNewPID"].(int64)
			}
			if mapRet[i]["retMsg"] != nil {
				_ = mapRet[i]["retMsg"].(string)
			}
			if mapRet[i]["returnNo"] != nil {
				retCode = int(mapRet[i]["returnNo"].(int64))
			}
		}
		//一个帐号，只有当初始DID=PID，然后添加别的帐号的时候，才会取新的DID，之后就不会再取（未取到则不调用云同步）
		if retCode >= 0 && rmvNewPID != 0 {
			//自选股云同步
			retFlag, err := mystocksyn.CopyAndWrite(strconv.FormatInt(curGID, 10), strconv.FormatInt(rmvNewPID, 10), strconv.FormatInt(rmvCID, 10))
			shareUserLogger.InfoFormat("添加绑定-自选股云同步开始 uidFrom:%d  uidTo:%d  uidOld:%d  返回结果:%s", curGID, rmvNewPID, rmvCID, retFlag)
			if err != nil {
				shareUserLogger.ErrorFormat(err, "移除绑定-自选股云同步失败 uidFrom:%d  uidTo:%d  uidOld:%d", curGID, rmvNewPID, rmvCID)
			}
			//客户端指标数据等同步
			softdataService := NewSoftDataSyncService()
			retCode, err = softdataService.SoftUserDataSync(rmvNewPID, curGID)
			shareUserLogger.InfoFormat("添加绑定-客户端指标数据 uidFrom:%d  uidTo:%d 返回结果:%d", curGID, rmvNewPID, retCode)
			if err != nil {
				shareUserLogger.ErrorFormat(err, "移除绑定-客户端指标数据等同步失败 uidFrom:%d  uidTo:%d", curGID, rmvNewPID)
			}

			//绑定手机号变更-通知投教实战体系
			if !validate.IsMobile(curUserName) && len(curUserName) < 20 {
				shareUserLogger.InfoFormat("绑定手机号变更-通知投教实战体系 emcard:%s ", curUserName)

				service.NoticeRedis_TouJiaoSZTX(curUserName)
			}
		}
	}

	return rmvNewPID, retCode, err
}

// GetLoginIDByName 根据账号获取对应的CID DID等信息
func (service *UserService) GetLoginIDByName(userName string, userPasswd string, createLogin int) (*model.AccountLoginIDInfo_Response, error) {
	var loginIDInfo = new(model.AccountLoginIDInfo_Response)
	var mapRet []map[string]interface{}
	mapRet, err := service.userRepo.GetLoginIDByName(userName, userPasswd, createLogin)

	shareUserLogger.InfoFormat("userName:%s userPasswd:%s createLogin:%d GetLoginIDByName_ResultSet返回：", userName, userPasswd, createLogin, mapRet)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "获取LoginID接口失败", mapRet)
		return nil, err
	}

	if len(mapRet) > 0 {
		for i, _ := range mapRet {
			if mapRet[i]["userType"] != nil {
				loginIDInfo.UserType = mapRet[i]["userType"].(int64)
			}
			if mapRet[i]["showName"] != nil {
				loginIDInfo.ShowName = mapRet[i]["showName"].(string)
			}
			//if mapRet[i]["guidCP"] != nil {
			//	loginIDInfo.GuidCP = mapRet[i]["guidCP"].(string)
			//}
			if mapRet[i]["uniqueID"] != nil {
				loginIDInfo.UniqueID = mapRet[i]["uniqueID"].(int64)
			}
			if mapRet[i]["CID"] != nil {
				loginIDInfo.CID = mapRet[i]["CID"].(int64)
			}
			if mapRet[i]["DID"] != nil {
				loginIDInfo.DID = mapRet[i]["DID"].(int64)
			}
			if mapRet[i]["PID"] != nil {
				loginIDInfo.PID = mapRet[i]["PID"].(int64)
			}
		}
	}

	return loginIDInfo, err
}

// GetEndDate 根据账号获取到期日信息
func (service *UserService) GetEndDate(userName string) (string, error) {
	var mapRet []map[string]interface{}
	mapRet, err := service.userRepo.EMGet_EndDate_New_Result(userName)

	shareUserLogger.InfoFormat("[查询到期日] userName:%s EMGet_EndDate_New_Result返回：", userName, mapRet)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "获取到期日接口失败", mapRet)
		return "", err
	}

	if len(mapRet) > 0 {
		for i, _ := range mapRet {
			if mapRet[i]["enddate"] != nil {
				enddate := mapRet[i]["enddate"].(string)
				return enddate, nil
			}
		}
	}

	return "", err
}

// SavePicOrName 修改昵称头像
func (service *UserService) SavePicOrName(username string, pictrue string, nickname string) (string, string, error) {
	url := config.CurrentConfig.SaveAccountPicOrNameApiUrl
	url = fmt.Sprintf(url, username, pictrue, nickname)

	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	shareUserLogger.Info(fmt.Sprintf("SavePicOrName调用记录=> url:%s response:%s", url, body))

	if errReturn != nil {
		shareUserLogger.Error(errReturn, "SavePicOrName Api接口异常=>url="+url)
		return "", "", errReturn
	}
	_ = contentType
	_ = intervalTime
	webApiResp := contract.WebApiResponse{}
	err := _json.Unmarshal(body, &webApiResp)
	if err != nil {
		return "", "", err
	}
	retMsg := webApiResp.RetMsg
	retCode := webApiResp.RetCode

	return retCode, retMsg, err
}

// GetAccountProfile 获取个人资料
func (service *UserService) GetAccountProfile(uid int64, username string) (*model.ProfileInfo, error) {
	//登录账号是手机号，获取手机号绑定的EM号
	if validate.IsMobile(username) {
		boundAccount_Response, err := service.BoundGroupQryLogin(uid)
		if err != nil {
			shareUserLogger.Error(err, "BoundGroupQryLogin 接口异常=>username:"+username+"")
			return nil, err
		}

		for i, _ := range boundAccount_Response {
			accounttype := boundAccount_Response[i].AccountType
			if accounttype == 0 {
				username = boundAccount_Response[i].AccountName
			}
		}
	}

	profile := new(model.ProfileInfo)
	url := config.CurrentConfig.GetAccountProfileApiUrl
	url = fmt.Sprintf(url, username)

	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	shareUserLogger.Info(fmt.Sprintf("GetAccountProfile调用记录=> url:%s response:%s", url, body))

	if errReturn != nil {
		shareUserLogger.Error(errReturn, "GetAccountProfile Api接口异常=>url="+url)
		return nil, errReturn
	}
	_ = contentType
	_ = intervalTime
	webApiResp := contract.WebApiResponse{}
	err := _json.Unmarshal(body, &webApiResp)
	if err != nil {
		return nil, err
	}
	_ = webApiResp.RetMsg
	retCode := webApiResp.RetCode
	message := webApiResp.Message

	if retCode == "0" && message != nil {
		_json.Unmarshal(message.(string), profile)
	}

	return profile, err
}

// SaveAccountProfile 修改个人资料
func (service *UserService) SaveAccountProfile(myPro *model.ProfileInfo) (string, string, error) {
	url := config.CurrentConfig.SaveAccountProfileNewApiUrl
	url = fmt.Sprintf(url, myPro.UserID, myPro.UserName, myPro.Name, myPro.Phone, myPro.Mobile, myPro.IDCard, myPro.Address, myPro.Email, myPro.Sex, myPro.CoName, myPro.ProvinceName, myPro.Description, myPro.Age, myPro.AgreementA, myPro.AgreementA, myPro.TZvalue, myPro.QQ, myPro.ServiceAgentId, myPro.ProvinceName1, myPro.WeixinID, myPro.Birth_year, myPro.Birth_month, myPro.Birth_day, myPro.Contact, myPro.Tzph, myPro.Zjgm, myPro.Profession, myPro.Degree, myPro.Dp_time, myPro.InternationalNumber, myPro.ProvinceNumber, myPro.SmallNumber)

	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	shareUserLogger.Info(fmt.Sprintf("GetAccountProfile调用记录=> url:%s response:%s", url, body))

	if errReturn != nil {
		shareUserLogger.Error(errReturn, "GetAccountProfile Api接口异常=>url="+url)
		return "-1", "", errReturn
	}
	_ = contentType
	_ = intervalTime
	webApiResp := contract.WebApiResponse{}
	err := _json.Unmarshal(body, &webApiResp)
	if err != nil {
		return "", "", err
	}
	retMsg := webApiResp.RetMsg
	retCode := webApiResp.RetCode

	return retCode, retMsg, err
}

// BindAccountSelect 根据uid和pid获取绑定的手机号
/* 返回：
{
RetCode: "-104200002",
RetMsg: "111****1111",
Message: {
guid: "oX2F16dtBtNfR2UvyRH24g=="
}
}
*/
func (service *UserService) BindAccountSelect(uid int64, pid int) (string, string, error) {
	url := config.CurrentConfig.BindAccountSelectApi
	url = fmt.Sprintf(url, uid, pid)

	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	shareUserLogger.Info(fmt.Sprintf("BindAccountSelect调用记录=> url:%s response:%s", url, body))

	if errReturn != nil {
		shareUserLogger.Error(errReturn, "BindAccountSelect Api接口异常=>url="+url)
		return "", "", errReturn
	}
	_ = contentType
	_ = intervalTime
	webApiResp := contract.WebApiResponse{}
	err := _json.Unmarshal(body, &webApiResp)
	if err != nil {
		return "", "", err
	}
	retmsg := webApiResp.RetMsg
	retobj := webApiResp.Message
	guidmap := retobj.(map[string]string)

	if len(retmsg) == 11 && len(guidmap) > 0 {
		//base64转换为十六进制字符串
		mobilex, _ := mobile.Base642Hex(guidmap["guid"])
		return retmsg, mobilex, nil
	} else {
		return "", "", nil
	}
}

// BindAccount 绑定手机号
/* 返回
{
RetCode: "-104200002",
RetMsg: "@syjs2011已与手机号码15618681991绑定！@0",
Message: null
}
*/
func (service *UserService) BindAccount(uid int64, pid int, mobile string) (string, error) {
	url := config.CurrentConfig.BindAccountApi
	url = fmt.Sprintf(url, uid, pid, mobile)

	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	shareUserLogger.Info(fmt.Sprintf("BindAccountSelect调用记录=> url:%s response:%s", url, body))

	if errReturn != nil {
		shareUserLogger.Error(errReturn, "BindAccount Api接口异常=>url="+url)
		return "", errReturn
	}
	_ = contentType
	_ = intervalTime
	webApiResp := contract.WebApiResponse{}
	err := _json.Unmarshal(body, &webApiResp)
	if err != nil {
		return "", err
	}
	retmsg := webApiResp.RetMsg

	strings_msg := strings.Split(retmsg, "@")
	if len(strings_msg) > 0 {
		if strings.ContainsAny(strings_msg[1], "绑定") {
			return "0", nil
		}
	}
	return "-1", nil
}

// LoginDaysAndProductByCID 获取用户连登天数以及产品信息
func (srv *UserService) LoginDaysAndProductByCID(cid int64) (loginDays int, continueLoginDays int, product string, activateTime string, endTime string, err error) {
	var mapRet []map[string]interface{}
	shareUserLogger.DebugFormat("LoginDaysAndProductByCID 获取用户连登天数以及产品信息 cid=%d", cid)
	mapRet, err = srv.userRepo.LoginDaysAndProduct(cid)

	if err != nil {
		shareUserLogger.ErrorFormat(err, "获取用户连登天数以及产品信息失败 mapRet=%s", _json.GetJsonString(mapRet))
		return
	}
	shareUserLogger.DebugFormat("LoginDaysAndProductByCID 获取用户连登天数以及产品信息 cid=%d mapRet=%s", cid, _json.GetJsonString(mapRet))
	if len(mapRet) > 0 {
		if mapRet[0]["TtlLoginDays"] != nil {
			loginDays = int(mapRet[0]["TtlLoginDays"].(int64))
		}
		if mapRet[0]["MaxLoginDays"] != nil {
			continueLoginDays = int(mapRet[0]["MaxLoginDays"].(int64))
		}
		if continueLoginDays > 3 {
			if continueLoginDays%3 == 0 {
				continueLoginDays = 3
			} else {
				continueLoginDays = continueLoginDays % 3
			}
		}

		if mapRet[0]["Product"] != nil {
			product = mapRet[0]["Product"].(string)
		}
		activateTimeTmp := mapRet[0]["ActiveTime"]
		if activateTimeTmp != nil {
			activateTime = activateTimeTmp.(time.Time).Format("2006-01-02")
		}
		endTimeTmp := mapRet[0]["EndDate"]
		if endTimeTmp != nil {
			endTime = endTimeTmp.(time.Time).Format("2006-01-02")
		}
		return
	} else {
		return
	}
}

// CidFindGid 根据Cid获取Gid
// gid<=0则未查询到相关的信息
func (srv *UserService) CidFindGid(cid int64) (gid int64, err error) {
	var mapRet []map[string]interface{}
	mapRet, err = srv.userRepo.CidFindGid(cid)

	if err != nil {
		shareUserLogger.ErrorFormat(err, "CidFindGid 根据Cid获取Gid mapRet=%s", _json.GetJsonString(mapRet))
		return
	}
	shareUserLogger.DebugFormat("CidFindGid 根据Cid获取Gid cid=%d mapRet=%s", cid, _json.GetJsonString(mapRet))
	if len(mapRet) > 0 {
		gid = mapRet[0]["parentid"].(int64)
		return gid, err
	} else {
		return
	}
}

// NoticeRedis_TouJiaoSZTX 绑定手机号变更-通知投教实战体系
func (srv *UserService) NoticeRedis_TouJiaoSZTX(emcard string) {
	key := TJSZTXkey_pre + "HasChangedEmcard"
	err := srv.RedisCache.HSet(key, emcard, emcard)
	shareUserLogger.InfoFormat("绑定手机号变更-通知投教实战体系 key:&s emcard:%s ", key, emcard)

	if err != nil {
		shareUserLogger.ErrorFormat(err, "绑定手机号变更-通知投教实战体系失败 emcard=%s", emcard)
		return
	}
}

func init() {
	protected.RegisterServiceLoader(userServiceName, userServiceLoader)
}

func userServiceLoader() {
	shareUserRepo = user.NewUserRepository(protected.DefaultConfig)
	shareUserLogger = dotlog.GetLogger(userServiceName)
}
