package live

import (
	"github.com/devfeel/mapper"
)

//直播问答
type LiveQuestionAnswer struct {
	Id                int
	LId               int
	Status            int
	AskContent        string
	AskUserId         int
	AskPassportId     int64
	AskUserName       string
	AskTime           mapper.JSONTime
	AnswerContent     string
	AnswerUserId      int
	AnswerAdminUserId int64
	AnswerPassportId  int64
	AnswerUserName    string
	AnswerTime        mapper.JSONTime
	IsRecommend       int
	IsEffectReward    int
	LiveName          string
	IsFromMobile      int
	IsPrivate         int
	ProductId         string

	// AdviserName 投资顾问
	AdviserName string

	// AdviserNo 投资编号
	AdviserNo string
}
