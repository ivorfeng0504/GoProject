package yqq

//益圈圈最新的统计信息
type YqqStatRet struct{
	RetCode string `json:RetCode,string`
	RetMsg string `json:RetMsg,string`
	Data	Stats 	`json:Data`
}

type Stats struct{
	LiveRoomNum int `json:LiveRoomNum,int`
	ContentNum int `json:ContentNum,int`
	QuestionNum int `json:QuestionNum,int`
	UserNum int `json:UserNum,int`
}