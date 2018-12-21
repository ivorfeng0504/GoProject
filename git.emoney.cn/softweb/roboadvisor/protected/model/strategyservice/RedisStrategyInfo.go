package strategyservice

type RedisStrategyInfo struct{
	StrategyGroupId int `json:"ID",int`//策略组ID
	StrategyGroupName string `json:"name",string`//策略组名称
	StrategyList []Strategy //策略列表
	StrategyGroup []InnerIds `json:"StrategyGroup",string`
}

type Strategy struct {
	StrategyId int  `json:"ID",int`//策略ID
	StrategyName string `json:"name",string`//`json:ID,int`策略名称
}


type InnerIds struct{
	StrategyId int `json:"id",int`
}