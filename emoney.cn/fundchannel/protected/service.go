package protected

import (
	"emoney.cn/fundchannel/config"
	"emoney.cn/fundchannel/const"
	"errors"
	"github.com/devfeel/dotlog"
)

type ServiceLoader func()

const (
	DefaultRedisID    = "default"
	defaultDatabaseID = "demodb"
	SessionRedisID                  = "session_redis"
)

var (
	DefaultConfig    *ServiceConfig
	serviceLoaderMap map[string]ServiceLoader
	ServiceLogger    dotlog.Logger
	//默认Redis
	DefaultRedisConfig *config.RedisInfo
	//Session会话专员专用Redis
	SessionRedisConfig *config.RedisInfo
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

	//获取数据库连接字符串
	dbInfo, exists := config.GetDataBaseInfo(defaultDatabaseID)
	if !exists || dbInfo.ServerUrl == "" {
		err = errors.New("no config " + defaultDatabaseID + " database config")
		ServiceLogger.Error(err, "no config "+defaultDatabaseID+" database config")
		return err
	}

	DefaultConfig = &ServiceConfig{
		DefaultDBConn: dbInfo.ServerUrl,
	}


	//初始化Redis配置信息
	redis, exists := config.GetRedisInfo(DefaultRedisID)
	if !exists {
		err = errors.New("no exists " + DefaultRedisID + " logger config")
		ServiceLogger.Error(err, "not exists "+DefaultRedisID+" logger config")
		return err
	}
	DefaultRedisConfig = redis


	//初始化SessionRedis配置信息
	sessionRedis, exists := config.GetRedisInfo(SessionRedisID)
	if !exists {
		err = errors.New("no exists " + SessionRedisID + " logger config")
		ServiceLogger.Error(err, "not exists "+SessionRedisID+" logger config")
		return err
	}
	SessionRedisConfig = sessionRedis

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
