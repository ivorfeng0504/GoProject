package myoptional

import "github.com/devfeel/mapper"

//评估结果
type EvaluationResult struct {
	//主键Id
	EvaluationResultId int64
	//用户UID
	UID string
	//创建时间
	CreateTime mapper.JSONTime
	//修改时间
	ModifyTime mapper.JSONTime
	//投资目标
	InvestTarget string
	//投资目标-描述
	InvestTargetDesc string
	//投资目标-提示
	InvestTargetTip string
	//选股主要考虑
	ChooseStockReason string
	//选股主要考虑-描述
	ChooseStockReasonDesc string
	//选股主要考虑-提示
	ChooseStockReasonTip string
	//股票一般拿多久
	HoldStockTime string
	//股票一般拿多久-描述
	HoldStockTimeDesc string
	//股票一般拿多久-提示
	HoldStockTimeTip string
	//操作买点
	BuyStyle string
	//操作买点-描述
	BuyStyleDesc string
	//操作买点-提示
	BuyStyleTip string
	//评测结果
	Result string
	//评测结果-描述
	ResultDesc string
}
