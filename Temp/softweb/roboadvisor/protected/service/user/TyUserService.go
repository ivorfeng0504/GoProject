package user

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/repository/user"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

type TyUserService struct {
	service.BaseService
	tyuserRepo *user.TyUserRepository
}

var (
	// shareTyUserRepo TyUserRepository
	shareTyUserRepo *user.TyUserRepository

	// shareTyUserLogger 共享的Logger实例
	shareTyUserLogger dotlog.Logger
)

const (
	tyuserServiceName = "tyUserServiceLogger"
)

func NewTyUserService() *TyUserService {
	tyuserService := &TyUserService{
		tyuserRepo: shareTyUserRepo,
	}
	tyuserService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return tyuserService
}


// QKTaste_RegistMobile_Web 获取手机号密码
func (service *TyUserService)  QKTaste_RegistMobile_Web(mobile string,hardinfo string,sid int ,tid int) (retCode int,retMsg string,passWord string,err error) {
	var mapRet []map[string]interface{}
	mapRet, err = service.tyuserRepo.QKTaste_RegistMobile_Web(mobile, hardinfo, sid, tid)

	shareUserLogger.InfoFormat("[获取手机号密码] mobile:%s sid:%d tid:%d  链接字符串:%s QKTaste_RegistMobile_Web返回：", mobile, sid, tid, protected.DefaultConfig.EmoneyDBConn, mapRet)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "获取手机号密码接口失败", mapRet)
		return -1, "", "", err
	}
	if len(mapRet) > 0 {
		for i, _ := range mapRet {
			retCode = int(mapRet[i]["returnNo"].(int64))
			retMsg = mapRet[i]["msg"].(string)
			passWord = mapRet[i]["passWord"].(string)
		}
	}

	return retCode, retMsg, passWord, err
}

/*
获取手机号密码
*/
func (service *TyUserService) GetMobilePwd_emoney(mobile string,hardware string,sid int, tid int) (int,string, error) {
	url := config.CurrentConfig.GetMobilePwdApiUrl
	url = fmt.Sprintf(url, mobile, hardware, sid, tid)

	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	shareUserLogger.Info(fmt.Sprintf("GetMobilePwd_emoney调用记录=> url:%s response:%s", url, body))

	if errReturn != nil {
		shareUserLogger.Error(errReturn, "GetMobilePwd_emoney Api接口异常=>url="+url)
		return -1, "", errReturn
	}
	_ = contentType
	_ = intervalTime
	webApiResp := contract.WebApiResponse{}
	err := _json.Unmarshal(body, &webApiResp)
	if err != nil {
		return -1, "", err
	}
	retcode := webApiResp.RetCode
	retobj := webApiResp.Message
	if retobj!=nil{
		pwdmap := retobj.(map[string]interface{})

		if retcode == "0" && len(pwdmap) > 0 {
			pwd := pwdmap["Msg"]
			return 0, pwd.(string), nil
		}
	}

	return -1, "", nil
}


func init() {
	protected.RegisterServiceLoader(tyuserServiceName, tyuserServiceLoader)
}

func tyuserServiceLoader() {
	shareTyUserRepo = user.NewTyUserRepository(protected.DefaultConfig)
	shareTyUserLogger = dotlog.GetLogger(tyuserServiceName)
}
