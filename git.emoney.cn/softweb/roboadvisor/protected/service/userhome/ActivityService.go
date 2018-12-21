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

type ActivityService struct {
	service.BaseService
	activityRepo *userhome_repo.ActivityRepository
}

var (
	shareActivityRepo   *userhome_repo.ActivityRepository
	shareActivityLogger dotlog.Logger
)

const (
	UserHomeActivityServiceName             = "ActivityService"
	EMNET_Activity_CachePreKey              = "EMNET:UserHome_ActivityBLL"
	EMNET_Activity_GetActivityById_CacheKey = EMNET_Activity_CachePreKey + ":GetActivityById:"
	EMNET_Activity_GetActivityList_CacheKey = EMNET_Activity_CachePreKey + ":GetActivityList:"
	EMNET_Activity_CacheSeconds             = 15 * 60
)

func NewActivityService() *ActivityService {
	srv := &ActivityService{
		activityRepo: shareActivityRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeActivityServiceName, userHomeActivityServiceLoader)
}

func userHomeActivityServiceLoader() {
	shareActivityRepo = userhome_repo.NewActivityRepository(protected.DefaultConfig)
	shareActivityLogger = dotlog.GetLogger(UserHomeActivityServiceName)
}

// GetActivityList 获取活动列表
func (srv *ActivityService) GetActivityList(state int) (activityList []*userhome_model.Activity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getActivityListDB(state)
	case config.ReadDB_CacheOrDB_UpdateCache:
		activityList, err = srv.getActivityListCache(state)
		if err == nil && activityList == nil {
			activityList, err = srv.refreshActivityList(state)
		}
		return activityList, err
	case config.ReadDB_RefreshCache:
		activityList, err = srv.refreshActivityList(state)
		return activityList, err
	default:
		return srv.getActivityListCache(state)
	}
}

// getActivityListDB 获取活动列表-查询数据库
func (srv *ActivityService) getActivityListDB(state int) (activityList []*userhome_model.Activity, err error) {
	activityList, err = srv.activityRepo.GetActivityList(state)
	if err != nil {
		shareActivityLogger.ErrorFormat(err, "查询活动列表异常 state=%d", state)
	}
	return activityList, err
}

// getActivityListCache 获取活动列表-查询缓存
func (srv *ActivityService) getActivityListCache(state int) (activityList []*userhome_model.Activity, err error) {
	cacheKey := EMNET_Activity_GetActivityList_CacheKey + strconv.Itoa(state)
	err = srv.RedisCache.GetJsonObj(cacheKey, &activityList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return activityList, err
}

// refreshActivityList 刷新缓存
func (srv *ActivityService) refreshActivityList(state int) (activityList []*userhome_model.Activity, err error) {
	cacheKey := EMNET_Activity_GetActivityList_CacheKey + strconv.Itoa(state)
	activityList, err = srv.getActivityListDB(state)
	if err != nil {
		return nil, err
	}
	if activityList != nil {
		srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(activityList), EMNET_Activity_CacheSeconds)
	}
	return activityList, err
}

// GetActivityById 获取活动详情
func (srv *ActivityService) GetActivityById(activityId int64) (activity *userhome_model.Activity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getActivityByIdDB(activityId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		activity, err = srv.getActivityByIdCache(activityId)
		if err == nil && activity == nil {
			activity, err = srv.refreshActivity(activityId)
		}
		return activity, err
	case config.ReadDB_RefreshCache:
		activity, err = srv.refreshActivity(activityId)
		return activity, err
	default:
		return srv.getActivityByIdCache(activityId)
	}
}

// getActivityListDB 获取活动列表-查询数据库
func (srv *ActivityService) getActivityByIdDB(activityId int64) (activity *userhome_model.Activity, err error) {
	activity, err = srv.activityRepo.GetActivityById(activityId)
	if err != nil {
		shareActivityLogger.ErrorFormat(err, "获取活动详情异常 getActivityByIdDB activityId=%d", activityId)
	}
	return activity, err
}

// getActivityByIdCache 获取活动列表-查询缓存
func (srv *ActivityService) getActivityByIdCache(activityId int64) (activity *userhome_model.Activity, err error) {
	cacheKey := EMNET_Activity_GetActivityById_CacheKey + strconv.FormatInt(activityId, 10)
	err = srv.RedisCache.GetJsonObj(cacheKey, &activity)
	if err == redis.ErrNil {
		return nil, nil
	}
	return activity, err
}

// refreshActivity 刷新缓存
func (srv *ActivityService) refreshActivity(activityId int64) (activity *userhome_model.Activity, err error) {
	cacheKey := EMNET_Activity_GetActivityById_CacheKey + strconv.FormatInt(activityId, 10)
	activity, err = srv.getActivityByIdDB(activityId)
	if err != nil {
		return nil, err
	}
	if activity != nil {
		srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(activity), EMNET_Activity_CacheSeconds)
	}
	return activity, err
}
