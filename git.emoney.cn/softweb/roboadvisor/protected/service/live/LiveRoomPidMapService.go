package live

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	liverepo "git.emoney.cn/softweb/roboadvisor/protected/repository/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type LiveRoomPidMapService struct {
	service.BaseService
	liveRoomPidMapRepo *liverepo.LiveRoomPidMapRepository
}

var (
	// shareRoomPidMapRepo 共享的仓储
	shareLiveRoomPidMapRepo *liverepo.LiveRoomPidMapRepository

	// shareLiveRoomPidMapLogger 共享的Logger实例
	shareLiveRoomPidMapLogger dotlog.Logger
)

const (
	liveRoomPidMapServiceName               = "LiveRoomPidMapService"
	EMNET_LiveRoomPidMap_CachePreKey        = "EMoney.ContentLivePlat.Business.LiveRoomPidMapBll."
	EMNET_LiveRoomPidMap_GetMapListCacheKey = EMNET_LiveRoomPidMap_CachePreKey + "All"
)

func NewLiveRoomPidMapService() *LiveRoomPidMapService {
	service := &LiveRoomPidMapService{
		liveRoomPidMapRepo: shareLiveRoomPidMapRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return service
}

// GetMapList 获取所有直播间与PID之间的映射关系
func (service *LiveRoomPidMapService) GetRoomId(pid int) (roomId int) {
	if pid <= 0 {
		return 0
	}
	mapList, err := service.GetMapList()
	if err != nil || mapList == nil {
		return 0
	}
	for _, mapInfo := range mapList {
		if mapInfo.PID == pid {
			return mapInfo.RoomId
		}
	}
	return 0
}

// GetMapList 获取所有直播间与PID之间的映射关系
func (service *LiveRoomPidMapService) GetMapList() (mapList []*livemodel.LiveRoomPidMap, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getMapListDB()
	case config.ReadDB_CacheOrDB_UpdateCache:
		mapList, err = service.getMapListCache()
		if err == nil && mapList == nil {
			mapList, err = service.refreshMapListCache()
		}
		return mapList, err
	default:
		return service.getMapListCache()
	}
}

// getMapListDB 获取所有直播间与PID之间的映射关系
func (service *LiveRoomPidMapService) getMapListDB() (mapList []*livemodel.LiveRoomPidMap, err error) {
	mapList, err = service.liveRoomPidMapRepo.GetMapList()
	return mapList, err
}

// getMapListCache 获取所有直播间与PID之间的映射关系
func (service *LiveRoomPidMapService) getMapListCache() (mapList []*livemodel.LiveRoomPidMap, err error) {
	cacheKey := EMNET_LiveRoomPidMap_GetMapListCacheKey
	err = service.RedisCache.GetJsonObj(cacheKey, &mapList)
	if err == redis.ErrNil {
		return mapList, nil
	}
	if err != nil {
		shareLiveRoomPidMapLogger.Error(err, "读取缓存-获取所有直播间与PID之间的映射关系 getMapListCache 异常")
	}
	return mapList, err
}

// refreshMapListCache 刷新所有直播间与PID之间的映射关系缓存
func (service *LiveRoomPidMapService) refreshMapListCache() (mapList []*livemodel.LiveRoomPidMap, err error) {
	mapList, err = service.liveRoomPidMapRepo.GetMapList()
	if err != nil || mapList == nil {
		return mapList, err
	}
	cacheKey := EMNET_LiveRoomPidMap_GetMapListCacheKey
	_, err = service.RedisCache.SetJsonObj(cacheKey, mapList)
	return mapList, err
}

func init() {
	protected.RegisterServiceLoader(liveRoomPidMapServiceName, liveRoomPidMapServiceLoader)
}

func liveRoomPidMapServiceLoader() {
	shareLiveRoomPidMapRepo = liverepo.NewLiveRoomPidMapRepository(protected.DefaultConfig)
	shareLiveRoomPidMapLogger = dotlog.GetLogger(liveRoomPidMapServiceName)
}
