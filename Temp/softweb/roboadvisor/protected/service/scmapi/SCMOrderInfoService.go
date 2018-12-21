package scmapi

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"strings"
	"unicode/utf8"
)

const (
	SCMOrderInfoServiceName = "SCMOrderInfoService"
)

var (
	shareSCMOrderInfoServiceLogger dotlog.Logger
)

func init() {
	protected.RegisterServiceLoader(SCMOrderInfoServiceName, func() {
		shareSCMOrderInfoServiceLogger = dotlog.GetLogger(SCMOrderInfoServiceName)
	})
}

// GetRefundSignDesc 获取退款状态描述
func GetRefundSignDesc(state int) string {
	dict := map[int]string{0: "已购买", -1: "已退货", 1: "全部"}
	return dict[state]
}

// QueryOrderProdListByParams 查询指定用户的订单 account与mobile传递其中之一即可 state为订单退款状态
func QueryOrderProdListByParams(request QueryOrderProdListByParamsRequest) (orderList []*OrderProdInfo, err error) {
	if len(request.EmCard) == 0 && len(request.MidPwd) == 0 {
		err = errors.New("EmCard与MidPwd不能同时为空")
		return orderList, err
	}
	apiUrl := config.CurrentConfig.SCMAPI_QueryOrderProdListByParamsApi
	if len(apiUrl) == 0 {
		err = errors.New("SCMAPI_QueryOrderProdListByParamsApi 配置不能为空")
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "QueryOrderProdListByParams 查询指定用户的订单 SCMAPI_QueryOrderProdListByParamsApi 配置不能为空")
		return nil, err
	}
	requestJson := _json.GetJsonString(request)
	apiUrl += "&jsonStr=" + requestJson
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(errReturn, "QueryOrderProdListByParams 查询指定用户的订单 网关请求异常 requestUrl=%s body=%s", apiUrl, body)
		return nil, errReturn
	}
	_ = contentType
	_ = intervalTime
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "QueryOrderProdListByParams 查询指定用户的订单 网关响应数据解析异常 requestUrl=%s body=%s", apiUrl, body)
		return nil, err
	}
	if apiGatewayResp.RetCode != 0 {
		err = errors.New(apiGatewayResp.RetMsg)
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "QueryOrderProdListByParams 查询指定用户的订单 网关响应状态码异常 requestUrl=%s apiGatewayResp=%s body=%s", apiUrl, _json.GetJsonString(apiGatewayResp), body)
		return nil, err
	}
	responseMsg := QueryOrderProdListByParamsResponse{}
	err = _json.Unmarshal(apiGatewayResp.Message, &responseMsg)
	if err != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "QueryOrderProdListByParams 查询指定用户的订单 apiGatewayResp.Message解析异常 requestUrl=%s apiGatewayResp=%s body=%s", apiUrl, _json.GetJsonString(apiGatewayResp), body)
		return nil, err
	}
	if responseMsg.Code != 1 {
		err = errors.New("查询指定用户的订单异常 " + responseMsg.Msg)
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "QueryOrderProdListByParams 查询指定用户的订单 接口异常 requestUrl=%s responseMsg=%s body=%s", apiUrl, _json.GetJsonString(responseMsg), body)
		return nil, err
	}
	orderList = responseMsg.Data
	shareSCMOrderInfoServiceLogger.DebugFormat("QueryOrderProdListByParams 查询指定用户的订单 request=%s orderList=%s apiUrl=%s body=%s", _json.GetJsonString(request), _json.GetJsonString(orderList), apiUrl, body)
	return orderList, nil
}

// ValidateReturnbackAndRefund 验证是否支持快速退货退款
func ValidateReturnbackAndRefund(orderId string) (result ValidateReturnbackAndRefundResponse, err error) {
	apiUrl := config.CurrentConfig.SCMAPI_ValidateReturnbackAndRefundApi
	if len(apiUrl) == 0 {
		err = errors.New("SCMAPI_ValidateReturnbackAndRefundApi 配置不能为空")
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ValidateReturnbackAndRefund 验证是否支持快速退货退款 SCMAPI_ValidateReturnbackAndRefundApi 配置不能为空")
		return result, err
	}
	apiUrl = fmt.Sprintf(apiUrl, orderId)
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(errReturn, "ValidateReturnbackAndRefund 验证是否支持快速退货退款 网关请求异常 requestUrl=%s body=%s", apiUrl, body)
		return result, errReturn
	}
	_ = contentType
	_ = intervalTime
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ValidateReturnbackAndRefund 验证是否支持快速退货退款 网关响应数据解析异常 requestUrl=%s body=%s", apiUrl, body)
		return result, err
	}
	if apiGatewayResp.RetCode != 0 {
		err = errors.New(apiGatewayResp.RetMsg)
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ValidateReturnbackAndRefund 验证是否支持快速退货退款 网关响应状态码异常 requestUrl=%s apiGatewayResp=%s body=%s", apiUrl, _json.GetJsonString(apiGatewayResp), body)
		return result, err
	}
	err = _json.Unmarshal(apiGatewayResp.Message, &result)
	if err != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ValidateReturnbackAndRefund 验证是否支持快速退货退款 apiGatewayResp.Message解析异常 requestUrl=%s apiGatewayResp=%s body=%s", apiUrl, _json.GetJsonString(apiGatewayResp), body)
		return result, err
	}
	if result.StatusCode != 0 {
		err = errors.New("验证是否支持快速退货退款 " + result.Message)
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ValidateReturnbackAndRefund 验证是否支持快速退货退款 接口异常 requestUrl=%s responseMsg=%s body=%s", apiUrl, _json.GetJsonString(result), body)
		return result, err
	}
	return result, nil
}

// ReturnbackAndRefund 快速退货退款
func ReturnbackAndRefund(request ReturnbackAndRefundRequest) (result ReturnbackAndRefundResponse, err error) {
	apiUrl := config.CurrentConfig.SCMAPI_ReturnbackAndRefundApi
	if len(apiUrl) == 0 {
		err = errors.New("SCMAPI_ReturnbackAndRefundApi 配置不能为空")
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 SCMAPI_ReturnbackAndRefundApi 配置不能为空")
		return result, err
	}
	postData := _json.GetJsonString(request)

	if len(request.OrderId) == 0 {
		err = errors.New("订单号不能为空")
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 OrderId 订单号不能为空 postData=%s", postData)
		return result, err
	}
	calcReason := request.Reason
	reasonIndex := strings.Index(calcReason, "|")
	if reasonIndex >= 0 {
		calcReason = calcReason[reasonIndex+1:]
	}
	if len(calcReason) == 0 {
		err = errors.New("退款理由不能为空")
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 Reason 退款理由不能为空 postData=%s", postData)
		return result, err
	}

	if utf8.RuneCountInString(calcReason) > 200 {
		err = errors.New("退款理由不能超过200个字符")
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 Reason 退款理由不能超过200个字符 postData=%s", postData)
		return result, err
	}

	if request.ReturnFlow != _const.ReturnFlow_Quick {
		err = errors.New("不支持的退货流程类型")
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 ReturnFlow 不支持的退货流程类型 postData=%s", postData)
		return result, err
	}

	if request.ReturnMode == _const.RefundMode_Bank {
		if len(request.OrderRetBankInfo.ReProName) == 0 {
			err = errors.New("收款人不能为空")
			shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 ReProName 收款人不能为空 postData=%s", postData)
			return result, err
		}

		if utf8.RuneCountInString(request.OrderRetBankInfo.ReProName) > 100 {
			err = errors.New("收款人不能大于100个字符")
			shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 ReProName 收款人不能大于100个字符 postData=%s", postData)
			return result, err
		}

		if len(request.OrderRetBankInfo.ReProCode) == 0 {
			err = errors.New("收款人账号不能为空")
			shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 ReProCode 收款人账号不能为空 postData=%s", postData)
			return result, err
		}
		if utf8.RuneCountInString(request.OrderRetBankInfo.ReProCode) > 30 {
			err = errors.New("收款人账号不能大于30个字符")
			shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 ReProCode 收款人账号不能大于30个字符 postData=%s", postData)
			return result, err
		}

		if len(request.OrderRetBankInfo.ReBankName) == 0 {
			err = errors.New("收款银行不能为空")
			shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 ReBankName 收款银行不能为空 postData=%s", postData)
			return result, err
		}

		if len(request.OrderRetBankInfo.ReBankArea) == 0 {
			err = errors.New("开户行不能为空")
			shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 ReBankArea 开户行不能为空 postData=%s", postData)
			return result, err
		}
		if utf8.RuneCountInString(request.OrderRetBankInfo.ReBankArea) > 100 {
			err = errors.New("开户行不能大于100个字符")
			shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 ReBankArea 开户行不能大于100个字符 postData=%s", postData)
			return result, err
		}
	}
	body, contentType, intervalTime, errReturn := _http.HttpPost(apiUrl, postData, "application/json")
	if errReturn != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(errReturn, "ReturnbackAndRefund 快速退货退款异常 网关请求异常 requestUrl=%s postData=%s", apiUrl, postData)
		return result, errReturn
	}
	_ = contentType
	_ = intervalTime
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 网关响应数据解析异常 requestUrl=%s body=%s", apiUrl, body)
		return result, err
	}
	if apiGatewayResp.RetCode != 0 {
		err = errors.New(apiGatewayResp.RetMsg)
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 网关响应状态码异常 requestUrl=%s apiGatewayResp=%s body=%s", apiUrl, _json.GetJsonString(apiGatewayResp), body)
		return result, err
	}
	err = _json.Unmarshal(apiGatewayResp.Message, &result)
	if err != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 apiGatewayResp.Message解析异常 requestUrl=%s apiGatewayResp=%s body=%s", apiUrl, _json.GetJsonString(apiGatewayResp), body)
		return result, err
	}
	if result.StatusCode != 0 {
		err = errors.New(result.Message)
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "ReturnbackAndRefund 快速退货退款异常 接口异常 requestUrl=%s responseMsg=%s body=%s", apiUrl, _json.GetJsonString(result), body)
		return result, err
	}
	return result, nil
}

// GetRefundStatus 查询退款状态
func GetRefundStatus(orderId string) (result []*OrderRefundState, err error) {
	apiUrl := config.CurrentConfig.SCMAPI_GetRefundStatusApi
	if len(apiUrl) == 0 {
		err = errors.New("SCMAPI_GetRefundStatusApi 配置不能为空")
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "GetRefundStatus 查询退款状态 SCMAPI_GetRefundStatusApi 配置不能为空")
		return result, err
	}
	apiUrl = fmt.Sprintf(apiUrl, orderId, _const.RefundAppId_Roboadvior)
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(errReturn, "GetRefundStatus 查询退款状态 网关请求异常 requestUrl=%s body=%s", apiUrl, body)
		return result, errReturn
	}
	_ = contentType
	_ = intervalTime
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "GetRefundStatus 查询退款状态 网关响应数据解析异常 requestUrl=%s body=%s", apiUrl, body)
		return result, err
	}
	if apiGatewayResp.RetCode != 0 {
		err = errors.New(apiGatewayResp.RetMsg)
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "GetRefundStatus 查询退款状态 网关响应状态码异常 requestUrl=%s apiGatewayResp=%s body=%s", apiUrl, _json.GetJsonString(apiGatewayResp), body)
		return result, err
	}

	var refundResponse OrderRefundStateResponse
	err = _json.Unmarshal(apiGatewayResp.Message, &refundResponse)
	if err != nil {
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "GetRefundStatus 查询退款状态 apiGatewayResp.Message解析异常 requestUrl=%s apiGatewayResp=%s body=%s", apiUrl, _json.GetJsonString(apiGatewayResp), body)
		return result, err
	}

	if refundResponse.StatusCode != 0 {
		err = errors.New("查询退款状态异常 " + refundResponse.Message)
		shareSCMOrderInfoServiceLogger.ErrorFormat(err, "GetRefundStatus 查询退款状态 接口异常 requestUrl=%s responseMsg=%s body=%s", apiUrl, _json.GetJsonString(result), body)
		return result, err
	}
	result = refundResponse.OrderRefundList
	return result, nil
}
