package stock3minute

type StockThreeMinuteInfoMain struct {
	//公司概要
	OverviewValue string
	//生命周期（单选，例：1）；（1初创期、2成长期、3成熟期、4衰退期）
	LifecycleValue int
}

func GetLifecycleValueDesc(lifecycleValue int) string {
	switch lifecycleValue {
	case 1:
		return "初创期"
	case 2:
		return "成长期"
	case 3:
		return "成熟期"
	case 4:
		return "衰退期"
	case 5:
		return "周期底部"
	case 6:
		return "周期顶部"
	case 7:
		return "周期向下"
	case 8:
		return "周期向上"
	}
	return ""
}
