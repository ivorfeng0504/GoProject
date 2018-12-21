package Strategy

type Strategy struct {
	Id int

	// 调整日期，格式为：yyyy-MM-dd
	date string

	// 策略代码
	code string

	// 策略名称
	name string

	// 策略风格，数字1~5对应保守到激进五种风格
	style int

	// 策略级别，数字1~3对应同一种风格中风险由低到高的三
	level int

	// 策略简介
	summary string

	// 策略分析
	memo string

	// 策略标签
	labels []string

	// 策略配置
	samples []Sample

	// 背景图片
	bgImg string

	//基金组合涨幅
	LatestInfo LatestInfo
}

type Sample struct {
	// 基金类型
	typeName string `json:type,string`

	// 基金代码
	ticker string

	// 基金名称
	name string

	// 基金占比，已乘以100
	weight float64
}

