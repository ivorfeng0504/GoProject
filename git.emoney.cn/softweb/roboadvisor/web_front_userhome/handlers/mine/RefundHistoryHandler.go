package mine

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/scmapi"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel/userhome"
	"github.com/devfeel/dotweb"
)

func MyOrder(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "my_order.html")
}

// OrderDetail 订单详情
func OrderDetail(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "userMyRefundList.html")
}

// QueryOrderList 订单查询
func QueryOrderList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	if user == nil {
		response.RetCode = -1
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}
	orderList, err := agent.QueryOrderList(*user)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "订单查询异常"
		global.InnerLogger.ErrorFormat(err, "订单查询异常 uid=%s", user.UID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = orderList
	return ctx.WriteJson(response)
}

// ValidateOrder 验证是否支持快速退货退款
func ValidateOrder(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	if user == nil {
		response.RetCode = -1
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}
	orderId := ctx.PostFormValue("OrderId")
	validateResp := viewmodel.ValidateOrderResponse{
		NickName: user.NickName,
	}
	result, err := agent.ValidateOrder(orderId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "验证是否支持快速退货退款异常"
		global.InnerLogger.ErrorFormat(err, "验证是否支持快速退货退款异常 uid=%s orderId=%s", user.UID, orderId)
		return ctx.WriteJson(response)
	}
	if result.IsRefundByCustomer {
		if result.IsOriginRefund {
			validateResp.ReFundMode = _const.RefundMode_Third
			validateResp.ReFundModeDesc = "微信/支付宝"
		} else {
			validateResp.ReFundMode = _const.RefundMode_Bank
			validateResp.ReFundModeDesc = "银行汇款"
		}
	} else {
		validateResp.ReFundMode = _const.RefundMode_NotSupport
		validateResp.ReFundModeDesc = "不支持退款"
	}

	//银行列表
	if validateResp.ReFundMode == _const.RefundMode_Bank {
		bankList := []string{"中国农业银行", "交通银行", "北京银行", "上海银行", "中国银行", "中国建设银行", "成都银行", "中国光大银行",
			"广发银行", "兴业银行", "中国民生银行", "招商银行", "城市商业银行", "浙商银行", "中信银行", "汉口银行", "徽商银行", "华夏银行",
			"杭州银行", "中国工商银行", "金华银行", "九江银行", "江苏银行", "锦州银行", "宁波银行", "南京银行", "平安银行", "中国邮政储蓄银行", "农村商业银行", "上海浦东发展银行"}
		for _, bankName := range bankList {
			validateResp.BankInfoList = append(validateResp.BankInfoList, &viewmodel.BankInfo{
				BankName:  bankName,
				BankValue: bankName,
			})
		}
	}
	//if result.PaymentDetailList != nil {
	//	for _, payment := range result.PaymentDetailList {
	//		validateResp.BankInfoList = append(validateResp.BankInfoList, &viewmodel.BankInfo{
	//			BankName:  payment.PaymentSource,
	//			BankValue: payment.PaymentSource,
	//		})
	//	}
	//}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = validateResp
	return ctx.WriteJson(response)
}

// RefundSubmit 用户退款提交
func RefundSubmit(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	if user == nil {
		response.RetCode = -1
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}
	submitData := agent.RefundSubmitData{}
	err := ctx.Bind(&submitData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "提交的参数不正确"
		return ctx.WriteJson(response)
	}
	if len(submitData.OrderId) == 0 {
		response.RetCode = -3
		response.RetMsg = "订单号不正确"
		return ctx.WriteJson(response)
	}

	if submitData.RefundMode <= 0 {
		response.RetCode = -4
		response.RetMsg = "退款方式不正确"
		return ctx.WriteJson(response)
	}

	if submitData.RefundMode == _const.RefundMode_Bank {
		if len(submitData.BankValue) == 0 {
			response.RetCode = -5
			response.RetMsg = "银行信息不正确"
			return ctx.WriteJson(response)
		}
		if len(submitData.Name) == 0 {
			response.RetCode = -6
			response.RetMsg = "用户姓名不正确"
			return ctx.WriteJson(response)
		}
		if len(submitData.BankAccount) == 0 {
			response.RetCode = -7
			response.RetMsg = "银行账号不正确"
			return ctx.WriteJson(response)
		}
		if len(submitData.BankDetail) == 0 {
			response.RetCode = -8
			response.RetMsg = "开户行信息不正确"
			return ctx.WriteJson(response)
		}
	}
	submitData.Account = user.Account
	submitData.Mobile = user.MobileX
	submitData.RefundAppId = _const.RefundAppId_Roboadvior
	result, err := agent.RefundSubmit(submitData)
	if err != nil {
		response.RetCode = -9
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	return ctx.WriteJson(response)
}

// GetOrderInfo 查询订单信息
func GetOrderInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	if user == nil {
		response.RetCode = -1
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}

	orderId := ctx.PostFormValue("OrderId")
	if len(orderId) == 0 {
		response.RetCode = -2
		response.RetMsg = "订单号不正确"
		return ctx.WriteJson(response)
	}
	orderList, err := agent.QueryOrderList(*user)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = "订单查询异常"
		global.InnerLogger.ErrorFormat(err, "订单查询异常 uid=%s", user.UID)
		return ctx.WriteJson(response)
	}
	if orderList == nil || len(orderList) == 0 {
		response.RetCode = -4
		response.RetMsg = "未查询到订单信息"
		return ctx.WriteJson(response)
	}
	var orderInfo *agent.RefundOrderInfo
	for _, item := range orderList {
		if item.OrderId == orderId {
			orderInfo = item
		}
	}
	if orderInfo == nil {
		response.RetCode = -5
		response.RetMsg = "未查询到订单信息"
		return ctx.WriteJson(response)
	}

	result, err := agent.GetRefundStatus(orderId)
	if err != nil {
		response.RetCode = -6
		response.RetMsg = "退款状态查询异常"
		global.InnerLogger.ErrorFormat(err, "退款状态查询异常 uid=%s orderId=%s", user.UID, orderId)
		return ctx.WriteJson(response)
	}
	refundOrderInfo := viewmodel.RefundOrderInfoDetail{}
	refundOrderInfo.OrderId = orderId
	refundOrderInfo.ProductList = orderInfo.ProductList

	if result != nil {
		for _, status := range result {
			detail := &viewmodel.RefundDetail{
				OrderId:          orderId,
				ReturnModeName:   status.ReturnModeName,
				RefundTypeDesc:   status.RefundType,
				RealbackPrice:    status.RealbackPrice,
				RefundStatusDesc: scmapi.GetRefundStateDesc(status.RefundState),
				RefundTime:       status.CreateTime,
			}
			refundOrderInfo.RefundDetailList = append(refundOrderInfo.RefundDetailList, detail)
		}
	}

	//设置退款时间
	if refundOrderInfo.ProductList != nil && len(refundOrderInfo.RefundDetailList) > 0 {
		for _, item := range refundOrderInfo.ProductList {
			item.RefundTime = refundOrderInfo.RefundDetailList[0].RefundTime
		}
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = refundOrderInfo
	return ctx.WriteJson(response)
}
