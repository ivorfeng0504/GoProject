package service

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dataapi"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/util/time"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type GuessChangeActivityService struct {
	service.BaseService
	guessChangeActivityRepo *userhome_repo.GuessChangeActivityRepository
}

var (
	shareGuessChangeActivityRepo   *userhome_repo.GuessChangeActivityRepository
	shareGuessChangeActivityLogger dotlog.Logger
)

const (
	UserHomeGuessChangeActivityServiceName                           = "GuessChangeActivityService"
	EMNET_GuessChangeActivity_CachePreKey                            = "EMNET:GuessChangeActivityService"
	EMNET_GuessChangeActivity_GetCurrentGuessChangeActivity_CacheKey = EMNET_GuessChangeActivity_CachePreKey + ":GetCurrentGuessChangeActivity:"
	EMNET_GuessChangeActivity_CacheSeconds                           = 15 * 60
)

func NewGuessChangeActivityService() *GuessChangeActivityService {
	srv := &GuessChangeActivityService{
		guessChangeActivityRepo: shareGuessChangeActivityRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeGuessChangeActivityServiceName, guessChangeActivityServiceLoader)
}

func guessChangeActivityServiceLoader() {
	shareGuessChangeActivityRepo = userhome_repo.NewGuessChangeActivityRepository(protected.DefaultConfig)
	shareGuessChangeActivityLogger = dotlog.GetLogger(UserHomeGuessChangeActivityServiceName)
}

// GetCurrentGuessChangeActivity 获取当前最新一期的猜涨跌活动
func (srv *GuessChangeActivityService) GetCurrentGuessChangeActivity() (activity *userhome_model.GuessChangeActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getCurrentGuessChangeActivityDB()
	case config.ReadDB_CacheOrDB_UpdateCache:
		activity, err = srv.getCurrentGuessChangeActivityCache()
		if err == nil && activity == nil {
			activity, err = srv.refreshCurrentGuessChangeActivity()
		}
		return activity, err
	case config.ReadDB_RefreshCache:
		activity, err = srv.refreshCurrentGuessChangeActivity()
		return activity, err
	default:
		return srv.getCurrentGuessChangeActivityCache()
	}
}

// getCurrentGuessChangeActivityCache 获取当前最新一期的猜涨跌活动-读取缓存
func (srv *GuessChangeActivityService) getCurrentGuessChangeActivityCache() (activity *userhome_model.GuessChangeActivity, err error) {
	cacheKey := EMNET_GuessChangeActivity_GetCurrentGuessChangeActivity_CacheKey + time.Now().Format("20060102")
	err = srv.RedisCache.GetJsonObj(cacheKey, &activity)
	if err == redis.ErrNil {
		return nil, nil
	}
	return activity, err
}

// getCurrentGuessChangeActivityDB 获取当前最新一期的猜涨跌活动-读取数据库
func (srv *GuessChangeActivityService) getCurrentGuessChangeActivityDB() (activity *userhome_model.GuessChangeActivity, err error) {
	activity, err = srv.guessChangeActivityRepo.GetNewstGuessChangeActivity()
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "GetCurrentGuessChangeActivity 获取当前最新一期的猜涨跌活动 异常")
	}
	return
}

// refreshCurrentGuessChangeActivity 获取当前最新一期的猜涨跌活动-刷新缓存
func (srv *GuessChangeActivityService) refreshCurrentGuessChangeActivity() (activity *userhome_model.GuessChangeActivity, err error) {
	cacheKey := EMNET_GuessChangeActivity_GetCurrentGuessChangeActivity_CacheKey + time.Now().Format("20060102")
	activity, err = srv.getCurrentGuessChangeActivityDB()
	if err != nil {
		return nil, err
	}
	if activity != nil {
		err = srv.RedisCache.Set(cacheKey, _json.GetJsonString(activity), EMNET_GuessChangeActivity_CacheSeconds)
	}
	return activity, err
}

// InsertGuessChangeActivity 新增一个猜涨跌活动
func (srv *GuessChangeActivityService) InsertGuessChangeActivity(issueNumber string, beginTime time.Time, endTime time.Time, activityCycle string) (err error) {
	err = srv.guessChangeActivityRepo.InsertGuessChangeActivity(issueNumber, beginTime, endTime, activityCycle)
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "InsertGuessChangeActivity 新增一个猜涨跌活动 异常 issueNumber=%s", issueNumber)
	} else {
		srv.refreshCurrentGuessChangeActivity()
	}
	return err
}

// getGuessChangeActivityByIssueNumber 获取指定一期的猜涨跌活动-读取数据库
func (srv *GuessChangeActivityService) getGuessChangeActivityByIssueNumberDB(issueNumber string) (activity *userhome_model.GuessChangeActivity, err error) {
	activity, err = srv.guessChangeActivityRepo.GetGuessChangeActivityByIssueNumber(issueNumber)
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "getGuessChangeActivityByIssueNumber 获取指定一期的猜涨跌活动-读取数据库 异常 issueNumber=%s", issueNumber)
	}
	return activity, err
}

// InitNextGuessChangeActivity 初始化下一期猜涨跌活动
func (srv *GuessChangeActivityService) InitNextGuessChangeActivity() (nextIssueNum string, err error) {
	now := time.Now()
	nextTradeDate, err := dataapi.GetNextTradeDay(now)
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "InitNextGuessChangeActivity 初始化下一期猜涨跌活动 异常 获取下一个交易日异常")
		return nextIssueNum, err
	} else {
		shareGuessChangeActivityLogger.DebugFormat("InitNextGuessChangeActivity 初始化下一期猜涨跌活动 nextTradeDate=%s", nextTradeDate.Format("2006-01-02"))
	}
	nextIssueNum = nextTradeDate.Format("20060102")
	activity, err := srv.getGuessChangeActivityByIssueNumberDB(nextIssueNum)
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "InitNextGuessChangeActivity 初始化下一期猜涨跌活动 获取下一期活动异常 nextIssueNum=%s", nextIssueNum)
		return nextIssueNum, err
	}
	//下期活动已经存在 无需再次初始化
	if activity != nil {
		shareGuessChangeActivityLogger.DebugFormat("InitNextGuessChangeActivity 初始化下一期猜涨跌活动 下期活动已经存在 无需再次初始化")
		return nextIssueNum, nil
	}
	beginTime, err := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02")+" 09:00:00")
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "InitNextGuessChangeActivity 初始化下一期猜涨跌活动 Parse beginTime 异常 nextIssueNum=%s", nextIssueNum)
		return nextIssueNum, err
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", nextTradeDate.Format("2006-01-02")+" 08:59:59")
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "InitNextGuessChangeActivity 初始化下一期猜涨跌活动 Parse endTime 异常 nextIssueNum=%s", nextIssueNum)
		return nextIssueNum, err
	}
	monday, friday := _time.GetMondayAndFriday(nextTradeDate)
	activityCycle := monday.Format("1.2") + "-" + friday.Format("1.2")
	err = srv.InsertGuessChangeActivity(nextIssueNum, beginTime, endTime, activityCycle)
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "InitNextGuessChangeActivity 初始化下一期猜涨跌活动 InsertGuessChangeActivity 异常 nextIssueNum=%s activityCycle=%s", nextIssueNum, activityCycle)
		return nextIssueNum, err
	}
	return nextIssueNum, err
}

// PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品
func (srv *GuessChangeActivityService) PublishCurrentGuessChangeActivity(activityId int64, issueDate time.Time) (issueNum string, err error) {
	issueNum = issueDate.Format("20060102")
	isTradeDate, err := dataapi.IsTradeDay(issueDate)
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 判断交易日异常 issueNum=%s", issueNum)
		return issueNum, err
	}
	//非交易日 无需开奖
	if isTradeDate == false {
		shareGuessChangeActivityLogger.DebugFormat("PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 非交易日 无需开奖 issueNum=%s", issueNum)
		return issueNum, nil
	}
	_, change, tradeDate, err := dataapi.GetSH000001ClosePrice()
	year, month, day := tradeDate.Date()
	year2, month2, day2 := issueDate.Date()
	if year != year2 || month != month2 || day != day2 {
		errMsg := fmt.Sprintf("SH000001 未获取到期望的收盘价格  期望的时间为 %s  获取到的收盘价格时间为 %s", issueDate, tradeDate)
		err = errors.New(errMsg)
		shareGuessChangeActivityLogger.ErrorFormat(err, errMsg)
		return issueNum, nil
	}
	//公布本期活动的竞猜结果
	f, err := strconv.ParseFloat(change, 32)
	if err != nil {
		shareGuessChangeActivityLogger.DebugFormat("PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 ParseFloat异常 issueNum=%s change=%s", issueNum, change)
		return issueNum, err
	}
	result := 0
	if f > 0 {
		result = 1
	} else if f < 0 {
		result = -1
	}
	err = srv.guessChangeActivityRepo.PublishGuessChangeActivity(issueNum, result)
	if err != nil {
		shareGuessChangeActivityLogger.DebugFormat("PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 PublishGuessChangeActivity 异常 issueNum=%s change=%s result=%d", issueNum, change, result)
		return issueNum, err
	}

	//修改用户的竞猜状态
	userInGuessChangeSrv := NewUserInGuessChangeActivityService()
	_, err = userInGuessChangeSrv.PublishUserGuessChangeResult(issueNum, result)
	if err != nil {
		shareGuessChangeActivityLogger.ErrorFormat(err, "PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 PublishUserGuessChangeResult异常 issueNum=%s result=%d", issueNum, result)
		return issueNum, err
	}

	//如果weekday大于等于3 发放研报奖励
	if issueDate.Weekday() >= 3 {
		//查询当前活动的奖品
		activityAwardSrv := NewAwardInActivityService()
		awardList, err := activityAwardSrv.GetActivityAwardListByActivityIdDB(activityId, true)
		if err != nil {
			shareGuessChangeActivityLogger.ErrorFormat(err, "PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 发放研报奖励异常 获取奖品失败 issueNum=%s result=%d activityId=%d", issueNum, result, activityId)
			return issueNum, err
		}
		if awardList == nil || len(awardList) == 0 {
			err = errors.New("未获取到有效的研报奖品")
			shareGuessChangeActivityLogger.ErrorFormat(err, "PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 发放研报奖励异常 获取奖品为空 issueNum=%s result=%d activityId=%d", issueNum, result, activityId)
			return issueNum, err
		}
		award := awardList[0]
		monday, friday := _time.GetMondayAndFriday(issueDate)
		if len(award.LinkUrl) == 0 {
			err = errors.New("研报地址不能为空")
			shareGuessChangeActivityLogger.ErrorFormat(err, "PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 发放研报奖励异常 研报地址不能为空 issueNum=%s result=%d activityId=%d awardInfo=%s", issueNum, result, activityId, _json.GetJsonString(award))
			return issueNum, err
		}
		_, err = userInGuessChangeSrv.UpdateUserGuessChangeAward(monday.Format("20060102"), friday.Format("20060102"), award.AwardId, award.AwardName, award.LinkUrl)
		if err != nil {
			shareGuessChangeActivityLogger.ErrorFormat(err, "PublishCurrentGuessChangeActivity 公布当期猜涨跌结果，并且发放奖品 UpdateUserGuessChangeAward 发放奖励失败 issueNum=%s result=%d activityId=%d awardInfo=%s", issueNum, result, activityId, _json.GetJsonString(award))
			return issueNum, err
		}
	}
	return issueNum, err
}
