package viewmodel

import "time"

type LiveQuestionAnswerList []*LiveQuestionAnswer

func (list LiveQuestionAnswerList) Len() int {
	return len(list)
}

func (list LiveQuestionAnswerList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list LiveQuestionAnswerList) Less(i, j int) bool {
	return time.Time(list[j].AskTime).After(time.Time(list[i].AskTime))
}
