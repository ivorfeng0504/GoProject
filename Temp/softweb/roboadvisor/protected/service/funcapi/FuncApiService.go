package funcapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"github.com/devfeel/dotlog"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	shareFuncApiLogger dotlog.Logger
)

const (
	funcApiServiceName = "FuncApiService"
	//大单天眼特权Id
	TeQuan_DDTY = "1"
)

func init() {
	protected.RegisterServiceLoader(funcApiServiceName, func() {
		shareFuncApiLogger = dotlog.GetLogger(funcApiServiceName)
	})
}

// OpenJRPTFunc 开通金融平台功能 funcId见JRPTFunc.go
func OpenJRPTFunc(funcId int, mobile string, enableDays int) (success bool, err error) {
	presentId := ""
	applyUserName := "智投-用户中心"
	switch funcId {
	case _const.JRPTFunc_DDTY:
		presentId = TeQuan_DDTY
		break
	}
	if presentId == "" {
		err = errors.New("无效的功能Id")
		shareFuncApiLogger.ErrorFormat(err, "OpenJRPTFunc 开通金融平台功能失败 无效的功能Id funcId=%d mobile=%s enableDays=%d", strconv.Itoa(funcId), mobile, enableDays)
		return false, err
	}
	result, err := SendTeQuanCore(mobile, presentId, enableDays, applyUserName)
	if err != nil {
		shareFuncApiLogger.ErrorFormat(err, "OpenJRPTFunc 开通金融平台功能失败 funcId=%d mobile=%s enableDays=%d", strconv.Itoa(funcId), mobile, enableDays)
		return false, err
	}
	if result == nil || result.Code != 0 {
		shareFuncApiLogger.ErrorFormat(nil, "OpenJRPTFunc 开通金融平台功能失败 funcId=%d mobile=%s enableDays=%d", strconv.Itoa(funcId), mobile, enableDays)
		return false, nil
	} else {
		shareFuncApiLogger.InfoFormat("OpenJRPTFunc 开通金融平台功能成功 funcId=%d mobile=%s enableDays=%d", strconv.Itoa(funcId), mobile, enableDays)
		return true, nil
	}
}

// SendTeQuan 特权开通
// mobile 手机号
// presentId 特权Id
// enableDays开通天数
// applyUserName 申请人
func SendTeQuanCore(mobile string, presentId string, enableDays int, applyUserName string) (result *FuncApiResponse, err error) {
	if len(config.CurrentConfig.SendTeQuanApi) == 0 {
		err = errors.New("未配置特权开通API地址")
		shareFuncApiLogger.ErrorFormat(err, "未配置特权开通API地址 mobile=%s presentId=%s enableDays=%d applyUserName=%s", mobile, presentId, enableDays, applyUserName)
		return nil, err
	}
	//手机号
	accountType := "M"
	apiUrl := fmt.Sprintf(config.CurrentConfig.SendTeQuanApi, mobile, accountType, presentId, enableDays, applyUserName)
	shareFuncApiLogger.InfoFormat("SendTeQuan 准备开通特权 请求地址为：%s", apiUrl)
	resp, err := http.Get(apiUrl)
	defer func() { resp.Body.Close() }()
	if err != nil {
		shareFuncApiLogger.ErrorFormat(err, "SendTeQuan 开通特权异常 请求地址为：%s", apiUrl)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("特权开通API请求异常，状态码不正确")
		shareFuncApiLogger.ErrorFormat(err, "SendTeQuan 开通特权异常  状态码不正确 状态码=%d 请求地址为：%s", resp.StatusCode, apiUrl)
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		shareFuncApiLogger.ErrorFormat(err, "SendTeQuan 开通特权异常 读取请求体异常 请求地址为：%s", apiUrl)
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		shareFuncApiLogger.ErrorFormat(err, "SendTeQuan 开通特权异常 结果反序列化失败 请求地址为：%s", apiUrl)
		return nil, err
	}
	if result.Code == 0 {
		shareFuncApiLogger.InfoFormat("SendTeQuan 开通特权成功 响应结果为%s 请求地址为：%s", string(data), apiUrl)
	} else {
		shareFuncApiLogger.ErrorFormat(err, "SendTeQuan 开通特权异常 响应结果为%s 请求地址为：%s", string(data), apiUrl)
	}
	return result, nil
}
