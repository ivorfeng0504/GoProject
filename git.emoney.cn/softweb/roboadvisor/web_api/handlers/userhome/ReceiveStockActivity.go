package userhome

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"github.com/devfeel/dotweb"
)

// GetNewstStockPool 获取最新的股票池
func GetNewstStockPool(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}

	srv := service.NewReceiveStockActivityService()
	stock, err := srv.GetNewstStockPool()
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = stock
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
