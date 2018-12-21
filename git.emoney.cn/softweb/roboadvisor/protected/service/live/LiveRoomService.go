package live

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	liverepo "git.emoney.cn/softweb/roboadvisor/protected/repository/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/dotweb/framework/json"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

var (
	//shareRoomRepo 共享的LiveRoomRepository实例
	shareRoomRepo *liverepo.LiveRoomRepository

	// shareRoomLogger 共享的Logger实例
	shareRoomLogger dotlog.Logger
)

const (
	liveRoomServiceName                = "LiveRoomService"
	EMNET_LiveRoom_CachePreKey         = "EMoney:Live:LiveRoomV2.0"
	EMNET_LiveRoom_GetLiveRoomCacheKey = EMNET_LiveRoom_CachePreKey + "LiveRoomHashId"
)

type LiveRoomService struct {
	service.BaseService
	roomRepo *liverepo.LiveRoomRepository
}

// NewLiveRoomService 创建LiveRoomService服务
func NewLiveRoomService() *LiveRoomService {
	liveRoomService := &LiveRoomService{
		roomRepo: shareRoomRepo,
	}
	liveRoomService.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return liveRoomService
}

// GetLiveRoom 根据直播间Id获取直播间信息
func (service *LiveRoomService) GetLiveRoom(roomId int) (*livemodel.LiveRoom, error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getLiveRoomDB(roomId)
	default:
		return service.getLiveRoomCache(roomId)
	}
}

// getLiveRoomCache 根据直播间Id获取直播间信息
func (service *LiveRoomService) getLiveRoomCache(roomId int) (*livemodel.LiveRoom, error) {
	json, err := service.RedisCache.HGet(EMNET_LiveRoom_GetLiveRoomCacheKey, strconv.Itoa(roomId))
	if err != nil || json == "" {
		return nil, err
	}
	room := new(livemodel.LiveRoom)
	err = jsonutil.Unmarshal(json, room)
	if err != nil {
		shareRoomLogger.Error(err, "根据直播间Id获取直播间信息异常 getLiveRoomCache roomId="+strconv.Itoa(roomId))
		return nil, err
	}
	return room, nil
}

// getLiveRoomDB 根据直播间Id获取直播间信息
func (service *LiveRoomService) getLiveRoomDB(roomId int) (*livemodel.LiveRoom, error) {
	room, err := service.roomRepo.GetLiveRoom(roomId)
	if err != nil {
		shareRoomLogger.Error(err, "根据直播间Id获取直播间信息异常 getLiveRoomDB roomId="+strconv.Itoa(roomId))
		return nil, err
	}
	return room, err
}

func init() {
	protected.RegisterServiceLoader(liveRoomServiceName, liveRoomServiceLoader)
}

func liveRoomServiceLoader() {
	shareRoomRepo = liverepo.NewLiveRoomRepository(protected.DefaultConfig)
	shareRoomLogger = dotlog.GetLogger(liveRoomServiceName)
}
