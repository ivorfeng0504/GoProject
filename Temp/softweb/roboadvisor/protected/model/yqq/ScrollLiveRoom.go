package yqq

//最新的直播信息列表
type ScrollLiveRoomRet struct{
	RetCode string `json:RetCode,string`
	RetMsg string `json:RetMsg,string`
	Data []Item
}

type Item struct{
	Room RoomInfo
	Live LiveInfo
	VisitAddrNum int `json:VisitAddrNum,int`
}

type RoomInfo struct{
	Id int `json:Id,int`
	LiveName string `json:LiveName,string`
	LiveImg string `json:LiveImg,string`
	LiveIntro string `json:LiveIntro,string`
	CreateTime string `json:CreateTime,string`
	LastModifyTime string `json:LastModifyTime,string`
	TagNameStr string `json:TagNameStr,string`
	TagIdStr string `json:TagIdStr,string`
	FansNum int `json:FansNum,int`
	VisitAddrNum int `json:VisitAddrNum,int`
	IsVIP int `json:IsVIP,int`
	HasVIP int `json:hasVIP,int`
	IsSoftAuth int `json:IsSoftAuth,int`
}


type LiveInfo struct{
	Id int `json:Id,int`
	Content string `json:Content,string`
	ContentImg string `json:ContentImg,string`
	CreateTime string `json:CreateTime,string`
	LastModifyTime string `json:LastModifyTime,string`
}