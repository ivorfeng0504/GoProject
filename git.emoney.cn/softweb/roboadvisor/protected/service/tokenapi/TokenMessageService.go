package tokenapi

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
)

const (
	tokenMessageServiceName = "TokenMessageService"
)

var (
	shareTokenMessageLogger dotlog.Logger
)

func init() {
	protected.RegisterServiceLoader(tokenMessageServiceName, func() {
		shareTokenMessageLogger = dotlog.GetLogger(tokenMessageServiceName)
	})
}

// CreateToken 创建Token
func CreateToken(appId string, lifeSeconds int, data interface{}) (response *HandlerResponse, err error) {
	response = &HandlerResponse{}
	url := config.CurrentConfig.TokenCreateUrl
	dataJson := _json.GetJsonString(data)
	tokenInfo := TokenInfo{
		TokenBody:   dataJson,
		AppID:       appId,
		LifeSeconds: lifeSeconds,
	}
	requestJson := _json.GetJsonString(tokenInfo)
	body, contentType, intervalTime, errReturn := _http.HttpPost(url, requestJson, "application/json")
	if errReturn != nil {
		shareTokenMessageLogger.ErrorFormat(errReturn, "CreateToken HTTP访问异常  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
		return nil, errReturn
	}
	_ = contentType
	_ = intervalTime
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "API网关请求结果反序列化失败"
		shareTokenMessageLogger.ErrorFormat(err, "CreateToken API网关请求结果反序列化失败  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
		return response, err
	}
	if apiGatewayResp.RetCode != 0 {
		response.RetCode = apiGatewayResp.RetCode
		response.RetMsg = apiGatewayResp.RetMsg
		shareTokenMessageLogger.WarnFormat("CreateToken 网关访问异常  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
		return response, err
	}

	err = _json.Unmarshal(apiGatewayResp.Message, &response)
	if err != nil {
		shareTokenMessageLogger.ErrorFormat(err, "CreateToken 真实接口反序列化失败  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
	} else {
		shareTokenMessageLogger.DebugFormat("CreateToken  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
	}
	return response, err
}

// QueryToken 查询Token Token可多次使用
func QueryToken(appId string, token string) (response *HandlerResponse, err error) {
	response = &HandlerResponse{}
	url := fmt.Sprintf("%s&appid=%s&token=%s", config.CurrentConfig.TokenQueryUrl, appId, token)
	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	shareTokenMessageLogger.DebugFormat("QueryToken接口调用 url=%s   result=【%s】", url, body)
	if errReturn != nil {
		shareTokenMessageLogger.ErrorFormat(errReturn, "QueryToken HTTP访问异常  请求地址为:%s  响应结果为:%s", url, body)
		return nil, errReturn
	}
	_ = contentType
	_ = intervalTime

	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "API网关请求结果反序列化失败"
		shareTokenMessageLogger.ErrorFormat(err, "QueryToken API网关请求结果反序列化失败  请求地址为:%s  响应结果为:%s", url, body)
		return response, err
	}
	if apiGatewayResp.RetCode != 0 {
		response.RetCode = apiGatewayResp.RetCode
		response.RetMsg = apiGatewayResp.RetMsg
		shareTokenMessageLogger.WarnFormat("QueryToken 网关访问异常  请求地址为:%s  响应结果为:%s", url, body)
		return response, err
	}

	err = _json.Unmarshal(apiGatewayResp.Message, &response)
	if err != nil {
		shareTokenMessageLogger.ErrorFormat(err, "QueryToken 请求结果反序列化失败  请求地址为:%s  响应结果为:%s", url, body)
		return nil, err
	}
	return response, err
}

// VerifyToken 校验Token，校验之后Token失效
func VerifyToken(appId string, token string, isCheckBody bool, tokenBody string) (response *HandlerResponse, err error) {
	response = &HandlerResponse{}
	url := config.CurrentConfig.TokenVerifyUrl
	request := VerifyTokenRequest{
		Token:       token,
		AppID:       appId,
		IsCheckBody: isCheckBody,
		TokenBody:   tokenBody,
	}
	requestJson := _json.GetJsonString(request)
	body, contentType, intervalTime, errReturn := _http.HttpPost(url, requestJson, "application/json")
	if errReturn != nil {
		shareTokenMessageLogger.ErrorFormat(errReturn, "VerifyToken HTTP访问异常  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
		return nil, errReturn
	}
	_ = contentType
	_ = intervalTime
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "API网关请求结果反序列化失败"
		shareTokenMessageLogger.ErrorFormat(err, "VerifyToken API网关请求结果反序列化失败  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
		return response, err
	}
	if apiGatewayResp.RetCode != 0 {
		response.RetCode = apiGatewayResp.RetCode
		response.RetMsg = apiGatewayResp.RetMsg
		shareTokenMessageLogger.WarnFormat("VerifyToken 网关访问异常  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
		return response, err
	}

	err = _json.Unmarshal(apiGatewayResp.Message, &response)
	if err != nil {
		shareTokenMessageLogger.ErrorFormat(err, "VerifyToken 真实接口反序列化失败  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
	} else {
		shareTokenMessageLogger.DebugFormat("VerifyToken  请求地址为:%s  请求参数为:%s  响应结果为:%s", url, requestJson, body)
	}
	return response, err
}
