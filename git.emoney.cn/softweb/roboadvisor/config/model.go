package config

import (
	"encoding/xml"
)

const (
	//只从数据库读取
	ReadDB_JustDB = "1"
	//只从缓存读取
	ReadDB_JustCache = "0"
	//先从缓存读，如果没有再从数据库读取,不更新缓存
	ReadDB_CacheOrDB_NotUpdateCache = "2"
	//先从缓存读，如果没有再从数据库读取,并更新缓存
	ReadDB_CacheOrDB_UpdateCache = "3"
	//读取数据库并刷新缓存
	ReadDB_RefreshCache = "4"
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
	//rsa解密接口地址
	RsaDecryptUrl string
	//移动端产品获取接口地址
	MobileProductApi string
	//市场价格行情接口地址
	MarketPriceApi string
	//获取绑定的手机号接口地址
	BindAccountSelectApi string
	//绑定接口
	BindAccountApi string
	//猜数字活动Id
	Activity_GuessNumber int64
	//连登-领双响炮活动
	Activity_ReceiveStock int64
	//猜涨跌活动Id
	Activity_GuessChange int64
	//特权开通API地址
	SendTeQuanApi string
	//股池接口
	CommonPoolApi string
	//盘中预警股池接口
	OnTradingPoolApi string
	//双响炮策略KEY
	ShuangXiangPaoStrategyKey string
	//A股码表接口
	StockTableApi string
	//个股行情接口
	StockQUOTESUrl string
	//个股PE接口
	StockPEUrl string
	//资讯接口
	NewsInformationApi string
	//获取直播间列表
	GetLiveRoomListUrl string
	//资讯模板API接口
	NewsInfomationTemplateApi string
	//资讯栏目ID
	ColumnID string
	//获取VIP和热门直播
	GetRecommandLiveRoomListUrl string
	//获取直播间锦囊
	GetSilkBagListUrl string
	//获取益圈圈直播统计信息
	GetExpertLiveDataUrl string
	//按标签查询直播列表
	GetTagLiveRoomInfoUrl string
	//获取滚动直播列表
	GetScrollLiveRoomListUrl string
	//获取益圈圈首页的数据
	GetYqqHomeDataUrl string
	//物流-查询订单明细接口
	SCMAPI_QueryOrderProdListByParamsApi string
	//物流-快速退货退款接口
	SCMAPI_ReturnbackAndRefundApi string
	//物流-验证是否支持快速退货退款接口
	SCMAPI_ValidateReturnbackAndRefundApi string
	//物流-获取退款状态接口
	SCMAPI_GetRefundStatusApi string
	//获取用户的个人资料
	GetAccountProfileApiUrl string
	//保存用户的个人资料
	SaveAccountProfileNewApiUrl string
	//个人资料头像昵称修改
	SaveAccountPicOrNameApiUrl string
	//策略资讯副窗口AppId
	StrategyServiceAppId string
	//配置资源接口地址
	ResourceConfigAPI string
	//策略对应的直播间关系配置Key
	ConfigKey_StrategyLiveRoom string
	//用户中心用户头像地址
	UserHome_Headportrait_UrlFormat string

	//百度OCR文字识别API（高精度版）
	Baidu_OCR_Accurate_Url string
	//百度OCR文字识别API（普通版）
	Baidu_OCR_General_Url string
	//获取百度API调用的AccessToken
	Baidu_OCR_AccessToken_Rul string
	//百度API_KEY
	Baidu_OCR_APIKey string
	//百度Secret_KEY
	Baidu_OCR_SecretKey string
	// 自选股云同步接口
	MyStockSynURL string
	// 客户端指标数据同步接口mysql地址
	SoftDataSyncMysqlURL string

	//获取手机号密码接口
	GetMobilePwdApiUrl string

	//关注益圈圈直播接口url
	FocusLiveUrl string

	//客户端策略ID和专家资讯策略对应关系
	ClientExpertStrategyRelation string

	//个股微股吧页面地址
	StockTalkPageUrl string

	//获取用户所在服务商地区
	ServiceAgentNameApi string

	//重要提示栏目ID
	ColumnIDImportantTips string

	//客户端策略和专家策略对应关系配置Key
	ConfigKey_ClientExpertStrategyRelation string

	// 个股三分钟 获取股票码表
	StockThreeMinuteAPI_GetStockListInfo string

	//个股三分钟 获取基本面信息
	StockThreeMinuteAPI_GetStockThreeMinuteInfo string

	//智盈产品ID
	ConfigKey_ZYPID string

	//用户培训策略显示配置信息
	ConfigKey_TrainClient string

	//用户培训标签显示配置信息
	ConfigKey_TrainTag string
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
