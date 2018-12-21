package config

import (
	"encoding/xml"
)

//配置信息
type AppConfig struct {
	XMLName    xml.Name        `xml:"config"`
	AppSets    []*AppSet       `xml:"appsettings>add"`
	Redises    []*RedisInfo    `xml:"redises>redis"`
	Databases  []*DataBaseInfo `xml:"databases>database"`
	AllowIps   []string        `xml:"allowips>ip"`
	ConfigPath string
	// StaticServerHost 静态资源域名占位符，域名末尾不包含/
	StaticServerHost string
	// StaticResEnv 静态资源环境配置后缀占位符
	StaticResEnv string
	// ResourceVersion 静态资源版本号
	ResourceVersion string
	// ReadDB 是否从数据库读取数据，0则只从Redis读取数据，1则从数据库读取，2则先从缓存读取，读取不到从数据库进行读取，3则先从缓存读取，读取不到从数据库进行读取，并更新缓存
	ReadDB string
	// ReadDB_UserHome 是否从数据库读取数据，0则只从Redis读取数据，1则从数据库读取，2则先从缓存读取，读取不到从数据库进行读取，3则先从缓存读取，读取不到从数据库进行读取，并更新缓存
	ReadDB_UserHome string
	// ReadDB_MyOptional 是否从数据库读取数据，0则只从Redis读取数据，1则从数据库读取，2则先从缓存读取，读取不到从数据库进行读取，3则先从缓存读取，读取不到从数据库进行读取，并更新缓存
	ReadDB_MyOptional string
	// SSOUrl SSO接口地址
	SSOUrl string
	//WebApi主机地址  如http://127.0.0.1
	WebApiHost string
	//StockTalk主机地址  如http://127.0.0.1:8091
	StockTalkWebApiHost string
	//图文直播服务接口主机地址  如http://127.0.0.1:8092
	LiveWebApiHost string
	//默认直播间权限，所有人都拥有该直播间权限，多个直播间权限以,分隔
	DefaultRoomList string
	//默认提问的直播间
	DefaultAskRoom string
	//token创建接口地址
	TokenCreateUrl string
	//token查询接口地址
	TokenQueryUrl string
	//token校验接口地址
	TokenVerifyUrl string
	//AppId
	AppId string
	//手机号加密接口
	EncryptMobileApi string
	// 数据库同步接口
	DBSyncApi string
	//服务器虚拟目录 例如/zt
	ServerVirtualPath string
	//发送短信接口地址
	SendMsgApi string
	//Resource目录 以/结尾
	ResourcePath string
	//禁用HTML缓存
	DisabledHtmlCache bool
	//获取加密手机号
	Boundgroupqrylogin string
	//获取用户等级
	QueryUserRiskInfo string
	//获取益圈圈 最新直播内容 接口
	YqqApiNewLiveInfoUrl string
	//获取益圈圈 全部直播内容 接口
	YqqApiAllLiveInfoUrl string
	//获取益圈圈 全部问答内容 接口
	YqqApiAllQuestionUrl string
	//获取益圈圈 我的问答内容 接口
	YqqApiMyQuestionUrl string
}

//AppSetting配置
type AppSet struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

//Redis信息
type RedisInfo struct {
	ID             string `xml:"id,attr"`
	ServerUrl      string `xml:"serverurl,attr"`
	ReadOnlyServer string `xml:"readonlyserver,attr"`
	BackupServer   string `xml:"backupserver,attr"`
	MaxIdle        int    `xml:"maxidle,attr"`
	MaxActive      int    `xml:"maxactive,attr"`
}

//DataBase信息
type DataBaseInfo struct {
	ID        string `xml:"id,attr"`
	ServerUrl string `xml:"serverurl,attr"`
}
