package stock3minute

import (
	"git.emoney.cn/softweb/roboadvisor/util/strings"
	"strings"
)

type StockThreeMinuteInfoBase struct {
	StockCode string
	StockName string
	//市值（1巨盘、2大盘、3中盘、4小盘、5微盘）
	TypeCap int
	//风格 逗号分隔（1成长、2价值、3周期、4题材、5高价值）
	TypeStyle string
	//估值（分值1~5，>=4 偏低， ＜=2 偏高，其他适中）
	ScoreTTM float32
	//成长指数（分值1~5，>=4 高， ＜=2 低，其他中)
	ScoreGrowing float32
	//盈利能力（分值1~5，>=4 高， ＜=2 低，其他中）
	ScoreProfit float32
}

func GetTypeCapDesc(typeCap int) string {
	switch typeCap {
	case 1:
		return "巨盘"
	case 2:
		return "大盘"
	case 3:
		return "中盘"
	case 4:
		return "小盘"
	case 5:
		return "微盘"
	}
	return ""
}

func GetTypeStyleDesc(typeStyle string) string {
	if len(typeStyle) == 0 {
		return ""
	}
	desc := ""
	styleList := strings.Split(typeStyle, ",")
	for _, style := range styleList {
		styleDesc := GetTypeStyleDescCore(style)
		if len(styleDesc) > 0 {
			desc += styleDesc + "、"
		}
	}
	if len(desc) > 0 && _strings.LastString(desc) == "、" {
		desc = _strings.SubString(desc, 0, _strings.StringLen(desc)-1)
	}
	return desc
}

func GetTypeStyleDescCore(typeStyle string) string {
	switch typeStyle {
	case "1":
		return "成长"
	case "2":
		return "价值"
	case "3":
		return "周期"
	case "4":
		return "题材"
	case "5":
		return "高价值"
	}
	return ""
}
func GetScoreTTMDesc(scoreTTM float32) string {
	if scoreTTM >= 4 {
		return "偏低"
	} else if scoreTTM <= 2 {
		return "偏高"
	} else {
		return "适中"
	}
}

func GetScoreGrowingDesc(scoreGrowing float32) string {
	if scoreGrowing >= 4 {
		return "高"
	} else if scoreGrowing <= 2 {
		return "低"
	} else {
		return "中"
	}
}

func GetScoreProfitDesc(scoreProfit float32) string {
	if scoreProfit >= 4 {
		return "高"
	} else if scoreProfit <= 2 {
		return "低"
	} else {
		return "中"
	}
}
