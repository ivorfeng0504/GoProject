package bk1minute

import (
	"github.com/devfeel/mapper"
	"html/template"
)

type BK1MinutesInfo struct{

	//唯一标识
	ID int

	//板块code
	BK_Code string

	//板块name
	BK_Name string

	//关键逻辑
	BK_KeyLogic string

	//关键逻辑
	BK_KeyLogicHTML template.HTML

	//最新动力
	BK_Impetus string

	//最新动力
	BK_ImpetusHTML template.HTML

	//上游
	UpperBK string

	UpperBKList []*BKInfo

	//下游
	LowerBK string

	LowerBKList []*BKInfo

	//行业内容
	IndustryContent string

	IndustryContentList []*IndustryInfo

	//是否删除
	IsDeleted bool

	//创建用户
	CreateUser string

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime
}

type BKInfo struct {
	BKName string
}

type IndustryInfo struct {
	StockDesc string
	StockName string
}