package service

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type UserInGuessChangeActivityService struct {
	service.BaseService
	userInGuessChangeActivityRepo *userhome_repo.UserInGuessChangeActivityRepository
}

var (
	shareUserInGuessChangeActivityRepo   *userhome_repo.UserInGuessChangeActivityRepository
	shareUserInGuessChangeActivityLogger dotlog.Logger
)

const (
	UserHomeUserInGuessChangeActivityServiceName                             = "UserInGuessChangeActivityService"
	EMNET_UserInGuessChangeActivity_CachePreKey                              = "EMNET:UserInGuessChangeActivityService"
	EMNET_UserInGuessChangeActivity_GetGuessTotal_CacheKey                   = EMNET_UserInGuessChangeActivity_CachePreKey + ":GetGuessTotal:"
	EMNET_UserInGuessChangeActivity_GetUserInGuessChangeActivity_CacheKey    = EMNET_UserInGuessChangeActivity_CachePreKey + ":GetUserInGuessChangeActivity:"
	EMNET_UserInGuessChangeActivity_GetMyGuessChangeInfoCurrentWeek_CacheKey = EMNET_UserInGuessChangeActivity_CachePreKey + ":GetMyGuessChangeInfoCurrentWeek:"
	EMNET_UserInGuessChangeActivity_GetMyGuessageAwardList_CacheKey          = EMNET_UserInGuessChangeActivity_CachePreKey + ":GetMyGuessageAwardList:"
	EMNET_UserInGuessChangeActivity_GetMyGuessChangeInfoNewst_CacheKey       = EMNET_UserInGuessChangeActivity_CachePreKey + ":GetMyGuessChangeInfoNewst:"
	EMNET_UserInGuessChangeActivity_GetUserJoinCount_CacheKey                = EMNET_UserInGuessChangeActivity_CachePreKey + ":GetUserJoinCount:"
	EMNET_UserInGuessChangeActivity_CacheSeconds                             = 15 * 60
)

func NewUserInGuessChangeActivityService() *UserInGuessChangeActivityService {
	srv := &UserInGuessChangeActivityService{
		userInGuessChangeActivityRepo: shareUserInGuessChangeActivityRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeUserInGuessChangeActivityServiceName, userInGuessChangeActivityServiceLoader)
}

func userInGuessChangeActivityServiceLoader() {
	shareUserInGuessChangeActivityRepo = userhome_repo.NewUserInGuessChangeActivityRepository(protected.DefaultConfig)
	shareUserInGuessChangeActivityLogger = dotlog.GetLogger(UserHomeUserInGuessChangeActivityServiceName)
}

// GetGuessTotal 获取指定一期的竞猜统计
func (srv *UserInGuessChangeActivityService) GetGuessTotal(issueNumber string) (total *GuessChangeTotal, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getGuessTotalDB(issueNumber)
	case config.ReadDB_CacheOrDB_UpdateCache:
		total, err := srv.getGuessTotalCache(issueNumber)
		if err == nil && total == nil {
			total, err = srv.refreshGuessTotal(issueNumber)
		}
		return total, err
	case config.ReadDB_RefreshCache:
		total, err = srv.refreshGuessTotal(issueNumber)
		return total, err
	default:
		return srv.getGuessTotalCache(issueNumber)
	}
}

// getGuessTotalDB 获取指定一期的竞猜统计-读取数据库
func (srv *UserInGuessChangeActivityService) getGuessTotalDB(issueNumber string) (total *GuessChangeTotal, err error) {
	upCount, downCount, err := srv.userInGuessChangeActivityRepo.GetGuessTotal(issueNumber)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "GetGuessTotal 获取指定一期的竞猜统计 异常")
	} else {
		total = &GuessChangeTotal{
			UpCount:   upCount,
			DownCount: downCount,
		}
	}
	return total, err
}

// getGuessTotalCache 获取指定一期的竞猜统计-读取缓存
func (srv *UserInGuessChangeActivityService) getGuessTotalCache(issueNumber string) (total *GuessChangeTotal, err error) {
	cacheKey := EMNET_UserInGuessChangeActivity_GetGuessTotal_CacheKey + issueNumber
	err = srv.RedisCache.GetJsonObj(cacheKey, &total)
	if err == redis.ErrNil {
		return nil, nil
	}
	return total, err
}

// refreshGuessTotal 获取指定一期的竞猜统计-刷新缓存
func (srv *UserInGuessChangeActivityService) refreshGuessTotal(issueNumber string) (total *GuessChangeTotal, err error) {
	cacheKey := EMNET_UserInGuessChangeActivity_GetGuessTotal_CacheKey + issueNumber
	total, err = srv.getGuessTotalDB(issueNumber)
	if err != nil {
		return nil, err
	}
	err = srv.RedisCache.Set(cacheKey, _json.GetJsonString(total), EMNET_UserInGuessChangeActivity_CacheSeconds)
	return total, err
}

// GetUserInGuessChangeActivity 获取指定用户在指定一期的参与记录
func (srv *UserInGuessChangeActivityService) GetUserInGuessChangeActivity(userInfoId int, issueNumber string) (userInGuessChangeActivity *userhome_model.UserInGuessChangeActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserInGuessChangeActivityDB(userInfoId, issueNumber)
	case config.ReadDB_CacheOrDB_UpdateCache:
		userInGuessChangeActivity, err := srv.getUserInGuessChangeActivityCache(userInfoId, issueNumber)
		if err == nil && userInGuessChangeActivity == nil {
			userInGuessChangeActivity, err = srv.refreshUserInGuessChangeActivity(userInfoId, issueNumber)
		}
		return userInGuessChangeActivity, err
	case config.ReadDB_RefreshCache:
		userInGuessChangeActivity, err = srv.refreshUserInGuessChangeActivity(userInfoId, issueNumber)
		return userInGuessChangeActivity, err
	default:
		return srv.getUserInGuessChangeActivityCache(userInfoId, issueNumber)
	}
}

// getUserInGuessChangeActivityDB 获取指定用户在指定一期的参与记录-读取数据库
func (srv *UserInGuessChangeActivityService) getUserInGuessChangeActivityDB(userInfoId int, issueNumber string) (userInGuessChangeActivity *userhome_model.UserInGuessChangeActivity, err error) {
	userInGuessChangeActivity, err = srv.userInGuessChangeActivityRepo.GetUserInGuessChangeActivity(userInfoId, issueNumber)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "getUserInGuessChangeActivityDB 获取指定用户在指定一期的参与记录-读取数据库 异常")
	}
	return userInGuessChangeActivity, err
}

// getUserInGuessChangeActivityCache 获取指定用户在指定一期的参与记录-读取缓存
func (srv *UserInGuessChangeActivityService) getUserInGuessChangeActivityCache(userInfoId int, issueNumber string) (userInGuessChangeActivity *userhome_model.UserInGuessChangeActivity, err error) {
	cacheKey := fmt.Sprintf("%s%d:%s", EMNET_UserInGuessChangeActivity_GetUserInGuessChangeActivity_CacheKey, userInfoId, issueNumber)
	err = srv.RedisCache.GetJsonObj(cacheKey, &userInGuessChangeActivity)
	if err == redis.ErrNil {
		return nil, nil
	}
	return userInGuessChangeActivity, err
}

// refreshUserInGuessChangeActivity 获取指定用户在指定一期的参与记录-刷新缓存
func (srv *UserInGuessChangeActivityService) refreshUserInGuessChangeActivity(userInfoId int, issueNumber string) (userInGuessChangeActivity *userhome_model.UserInGuessChangeActivity, err error) {
	cacheKey := fmt.Sprintf("%s%d:%s", EMNET_UserInGuessChangeActivity_GetUserInGuessChangeActivity_CacheKey, userInfoId, issueNumber)
	userInGuessChangeActivity, err = srv.getUserInGuessChangeActivityDB(userInfoId, issueNumber)
	if err != nil {
		return nil, err
	}
	if userInGuessChangeActivity != nil {
		err = srv.RedisCache.Set(cacheKey, _json.GetJsonString(userInGuessChangeActivity), EMNET_UserInGuessChangeActivity_CacheSeconds)
	}
	return userInGuessChangeActivity, err
}

// GetMyGuessChangeInfoCurrentWeek 获取用户本周的参与记录
func (srv *UserInGuessChangeActivityService) GetMyGuessChangeInfoCurrentWeek(userInfoId int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getMyGuessChangeInfoCurrentWeekDB(userInfoId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		myGuessChangeList, err := srv.getMyGuessChangeInfoCurrentWeekCache(userInfoId)
		if err == nil && myGuessChangeList == nil {
			myGuessChangeList, err = srv.refreshMyGuessChangeInfoCurrentWeek(userInfoId)
		}
		return myGuessChangeList, err
	case config.ReadDB_RefreshCache:
		myGuessChangeList, err = srv.refreshMyGuessChangeInfoCurrentWeek(userInfoId)
		return myGuessChangeList, err
	default:
		return srv.getMyGuessChangeInfoCurrentWeekCache(userInfoId)
	}
}

// getMyGuessChangeInfoCurrentWeekDB 获取用户本周的参与记录-读取数据库
func (srv *UserInGuessChangeActivityService) getMyGuessChangeInfoCurrentWeekDB(userInfoId int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	date := time.Now()
	myGuessChangeList, err = srv.userInGuessChangeActivityRepo.GetUserInGuessChangeActivityListCurrentWeek(userInfoId, date)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "getMyGuessChangeInfoCurrentWeekDB 获取用户本周的参与记录-读取数据库 异常")
	}
	return myGuessChangeList, err
}

// getMyGuessChangeInfoCurrentWeekCache 获取用户本周的参与记录-读取缓存
func (srv *UserInGuessChangeActivityService) getMyGuessChangeInfoCurrentWeekCache(userInfoId int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	cacheKey := EMNET_UserInGuessChangeActivity_GetMyGuessChangeInfoCurrentWeek_CacheKey + strconv.Itoa(userInfoId)
	err = srv.RedisCache.GetJsonObj(cacheKey, &myGuessChangeList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return myGuessChangeList, err
}

// refreshMyGuessChangeInfoCurrentWeek 获取用户本周的参与记录-刷新缓存
func (srv *UserInGuessChangeActivityService) refreshMyGuessChangeInfoCurrentWeek(userInfoId int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	cacheKey := EMNET_UserInGuessChangeActivity_GetMyGuessChangeInfoCurrentWeek_CacheKey + strconv.Itoa(userInfoId)
	myGuessChangeList, err = srv.getMyGuessChangeInfoCurrentWeekDB(userInfoId)
	if err != nil {
		return nil, err
	}
	if myGuessChangeList != nil {
		err = srv.RedisCache.Set(cacheKey, _json.GetJsonString(myGuessChangeList), EMNET_UserInGuessChangeActivity_CacheSeconds)
	}
	return myGuessChangeList, err
}

// GetMyGuessageAwardList 获取用户得到的奖品列表
func (srv *UserInGuessChangeActivityService) GetMyGuessageAwardList(userInfoId int, top int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getMyGuessageAwardListDB(userInfoId, top)
	case config.ReadDB_CacheOrDB_UpdateCache:
		myGuessChangeList, err := srv.getMyGuessageAwardListCache(userInfoId, top)
		if err == nil && myGuessChangeList == nil {
			myGuessChangeList, err = srv.refreshMyGuessageAwardList(userInfoId, top)
		}
		return myGuessChangeList, err
	case config.ReadDB_RefreshCache:
		myGuessChangeList, err = srv.refreshMyGuessageAwardList(userInfoId, top)
		return myGuessChangeList, err
	default:
		return srv.getMyGuessageAwardListCache(userInfoId, top)
	}
}

// getMyGuessageAwardListDB 获取用户得到的奖品列表-读取数据库
func (srv *UserInGuessChangeActivityService) getMyGuessageAwardListDB(userInfoId int, top int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	myGuessChangeList, err = srv.userInGuessChangeActivityRepo.GetUserInGuessChangeActivityAwardList(userInfoId, top)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "getMyGuessageAwardListDB 获取用户得到的奖品列表-读取数据库 异常")
	}
	return myGuessChangeList, err
}

// getMyGuessageAwardListCache 获取用户得到的奖品列表-读取缓存
func (srv *UserInGuessChangeActivityService) getMyGuessageAwardListCache(userInfoId int, top int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	cacheKey := fmt.Sprintf("%s%d:%d", EMNET_UserInGuessChangeActivity_GetMyGuessageAwardList_CacheKey, userInfoId, top)
	err = srv.RedisCache.GetJsonObj(cacheKey, &myGuessChangeList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return myGuessChangeList, err
}

// refreshMyGuessageAwardList 获取用户得到的奖品列表-刷新缓存
func (srv *UserInGuessChangeActivityService) refreshMyGuessageAwardList(userInfoId int, top int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	cacheKey := fmt.Sprintf("%s%d:%d", EMNET_UserInGuessChangeActivity_GetMyGuessageAwardList_CacheKey, userInfoId, top)
	myGuessChangeList, err = srv.getMyGuessageAwardListDB(userInfoId, top)
	if err != nil {
		return nil, err
	}
	if myGuessChangeList != nil {
		err = srv.RedisCache.Set(cacheKey, _json.GetJsonString(myGuessChangeList), EMNET_UserInGuessChangeActivity_CacheSeconds)
	}
	return myGuessChangeList, err
}

// GetMyGuessChangeInfoNewst 获取用户近期的竞猜记录
func (srv *UserInGuessChangeActivityService) GetMyGuessChangeInfoNewst(userInfoId int, top int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getMyGuessChangeInfoNewstDB(userInfoId, top)
	case config.ReadDB_CacheOrDB_UpdateCache:
		myGuessChangeList, err := srv.getMyGuessChangeInfoNewstCache(userInfoId, top)
		if err == nil && myGuessChangeList == nil {
			myGuessChangeList, err = srv.refreshMyGuessChangeInfoNewst(userInfoId, top)
		}
		return myGuessChangeList, err
	case config.ReadDB_RefreshCache:
		myGuessChangeList, err = srv.refreshMyGuessChangeInfoNewst(userInfoId, top)
		return myGuessChangeList, err
	default:
		return srv.getMyGuessChangeInfoNewstCache(userInfoId, top)
	}
}

// getMyGuessChangeInfoNewstDB 获取用户近期的竞猜记录-读取数据库
func (srv *UserInGuessChangeActivityService) getMyGuessChangeInfoNewstDB(userInfoId int, top int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	myGuessChangeList, err = srv.userInGuessChangeActivityRepo.GetUserInGuessChangeActivityList(userInfoId, top)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "GetMyGuessChangeInfoNewst 获取用户近期的竞猜记录-读取数据库 异常")
	}
	return myGuessChangeList, err
}

// getMyGuessChangeInfoNewstCache 获取用户近期的竞猜记录-读取缓存
func (srv *UserInGuessChangeActivityService) getMyGuessChangeInfoNewstCache(userInfoId int, top int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	cacheKey := fmt.Sprintf("%s%d:%d", EMNET_UserInGuessChangeActivity_GetMyGuessChangeInfoNewst_CacheKey, userInfoId, top)
	err = srv.RedisCache.GetJsonObj(cacheKey, &myGuessChangeList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return myGuessChangeList, err
}

// refreshMyGuessChangeInfoNewst 获取用户近期的竞猜记录-刷新缓存
func (srv *UserInGuessChangeActivityService) refreshMyGuessChangeInfoNewst(userInfoId int, top int) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	cacheKey := fmt.Sprintf("%s%d:%d", EMNET_UserInGuessChangeActivity_GetMyGuessChangeInfoNewst_CacheKey, userInfoId, top)
	myGuessChangeList, err = srv.getMyGuessChangeInfoNewstDB(userInfoId, top)
	if err != nil {
		return nil, err
	}
	if myGuessChangeList != nil {
		err = srv.RedisCache.Set(cacheKey, _json.GetJsonString(myGuessChangeList), EMNET_UserInGuessChangeActivity_CacheSeconds)
	}
	return myGuessChangeList, err
}

// GuessChangeSubmit 新增用户竞猜记录
func (srv *UserInGuessChangeActivityService) GuessChangeSubmit(issueNumber string, userInfoId int, uid int64, nickName string, result int, activityCycle string, activityId int64) (err error) {
	err = srv.userInGuessChangeActivityRepo.InsertUserInGuessChangeActivity(issueNumber, userInfoId, uid, nickName, result, activityCycle)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "GuessChangeSubmit 新增用户竞猜记录 异常 issueNumber=%s userInfoId=%d uid=%d nickName=%d result=%d activityCycle=%s activityId=%d", issueNumber, userInfoId, uid, nickName, result, activityCycle, activityId)
		return err
	} else {
		srv.refreshGuessTotal(issueNumber)
		go srv.refreshMyGuessageAwardList(userInfoId, 5)
		go srv.refreshMyGuessChangeInfoCurrentWeek(userInfoId)
		go srv.refreshUserInGuessChangeActivity(userInfoId, issueNumber)
		go srv.refreshMyGuessChangeInfoNewst(userInfoId, 10)
	}
	//添加用户活动参与记录
	usrInActivitySrv := NewUserInActivityService()
	_, err = usrInActivitySrv.InsertUserInActivity(userInfoId, activityId)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "GuessChangeSubmit 新增用户竞猜记录 写入用户活动参与记录【异常】")
	}
	return err
}

// PublishUserGuessChangeResult 公布用户竞猜结果
func (srv *UserInGuessChangeActivityService) PublishUserGuessChangeResult(issueNumber string, result int) (count int64, err error) {
	count, err = srv.userInGuessChangeActivityRepo.PublishUserGuessChangeResult(issueNumber, result)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "PublishUserGuessChangeResult 公布用户竞猜结果 异常 issueNumber=%s result=%d", issueNumber, result)
	}
	return count, err
}

// UpdateUserGuessChangeAward 发放奖品
func (srv *UserInGuessChangeActivityService) UpdateUserGuessChangeAward(startIssueNumber string, endIssueNumber string, awardId int64, awardName string, reportUrl string) (count int64, err error) {
	count, err = srv.userInGuessChangeActivityRepo.UpdateUserGuessChangeAward(startIssueNumber, endIssueNumber, awardId, awardName, reportUrl)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "UpdateUserGuessChangeAward 发放奖品 异常 startIssueNumber=%s endIssueNumber=%s awardId=%d awardName=%s reportUrl=%s", startIssueNumber, endIssueNumber, awardId, awardName, reportUrl)
	}
	return count, err
}

// GetUserJoinCount 获取用户参与次数
func (srv *UserInGuessChangeActivityService) GetUserJoinCount(userInfoId int) (count int64, err error) {
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

// getUserJoinCountDB 获取用户参与次数-查询数据库
func (srv *UserInGuessChangeActivityService) getUserJoinCountDB(userInfoId int) (count int64, err error) {
	count, err = srv.userInGuessChangeActivityRepo.GetUserJoinCount(userInfoId)
	if err != nil {
		shareUserInGuessChangeActivityLogger.ErrorFormat(err, "getUserJoinCountDB 获取用户参与次数-查询数据库 异常 userInfoId=%d", userInfoId)
	}
	return count, err
}

// GetUserJoinCount 获取用户参与次数
func (srv *UserInGuessChangeActivityService) getUserJoinCountCache(userInfoId int) (count int64, err error) {
	cacheKey := EMNET_UserInGuessChangeActivity_GetUserJoinCount_CacheKey + strconv.Itoa(userInfoId)
	err = srv.RedisCache.GetJsonObj(cacheKey, &count)
	return count, err
}

// GetUserJoinCount 获取用户参与次数
func (srv *UserInGuessChangeActivityService) refreshUserJoinCount(userInfoId int) (count int64, err error) {
	cacheKey := EMNET_UserInGuessChangeActivity_GetUserJoinCount_CacheKey + strconv.Itoa(userInfoId)
	count, err = srv.getUserJoinCountDB(userInfoId)
	if err != nil {
		return 0, err
	}
	err = srv.RedisCache.Set(cacheKey, _json.GetJsonString(count), EMNET_UserInGuessChangeActivity_CacheSeconds)
	return count, err
}

//猜涨跌统计
type GuessChangeTotal struct {
	UpCount   int64
	DownCount int64
}
