package strategyservice

//客户端策略信息
type ClientStrategyInfo struct {
	//栏目Id
	ColumnInfoId int
	//策略Id
	ClientStrategyId int
	//策略名称
	ClientStrategyName string
	//父级Id
	ParentId int
	//父级名称
	ParentName string
}
