package yqq

type YqqMyQuestionResult struct {
	RetCode string
	RetMsg  string
	Data []*LiveQuestionAnswer
}

type LiveQuestionAnswer struct {
	AskContent     string
	AskUserName    string
	AnswerContent  string
	AnswerUserName string
	// 投资顾问名称
	MasterName string
	// 投资顾问执业证书编号
	MasterNumber string
}




