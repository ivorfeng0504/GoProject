package service

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type UserInActivityService struct {
	service.BaseService
	userInActivityRepo *userhome_repo.UserInActivityRepository
}

var (
	shareUserInActivityRepo   *userhome_repo.UserInActivityRepository
	shareUserInActivityLogger dotlog.Logger
)

const (
	UserHomeUserInActivityServiceName                   = "UserInActivityService"
	EMNET_UserInActivity_CachePreKey                    = "EMNET:UserInActivityService:"
	EMNET_UserInActivity_GetUserInActivityList_CacheKey = EMNET_UserInActivity_CachePreKey + "GetUserInActivityList:"
	EMNET_UserInActivity_CacheSeconds                   = 60 * 60
)

func NewUserInActivityService() *UserInActivityService {
	srv := &UserInActivityService{
		userInActivityRepo: shareUserInActivityRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeUserInActivityServiceName, userHomeUserInActivityServiceLoader)
}

func userHomeUserInActivityServiceLoader() {
	shareUserInActivityRepo = userhome_repo.NewUserInActivityRepository(protected.DefaultConfig)
	shareUserInActivityLogger = dotlog.GetLogger(UserHomeUserInActivityServiceName)
}

// GetUserInActivityMap 获取用户参与活动Map
func (srv *UserInActivityService) GetUserInActivityMap(userInfoId int) (userActivityMap map[int64]bool, err error) {
	userActivityMap = make(map[int64]bool)
	recordList, err := srv.GetUserInActivityList(userInfoId)
	if err != nil {
		return userActivityMap, err
	}
	for _, record := range recordList {
		userActivityMap[record.ActivityId] = true
	}
	return userActivityMap, nil
}

// GetUserInActivityList 获取用户的活动参与记录
func (srv *UserInActivityService) GetUserInActivityList(userInfoId int) (recordList []*userhome_model.UserInActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserInActivityListDB(userInfoId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		recordList, err = srv.getUserInActivityListCache(userInfoId)
		if err == nil && recordList == nil {
			recordList, err = srv.refreshUserInActivityList(userInfoId)
		}
		return recordList, err
	case config.ReadDB_RefreshCache:
		recordList, err = srv.refreshUserInActivityList(userInfoId)
		return recordList, err
	default:
		return srv.getUserInActivityListCache(userInfoId)
	}
}

// getUserInActivityListDB 获取用户的活动参与记录-查询数据库
func (srv *UserInActivityService) getUserInActivityListDB(userInfoId int) (recordList []*userhome_model.UserInActivity, err error) {
	recordList, err = srv.userInActivityRepo.GetUserInActivityList(userInfoId)
	if err != nil {
		shareUserInActivityLogger.ErrorFormat(err, "获取用户的活动参与记录异常 getUserInActivityListDB userInfoId=%d", userInfoId)
	}
	return recordList, err
}

// getUserInActivityListCache 获取用户的活动参与记录-查询缓存
func (srv *UserInActivityService) getUserInActivityListCache(userInfoId int) (recordList []*userhome_model.UserInActivity, err error) {
	cacheKey := EMNET_UserInActivity_GetUserInActivityList_CacheKey + strconv.Itoa(userInfoId)
	err = srv.RedisCache.GetJsonObj(cacheKey, &recordList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return recordList, err
}

// refreshUserInActivityList 获取用户的活动参与记录-刷新缓存
func (srv *UserInActivityService) refreshUserInActivityList(userInfoId int) (recordList []*userhome_model.UserInActivity, err error) {
	cacheKey := EMNET_UserInActivity_GetUserInActivityList_CacheKey + strconv.Itoa(userInfoId)
	recordList, err = srv.getUserInActivityListDB(userInfoId)
	if err != nil {
		return nil, err
	}
	if recordList != nil {
		srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(recordList), EMNET_UserInActivity_CacheSeconds)
	}
	return recordList, err
}

// InsertUserInActivity 插入用户活动参与记录
func (srv *UserInActivityService) InsertUserInActivity(userInfoId int, activityId int64) (id int64, err error) {
	id, err = srv.userInActivityRepo.InsertUserInActivity(userInfoId, activityId)
	if err != nil {
		shareUserInActivityLogger.ErrorFormat(err, "插入用户活动参与记录异常 InsertUserInActivity userInfoId=%d activityId=%d", userInfoId, activityId)
	} else {
		//刷新缓存
		go srv.refreshUserInActivityList(userInfoId)
	}
	return id, err
}
