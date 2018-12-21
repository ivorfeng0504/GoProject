package click

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"github.com/devfeel/cache"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
)

var (
	shareClickServiceLogger dotlog.Logger
	shareClickServiceRedis  cache.RedisCache
	//点击处理器
	clickProvider map[string]func(identity string, clickNum int64) (err error)
)

const (
	ClickServiceName                         = "ClickService"
	EMNET_ClickService_CachePreKey           = "EMNET:ClickService:"
	EMNET_ClickService_ProcessQueue_CacheKey = EMNET_ClickService_CachePreKey + "ProcessQueue:"
	EMNET_ClickService_ClickStat_CacheKey    = EMNET_ClickService_CachePreKey + "ClickStat:"
)

func init() {
	clickProvider = make(map[string]func(identity string, clickNum int64) (err error))
	protected.RegisterServiceLoader(ClickServiceName, func() {
		shareClickServiceLogger = dotlog.GetLogger(ClickServiceName)
		shareClickServiceRedis = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	})
}

// Register 注册指定的点击统计服务
// clickType 点击服务名称 详见const/ClickType.go
// process 点击数据处理（数据库落库）  identity为资讯等的唯一标识  clickNum为点击数
func Register(clickType string, process func(identity string, clickNum int64) (err error)) (err error) {
	_, success := clickProvider[clickType]
	if success {
		err = errors.New("指定的Click处理器已存在 clickType=" + clickType)
	} else {
		clickProvider[clickType] = process
	}
	return err
}

// 手动处理click数据
func HandleClick(clickType string, identity string,clicknum int)(count int64, err error) {
	cacheKey := fmt.Sprintf("%s%s:%s", EMNET_ClickService_ClickStat_CacheKey, clickType, identity)
	_,err = shareClickServiceRedis.SetJsonObj(cacheKey, clicknum)

	count, err = shareClickServiceRedis.Incr(cacheKey)
	fmt.Println(count)
	return count, err
}

// AddClick 添加点击数
// clickType 点击服务名称 详见const/ClickType.go
// identity为资讯等的唯一标识
func AddClick(clickType string, identity string) (count int64, err error) {
	cacheKey := fmt.Sprintf("%s%s:%s", EMNET_ClickService_ClickStat_CacheKey, clickType, identity)
	count, err = shareClickServiceRedis.Incr(cacheKey)
	queueCacheKey := EMNET_ClickService_ProcessQueue_CacheKey + clickType
	shareClickServiceRedis.HSet(queueCacheKey, identity, "1")
	return count, err
}

// GetClick 获取点击数
func GetClick(clickType string, identity string) (count int64, err error) {
	cacheKey := fmt.Sprintf("%s%s:%s", EMNET_ClickService_ClickStat_CacheKey, clickType, identity)
	count, err = shareClickServiceRedis.GetInt64(cacheKey)
	return count, err
}

// ProcessQueue 处理点击队列
func ProcessQueue() error {
	for clickType, process := range clickProvider {
		cacheKey := EMNET_ClickService_ProcessQueue_CacheKey + clickType
		allIdentity, err := shareClickServiceRedis.HGetAll(cacheKey)
		if err == redis.ErrNil {
			shareClickServiceLogger.DebugFormat("ProcessQueue 处理点击队列 没有需要处理的队列")
			return nil
		}
		if err != nil {
			shareClickServiceLogger.ErrorFormat(err, "ProcessQueue 处理点击队列 获取所有identity失败  cacheKey=%s", cacheKey)
			return err
		}
		for identity, _ := range allIdentity {
			//清除该项任务
			shareClickServiceRedis.HDel(cacheKey, identity)
			clickNum, err := GetClick(clickType, identity)
			if err != nil {
				if err == redis.ErrNil {
					shareClickServiceLogger.DebugFormat("ProcessQueue 处理点击队列 获取点击数失败 Redis为空 clickType=%s   identity=%s", clickType, identity)
					continue
				}
				shareClickServiceLogger.ErrorFormat(err, "ProcessQueue 处理点击队列 获取点击数失败  clickType=%s   identity=%s", clickType, identity)
				//如果处理异常 将该任务重新添加到redis集合中
				shareClickServiceRedis.HSet(cacheKey, identity, "1")
				continue
			}
			err = process(identity, clickNum)
			if err != nil {
				shareClickServiceLogger.ErrorFormat(err, "ProcessQueue 处理点击队列 Process处理异常  clickType=%s   identity=%s clickNum=%d", clickType, identity, clickNum)
				//如果处理异常 将该任务重新添加到redis集合中
				//垃圾数据 可能导致一直循环处理  暂时忽略错误
				//shareClickServiceRedis.HSet(cacheKey, identity, "1")
				continue
			}
		}
	}
	return nil
}
