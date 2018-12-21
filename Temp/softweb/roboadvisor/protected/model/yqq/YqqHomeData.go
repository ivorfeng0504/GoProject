package yqq

//益圈圈首页数据
type YqqHomeData struct{
	RetCode string `json:RetCode,string`
	RetMsg string `json:RetMsg,string`
	Data HomeData `json:Data`
}

type HomeData struct{
	VipList []LiveRoom
	HotList []LiveRoom
	ScrollList []LiveContent
}

type LiveRoom struct{
	ColumnName string `json:ColumnName,string`
	LId int `json:LId,int`
	Sort int `json:Sort,int`
	ColumnSort int `json:ColumnSort,int`
	LiveName string `json:LiveName,string`
	LiveImg string `json:LiveImg,string`
	LiveIntro string `json:LiveIntro,string`
}

type LiveContent struct{
	LId int `json:LId,int`
	LiveName string `json:LiveName,string`
	LiveImg string `json:LiveImg,string`
	LiveNewContent string `json:LiveNewContent,string`
	NewContentTime string `json:NewContentTime,string`
}