package viewmodel

import "time"

type LiveContentList []*LiveContent

func (list LiveContentList) Len() int {
	return len(list)
}

func (list LiveContentList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list LiveContentList) Less(i, j int) bool {
	return time.Time(list[j].CreateTime).After(time.Time(list[i].CreateTime))
}
