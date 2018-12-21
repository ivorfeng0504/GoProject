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

type UserInGuessMarketActivityService struct {
	service.BaseService
	userInGuessMarketActivityRepo *userhome_repo.UserInGuessMarketActivityRepository
}

var (
	shareUserInGuessMarketActivityRepo   *userhome_repo.UserInGuessMarketActivityRepository
	shareUserInGuessMarketActivityLogger dotlog.Logger
)

const (
	UserInGuessMarketActivityServiceName                                             = "UserInGuessMarketActivityService"
	EMNET_UserInGuessMarketActivity_CachePreKey                                      = "EMNET:UserInGuessMarketActivityService:"
	EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivityGuessedList_CacheKey = EMNET_UserInGuessMarketActivity_CachePreKey + ":GetUserInGuessMarketActivityGuessedList:"
	EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivityList_CacheKey        = EMNET_UserInGuessMarketActivity_CachePreKey + ":GetUserInGuessMarketActivityList:"
	EMNET_UserInGuessMarketActivity_GetUserJoinCount_CacheKey                        = EMNET_UserInGuessMarketActivity_CachePreKey + ":GetUserJoinCount:"
	EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivity_CacheKey            = EMNET_UserInGuessMarketActivity_CachePreKey + ":GetUserInGuessMarketActivity:"
	EMNET_UserInGuessMarketActivity_CacheSeconds                                     = 15 * 60
)

func NewUserInGuessMarketActivityService() *UserInGuessMarketActivityService {
	srv := &UserInGuessMarketActivityService{
		userInGuessMarketActivityRepo: shareUserInGuessMarketActivityRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserInGuessMarketActivityServiceName, userInGuessMarketActivityServiceLoader)
}

func userInGuessMarketActivityServiceLoader() {
	shareUserInGuessMarketActivityRepo = userhome_repo.NewUserInGuessMarketActivityRepository(protected.DefaultConfig)
	shareUserInGuessMarketActivityLogger = dotlog.GetLogger(UserInGuessMarketActivityServiceName)
}

// GetUserInGuessMarketActivityGuessedList 查询最新的获奖记录列表
func (srv *UserInGuessMarketActivityService) GetUserInGuessMarketActivityGuessedList(top int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserInGuessMarketActivityGuessedListDB(top)
	case config.ReadDB_CacheOrDB_UpdateCache:
		resultList, err = srv.getUserInGuessMarketActivityGuessedListCache(top)
		if err == nil && resultList == nil {
			resultList, err = srv.refreshUserInGuessMarketActivityGuessedList(top)
		}
		return resultList, err
	case config.ReadDB_RefreshCache:
		resultList, err = srv.refreshUserInGuessMarketActivityGuessedList(top)
		return resultList, err
	default:
		return srv.getUserInGuessMarketActivityGuessedListCache(top)
	}
}

// getUserInGuessMarketActivityGuessedListDB 查询最新的获奖记录列表-读取数据库
func (srv *UserInGuessMarketActivityService) getUserInGuessMarketActivityGuessedListDB(top int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	resultList, err = srv.userInGuessMarketActivityRepo.GetUserInGuessMarketActivityGuessedList(top)
	if err != nil {
		shareUserInGuessMarketActivityLogger.ErrorFormat(err, "查询最新的获奖记录列表异常")
	}
	return resultList, err
}

// getUserInGuessMarketActivityGuessedListCache 查询最新的获奖记录列表-读取缓存
func (srv *UserInGuessMarketActivityService) getUserInGuessMarketActivityGuessedListCache(top int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	cacheKey := EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivityGuessedList_CacheKey + strconv.Itoa(top)
	err = srv.RedisCache.GetJsonObj(cacheKey, &resultList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return resultList, err
}

// refreshUserInGuessMarketActivityGuessedList 查询最新的获奖记录列表-刷新缓存
func (srv *UserInGuessMarketActivityService) refreshUserInGuessMarketActivityGuessedList(top int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	cacheKey := EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivityGuessedList_CacheKey + strconv.Itoa(top)
	resultList, err = srv.getUserInGuessMarketActivityGuessedListDB(top)
	if err != nil {
		return nil, err
	}
	if resultList != nil {
		srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(resultList), EMNET_UserInGuessMarketActivity_CacheSeconds)
	}
	return resultList, err
}

// GetUserInGuessMarketActivityList 查询用户最新的参与记录
func (srv *UserInGuessMarketActivityService) GetUserInGuessMarketActivityList(top int, userInfoId int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserInGuessMarketActivityListDB(top, userInfoId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		resultList, err = srv.getUserInGuessMarketActivityListCache(top, userInfoId)
		if err == nil && resultList == nil {
			resultList, err = srv.refreshUserInGuessMarketActivityList(top, userInfoId)
		}
		return resultList, err
	case config.ReadDB_RefreshCache:
		resultList, err = srv.refreshUserInGuessMarketActivityList(top, userInfoId)
		return resultList, err
	default:
		return srv.getUserInGuessMarketActivityListCache(top, userInfoId)
	}
}

// getUserInGuessMarketActivityListDB 查询用户最新的参与记录-读取数据库
func (srv *UserInGuessMarketActivityService) getUserInGuessMarketActivityListDB(top int, userInfoId int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	resultList, err = srv.userInGuessMarketActivityRepo.GetUserInGuessMarketActivityList(top, userInfoId)
	if err != nil {
		shareUserInGuessMarketActivityLogger.ErrorFormat(err, "查询用户最新的参与记录 userInfoId=%d", userInfoId)
	}
	return resultList, err
}

// getUserInGuessMarketActivityListCache 查询用户最新的参与记录-读取缓存
func (srv *UserInGuessMarketActivityService) getUserInGuessMarketActivityListCache(top int, userInfoId int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	cacheKey := EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivityList_CacheKey + strconv.Itoa(top) + ":" + strconv.Itoa(userInfoId)
	err = srv.RedisCache.GetJsonObj(cacheKey, &resultList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return resultList, err
}

// refreshUserInGuessMarketActivityList 查询用户最新的参与记录-刷新缓存
func (srv *UserInGuessMarketActivityService) refreshUserInGuessMarketActivityList(top int, userInfoId int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	cacheKey := EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivityList_CacheKey + strconv.Itoa(top) + ":" + strconv.Itoa(userInfoId)
	resultList, err = srv.getUserInGuessMarketActivityListDB(top, userInfoId)
	if err != nil {
		return nil, err
	}
	if resultList != nil {
		srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(resultList), EMNET_UserInGuessMarketActivity_CacheSeconds)
	}
	return resultList, err
}

// GetUserJoinCount 查询用户参与总数
func (srv *UserInGuessMarketActivityService) GetUserJoinCount(userInfoId int) (count int64, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserJoinCountDB(userInfoId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		count, err = srv.getUserJoinCountCache(userInfoId)
		if err == redis.ErrNil {
			count, err = srv.refreshUserJoinCount(userInfoId)
		}
		return count, err
	case config.ReadDB_RefreshCache:
		count, err = srv.refreshUserJoinCount(userInfoId)
		return count, err
	default:
		return srv.getUserJoinCountCache(userInfoId)
	}
}

// getUserJoinCountDB 查询用户参与总数-读取数据库
func (srv *UserInGuessMarketActivityService) getUserJoinCountDB(userInfoId int) (count int64, err error) {
	count, err = srv.userInGuessMarketActivityRepo.GetUserJoinCount(userInfoId)
	if err != nil {
		shareUserInGuessMarketActivityLogger.ErrorFormat(err, "查询用户参与总数 userInfoId=%d", userInfoId)
	}
	return count, err
}

// getUserJoinCountCache 查询用户参与总数-读取缓存
func (srv *UserInGuessMarketActivityService) getUserJoinCountCache(userInfoId int) (count int64, err error) {
	cacheKey := EMNET_UserInGuessMarketActivity_GetUserJoinCount_CacheKey + strconv.Itoa(userInfoId)
	err = srv.RedisCache.GetJsonObj(cacheKey, &count)
	return count, err
}

// refreshUserJoinCount 查询用户参与总数-刷新缓存
func (srv *UserInGuessMarketActivityService) refreshUserJoinCount(userInfoId int) (count int64, err error) {
	cacheKey := EMNET_UserInGuessMarketActivity_GetUserJoinCount_CacheKey + strconv.Itoa(userInfoId)
	count, err = srv.getUserJoinCountDB(userInfoId)
	if err != nil {
		return 0, err
	}
	srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(count), EMNET_UserInGuessMarketActivity_CacheSeconds)
	return count, err
}

// Insert 写入用户领数字记录
func (srv *UserInGuessMarketActivityService) Insert(uid int64, userInfoId int, nickName string, issueNum string, result string, guessMarketActivityId int64) (id int64, err error) {
	id, err = srv.userInGuessMarketActivityRepo.InsertUserInGuessMarketActivity(uid, userInfoId, nickName, issueNum, result, guessMarketActivityId)
	if err != nil {
		shareUserInGuessMarketActivityLogger.ErrorFormat(err, "写入用户领数字记录异常")
	} else {
		//刷新缓存
		go srv.refreshUserJoinCount(userInfoId)
		//中奖名单4条
		go srv.refreshUserInGuessMarketActivityGuessedList(4)
		//往期中奖名单6条
		go srv.refreshUserInGuessMarketActivityList(6, userInfoId)
		go srv.refreshUserInGuessMarketActivity(userInfoId, issueNum)
	}
	return id, err
}

// GetUserInGuessMarketActivity 获取用户指定期号的猜数字记录
func (srv *UserInGuessMarketActivityService) GetUserInGuessMarketActivity(userInfoId int, issueNum string) (result *userhome_model.UserInGuessMarketActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserInGuessMarketActivityDB(userInfoId, issueNum)
	case config.ReadDB_CacheOrDB_UpdateCache:
		result, err = srv.getUserInGuessMarketActivityCache(userInfoId, issueNum)
		if err == nil && result == nil {
			result, err = srv.refreshUserInGuessMarketActivity(userInfoId, issueNum)
		}
		return result, err
	case config.ReadDB_RefreshCache:
		result, err = srv.refreshUserInGuessMarketActivity(userInfoId, issueNum)
		return result, err
	default:
		return srv.getUserInGuessMarketActivityCache(userInfoId, issueNum)
	}
}

// getUserInGuessMarketActivityDB 获取用户指定期号的猜数字记录-读取数据库
func (srv *UserInGuessMarketActivityService) getUserInGuessMarketActivityDB(userInfoId int, issueNum string) (result *userhome_model.UserInGuessMarketActivity, err error) {
	result, err = srv.userInGuessMarketActivityRepo.GetUserInGuessMarketActivity(userInfoId, issueNum)
	if err != nil {
		shareUserInGuessMarketActivityLogger.ErrorFormat(err, "获取用户指定期号的猜数字记录 userInfoId=%d  issueNum=%s", userInfoId, issueNum)
	}
	return result, err
}

// getUserInGuessMarketActivityCache 获取用户指定期号的猜数字记录-读取缓存
func (srv *UserInGuessMarketActivityService) getUserInGuessMarketActivityCache(userInfoId int, issueNum string) (result *userhome_model.UserInGuessMarketActivity, err error) {
	cacheKey := EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivity_CacheKey + strconv.Itoa(userInfoId) + ":" + issueNum
	err = srv.RedisCache.GetJsonObj(cacheKey, &result)
	if err == redis.ErrNil {
		return nil, nil
	}
	return result, err
}

// refreshUserInGuessMarketActivity 获取用户指定期号的猜数字记录-刷新缓存
func (srv *UserInGuessMarketActivityService) refreshUserInGuessMarketActivity(userInfoId int, issueNum string) (result *userhome_model.UserInGuessMarketActivity, err error) {
	cacheKey := EMNET_UserInGuessMarketActivity_GetUserInGuessMarketActivity_CacheKey + strconv.Itoa(userInfoId) + ":" + issueNum
	result, err = srv.userInGuessMarketActivityRepo.GetUserInGuessMarketActivity(userInfoId, issueNum)
	if err != nil {
		return nil, err
	}
	if result != nil {
		err = srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(result), EMNET_UserInGuessMarketActivity_CacheSeconds)
	}
	return result, err
}

// GetWinner 获取获奖用户
func (srv *UserInGuessMarketActivityService) GetWinner(luckNum string, issueNum string) (result *userhome_model.UserInGuessMarketActivity, err error) {
	result, err = srv.userInGuessMarketActivityRepo.GetWinner(luckNum, issueNum)
	if err != nil {
		shareUserInGuessMarketActivityLogger.ErrorFormat(err, "GetWinner 获取获奖用户异常 luckNum=%s  issueNum=%s", luckNum, issueNum)
	}
	return result, err
}

// UpdateUserAward 更新用户奖品
func (srv *UserInGuessMarketActivityService) UpdateUserAward(userInfoId int, issueNum string, userInGuessMarketActivityId int64, awardId int64, awardName string) (n int64, err error) {
	n, err = srv.userInGuessMarketActivityRepo.UpdateUserAward(userInGuessMarketActivityId, awardId, awardName)
	if err != nil {
		shareUserInGuessMarketActivityLogger.ErrorFormat(err, "UpdateUserAward 更新用户奖品异常 userInGuessMarketActivityId=%d  awardId=%s awardName=%s", userInGuessMarketActivityId, awardId, awardName)
	} else {
		//刷新缓存
		go srv.refreshUserJoinCount(userInfoId)
		//中奖名单4条
		go srv.refreshUserInGuessMarketActivityGuessedList(4)
		//往期中奖名单6条
		go srv.refreshUserInGuessMarketActivityList(6, userInfoId)
		go srv.refreshUserInGuessMarketActivity(userInfoId, issueNum)
	}
	return n, err
}
