package userhome

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"github.com/devfeel/dotweb"
)

// RefreshSeedUserInfoList 刷新所有种子用户信息
func RefreshSeedUserInfoList(ctx dotweb.Context) error {
	response := contract.NewApiResponse()
	seedUserSrv := service.NewSeedUserInfoService()
	total, err := seedUserSrv.RefreshSeedUserInfoList()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.TotalCount = total
	}
	return ctx.WriteJson(response)
}
