package yqq


type YqqRet struct{
	RetCode string `json:RetCode,string`
	RetMsg string `json:RetMsg,string`
	Data []YqqRoom
}

type YqqRetSimple struct{
	RetCode string `json:RetCode,string`
	RetMsg string `json:RetMsg,string`
	Data []YqqRoomSimple
}

type YqqRoom struct {
	Id int `json:Id`
	Fid int `json:-`
	LiveName string  `json:LiveName,string`
	FirstLetter string  `json:-`
	PassportId int  `json:PassportId,string`
	UserName string  `json:UserName,string`
	LiveStatus int `json:LiveStatus,int`
	LiveImg string  `json:LiveImg,string`
	LiveIntro string  `json:LiveIntro,string`
	CreateTime string `json:CreateTime,string`
	LastModifyTime string `json:LastModifyTime,string`
	TagNameStr string  `json:TagNameStr,string`
	TagIdStr string  `json:TagIdStr,string`
	FansNum int `json:FansNum,int`
	BagNum int `json:BagNum,int`
	VisitAddrNum int `json:VisitAddrNum,int`
	IsFocus int `json:-`
	PId int `json:PId,int`
	IsVIP int `json:IsVIP,int`
	hasVIP int `json:hasVIP,int`
	IsHide int `json:IsHide,int`
	IsAwardVIP int `json:-`
	Content string  `json:Content,string`
	NewContentTime string `json:NewContentTime,string`
	ExpiredTime string `json:ExpiredTime,string`
	AdviserName string  `json:AdviserName,string`
	AdviserNo string  `json:AdviserNo,string`
	TopicName string  `json:TopicName,string`
	IsDelete int  `json:IsDelete,int`
	DeleteTime string `json:-`
	IsSoftAuth int 	`json:IsSoftAuth,int`
}

type YqqRoomSimple struct {
	Id int `json:Id`
	LiveName string  `json:LiveName,string`
	LiveStatus int `json:LiveStatus,int`
	LiveImg string  `json:LiveImg,string`
	LiveIntro string  `json:LiveIntro,string`
	CreateTime string `json:CreateTime,string`
	LastModifyTime string `json:LastModifyTime,string`
	TagNameStr string  `json:TagNameStr,string`
	TagIdStr string  `json:TagIdStr,string`
	FansNum int `json:FansNum,int`
	BagNum int `json:BagNum,int`
	VisitAddrNum int `json:VisitAddrNum,int`
	Content string  `json:Content,string`
	NewContentTime string `json:NewContentTime,string`
	AdviserName string  `json:AdviserName,string`
	AdviserNo string  `json:AdviserNo,string`
	TopicName string  `json:TopicName,string`
	IsSoftAuth int 	`json:IsSoftAuth,int`
}