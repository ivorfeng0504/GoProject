package Strategy

type LatestInfo struct {
	Id int

	// 策略代码
	ticker string

	// 最新净值
	nav float64

	// 最新净值日期(yyyy-MM-dd)
	navDate string

	// 净值涨跌额
	navDiff float64

	// 净值涨跌幅
	navChange float64

	// 组合涨幅
	percents  []Percent
}

type Percent struct {
	// 标识：策略代码/""/000300
	ticker string

	// 名称：本产品/比较基准/沪深300
	name string

	// 近一周涨跌幅，已乘以100；下同
	percentW1 float64

	// 近一月涨跌幅
	percentM1 float64

	// 近三月涨跌幅
	percentM3 float64

	// 近六月涨跌幅
	percentM6 float64

	// 今年以来涨跌幅
	percentYtd float64

	// 近一年涨跌幅
	percentY1 float64

	// 近三年涨跌幅
	percentY3 float64
}
