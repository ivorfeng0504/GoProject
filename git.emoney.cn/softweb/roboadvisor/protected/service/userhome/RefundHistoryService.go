package service

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/rsa"
	"git.emoney.cn/softweb/roboadvisor/protected/service/scmapi"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type RefundHistoryService struct {
	service.BaseService
	refundHistoryRepo *userhome_repo.RefundHistoryRepository
}

var (
	shareRefundHistoryRepo   *userhome_repo.RefundHistoryRepository
	shareRefundHistoryLogger dotlog.Logger
)

const (
	UserHomeRefundHistoryServiceName = "RefundHistoryService"
)

func NewRefundHistoryService() *RefundHistoryService {
	srv := &RefundHistoryService{
		refundHistoryRepo: shareRefundHistoryRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeRefundHistoryServiceName, refundHistoryServiceLoader)
}

func refundHistoryServiceLoader() {
	shareRefundHistoryRepo = userhome_repo.NewRefundHistoryRepository(protected.DefaultConfig)
	shareRefundHistoryLogger = dotlog.GetLogger(UserHomeRefundHistoryServiceName)
}

// InsertRefundHistory 新增一条退款记录
func (srv *RefundHistoryService) InsertRefundHistory(model *userhome_model.RefundHistory) (id int64, err error) {
	id, err = srv.refundHistoryRepo.InsertRefundHistory(model)
	if err != nil {
		shareRefundHistoryLogger.ErrorFormat(err, "InsertRefundHistory 新增一条退款记录 异常 model=%s", _json.GetJsonString(model))
	}
	return id, err
}

// QueryOrderList 订单查询
func (srv *RefundHistoryService) QueryOrderList(request scmapi.QueryOrderProdListByParamsRequest) (orderList []*agent.RefundOrderInfo, err error) {
	orderListSrc, err := scmapi.QueryOrderProdListByParams(request)
	if err != nil {
		shareRefundHistoryLogger.ErrorFormat(err, "QueryOrderList 订单查询 异常 request=%s", _json.GetJsonString(request))
		return nil, err
	}
	if orderListSrc == nil || len(orderListSrc) == 0 {
		return nil, nil
	}
	orderMap := make(map[string]int)
	for _, orderInfoSrc := range orderListSrc {
		orderIndex, success := orderMap[orderInfoSrc.OrderId]
		orderInfo := &agent.RefundOrderInfo{
			CanRefund: true,
		}
		if success {
			//订单已经存在
			orderInfo = orderList[orderIndex]
		} else {
			//新的订单
			orderInfo.OrderId = orderInfoSrc.OrderId
			orderList = append(orderList, orderInfo)
			orderMap[orderInfoSrc.OrderId] = len(orderList) - 1
		}
		orderInfo.Price += orderInfoSrc.SPrice
		orderInfo.ProductList = append(orderInfo.ProductList, &agent.RefundProductInfo{
			ProductName:   orderInfoSrc.ProdName,
			Price:         orderInfoSrc.SPrice,
			Count:         1,
			CreateTime:    orderInfoSrc.OrdAddTime,
			State:         orderInfoSrc.RefundSign,
			StateDesc:     scmapi.GetRefundSignDesc(orderInfoSrc.RefundSign),
			IsQuickReturn: orderInfoSrc.IsRefundByCustomer,
		})
		//如果有一个产品不能退款 则整个订单不能退
		if orderInfoSrc.RefundSign != 0 || orderInfoSrc.IsRefundByCustomer != 1 {
			orderInfo.CanRefund = false
		}
		if orderInfoSrc.RefundSign == -1 {
			orderInfo.IsRefund = true
		}
	}
	return orderList, err
}

// ValidateOrder 验证是否支持快速退货退款
func (srv *RefundHistoryService) ValidateOrder(orderId string) (result scmapi.ValidateReturnbackAndRefundResponse, err error) {
	result, err = scmapi.ValidateReturnbackAndRefund(orderId)
	if err != nil {
		shareRefundHistoryLogger.ErrorFormat(err, "ValidateOrder 验证是否支持快速退货退款 异常 orderId=%s", orderId)
	}
	return result, err
}

// RefundSubmit 提交退款
func (srv *RefundHistoryService) RefundSubmit(submitData agent.RefundSubmitData) (result scmapi.ReturnbackAndRefundResponse, err error) {
	request := scmapi.ReturnbackAndRefundRequest{
		OrderId:     submitData.OrderId,
		AppId:       submitData.RefundAppId,
		Reason:      submitData.Reason,
		ReturnFlow:  _const.ReturnFlow_Quick,
		ReturnMode:  submitData.RefundMode,
		IsAllReturn: 1,
	}
	//银行卡退款
	if submitData.RefundMode == _const.RefundMode_Bank {
		bankAccount, err := rsa.DecryptRSA(submitData.BankAccount)
		if err != nil {
			shareRefundHistoryLogger.ErrorFormat(err, "RefundSubmit 提交退款 银行卡解密失败 submitData=%s", _json.GetJsonString(submitData))
			return result, err
		}
		request.OrderRetBankInfo = scmapi.OrderRetBankInfoDo{
			ReProCode:  bankAccount,
			ReProName:  submitData.Name,
			ReBankArea: submitData.BankDetail,
			ReBankName: submitData.BankValue,
		}
	}

	result, err = scmapi.ReturnbackAndRefund(request)
	if err != nil {
		shareRefundHistoryLogger.ErrorFormat(err, "RefundSubmit 提交退款 异常 request=%s", _json.GetJsonString(request))
	} else {
		//数据库存储加密的银行卡号
		request.OrderRetBankInfo.ReProCode = submitData.BankAccount
		history := &userhome_model.RefundHistory{
			OrderId:    request.OrderId,
			Reason:     request.Reason,
			RefundMode: request.ReturnMode,
			SubmitData: _json.GetJsonString(request),
			Account:    submitData.Account,
			Mobile:     submitData.Mobile,
		}
		_, err = srv.InsertRefundHistory(history)
		if err != nil {
			shareRefundHistoryLogger.ErrorFormat(err, "RefundSubmit 提交退款 写入退款历史记录异常 request=%s history=%s", _json.GetJsonString(request), _json.GetJsonString(history))
		}
	}
	return result, err
}

// GetRefundStatus 查询退款状态
func (srv *RefundHistoryService) GetRefundStatus(orderId string) (result []*scmapi.OrderRefundState, err error) {
	result, err = scmapi.GetRefundStatus(orderId)
	if err != nil {
		shareRefundHistoryLogger.ErrorFormat(err, "GetRefundStatus 查询退款状态 异常 orderId=%s", orderId)
	}
	return result, err
}
