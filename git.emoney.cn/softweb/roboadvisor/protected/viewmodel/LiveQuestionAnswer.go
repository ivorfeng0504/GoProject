package viewmodel

import (
	"github.com/devfeel/mapper"
)

//直播问答
type LiveQuestionAnswer struct {
	AskContent     string
	AskUserName    string
	AskTime        mapper.JSONTime
	AnswerContent  string
	AnswerUserName string
	AnswerTime     mapper.JSONTime
	AskTimeStr     string
	//是否有投资编号等信息
	HasAdviserInfo bool
	//是否已经回答
	IsAnswered bool
	// AdviserName 投资顾问
	AdviserName string

	// AdviserNo 投资编号
	AdviserNo string
}
