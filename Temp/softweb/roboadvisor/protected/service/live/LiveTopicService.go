package live

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	liverepo "git.emoney.cn/softweb/roboadvisor/protected/repository/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var (
	//shareTopicRepo 共享的LiveTopicRepository实例
	shareTopicRepo *liverepo.LiveTopicRepository

	// shareTopicLogger 共享的Logger实例
	shareTopicLogger dotlog.Logger
)

const (
	liveTopicServiceName                  = "LiveTopicService"
	EMNET_LiveTopic_CachePreKey           = "EMoney:Live:LiveTopic_v5"
	EMNET_LiveTopic_CacheKeyLiveTopicList = EMNET_LiveTopic_CachePreKey + ":LiveTopicList"
)

type LiveTopicService struct {
	service.BaseService
	topicRepo *liverepo.LiveTopicRepository
}

// NewLiveTopicService 创建LiveTopicService服务
func NewLiveTopicService() *LiveTopicService {
	liveTopicService := &LiveTopicService{
		topicRepo: shareTopicRepo,
	}
	liveTopicService.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return liveTopicService
}

// HaveTopicOpen 直播间列表中是否有一个或以上的直播间开播
func (service *LiveTopicService) HaveTopicOpen(date time.Time, roomIdList ...int) bool {
	for _, roomId := range roomIdList {
		topic, err := service.GetTopic(roomId, date)
		if err == nil && topic != nil {
			return true
		}
	}
	return false
}

// GetTopic 根据直播间Id获取主题信息
func (service *LiveTopicService) GetTopic(roomId int, date time.Time) (topic *livemodel.LiveTopic, err error) {
	if roomId <= 0 {
		return nil, errors.New("直播间号不正确！")
	}
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getTopicDB(roomId, date)
	default:
		return service.getTopicCache(roomId, date)
	}
}

// getTopicCache 根据直播间Id获取主题信息
func (service *LiveTopicService) getTopicCache(roomId int, date time.Time) (topic *livemodel.LiveTopic, err error) {
	cacheKey := EMNET_LiveTopic_CacheKeyLiveTopicList + strconv.Itoa(roomId) + date.Format("2006/01/02")
	err = service.RedisCache.GetJsonObj(cacheKey, &topic)
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		shareTopicLogger.Error(err, "根据直播间Id获取主题信息异常 getTopicCache roomId="+strconv.Itoa(roomId))
		return nil, err
	}
	return topic, nil
}

// getTopicDB 根据直播间Id获取主题信息
func (service *LiveTopicService) getTopicDB(roomId int, date time.Time) (topic *livemodel.LiveTopic, err error) {
	topic, err = service.topicRepo.GetTopic(roomId, date)
	if err != nil {
		shareTopicLogger.Error(err, "根据直播间Id获取主题信息异常 getTopicCache roomId="+strconv.Itoa(roomId))
	}
	return topic, err
}

// GetNewestTopic 根据直播间Id获取最新主题信息
func (service *LiveTopicService) GetNewestTopic(roomId int) (topic *livemodel.LiveTopic, err error) {
	if roomId <= 0 {
		return nil, errors.New("直播间号不正确！")
	}
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getNewestTopicDB(roomId)
	default:
		return service.getNewestTopicCache(roomId)
	}
}

// getNewestTopicCache 根据直播间Id获取最新主题信息
func (service *LiveTopicService) getNewestTopicCache(roomId int) (topic *livemodel.LiveTopic, err error) {
	cacheKey := EMNET_LiveTopic_CacheKeyLiveTopicList + strconv.Itoa(roomId)
	err = service.RedisCache.GetJsonObj(cacheKey, &topic)
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		shareTopicLogger.Error(err, "根据直播间Id获取最新主题信息异常 getNewestTopicCache roomId="+strconv.Itoa(roomId))
		return nil, err
	}
	return topic, nil
}

// getNewestTopicDB 根据直播间Id获取最新主题信息
func (service *LiveTopicService) getNewestTopicDB(roomId int) (topic *livemodel.LiveTopic, err error) {
	topic, err = service.topicRepo.GetNewestTopic(roomId)
	if err != nil {
		shareTopicLogger.Error(err, "根据直播间Id获取最新主题信息异常 getNewestTopicDB roomId="+strconv.Itoa(roomId))
	}
	return topic, err
}

func init() {
	protected.RegisterServiceLoader(liveTopicServiceName, liveTopicServiceLoader)
}

func liveTopicServiceLoader() {
	shareTopicRepo = liverepo.NewLiveTopicRepository(protected.DefaultConfig)
	shareTopicLogger = dotlog.GetLogger(liveTopicServiceName)
}
