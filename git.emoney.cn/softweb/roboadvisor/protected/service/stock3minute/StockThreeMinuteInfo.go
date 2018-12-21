package stock3minute

// StockThreeMinuteInfo 个股三分钟信息
type StockThreeMinuteInfo struct {
	StockCode string
	StockName string
	//公司概要
	OverviewValue string
	//市值（1巨盘、2大盘、3中盘、4小盘、5微盘）
	TypeCap string
	//风格（1成长、2价值、3周期、4题材、5高价值）
	TypeStyle string
	//生命周期（单选，例：1）；（1初创期、2成长期、3成熟期、4衰退期）
	LifecycleValue string
	//估值（分值1~5，>=4 偏低， ＜=2 偏高，其他适中）
	ScoreTTM string
	//成长指数（分值1~5，>=4 高， ＜=2 低，其他中)
	ScoreGrowing string
	//盈利能力（分值1~5，>=4 高， ＜=2 低，其他中）
	ScoreProfit string
}
