package mobile

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"errors"
	"time"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	//"net/url"
	"net/url"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type MessageService struct {
	service.BaseService
}
var (
	shareMessageLogger dotlog.Logger
)

const (
	messageServiceName = "MessageService"
)

func NewMessageService() *MessageService {
	service := &MessageService{

	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return service
}

// 发送短信验证码
/*
departmentId：部门编号
description：30|999999999|短信用途
expire：过期时间
nextTime：下次发送短信时间间隔
uniqueTag：唯一标识
applicationId：事先分配好
*/
func (service *MessageService) SendValidateCode(mobile string, content string, code string, departmentId string, description string, expire time.Time, nextTime time.Time, uniqueTag string, applicationId string) (retMsg string,err error) {
	msgkey := uniqueTag + "_" + mobile
	nowtime := time.Now()
	msg := new(contract.SendCodeMessage)

	err = service.RedisCache.GetJsonObj(msgkey, msg)

	if err != nil {
		shareMessageLogger.Error(err, "发送验证码获取redis异常 mobile="+mobile)
	}

	if msg != nil && nowtime.Before(msg.NextTime) {
		return "您的操作太频繁，请稍后再试", nil
	}

	retCode, err := SendMessage(content, mobile, departmentId, description, applicationId)

	//retCode := "0"
	if retCode != "0" {
		return "系统异常，请稍后再试", nil
	}

	msg = new(contract.SendCodeMessage)
	msg.Msg = code
	msg.Expire = expire
	msg.NextTime = nextTime

	_, err = service.RedisCache.SetJsonObj(msgkey, msg)

	if err != nil {
		shareMessageLogger.Error(err, "发送验证码存入redis异常 mobile="+mobile)
	}

	return "0", nil
}

func (service *MessageService) CheckValidateCode(mobile string, code string, uniqueTag string) (retMsg string,err error) {
	msgkey := uniqueTag + "_" + mobile
	nowtime := time.Now()
	msg := new(contract.SendCodeMessage)
	fmt.Println(msgkey)
	err = service.RedisCache.GetJsonObj(msgkey, msg)
	fmt.Println(msg)
	if err != nil {
		shareMessageLogger.Error(err, "验证验证码获取redis异常 mobile="+mobile)
		//return "系统异常，请稍后再试", err
	}

	if msg == nil || msg.Expire.Before(nowtime) {
		return "您的验证码无效，请重新发送验证码", nil
	}

	if code == msg.Msg {
		return "0", nil
	} else {
		return "验证码错误", nil
	}
}

func SendMessage(content string,mobile string,departmentId string,description string,applicationId string) (retCode string, err error) {
	sendmsgurl := config.CurrentConfig.SendMsgApi
	appid := config.CurrentConfig.AppId

	//明文手机号加密
	if len(mobile) == 11 {
		mobile, err = EncryptMobileHex(mobile)

		if err != nil {
			shareMessageLogger.Error(err, "发送短信加密手机异常 mobile="+mobile)
			return "-1", err
		}
	}
	urlparms := "appId=%s&ApplicationId=%s&content=%s&departmentId=%s&description=%s&level=%s&remark=%s"
	urlparms = fmt.Sprintf(urlparms, appid, applicationId, content, departmentId, description, "1", "")

	sendmsgurl = sendmsgurl + "?" + url.PathEscape(urlparms) + "&phone=" + mobile
	fmt.Println(sendmsgurl)

	body, contentType, intervalTime, errReturn := _http.HttpGet(sendmsgurl)
	fmt.Println(body)
	if errReturn != nil {
		return retCode, errReturn
	}
	_ = contentType
	_ = intervalTime
	smsApiResp := contract.WebApiResponse{}
	err = _json.Unmarshal(body, &smsApiResp)
	if err != nil {
		return retCode, err
	}
	retCode = smsApiResp.RetCode

	if retCode != "0" {
		return retCode, errors.New(smsApiResp.RetMsg)
	} else {
		return retCode, nil
	}
}

func init() {
	protected.RegisterServiceLoader(messageServiceName, messageServiceLoader)
}

func messageServiceLoader() {
	shareMessageLogger = dotlog.GetLogger(messageServiceName)
}