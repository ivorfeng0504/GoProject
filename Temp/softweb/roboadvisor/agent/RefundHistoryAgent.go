package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/scmapi"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"time"
)

// QueryOrderList 订单查询
func QueryOrderList(userInfo userhome_model.UserInfo) (orderList []*RefundOrderInfo, err error) {
	req := contract.NewApiRequest()
	now := time.Now()
	request := scmapi.QueryOrderProdListByParamsRequest{
		EmCard:          userInfo.Account,
		MidPwd:          userInfo.MobileX,
		OrdAddTimeStart: now.AddDate(0, -6, 0).Format("2006-01-02"),
		OrdAddTimeEnd:   now.AddDate(0, 0, 1).Format("2006-01-02"),
		Refund_Sign:     "1",
		WitchState:      "A9002",
	}

	//如果MidPwd与EmCard都为空 则直接返回
	if len(request.MidPwd) == 0 && len(request.EmCard) == 0 {
		return orderList, nil
	}

	//如果获取到手机号 则不传递EM账号
	if len(request.MidPwd) > 0 {
		request.EmCard = ""
	}
	req.RequestData = request
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/refundorder/queryorderlist", req)
	if err != nil {
		return orderList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return orderList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &orderList)
	return orderList, nil
}

// ValidateOrder 验证是否支持快速退货退款
func ValidateOrder(orderId string) (result scmapi.ValidateReturnbackAndRefundResponse, err error) {
	if len(orderId) == 0 {
		err = errors.New("订单号不正确")
		return result, err
	}
	req := contract.NewApiRequest()
	req.RequestData = orderId
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/refundorder/validateorder", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return result, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

// RefundSubmit 提交退款
func RefundSubmit(submitData RefundSubmitData) (result scmapi.ReturnbackAndRefundResponse, err error) {
	req := contract.NewApiRequest()
	req.RequestData = submitData
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/refundorder/refundsubmit", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return result, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

// GetRefundStatus 查询退款状态
func GetRefundStatus(orderId string) (result []*scmapi.OrderRefundState, err error) {
	if len(orderId) == 0 {
		err = errors.New("订单号不正确")
		return result, err
	}
	req := contract.NewApiRequest()
	req.RequestData = orderId
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/refundorder/getrefundstatus", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return result, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

//退款订单列表
type RefundOrderInfo struct {
	//订单号
	OrderId string
	//产品列表
	ProductList []*RefundProductInfo
	//产品总价
	Price float32
	//是否可以退款
	CanRefund bool
	//是否已经申请退款
	IsRefund bool
}

//产品信息
type RefundProductInfo struct {
	//产品名称
	ProductName string
	//价格
	Price float32
	//数量
	Count int
	//购买时间
	CreateTime string
	//状态
	State int
	//状态描述
	StateDesc string
	//是否支持快速退货退款0：不支持；1：支持
	IsQuickReturn int

	//退款时间 前端用
	RefundTime string
}

//用户退款提交参数
type RefundSubmitData struct {
	//订单号
	OrderId string
	//退款理由
	Reason string
	//用户姓名
	Name string
	//银行账号
	BankAccount string
	//收款银行
	BankValue string
	//银行详情，支行信息
	BankDetail string
	//退款方式
	RefundMode int
	//EM账号
	Account string
	//手机号
	Mobile string
	//调用方
	RefundAppId string
}
