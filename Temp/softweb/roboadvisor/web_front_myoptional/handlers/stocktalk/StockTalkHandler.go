package stocktalk

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/myoptional"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"github.com/devfeel/dotweb"
	"strconv"
	"strings"
)

// InsertStockTalk 插入一个微股吧评论
func InsertStockTalk(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_myoptional.UserInfo(ctx)
	stockCode := ctx.PostFormValue("StockCode")
	if len(stockCode) == 0 {
		response.RetCode = -1
		response.RetMsg = "股票代码不能为空"
		return ctx.WriteJson(response)
	}

	stockName := ctx.PostFormValue("StockName")
	if len(stockName) == 0 {
		response.RetCode = -2
		response.RetMsg = "股票名称不能为空"
		return ctx.WriteJson(response)
	}

	content := ctx.PostFormValue("Content")
	if len(content) == 0 {
		response.RetCode = -3
		response.RetMsg = "评论内容不能为空"
		return ctx.WriteJson(response)
	}

	requestData := myoptional_model.StockTalk{
		Content:   content,
		StockCode: stockCode,
		StockName: stockName,
		NickName:  user.NickName,
		UID:       strconv.FormatInt(user.UID, 10),
	}
	err := agent.InsertStockTalk(requestData)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
	}
	return ctx.WriteJson(response)
}

// GetStockTalkListByDate 获取指定日期的微股吧评论
func GetStockTalkListByDate(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	dateStr := _http.GetRequestValue(ctx, "Date")
	incre := _http.GetRequestValue(ctx, "Incre")
	IsIncre, _ := strconv.ParseBool(incre)
	list, err := agent.GetStockTalkListByDate(IsIncre, dateStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = list
	}
	return ctx.WriteJson(response)
}

// GetStockTalkListByPage 分页获取微股吧评论
func GetStockTalkListByPage(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	stockCodeListStr := _http.GetRequestValue(ctx, "StockCodeList")
	pageIndexStr := _http.GetRequestValue(ctx, "PageIndex")
	pageSizeStr := _http.GetRequestValue(ctx, "PageSize")
	pageIndex, _ := strconv.Atoi(pageIndexStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	list, err := agent.GetStockTalkListByPage(strings.Split(stockCodeListStr, ","), pageIndex, pageSize)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = list
	}
	return ctx.WriteJson(response)
}
