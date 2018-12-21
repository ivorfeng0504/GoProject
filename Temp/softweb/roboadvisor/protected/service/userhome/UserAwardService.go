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

type UserAwardService struct {
	service.BaseService
	userAwardRepo *userhome_repo.UserAwardRepository
}

var (
	shareUserAwardRepo   *userhome_repo.UserAwardRepository
	shareUserAwardLogger dotlog.Logger
)

const (
	UserHomeUserAwardServiceName              = "UserAwardService"
	EMNET_UserAward_CachePreKey               = "EMNET:UserAwardService"
	EMNET_UserAward_GetUserAwardList_CacheKey = EMNET_UserAward_CachePreKey + ":GetUserAwardList:"
	EMNET_UserAward_CacheSeconds              = 15 * 60
)

func NewUserAwardService() *UserAwardService {
	srv := &UserAwardService{
		userAwardRepo: shareUserAwardRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeUserAwardServiceName, userHomeUserAwardServiceLoader)
}

func userHomeUserAwardServiceLoader() {
	shareUserAwardRepo = userhome_repo.NewUserAwardRepository(protected.DefaultConfig)
	shareUserAwardLogger = dotlog.GetLogger(UserHomeUserAwardServiceName)
}

// GetUserAwardList 获取用户的奖品列表
func (srv *UserAwardService) GetUserAwardList(userInfoId int, state int) (userAwardList []*userhome_model.UserAward, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserAwardListDB(userInfoId, state)
	case config.ReadDB_CacheOrDB_UpdateCache:
		userAwardList, err = srv.getUserAwardListCache(userInfoId, state)
		if err == nil && userAwardList == nil {
			userAwardList, err = srv.refreshUserAwardList(userInfoId, state)
		}
		return userAwardList, err
	case config.ReadDB_RefreshCache:
		userAwardList, err = srv.refreshUserAwardList(userInfoId, state)
		return userAwardList, err
	default:
		return srv.getUserAwardListCache(userInfoId, state)
	}
}

// getUserAwardListDB 获取用户的奖品列表-查询数据库
func (srv *UserAwardService) getUserAwardListDB(userInfoId int, state int) (userAwardList []*userhome_model.UserAward, err error) {
	userAwardList, err = srv.userAwardRepo.GetUserAwardList(userInfoId, state)
	if err != nil {
		shareUserAwardLogger.ErrorFormat(err, "获取用户的奖品列表异常 userInfoId=%d", userInfoId)
	}
	return userAwardList, err
}

// getUserAwardListCache 获取用户的奖品列表-查询缓存
func (srv *UserAwardService) getUserAwardListCache(userInfoId int, state int) (userAwardList []*userhome_model.UserAward, err error) {
	cacheKey := EMNET_UserAward_GetUserAwardList_CacheKey + strconv.Itoa(userInfoId) + ":" + strconv.Itoa(state)
	err = srv.RedisCache.GetJsonObj(cacheKey, &userAwardList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return userAwardList, err
}

// refreshUserAwardList 获取用户的奖品列表-刷新缓存
func (srv *UserAwardService) refreshUserAwardList(userInfoId int, state int) (userAwardList []*userhome_model.UserAward, err error) {
	cacheKey := EMNET_UserAward_GetUserAwardList_CacheKey + strconv.Itoa(userInfoId) + ":" + strconv.Itoa(state)
	userAwardList, err = srv.getUserAwardListDB(userInfoId, state)
	if err != nil {
		return nil, err
	}
	if userAwardList != nil {
		srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(userAwardList), EMNET_UserAward_CacheSeconds)
	}
	return userAwardList, err
}

// InsertUserAward 插入用户奖品记录
func (srv *UserAwardService) InsertUserAward(userInfoId int, awardName string, awardId int64, introduceVideo string, activityId int64, activityName string, state int, awardImg string, avaDay int, awardType int) (id int64, err error) {
	id, err = srv.userAwardRepo.InsertUserAward(userInfoId, awardName, awardId, introduceVideo, activityId, activityName, state, awardImg, avaDay, awardType)
	if err != nil {
		shareUserAwardLogger.ErrorFormat(err, "插入用户奖品记录异常 InsertUserAward")
	} else {
		//刷新缓存
		go srv.refreshUserAwardList(userInfoId, state)
	}
	return id, err
}
