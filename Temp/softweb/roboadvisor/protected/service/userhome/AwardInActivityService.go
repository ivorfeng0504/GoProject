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

type AwardInActivityService struct {
	service.BaseService
	awardInActivityRepo *userhome_repo.AwardInActivityRepository
}

var (
	shareAwardInActivityRepo   *userhome_repo.AwardInActivityRepository
	shareAwardInActivityLogger dotlog.Logger
)

const (
	AwardInActivityServiceName                                      = "AwardInActivityService"
	EMNET_AwardInActivity_CachePreKey                               = "EMNET:AwardInActivityService:"
	EMNET_AwardInActivity_GetActivityAwardListByActivityId_CacheKey = EMNET_AwardInActivity_CachePreKey + "GetActivityAwardListByActivityId:"
	EMNET_AwardInActivity_CacheSeconds                              = 60 * 30
)

func NewAwardInActivityService() *AwardInActivityService {
	srv := &AwardInActivityService{
		awardInActivityRepo: shareAwardInActivityRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(AwardInActivityServiceName, awardInActivityServiceLoader)
}

func awardInActivityServiceLoader() {
	shareAwardInActivityRepo = userhome_repo.NewAwardInActivityRepository(protected.DefaultConfig)
	shareAwardInActivityLogger = dotlog.GetLogger(AwardInActivityServiceName)
}

//获取指定活动的奖品列表  activityId活动Id  ignoreExpired是否过滤掉过期的奖品设置
func (srv *AwardInActivityService) GetActivityAwardListByActivityId(activityId int64, ignoreExpired bool) (awardList []*userhome_model.ActivityAward, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.GetActivityAwardListByActivityIdDB(activityId, ignoreExpired)
	case config.ReadDB_CacheOrDB_UpdateCache:
		awardList, err = srv.getActivityAwardListByActivityIdCache(activityId, ignoreExpired)
		if err == nil && awardList == nil {
			awardList, err = srv.refreshetActivityAwardListByActivityId(activityId, ignoreExpired)
		}
		return awardList, err
	case config.ReadDB_RefreshCache:
		awardList, err = srv.refreshetActivityAwardListByActivityId(activityId, ignoreExpired)
		return awardList, err
	default:
		return srv.getActivityAwardListByActivityIdCache(activityId, ignoreExpired)
	}
}

//获取指定活动的奖品列表-查询数据库  activityId活动Id  ignoreExpired是否过滤掉过期的奖品设置
func (srv *AwardInActivityService) GetActivityAwardListByActivityIdDB(activityId int64, ignoreExpired bool) (awardList []*userhome_model.ActivityAward, err error) {
	awardList, err = srv.awardInActivityRepo.GetActivityAwardListByActivityId(activityId, ignoreExpired)
	if err != nil {
		shareAwardInActivityLogger.ErrorFormat(err, "GetActivityAwardListByActivityId 获取指定活动的奖品列表 异常 activityId=%d ignoreExpired=%s", activityId, ignoreExpired)
	}
	return awardList, err
}

//获取指定活动的奖品列表-读取缓存  activityId活动Id  ignoreExpired是否过滤掉过期的奖品设置
func (srv *AwardInActivityService) getActivityAwardListByActivityIdCache(activityId int64, ignoreExpired bool) (awardList []*userhome_model.ActivityAward, err error) {
	cacheKey := EMNET_AwardInActivity_GetActivityAwardListByActivityId_CacheKey + strconv.FormatInt(activityId, 10) + ":" + strconv.FormatBool(ignoreExpired)
	err = srv.RedisCache.GetJsonObj(cacheKey, &awardList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return awardList, err
}

//获取指定活动的奖品列表-刷新缓存  activityId活动Id  ignoreExpired是否过滤掉过期的奖品设置
func (srv *AwardInActivityService) refreshetActivityAwardListByActivityId(activityId int64, ignoreExpired bool) (awardList []*userhome_model.ActivityAward, err error) {
	cacheKey := EMNET_AwardInActivity_GetActivityAwardListByActivityId_CacheKey + strconv.FormatInt(activityId, 10) + ":" + strconv.FormatBool(ignoreExpired)
	awardList, err = srv.GetActivityAwardListByActivityIdDB(activityId, ignoreExpired)
	if err != nil {
		return nil, err
	}
	if awardList != nil {
		err = srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(awardList), EMNET_AwardInActivity_CacheSeconds)
	}
	return awardList, err
}
