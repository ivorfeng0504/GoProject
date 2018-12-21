package protected

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"github.com/devfeel/dotlog"
)

type ServiceLoader func()

const (
	DefaultRedisID                  = "default"
	ContentLivePlatRedisID          = "content_live_plat"
	SessionRedisID                  = "session_redis"
	MarketDataRedisID               = "market_redis"
	StrategyInnerOutRedisID         = "strategyout_redis"
	OnTradingMonitorPoolRedisID     = "strategy_ontradingpool_redis"
	AfterTradingStrategyPoolRedisID = "strategy_aftertradingpool_redis"
	UserHomeRedisID                 = "userhome_redis"
	UserHomeStatRedisID             = "userhomestat_redis"
	ContentLivePlatDatabaseID       = "content_live_plat_db"
	EMoneyRoboAdvisorID             = "emoney_roboadvisor_db"
	FeedbDatabaseID                 = "feedb_db"
	ExpertNewsRedisID               = "expert_redis"
	ExpertNewsStatRedisID           = "expertstat_redis"
	CaptchaStoreRedisID             = "captchastore_redis"
	EmoneyDatabaseID                = "emoney_db"
	TJSZTXRedisID                = "tjsztx_redis"
	UserTrainDatabaseID				= "usertrain_db"
)

var (
	DefaultConfig *ServiceConfig
	//默认Redis
	DefaultRedisConfig *config.RedisInfo
	//图文直播专用Redis
	LiveRedisConfig *config.RedisInfo
	//Session会话专员专用Redis
	SessionRedisConfig *config.RedisInfo
	//行情Redis数据
	MarketRedisConfig *config.RedisInfo
	//用户中心Redis数据
	UserHomeRedisConfig *config.RedisInfo
	//用户中心统计Redis数据
	UserHomeStatRedisConfig *config.RedisInfo

	//专家资讯Redis数据
	ExpertNewsRedisConfig *config.RedisInfo
	//专家资讯统计Redis数据
	ExpertNewsStatRedisConfig *config.RedisInfo

	//验证码Redis
	CaptchaStoreRedisConfig *config.RedisInfo

	//软件内部输出的策略Redis
	StrategyInnerOutRedisConfig *config.RedisInfo

	//盘中预警股池redis
	OnTradingMonitorPoolRedisConfig *config.RedisInfo

	//盘后选股策略股池redis
	AfterTradingStrategyPoolRedisConfig *config.RedisInfo

	//绑定解绑通知投教实战体系
	TJSZTXRedisConfig *config.RedisInfo

	serviceLoaderMap map[string]ServiceLoader
	ServiceLogger    dotlog.Logger
)

func init() {
	serviceLoaderMap = make(map[string]ServiceLoader)
}

func InitLogger() {
	ServiceLogger = dotlog.GetLogger(_const.LoggerName_Service)
}

func Init() error {
	InitLogger()
	var err error

	DefaultConfig = &ServiceConfig{}

	//获取ContentLivePlat数据库连接字符串
	liveDbInfo, exists := config.GetDataBaseInfo(ContentLivePlatDatabaseID)
	if !exists || liveDbInfo.ServerUrl == "" {
		err = errors.New("no config " + ContentLivePlatDatabaseID + " database config")
		ServiceLogger.Error(err, "no config "+ContentLivePlatDatabaseID+" database config")
		return err
	}

	DefaultConfig.ContentLivePlatDBConn = liveDbInfo.ServerUrl

	//获取EMoney_RoboAdvisor数据库连接字符串
	roboAdvisor, exists := config.GetDataBaseInfo(EMoneyRoboAdvisorID)
	if !exists || roboAdvisor.ServerUrl == "" {
		err = errors.New("no config " + EMoneyRoboAdvisorID + " database config")
		ServiceLogger.Error(err, "no config "+EMoneyRoboAdvisorID+" database config")
		return err
	}

	DefaultConfig.EMoney_RoboAdvisorDBConn = roboAdvisor.ServerUrl

	//获取feedb数据库连接字符串
	feedb, exists := config.GetDataBaseInfo(FeedbDatabaseID)
	if !exists || feedb.ServerUrl == "" {
		err = errors.New("no config " + FeedbDatabaseID + " database config")
		ServiceLogger.DebugFormat("no config " + FeedbDatabaseID + " database config")
		//return err
	} else {
		DefaultConfig.FeedbDBConn = feedb.ServerUrl
	}

	//获取emoney数据库连接字符串
	emoney, exists := config.GetDataBaseInfo(EmoneyDatabaseID)
	if !exists || emoney.ServerUrl == "" {
		err = errors.New("no config " + EmoneyDatabaseID + " database config")
		ServiceLogger.DebugFormat("no config " + EmoneyDatabaseID + " database config")
		//return err
	} else {
		DefaultConfig.EmoneyDBConn = emoney.ServerUrl
	}

	//获取usertrain数据库连接字符串
	usertrain, exists := config.GetDataBaseInfo(UserTrainDatabaseID)
	if !exists || usertrain.ServerUrl == "" {
		err = errors.New("no config " + UserTrainDatabaseID + " database config")
		ServiceLogger.DebugFormat("no config " + UserTrainDatabaseID + " database config")
		//return err
	} else {
		DefaultConfig.UserTrainDBConn = usertrain.ServerUrl
	}

	//初始化Redis配置信息
	redis, exists := config.GetRedisInfo(DefaultRedisID)
	if !exists {
		err = errors.New("no exists " + DefaultRedisID + " logger config")
		ServiceLogger.Error(err, "not exists "+DefaultRedisID+" logger config")
		return err
	}
	DefaultRedisConfig = redis

	//初始化内容直播Redis配置信息
	liveRedis, exists := config.GetRedisInfo(ContentLivePlatRedisID)
	if !exists {
		err = errors.New("no exists " + ContentLivePlatRedisID + " logger config")
		ServiceLogger.Error(err, "not exists "+ContentLivePlatRedisID+" logger config")
		return err
	}
	LiveRedisConfig = liveRedis

	//初始化SessionRedis配置信息
	sessionRedis, exists := config.GetRedisInfo(SessionRedisID)
	if !exists {
		err = errors.New("no exists " + SessionRedisID + " logger config")
		ServiceLogger.Error(err, "not exists "+SessionRedisID+" logger config")
		return err
	}
	SessionRedisConfig = sessionRedis

	//初始化Market行情Redis信息
	marketRedis, exists := config.GetRedisInfo(MarketDataRedisID)
	if !exists {
		err = errors.New("no exists " + MarketDataRedisID + " logger config")
		ServiceLogger.Error(err, "not exists "+MarketDataRedisID+" logger config")
		return err
	}
	MarketRedisConfig = marketRedis

	//初始化StrategyInnerOutRedis信息
	strategyInnerOutRedis, exists := config.GetRedisInfo(StrategyInnerOutRedisID)
	if !exists {
		err = errors.New("no exists " + StrategyInnerOutRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + StrategyInnerOutRedisID + " logger config")
		//return err
	}
	StrategyInnerOutRedisConfig = strategyInnerOutRedis

	//初始化OnTradingMonitorPoolRedisConfig信息
	onTradingPoolRedis, exists := config.GetRedisInfo(OnTradingMonitorPoolRedisID)
	if !exists {
		err = errors.New("no exists " + StrategyInnerOutRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + StrategyInnerOutRedisID + " logger config")
		//return err
	}
	OnTradingMonitorPoolRedisConfig = onTradingPoolRedis

	//初始化AfterTradingStrategyPoolRedisConfig信息
	afterTradingPoolRedis, exists := config.GetRedisInfo(AfterTradingStrategyPoolRedisID)
	if !exists {
		err = errors.New("no exists " + StrategyInnerOutRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + StrategyInnerOutRedisID + " logger config")
		//return err
	}
	AfterTradingStrategyPoolRedisConfig = afterTradingPoolRedis

	//初始化用户中心Redis信息
	userHomeRedis, exists := config.GetRedisInfo(UserHomeRedisID)
	if !exists {
		err = errors.New("no exists " + UserHomeRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + UserHomeRedisID + " logger config")
		//return err
	}
	UserHomeRedisConfig = userHomeRedis

	//初始化用户中心Stat Redis信息
	userHomeRedisStat, exists := config.GetRedisInfo(UserHomeStatRedisID)
	if !exists {
		err = errors.New("no exists " + UserHomeStatRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + UserHomeStatRedisID + " logger config")
		//return err
	}
	UserHomeStatRedisConfig = userHomeRedisStat

	//初始化专家资讯Redis信息
	expertnewsRedis, exists := config.GetRedisInfo(ExpertNewsRedisID)
	if !exists {
		err = errors.New("no exists " + ExpertNewsRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + ExpertNewsRedisID + " logger config")
		//return err
	}

	ExpertNewsRedisConfig = expertnewsRedis

	//初始化专家资讯Stat Redis信息
	expertnewsRedisStat, exists := config.GetRedisInfo(ExpertNewsStatRedisID)
	if !exists {
		err = errors.New("no exists " + ExpertNewsStatRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + ExpertNewsStatRedisID + " logger config")
		//return err
	}
	ExpertNewsStatRedisConfig = expertnewsRedisStat

	//初始化Captcha验证码StoreRedis
	captchaStoreRedis, exists := config.GetRedisInfo(CaptchaStoreRedisID)
	if !exists {
		err = errors.New("no exists " + ExpertNewsStatRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + ExpertNewsStatRedisID + " logger config")
		//return err
	}
	CaptchaStoreRedisConfig = captchaStoreRedis

	//初始化投教实战体系相关redis
	TJSZTXRedis, exists := config.GetRedisInfo(TJSZTXRedisID)
	if !exists {
		err = errors.New("no exists " + TJSZTXRedisID + " logger config")
		ServiceLogger.DebugFormat("not exists " + TJSZTXRedisID + " logger config")
		//return err
	}
	TJSZTXRedisConfig = TJSZTXRedis

	//执行已注册的配置初始化接口
	for _, loader := range serviceLoaderMap {
		loader()
	}
	return nil
}

// RegisterServiceLoader 注册服务加载接口
func RegisterServiceLoader(serviceName string, service ServiceLoader) {
	serviceLoaderMap[serviceName] = service
}
