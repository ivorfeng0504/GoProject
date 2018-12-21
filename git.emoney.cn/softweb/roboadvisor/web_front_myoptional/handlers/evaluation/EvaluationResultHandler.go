package evaluation

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/myoptional"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"strconv"
)

// Index 测评首页
func Index(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "invest_type_testing.html")
}

func DiscernPic(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "discernpic.html")
}

func Upfile(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "discern.html")
}

// GoTest 测评弹窗
func GoTest(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "zixuan_testing.html")
}

// GetResult 查看当前用户的评测结果，如果Message为null，则未参加测评
func GetResult(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_myoptional.UserInfo(ctx)
	uid := strconv.FormatInt(user.UID, 10)
	result, err := agent.GetEvaluationResult(uid)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "查看当前用户的评测结果异常，请稍后再试！"
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	return ctx.WriteJson(response)
}

// GetResultForClient 查看当前用户的评测结果 客户端调用接口
func GetResultForClient(ctx dotweb.Context) error {
	uid := ctx.QueryString("uid")
	result, err := agent.GetEvaluationResult(uid)
	if err != nil {
		//异常
		return ctx.WriteString("{|*-1*|}")
	} else {
		if result == nil {
			//未测评
			return ctx.WriteString("{|*0*|}")
		} else {
			//已测评
			return ctx.WriteString("{|*1*|}")
		}
	}
}

// SubmitResult 提交评测
func SubmitResult(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_myoptional.UserInfo(ctx)
	result := myoptional_model.EvaluationResult{
		InvestTarget:          ctx.FormValue("InvestTarget"),
		InvestTargetDesc:      ctx.FormValue("InvestTargetDesc"),
		InvestTargetTip:       ctx.FormValue("InvestTargetTip"),
		ChooseStockReason:     ctx.FormValue("ChooseStockReason"),
		ChooseStockReasonDesc: ctx.FormValue("ChooseStockReasonDesc"),
		ChooseStockReasonTip:  ctx.FormValue("ChooseStockReasonTip"),
		HoldStockTime:         ctx.FormValue("HoldStockTime"),
		HoldStockTimeDesc:     ctx.FormValue("HoldStockTimeDesc"),
		HoldStockTimeTip:      ctx.FormValue("HoldStockTimeTip"),
		BuyStyle:              ctx.FormValue("BuyStyle"),
		BuyStyleDesc:          ctx.FormValue("BuyStyleDesc"),
		BuyStyleTip:           ctx.FormValue("BuyStyleTip"),
		Result:                ctx.FormValue("Result"),
		ResultDesc:            ctx.FormValue("ResultDesc"),
	}
	result.UID = strconv.FormatInt(user.UID, 10)
	err := agent.SubmitResult(result)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "提交评测异常，请稍后再试！"
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// AnalyzeStock 上传图片，识别股票
func AnalyzeStock(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	upfile, err := ctx.Request().FormFile("StockImage")
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "上传图片异常，请稍后再试！"
		return ctx.WriteHtml(_json.GetJsonString(response))
	}

	//匹配股票代码的股票名称
	stockList, err := agent.AnalyzeStock(upfile.File)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "识别股票异常 " + err.Error()
		return ctx.WriteHtml(_json.GetJsonString(response))
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = stockList
	return ctx.WriteHtml(_json.GetJsonString(response))
}

// GetStrategyList 获取推荐策略及股票列表
func GetStrategyList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	strategyList, err := agent.GetStrategyList()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取推荐策略及股票列表异常，请稍后再试！"
		return ctx.WriteJson(response)
	} else {
		//股票取前三个
		if strategyList != nil {
			for _, strategy := range strategyList {
				if strategy.StockList != nil && len(strategy.StockList) > 3 {
					strategy.StockList = strategy.StockList[:3]
				}
			}
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyList
	return ctx.WriteJson(response)
}
