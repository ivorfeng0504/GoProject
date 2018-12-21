package service

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type UserReceivedStockService struct {
	service.BaseService
	userReceivedStockRepo *userhome_repo.UserReceivedStockRepository
}

var (
	shareUserReceivedStockRepo   *userhome_repo.UserReceivedStockRepository
	shareUserReceivedStockLogger dotlog.Logger
)

const (
	UserReceivedStockServiceName                         = "UserReceivedStockService"
	EMNET_UserReceivedStock_CachePreKey                  = "EMNET:UserReceivedStockService:"
	EMNET_UserReceivedStock_GetUserStockToday_CacheKey   = EMNET_UserReceivedStock_CachePreKey + "GetUserStockToday:"
	EMNET_UserReceivedStock_GetUserStockHistory_CacheKey = EMNET_UserReceivedStock_CachePreKey + "GetUserStockHistory:"
	EMNET_UserReceivedStock_CacheSeconds                 = 15 * 60
)

func NewUserReceivedStockService() *UserReceivedStockService {
	srv := &UserReceivedStockService{
		userReceivedStockRepo: shareUserReceivedStockRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserReceivedStockServiceName, userReceivedStockServiceLoader)
}

func userReceivedStockServiceLoader() {
	shareUserReceivedStockRepo = userhome_repo.NewUserReceivedStockRepository(protected.DefaultConfig)
	shareUserReceivedStockLogger = dotlog.GetLogger(UserReceivedStockServiceName)
}

// InsertUserStock 用户领取股票
func (srv *UserReceivedStockService) InsertUserStock(userStock *userhome_model.UserReceivedStock, activityId int64) (id int64, err error) {
	id, err = srv.userReceivedStockRepo.InsertUserStock(userStock)
	if err != nil {
		shareUserReceivedStockLogger.ErrorFormat(err, "InsertUserStock 用户领取股票 异常 userStock=%s", _json.GetJsonString(userStock))
	} else {
		//添加用户活动参与记录
		go func() {
			usrInActivitySrv := NewUserInActivityService()
			_, err2 := usrInActivitySrv.InsertUserInActivity(userStock.UserInfoId, activityId)
			if err2 != nil {
				shareUserReceivedStockLogger.ErrorFormat(err, "InsertUserStock 写入用户活动参与记录【异常】")
			}
		}()
		//刷新缓存
		go srv.refreshUserStockToday(userStock.UserInfoId)
		//前端需求 历史领取显示10条
		go srv.refreshUserStockHistory(userStock.UserInfoId, 10)
	}
	return id, err
}

// GetUserStockToday 获取用户当日领取的股票
func (srv *UserReceivedStockService) GetUserStockToday(userInfoId int) (userStock *userhome_model.UserReceivedStock, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserStockTodayDB(userInfoId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		userStock, err = srv.getUserStockTodayCache(userInfoId)
		if err == nil && userStock == nil {
			userStock, err = srv.refreshUserStockToday(userInfoId)
		}
		return userStock, err
	case config.ReadDB_RefreshCache:
		userStock, err = srv.refreshUserStockToday(userInfoId)
		return userStock, err
	default:
		return srv.getUserStockTodayCache(userInfoId)
	}
}

// getUserStockTodayDB 获取用户当日领取的股票-读取数据库
func (srv *UserReceivedStockService) getUserStockTodayDB(userInfoId int) (userStock *userhome_model.UserReceivedStock, err error) {
	userStock, err = srv.userReceivedStockRepo.GetUserStockToday(userInfoId)
	if err != nil {
		shareUserReceivedStockLogger.ErrorFormat(err, "获取用户当日领取的股票 异常 userInfoId=%d", userInfoId)
	}
	return userStock, err
}

// UpdateUserReceivedStock 更新用户领取的研报
func (srv *UserReceivedStockService) UpdateUserReceivedStock(userInfoId int, userReceivedStockId int64, reportUrl string) (err error) {
	err = srv.userReceivedStockRepo.UpdateUserReceivedStock(userReceivedStockId, reportUrl)
	if err != nil {
		shareUserReceivedStockLogger.ErrorFormat(err, "UpdateUserReceivedStock 更新用户领取的研报 异常 userInfoId=%d userReceivedStockId=%d reportUrl=%s", userInfoId, userReceivedStockId, reportUrl)
	} else {
		srv.refreshUserStockToday(userInfoId)
	}
	return err
}

// getUserStockTodayCache 获取用户当日领取的股票-读取缓存
func (srv *UserReceivedStockService) getUserStockTodayCache(userInfoId int) (userStock *userhome_model.UserReceivedStock, err error) {
	cacheKey := EMNET_UserReceivedStock_GetUserStockToday_CacheKey + ":" + strconv.Itoa(userInfoId) + ":" + time.Now().Format("20060102")
	err = srv.RedisCache.GetJsonObj(cacheKey, &userStock)
	if err == redis.ErrNil {
		return nil, nil
	}
	return userStock, err
}

// refreshUserStockToday 获取用户当日领取的股票-刷新缓存
func (srv *UserReceivedStockService) refreshUserStockToday(userInfoId int) (userStock *userhome_model.UserReceivedStock, err error) {
	cacheKey := EMNET_UserReceivedStock_GetUserStockToday_CacheKey + ":" + strconv.Itoa(userInfoId) + ":" + time.Now().Format("20060102")
	userStock, err = srv.getUserStockTodayDB(userInfoId)
	if err != nil {
		return nil, err
	}
	if userStock != nil {
		srv.RedisCache.Set(cacheKey, _json.GetJsonString(userStock), EMNET_UserReceivedStock_CacheSeconds)
	}
	return userStock, err
}

// GetUserStockHistory 获取用户领取的股票历史
func (srv *UserReceivedStockService) GetUserStockHistory(userInfoId int, top int) (userStockList []*userhome_model.UserReceivedStock, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserStockHistoryDB(userInfoId, top)
	case config.ReadDB_CacheOrDB_UpdateCache:
		userStockList, err = srv.getUserStockHistoryCache(userInfoId, top)
		if err == nil && userStockList == nil {
			userStockList, err = srv.refreshUserStockHistory(userInfoId, top)
		}
		return userStockList, err
	case config.ReadDB_RefreshCache:
		userStockList, err = srv.refreshUserStockHistory(userInfoId, top)
		return userStockList, err
	default:
		return srv.getUserStockHistoryCache(userInfoId, top)
	}
}

// getUserStockHistoryDB 获取用户领取的股票历史-读取数据库
func (srv *UserReceivedStockService) getUserStockHistoryDB(userInfoId int, top int) (userStockList []*userhome_model.UserReceivedStock, err error) {
	userStockList, err = srv.userReceivedStockRepo.GetUserStockHistory(userInfoId, top)
	if err != nil {
		shareUserReceivedStockLogger.ErrorFormat(err, "获取用户领取的股票历史-读取数据库 异常 userInfoId=%d top=%d", userInfoId, top)
	}
	return userStockList, err
}

// getUserStockHistoryCache 获取用户领取的股票历史-读取缓存
func (srv *UserReceivedStockService) getUserStockHistoryCache(userInfoId int, top int) (userStockList []*userhome_model.UserReceivedStock, err error) {
	cacheKey := fmt.Sprintf("%s%d:%d", EMNET_UserReceivedStock_GetUserStockHistory_CacheKey, userInfoId, top)
	err = srv.RedisCache.GetJsonObj(cacheKey, &userStockList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return userStockList, err
}

// refreshUserStockHistory 获取用户领取的股票历史-刷新缓存
func (srv *UserReceivedStockService) refreshUserStockHistory(userInfoId int, top int) (userStockList []*userhome_model.UserReceivedStock, err error) {
	cacheKey := fmt.Sprintf("%s%d:%d", EMNET_UserReceivedStock_GetUserStockHistory_CacheKey, userInfoId, top)
	userStockList, err = srv.getUserStockHistoryDB(userInfoId, top)
	if err != nil {
		return nil, err
	}
	if userStockList != nil {
		srv.RedisCache.Set(cacheKey, _json.GetJsonString(userStockList), EMNET_UserReceivedStock_CacheSeconds)
	}
	return userStockList, err
}
