package myoptional

// 策略入选消息
type StrategyStockTalk struct {
	//发布人昵称
	NickName string
	//股票代码
	StockCode string
	//股票名称
	StockName string
	//入选策略描述集合
	StrategyDescList []string
}
