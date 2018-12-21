package config

import (
	"encoding/xml"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/core"
	"git.emoney.cn/softweb/roboadvisor/util/file"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"io/ioutil"
	"strconv"
)

var (
	CurrentConfig  *AppConfig
	CurrentBaseDir string
	innerLogger    dotlog.Logger
	appSetMap      *core.CMap
	allowIPMap     *core.CMap
	redisMap       *core.CMap
	databaseMap    *core.CMap
)

func SetBaseDir(baseDir string) {
	CurrentBaseDir = baseDir
}

//初始化配置文件
func InitConfig(configFile string) *AppConfig {
	innerLogger = dotlog.GetLogger(_const.LoggerName_Inner)
	CurrentBaseDir = _file.GetCurrentDirectory()
	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]开始...")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]无法解析 - " + err.Error())
		panic(err)
	}

	var result AppConfig
	err = xml.Unmarshal(content, &result)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		panic(err)
	}

	//init config base
	CurrentConfig = &result

	//init AppConfig
	innerLogger.Info("AppConfig::InitConfig Load AppSet Start")
	tmpAppSetMap := core.NewCMap()
	for _, v := range result.AppSets {
		tmpAppSetMap.Set(v.Key, v.Value)
		innerLogger.Info("AppConfig::InitConfig Load AppSet => " + _json.GetJsonString(&v))
	}
	appSetMap = tmpAppSetMap
	innerLogger.Info("AppConfig::InitConfig Load AppSet Finished [" + strconv.Itoa(appSetMap.Len()) + "]")

	//init redisConfig
	innerLogger.Info("AppConfig::InitConfig Start Load RedisInfo")
	tmpRedisMap := core.NewCMap()
	for k, v := range result.Redises {
		tmpRedisMap.Set(v.ID, result.Redises[k])
		innerLogger.Info("AppConfig::InitConfig Load RedisInfo => " + _json.GetJsonString(v))
	}
	redisMap = tmpRedisMap
	innerLogger.Info("AppConfig::InitConfig Finish Load RedisInfo")

	//init databaseConfig
	innerLogger.Info("AppConfig::InitConfig Start Load DataBaseInfo")
	temDataBaseMap := core.NewCMap()
	for k, v := range result.Databases {
		temDataBaseMap.Set(v.ID, result.Databases[k])
		innerLogger.Info("AppConfig::InitConfig Load DataBaseInfo => " + _json.GetJsonString(v))
	}
	databaseMap = temDataBaseMap
	innerLogger.Info("AppConfig::InitConfig Finish Load DataBaseInfo")

	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]完成")

	CurrentConfig.ConfigPath = GetAppConfig("ConfigPath")
	CurrentConfig.StaticServerHost = GetAppConfig("StaticServerHost")
	CurrentConfig.StaticResEnv = GetAppConfig("StaticResEnv")
	CurrentConfig.ReadDB = GetAppConfig("ReadDB")
	CurrentConfig.ReadDB_UserHome = GetAppConfig("ReadDB_UserHome")
	CurrentConfig.ReadDB_MyOptional = GetAppConfig("ReadDB_MyOptional")
	CurrentConfig.SSOUrl = GetAppConfig("SSOUrl")
	CurrentConfig.WebApiHost = GetAppConfig("WebApiHost")
	CurrentConfig.StockTalkWebApiHost = GetAppConfig("StockTalkWebApiHost")
	CurrentConfig.LiveWebApiHost = GetAppConfig("LiveWebApiHost")
	CurrentConfig.DefaultRoomList = GetAppConfig("DefaultRoomList")
	CurrentConfig.DefaultAskRoom = GetAppConfig("DefaultAskRoom")
	CurrentConfig.TokenCreateUrl = GetAppConfig("TokenCreateUrl")
	CurrentConfig.TokenQueryUrl = GetAppConfig("TokenQueryUrl")
	CurrentConfig.TokenVerifyUrl = GetAppConfig("TokenVerifyUrl")
	CurrentConfig.AppId = GetAppConfig("AppId")
	CurrentConfig.EncryptMobileApi = GetAppConfig("EncryptMobileApi")
	CurrentConfig.ResourceVersion = GetAppConfig("ResourceVersion")
	CurrentConfig.DBSyncApi = GetAppConfig("DBSyncApi")
	CurrentConfig.ServerVirtualPath = GetAppConfig("ServerVirtualPath")
	CurrentConfig.SendMsgApi = GetAppConfig("SendMsgApi")
	CurrentConfig.ResourcePath = GetAppConfig("ResourcePath")
	CurrentConfig.DisabledHtmlCache = false
	disabledHtmlCacheStr := GetAppConfig("DisabledHtmlCache")

	if disabledHtmlCacheStr == "1" {
		CurrentConfig.DisabledHtmlCache = true
	}
	CurrentConfig.RsaDecryptUrl = GetAppConfig("RsaDecryptUrl")

	CurrentConfig.MobileProductApi = GetAppConfig("MobileProductApi")
	CurrentConfig.MarketPriceApi = GetAppConfig("MarketPriceApi")
	CurrentConfig.BindAccountSelectApi = GetAppConfig("BindAccountSelectApi")
	CurrentConfig.BindAccountApi = GetAppConfig("BindAccountApi")
	activityGuessNumber := GetAppConfig("Activity_GuessNumber")
	if len(activityGuessNumber) > 0 {
		CurrentConfig.Activity_GuessNumber, _ = strconv.ParseInt(activityGuessNumber, 10, 64)
	}
	activityReceiveStock := GetAppConfig("Activity_ReceiveStock")
	if len(activityReceiveStock) > 0 {
		CurrentConfig.Activity_ReceiveStock, _ = strconv.ParseInt(activityReceiveStock, 10, 64)
	}
	activityGuessChange := GetAppConfig("Activity_GuessChange")
	if len(activityGuessChange) > 0 {
		CurrentConfig.Activity_GuessChange, _ = strconv.ParseInt(activityGuessChange, 10, 64)
	}
	CurrentConfig.SendTeQuanApi = GetAppConfig("SendTeQuanApi")
	CurrentConfig.CommonPoolApi = GetAppConfig("CommonPoolApi")
	CurrentConfig.OnTradingPoolApi = GetAppConfig("OnTradingPoolApi")
	CurrentConfig.ShuangXiangPaoStrategyKey = GetAppConfig("ShuangXiangPaoStrategyKey")
	CurrentConfig.StockTableApi = GetAppConfig("StockTableApi")
	CurrentConfig.StockQUOTESUrl = GetAppConfig("StockQUOTESUrl")
	CurrentConfig.StockPEUrl = GetAppConfig("StockPEUrl")
	CurrentConfig.NewsInformationApi = GetAppConfig("NewsInformationApi")
	CurrentConfig.NewsInfomationTemplateApi = GetAppConfig("NewsInfomationTemplateApi")
	CurrentConfig.ColumnID = GetAppConfig("ColumnID")
	CurrentConfig.GetLiveRoomListUrl = GetAppConfig("GetLiveRoomListUrl")
	CurrentConfig.GetRecommandLiveRoomListUrl = GetAppConfig("GetRecommandLiveRoomListUrl")
	CurrentConfig.GetSilkBagListUrl = GetAppConfig("GetSilkBagListUrl")
	CurrentConfig.GetExpertLiveDataUrl = GetAppConfig("GetExpertLiveDataUrl")
	CurrentConfig.GetTagLiveRoomInfoUrl = GetAppConfig("GetTagLiveRoomInfoUrl")
	CurrentConfig.GetScrollLiveRoomListUrl = GetAppConfig("GetScrollLiveRoomListUrl")
	CurrentConfig.GetYqqHomeDataUrl = GetAppConfig("GetYqqHomeDataUrl")
	CurrentConfig.SCMAPI_QueryOrderProdListByParamsApi = GetAppConfig("SCMAPI_QueryOrderProdListByParamsApi")
	CurrentConfig.SCMAPI_ReturnbackAndRefundApi = GetAppConfig("SCMAPI_ReturnbackAndRefundApi")
	CurrentConfig.SCMAPI_ValidateReturnbackAndRefundApi = GetAppConfig("SCMAPI_ValidateReturnbackAndRefundApi")
	CurrentConfig.SCMAPI_GetRefundStatusApi = GetAppConfig("SCMAPI_GetRefundStatusApi")
	CurrentConfig.GetAccountProfileApiUrl = GetAppConfig("GetAccountProfileApiUrl")
	CurrentConfig.SaveAccountProfileNewApiUrl = GetAppConfig("SaveAccountProfileNewApiUrl")
	CurrentConfig.SaveAccountPicOrNameApiUrl = GetAppConfig("SaveAccountPicOrNameApiUrl")
	CurrentConfig.StrategyServiceAppId = GetAppConfig("StrategyServiceAppId")
	CurrentConfig.ResourceConfigAPI = GetAppConfig("ResourceConfigAPI")
	CurrentConfig.ConfigKey_StrategyLiveRoom = GetAppConfig("ConfigKey_StrategyLiveRoom")
	CurrentConfig.UserHome_Headportrait_UrlFormat = GetAppConfig("UserHome_Headportrait_UrlFormat")

	//百度OCR文字识别API相关配置
	CurrentConfig.Baidu_OCR_APIKey = GetAppConfig("Baidu_OCR_APIKey")
	CurrentConfig.Baidu_OCR_SecretKey = GetAppConfig("Baidu_OCR_SecretKey")
	CurrentConfig.Baidu_OCR_AccessToken_Rul = GetAppConfig("Baidu_OCR_AccessToken")
	CurrentConfig.Baidu_OCR_Accurate_Url = GetAppConfig("Baidu_OCR_Accurate")
	CurrentConfig.Baidu_OCR_General_Url = GetAppConfig("Baidu_OCR_General")

	//自选股云同步地址
	CurrentConfig.MyStockSynURL = GetAppConfig("MyStockSynURL")

	//客户端指标等数据同步接口mysql地址
	CurrentConfig.SoftDataSyncMysqlURL = GetAppConfig("SoftDataSyncMysqlURL")

	//获取手机号密码api接口
	CurrentConfig.GetMobilePwdApiUrl = GetAppConfig("GetMobilePwdApiUrl")

	//关注益圈圈直播接口url
	CurrentConfig.FocusLiveUrl = GetAppConfig("FocusLiveUrl")

	//客户端策略ID和专家资讯策略对应关系
	CurrentConfig.ClientExpertStrategyRelation = GetAppConfig("ClientExpertStrategyRelation")

	//个股微股吧页面地址
	CurrentConfig.StockTalkPageUrl = GetAppConfig("StockTalkPageUrl")

	//获取用户所在地区api地址
	CurrentConfig.ServiceAgentNameApi = GetAppConfig("ServiceAgentNameApi")

	//重要提示栏目ID
	CurrentConfig.ColumnIDImportantTips = GetAppConfig("ColumnIDImportantTips")

	//客户端策略和专家策略对应关系配置Key
	CurrentConfig.ConfigKey_ClientExpertStrategyRelation = GetAppConfig("ConfigKey_ClientExpertStrategyRelation")

	CurrentConfig.StockThreeMinuteAPI_GetStockListInfo = GetAppConfig("StockThreeMinuteAPI_GetStockListInfo")
	CurrentConfig.StockThreeMinuteAPI_GetStockThreeMinuteInfo = GetAppConfig("StockThreeMinuteAPI_GetStockThreeMinuteInfo")

	//智盈产品PID
	CurrentConfig.ConfigKey_ZYPID = GetAppConfig("ConfigKey_ZYPID")

	//用户培训策略显示配置信息
	CurrentConfig.ConfigKey_TrainClient = GetAppConfig("ConfigKey_TrainClient")

	//用户培训标签显示配置信息
	CurrentConfig.ConfigKey_TrainTag = GetAppConfig("ConfigKey_TrainTag")


	//生产环境下默认为当前目录
	if len(CurrentConfig.ResourcePath) == 0 {
		CurrentConfig.ResourcePath = "./"
	}
	return CurrentConfig
}

func GetAppConfig(key string) string {
	return appSetMap.GetString(key)
}

func GetAppSetMap() *core.CMap {
	return appSetMap
}

func GetRedisInfo(redisID string) (*RedisInfo, bool) {
	info, exists := redisMap.Get(redisID)
	if exists {
		return info.(*RedisInfo), exists
	} else {
		return nil, false
	}
}

func GetDataBaseInfo(databaseId string) (*DataBaseInfo, bool) {
	info, exists := databaseMap.Get(databaseId)
	if exists {
		return info.(*DataBaseInfo), exists
	} else {
		return nil, false
	}
}

//检测IP是否被允许访问
func CheckAllowIP(ip string) bool {
	return allowIPMap.Exists(ip)
}
