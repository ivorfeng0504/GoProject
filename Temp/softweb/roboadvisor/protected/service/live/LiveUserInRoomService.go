package live

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	liverepo "git.emoney.cn/softweb/roboadvisor/protected/repository/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dbsync"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type LiveUserInRoomService struct {
	service.BaseService
	liveUserInRoomRepo *liverepo.LiveUserInRoomRepository
}

var (
	// shareLiveUserInRoomRepo 共享的仓储
	shareLiveUserInRoomRepo *liverepo.LiveUserInRoomRepository

	// shareLiveUserInRoomLogger 共享的Logger实例
	shareLiveUserInRoomLogger dotlog.Logger
)

const (
	liveUserInRoomServiceName     = "LiveUserInRoomService"
	GetLiveUserInRoomListCacheKey = "LiveUserInRoomService.GetLiveUserInRoomList:"
)

func NewLiveUserInRoomService() *LiveUserInRoomService {
	service := &LiveUserInRoomService{
		liveUserInRoomRepo: shareLiveUserInRoomRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return service
}

// AddLiveUserInRoom 添加用户图文直播权限
func (service *LiveUserInRoomService) AddLiveUserInRoom(mobile string, roomId int, day int, orderId string, source string) (id int64, err error) {
	id, err = service.liveUserInRoomRepo.AddLiveUserInRoom(mobile, roomId, day, orderId, source)
	if err != nil {
		shareLiveUserInRoomLogger.Error(err, "添加用户图文直播权限 AddLiveUserInRoom 异常 mobile="+mobile+" roomId="+strconv.Itoa(roomId)+" day="+strconv.Itoa(day))
	} else {
		//写入成功刷新缓存
		service.refreshLiveUserInRoomListCache(mobile)
		//数据库同步
		dbsync.Sync(_const.SyncTable_LiveUserInRoom)
	}
	return id, err
}

// GetLiveUserInRoomList 根据手机号获取用户的直播间权限
func (service *LiveUserInRoomService) GetLiveUserInRoomList(mobile string) (rooms []*livemodel.LiveUserInRoom, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getLiveUserInRoomListDB(mobile)
	default:
		return service.getLiveUserInRoomListCache(mobile)
	}
}

// getLiveUserInRoomListCache 根据手机号获取用户的直播间权限
func (service *LiveUserInRoomService) getLiveUserInRoomListCache(mobile string) (rooms []*livemodel.LiveUserInRoom, err error) {
	cacheKey := GetLiveUserInRoomListCacheKey + mobile
	err = service.RedisCache.GetJsonObj(cacheKey, &rooms)
	if err == redis.ErrNil {
		return rooms, nil
	}
	if err != nil {
		shareLiveUserInRoomLogger.Error(err, "读取缓存-根据手机号获取用户的直播间权限 getLiveUserInRoomListCache 异常 mobile="+mobile)
	}
	return rooms, err
}

// refreshLiveUserInRoomListCache 根据手机号刷新用户的直播间权限
func (service *LiveUserInRoomService) refreshLiveUserInRoomListCache(mobile string) (rooms []*livemodel.LiveUserInRoom, err error) {
	cacheKey := GetLiveUserInRoomListCacheKey + mobile
	rooms, err = service.getLiveUserInRoomListDB(mobile)
	if err != nil {
		shareLiveUserInRoomLogger.Error(err, "刷新缓存-根据手机号刷新用户的直播间权限 refreshLiveUserInRoomListCache 异常 mobile="+mobile)
	}
	if rooms != nil {
		service.RedisCache.SetJsonObj(cacheKey, rooms)
	}
	return rooms, err
}

// getLiveUserInRoomListDB 根据手机号获取用户的直播间权限
func (service *LiveUserInRoomService) getLiveUserInRoomListDB(mobile string) (rooms []*livemodel.LiveUserInRoom, err error) {
	rooms, err = service.liveUserInRoomRepo.GetLiveUserInRoomList(mobile)
	if err != nil {
		shareLiveUserInRoomLogger.Error(err, "查询数据库-根据手机号获取用户的直播间权限 getLiveUserInRoomListDB 异常 mobile="+mobile)
	}
	return rooms, err
}

// GetRoomList 根据手机号获取用户当前的直播间列表并去重
func (service *LiveUserInRoomService) GetRoomList(mobile string) (rooms []int, err error) {
	roomList, err := service.GetLiveUserInRoomList(mobile)
	roomMap := make(map[int]int)
	if err != nil {
		shareLiveUserInRoomLogger.Error(err, "根据手机号获取用户的直播间列表并去重 GetRoomList 异常 mobile="+mobile)
	}
	now := time.Now()
	for _, room := range roomList {
		_, exist := roomMap[room.RoomId]
		if !exist && time.Time(room.ExpireTime).After(now) {
			rooms = append(rooms, room.RoomId)
			roomMap[room.RoomId] = 1
		}
	}
	return rooms, err
}

// GetDefaultRoomList 获取默认直播间权限
func (service *LiveUserInRoomService) GetDefaultRoomList() (rooms []int, err error) {
	roomListStr := config.CurrentConfig.DefaultRoomList
	if roomListStr == "" {
		return rooms, nil
	}
	roomList := strings.Split(roomListStr, ",")
	if roomList == nil || len(roomList) == 0 {
		return rooms, nil
	}
	for _, roomId := range roomList {
		id, err := strconv.Atoi(roomId)
		if err == nil {
			rooms = append(rooms, id)
		}
	}
	return rooms, err
}

// DeleteById 删除指定Id的数据（逻辑删除）
func (service *LiveUserInRoomService) DeleteById(liveUserInRoomId int) (err error) {
	if liveUserInRoomId <= 0 {
		return errors.New("Id不正确")
	}
	room, err := service.liveUserInRoomRepo.GetLiveUserInRoomById(liveUserInRoomId)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("该记录不存在")
	}
	err = service.liveUserInRoomRepo.DeleteById(liveUserInRoomId)
	//逻辑删除成功则更新缓存
	if err == nil {
		//写入成功刷新缓存
		service.refreshLiveUserInRoomListCache(room.Mobile)
		//数据库同步
		dbsync.Sync(_const.SyncTable_LiveUserInRoom)
	}
	return err
}

func init() {
	protected.RegisterServiceLoader(liveUserInRoomServiceName, liveUserInRoomServiceLoader)
}

func liveUserInRoomServiceLoader() {
	shareLiveUserInRoomRepo = liverepo.NewLiveUserInRoomRepository(protected.DefaultConfig)
	shareLiveUserInRoomLogger = dotlog.GetLogger(liveUserInRoomServiceName)
}
