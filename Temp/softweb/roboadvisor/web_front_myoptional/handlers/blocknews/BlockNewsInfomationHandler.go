package blocknews

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"github.com/devfeel/dotweb"
	"strconv"
)

// GetBlockNewsInfomation 获取板块相关资讯
func GetBlockNewsInfomation(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	blockCode := _http.GetRequestValue(ctx, "BlockCode")
	pageSizeStr := _http.GetRequestValue(ctx, "PageSize")
	if len(blockCode) < 4 {
		response.RetCode = -1
		response.RetMsg = "板块代码不正确"
		return ctx.WriteJson(response)
	}
	//板块代码可能是2006246或者BK6246
	if len(blockCode) > 4 {
		blockCode = blockCode[len(blockCode)-4:]
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 0
	}
	newsList, err := agent.GetBlockNewsInfomation(blockCode, pageSize)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	return ctx.WriteJson(response)
}
