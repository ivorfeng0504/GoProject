package service

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dataapi"
	"git.emoney.cn/softweb/roboadvisor/protected/service/funcapi"
	"git.emoney.cn/softweb/roboadvisor/protected/service/user"
	"git.emoney.cn/softweb/roboadvisor/util/array"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
)

type GuessMarketActivityService struct {
	service.BaseService
	guessMarketActivityRepo *userhome_repo.GuessMarketActivityRepository
}

var (
	shareGuessMarketActivityRepo   *userhome_repo.GuessMarketActivityRepository
	shareGuessMarketActivityLogger dotlog.Logger
)

const (
	GuessMarketActivityServiceName                                = "GuessMarketActivityService"
	EMNET_GuessMarketActivity_CachePreKey                         = "EMNET:GuessMarketActivityService:"
	EMNET_GuessMarketActivity_GetGuessMarketActivityList_CacheKey = EMNET_GuessMarketActivity_CachePreKey + "GetGuessMarketActivityList:"
	EMNET_GuessMarketActivity_GetCurrentActivity_CacheKey         = EMNET_GuessMarketActivity_CachePreKey + "GetCurrentActivity:"
	EMNET_GuessMarketActivity_NumQueue_CacheKey                   = EMNET_GuessMarketActivity_CachePreKey + "NumQueue"
	EMNET_GuessMarketActivity_CacheSeconds                        = 15 * 60
)

func NewGuessMarketActivityService() *GuessMarketActivityService {
	srv := &GuessMarketActivityService{
		guessMarketActivityRepo: shareGuessMarketActivityRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(GuessMarketActivityServiceName, guessMarketResultServiceLoader)
}

func guessMarketResultServiceLoader() {
	shareGuessMarketActivityRepo = userhome_repo.NewGuessMarketActivityRepository(protected.DefaultConfig)
	shareGuessMarketActivityLogger = dotlog.GetLogger(GuessMarketActivityServiceName)
}

// GetGuessMarketActivityList 查询最新的N条开奖信息，top小于等于0则查询所有
func (srv *GuessMarketActivityService) GetGuessMarketActivityList(top int) (resultList []*userhome_model.GuessMarketActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getGuessMarketActivityListDB(top)
	case config.ReadDB_CacheOrDB_UpdateCache:
		resultList, err = srv.getGuessMarketActivityListCache(top)
		if err == nil && resultList == nil {
			resultList, err = srv.refreshGuessMarketActivityList(top)
		}
		return resultList, err
	case config.ReadDB_RefreshCache:
		resultList, err = srv.refreshGuessMarketActivityList(top)
		return resultList, err
	default:
		return srv.getGuessMarketActivityListCache(top)
	}
}

// getGuessMarketActivityListDB 查询最新的N条开奖信息，top小于等于0则查询所有-读取数据库
func (srv *GuessMarketActivityService) getGuessMarketActivityListDB(top int) (resultList []*userhome_model.GuessMarketActivity, err error) {
	resultList, err = srv.guessMarketActivityRepo.GetGuessMarketActivityList(top)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "查询最新的N条开奖信息异常")
	}
	return resultList, err
}

// getGuessMarketActivityListCache 查询最新的N条开奖信息，top小于等于0则查询所有-读取缓存
func (srv *GuessMarketActivityService) getGuessMarketActivityListCache(top int) (resultList []*userhome_model.GuessMarketActivity, err error) {
	cacheKey := EMNET_GuessMarketActivity_GetGuessMarketActivityList_CacheKey + strconv.Itoa(top)
	err = srv.RedisCache.GetJsonObj(cacheKey, &resultList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return resultList, err
}

// refreshGuessMarketActivityList 查询最新的N条开奖信息，top小于等于0则查询所有-刷新缓存
func (srv *GuessMarketActivityService) refreshGuessMarketActivityList(top int) (resultList []*userhome_model.GuessMarketActivity, err error) {
	cacheKey := EMNET_GuessMarketActivity_GetGuessMarketActivityList_CacheKey + strconv.Itoa(top)
	resultList, err = srv.getGuessMarketActivityListDB(top)
	if err != nil {
		return nil, err
	}
	if resultList != nil {
		err = srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(resultList), EMNET_GuessMarketActivity_CacheSeconds)
	}
	return resultList, err
}

// GetCurrentActivity 查询当前的活动，正在进行或者已经结束
func (srv *GuessMarketActivityService) GetCurrentActivity() (activity *userhome_model.GuessMarketActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getCurrentActivityDB()
	case config.ReadDB_CacheOrDB_UpdateCache:
		activity, err = srv.getCurrentActivityCache()
		if err == nil && activity == nil {
			activity, err = srv.refreshCurrentActivity()
		}
		return activity, err
	case config.ReadDB_RefreshCache:
		activity, err = srv.refreshCurrentActivity()
		return activity, err
	default:
		return srv.getCurrentActivityCache()
	}
}

// getCurrentActivityDB 查询当前的活动，正在进行或者已经结束-读取数据库
func (srv *GuessMarketActivityService) getCurrentActivityDB() (activity *userhome_model.GuessMarketActivity, err error) {
	activity, err = srv.guessMarketActivityRepo.GetCurrentBeginingActivity()
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "查询当前的活动异常getCurrentActivityDB->GetCurrentBeginingActivity")
		return nil, err
	}
	if activity == nil {
		activity, err = srv.guessMarketActivityRepo.GetCurrentFinishedActivity()
	}
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "查询当前的活动异常getCurrentActivityDB->GetCurrentFinishedActivity")
		return nil, err
	}
	return activity, err
}

// getCurrentActivityCache 查询当前的活动，正在进行或者已经结束-读取缓存
func (srv *GuessMarketActivityService) getCurrentActivityCache() (activity *userhome_model.GuessMarketActivity, err error) {
	cacheKey := EMNET_GuessMarketActivity_GetCurrentActivity_CacheKey + time.Now().Format("20060102")
	err = srv.RedisCache.GetJsonObj(cacheKey, &activity)
	if err == redis.ErrNil {
		return nil, nil
	}
	return activity, err
}

// refreshCurrentActivity 查询当前的活动，正在进行或者已经结束-刷新缓存
func (srv *GuessMarketActivityService) refreshCurrentActivity() (activity *userhome_model.GuessMarketActivity, err error) {
	cacheKey := EMNET_GuessMarketActivity_GetCurrentActivity_CacheKey + time.Now().Format("20060102")
	activity, err = srv.getCurrentActivityDB()
	if err != nil {
		return nil, err
	}
	if activity != nil {
		err = srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(activity), EMNET_GuessMarketActivity_CacheSeconds)
	}
	return activity, err
}

// InsertActivity 插入新一期的活动
func (srv *GuessMarketActivityService) InsertActivity(issueNumber string, beginTime string, endTime string) (id int64, err error) {
	id, err = srv.guessMarketActivityRepo.InsertActivity(issueNumber, beginTime, endTime)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "InsertActivity 插入新一期领数字活动异常 issueNumber=%s", issueNumber)
	} else {
		//刷新缓存
		go srv.refreshCurrentActivity()
		//中奖名单
		go srv.refreshGuessMarketActivityList(6)
	}
	return id, err
}

// GrantAward 自动发放奖励
// 每日盘后判断是否是最后一个交易日 如果是则发放奖励
func (srv *GuessMarketActivityService) GrantAward() (issueNumber string, err error) {
	isTradeDay := false
	isTradeDay, err = dataapi.IsTradeDayToDay()
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 IsTradeDayToDay判断交易异常")
		return issueNumber, err
	}
	//如果不是交易日 则不进行任何处理
	if isTradeDay == false {
		shareGuessMarketActivityLogger.Debug("GrantAward 当前不是交易日 不需要发放奖励")
		return issueNumber, nil
	}
	now := time.Now()
	//获取本周最后一个交易日
	lastTradeDayThisWeek, err := dataapi.GetLastTradeWeekDay(now)
	if err == dataapi.ErrNoTradeDay {
		shareGuessMarketActivityLogger.Debug("GrantAward 本周没有交易日 不需要发放奖励")
		return issueNumber, nil
	}
	if lastTradeDayThisWeek.Weekday() != now.Weekday() {
		shareGuessMarketActivityLogger.Debug("GrantAward 今天不是最后一个交易日 不需要发放奖励")
		return issueNumber, nil
	}

	issueNumber, _, _, err = srv.GetCurrentIssueNumber()
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 GetCurrentIssueNumber获取当前期号异常")
		return issueNumber, err
	}
	if issueNumber != now.Format("20060102") {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 当前期号不是今日")
		return issueNumber, err
	}
	//获取活动信息
	activityId := config.CurrentConfig.Activity_GuessNumber
	if activityId <= 0 {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 Activity_GuessNumber  未配置 activityId=%d", activityId)
		return issueNumber, err
	}
	activitySrv := NewActivityService()
	activity, err := activitySrv.GetActivityById(activityId)
	if activity == nil {
		err = errors.New("未查询到活动信息 活动Id=" + strconv.FormatInt(activityId, 10))
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 GetActivityById  活动不存在 activityId=%d", activityId)
		return issueNumber, err
	}
	//获取猜数字活动信息
	guessMarketActivitySrv := NewGuessMarketActivityService()
	guessMarketActivity, err := guessMarketActivitySrv.guessMarketActivityRepo.GetGuessMarketActivityByIssueNumber(issueNumber)
	if guessMarketActivity == nil {
		err = errors.New("未查询到猜数字活动信息issueNumber =" + issueNumber)
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 GetGuessMarketActivityByIssueNumber  猜数字活动信息不存在 issueNumber=%s", issueNumber)
		return issueNumber, err
	}

	//判断是否已开奖
	if guessMarketActivity.IsPublish {
		err = errors.New("该期活动已开奖 无需再次开奖 issueNumber =" + issueNumber)
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 该期活动已开奖  issueNumber=%s", issueNumber)
		return issueNumber, err
	}

	//今天是最后的交易日 发放奖励
	//计算得奖号码
	luckNum, err := srv.getLuckNum(time.Time(guessMarketActivity.EndTime))
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 getLuckNum获取开奖数字失败")
		return issueNumber, err
	}

	//公布开奖号码
	_, err = srv.PublishActivityResult(issueNumber, luckNum)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 PublishActivityResult 公布开奖号码异常")
		return issueNumber, err
	}
	userActivitySrv := NewUserInGuessMarketActivityService()
	//获取得奖的用户
	winner, err := userActivitySrv.GetWinner(luckNum, issueNumber)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 GetWinner获取得奖的用户")
		return issueNumber, err
	}
	if winner == nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 GetWinner  未获取到得奖用户")
		return issueNumber, err
	}

	//获取奖励内容
	awardInActivitySrv := NewAwardInActivityService()
	activityAwardList, err := awardInActivitySrv.GetActivityAwardListByActivityId(activityId, true)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 GetActivityAwardListByActivityId 获取活动奖品失败")
		return issueNumber, err
	}
	if activityAwardList == nil || len(activityAwardList) == 0 {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 GetActivityAwardListByActivityId 活动奖品不存在")
		return issueNumber, err
	}

	//记录用户奖励信息
	awardInfo := activityAwardList[0]
	_, err = userActivitySrv.UpdateUserAward(winner.UserInfoId, winner.IssueNumber, winner.UserInGuessMarketActivityId, awardInfo.AwardId, awardInfo.AwardName)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 UpdateUserAward 设置更新用户奖品异常 winner=%s awardInfo=%s", jsonutil.GetJsonString(winner), jsonutil.GetJsonString(awardInfo))
		return issueNumber, err
	}

	userAwardSrv := NewUserAwardService()
	userAwardSrv.InsertUserAward(winner.UserInfoId, awardInfo.AwardName, awardInfo.AwardId, awardInfo.IntroduceVideo, awardInfo.ActivityId, activity.Title, _const.UserAwardState_Receive, awardInfo.AwardImg, awardInfo.AvailableDay, awardInfo.AwardType)

	//如果是特权 自动发放
	if awardInfo.AwardType == _const.AwardType_Function {
		shareGuessMarketActivityLogger.DebugFormat("GrantAward 自动发放奖励【开始发放特权】 activityId=%s issueNum=%s funcId=%d winner.UID=%d enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, winner.UID, awardInfo.AvailableDay)
		//根据GID查询用户的手机号
		userService := user.NewUserService()
		gid, err := userService.CidFindGid(winner.UID)
		if err != nil {
			shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 CidFindGid  特权发放异常  根据Cid获取Gid异常 activityId=%s issueNum=%s funcId=%d uid=%s enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, strconv.FormatInt(winner.UID, 10), awardInfo.AvailableDay)
			return issueNumber, err
		}
		if gid <= 0 {
			gid = winner.UID
		}
		shareGuessMarketActivityLogger.DebugFormat("GrantAward 自动发放奖励【开始发放特权】【gid已处理】 activityId=%s issueNum=%s funcId=%d winner.UID=%d  gid=%d enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, winner.UID, gid, awardInfo.AvailableDay)
		boundAccount, err := userService.BoundGroupQryLogin(gid)
		if err != nil {
			shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 BoundGroupQryLogin  特权发放异常  查询用户手机号异常 activityId=%s issueNum=%s funcId=%d gid=%d enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, gid, awardInfo.AvailableDay)
			return issueNumber, err
		}

		//如果用户没有手机号 则不进行特权发放 记录日志
		if boundAccount == nil || len(boundAccount) == 0 {
			shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 BoundGroupQryLogin  特权发放异常  未查询到用户手机号 activityId=%s issueNum=%s funcId=%d gid=%d enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, gid, awardInfo.AvailableDay)
			return issueNumber, err
		}
		encryptMobile := ""
		for _, boundInfo := range boundAccount {
			if boundInfo.AccountType == 1 && boundInfo.AccountName != "" {
				encryptMobile = boundInfo.EncryptMobile
			}
		}
		if len(encryptMobile) == 0 {
			shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 BoundGroupQryLogin  特权发放异常  未查询到用户手机号 activityId=%s issueNum=%s funcId=%d gid=%d enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, gid, awardInfo.AvailableDay)
			return issueNumber, err
		}
		//查询到加密手机号  进行特权发放
		success, err := funcapi.OpenJRPTFunc(awardInfo.JRPTFunc, encryptMobile, awardInfo.AvailableDay)
		if err != nil {
			shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 OpenJRPTFunc  特权发放异常 activityId=%s issueNum=%s funcId=%d gid=%d enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, gid, awardInfo.AvailableDay)
			return issueNumber, err
		}
		if success == false {
			err = errors.New("特权发放失败")
			shareGuessMarketActivityLogger.ErrorFormat(err, "GrantAward 自动发放奖励【异常】 OpenJRPTFunc  特权发放异常 activityId=%s issueNum=%s funcId=%d gid=%d enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, gid, awardInfo.AvailableDay)
			return issueNumber, err
		} else {
			shareGuessMarketActivityLogger.DebugFormat("GrantAward 自动发放奖励【发放特权成功】 activityId=%s issueNum=%s funcId=%d winner.UID=%d  gid=%d enableDays=%d", strconv.FormatInt(activityId, 10), issueNumber, awardInfo.JRPTFunc, winner.UID, gid, awardInfo.AvailableDay)
		}
	}
	return issueNumber, nil
}

// InitGuessMarketActivity 初始化新一期的领数字活动
// 如果当前不存在新一期的活动则创建新的活动
// 并且清除上一期的领号队列，初始化新的领号队列（0000-9999）
// isNext 是否初始化下期活动  true则初始化下期活动  false则初始化本周当期活动（一般用于第一次初始化）
func (srv *GuessMarketActivityService) InitGuessMarketActivity(isNext bool) (issueNumber string, err error) {
	activity, err := srv.guessMarketActivityRepo.GetCurrentBeginingActivity()
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "初始化新一期的领数字活动【异常】 InitGuessMarketActivity->GetCurrentBeginingActivity")
		return issueNumber, err
	}
	if activity != nil {
		shareGuessMarketActivityLogger.Debug("初始化新一期的领数字活动 InitGuessMarketActivity 当前活动正在进行，无需初始化")
		issueNumber = activity.IssueNumber
		return issueNumber, nil
	}
	//插入新一期活动记录
	nextIssueNumber, nextIssueDate, beginDate, err := srv.GetNextIssueNumber()

	if isNext == false {
		nextIssueNumber, nextIssueDate, beginDate, err = srv.GetCurrentIssueNumber()
	}
	issueNumber = nextIssueNumber
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "初始化新一期的领数字活动【异常】 InitGuessMarketActivity->GetNextIssueNumber")
		return issueNumber, err
	}
	//查询新一期的活动是否已经初始化
	activity, err = srv.guessMarketActivityRepo.GetGuessMarketActivityByIssueNumber(nextIssueNumber)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "初始化新一期的领数字活动【异常】 InitGuessMarketActivity->guessMarketActivityRepo.GetGuessMarketActivityByIssueNumber")
		return issueNumber, err
	}

	//如果记录已经存在 则不再进行处理
	if activity != nil {
		shareGuessMarketActivityLogger.DebugFormat("InitGuessMarketActivity 新一期的活动已经创建完毕 nextIssueNumber=%s", nextIssueNumber)
		return issueNumber, nil
	}

	week := int(nextIssueDate.Weekday())
	if week == 0 {
		week = 7
	}
	//开始时间应该为这周的第一个交易日的九点
	_, err = srv.InsertActivity(nextIssueNumber, beginDate.Format("2006-01-02")+" 09:00", nextIssueDate.Format("2006-01-02")+" 15:30")
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "初始化新一期的领数字活动【异常】 InitGuessMarketActivity->InsertActivity")
		return issueNumber, err
	}
	//删除上一期的队列
	srv.RedisCache.Delete(EMNET_GuessMarketActivity_NumQueue_CacheKey)

	//初始化新一期的队列
	//队列最大数量
	queueLen := 10000
	//每个数字的长度
	numLen := 4
	shuffleArr := _array.GetShuffleArray(queueLen)
	var shuffleStrArr []interface{}
	for i := 0; i < queueLen; i++ {
		item := strconv.Itoa(shuffleArr[i])
		//不足四位补0
		item = strings.Repeat("0", numLen-len(item)) + item
		shuffleStrArr = append(shuffleStrArr, item)
	}
	//每次500个，批量写入redis
	for i := 0; i < queueLen; i += 500 {
		_, err = srv.RedisCache.RPush(EMNET_GuessMarketActivity_NumQueue_CacheKey, shuffleStrArr[i:i+500]...)
		if err != nil {
			shareGuessMarketActivityLogger.ErrorFormat(err, "初始化新一期的领数字活动【异常】 写入队列出现异常 InitGuessMarketActivity->RedisCache.RPush")
			return issueNumber, err
		}
	}
	return issueNumber, err
}

// GetCurrentIssueNumber 获取当前的期号
func (srv *GuessMarketActivityService) GetCurrentIssueNumber() (issueNumber string, issueDate time.Time, beginDate time.Time, err error) {
	//获取本周5的日期
	now := time.Now()
	issueDate = time.Now()
	week := int(now.Weekday())
	if week == 0 {
		week = 7
	}
	//获取本周五的日期
	issueDate = now.AddDate(0, 0, 5-week)
	isTradeDay := false
	maxCall := 14
	for isTradeDay == false {
		if maxCall <= 0 {
			break
		}
		isTradeDay, err = dataapi.IsTradeDay(issueDate)
		if err != nil {
			shareGuessMarketActivityLogger.ErrorFormat(err, "获取当前的期号异常 判断交易日异常 GetCurrentIssueNumber-->IsTradeDay  date=%s", issueDate)
			return issueNumber, issueDate, beginDate, err
		}
		if isTradeDay == false {
			issueDate = issueDate.AddDate(0, 0, -1)
		}
		maxCall--
	}

	if isTradeDay {
		issueNumber = issueDate.Format("20060102")

		//获取到下一期的开奖日后 查询这周的第一个交易日作为活动的开始日期
		beginDate = issueDate.AddDate(0, 0, -(int(issueDate.Weekday()) - 1))
		findFirstTrade := false
		for i := 0; i < int(issueDate.Weekday())-1; i++ {
			isTradeDay, err = dataapi.IsTradeDay(beginDate)
			if err != nil {
				shareGuessMarketActivityLogger.ErrorFormat(err, "GetNextIssueNumber 获取当前第一个交易日异常 判断交易日异常-->IsTradeDay  date=%s", issueDate)
				return issueNumber, issueDate, beginDate, err
			}
			if isTradeDay {
				findFirstTrade = true
				break
			}
			beginDate = beginDate.AddDate(0, 0, -1)
		}
		if findFirstTrade == false {
			err = errors.New("获取当前的期号异常 未找到第一个交易日 findFirstTrade=false")
			shareGuessMarketActivityLogger.ErrorFormat(err, err.Error())
			return issueNumber, issueDate, beginDate, err
		}
		return issueNumber, issueDate, beginDate, nil
	} else {
		return issueNumber, issueDate, beginDate, errors.New("未获取到交易日，查询领数字活动期号失败")
	}
}

// GetNextIssueNumber 获取下一期的期号
func (srv *GuessMarketActivityService) GetNextIssueNumber() (issueNumber string, issueDate time.Time, beginDate time.Time, err error) {
	//获取本周5的日期
	now := time.Now()
	issueDate = time.Now()
	week := int(now.Weekday())
	if week == 0 {
		week = 7
	}
	//获取下周五的日期
	issueDate = now.AddDate(0, 0, 5-week+7)
	//最小日期为下周一
	minDate := issueDate.AddDate(0, 0, 5-1)
	isTradeDay := false
	maxCall := 14
	for isTradeDay == false {
		if maxCall <= 0 {
			break
		}
		isTradeDay, err = dataapi.IsTradeDay(issueDate)
		if err != nil {
			shareGuessMarketActivityLogger.ErrorFormat(err, "GetNextIssueNumber 获取当前第一个交易日异常 判断交易日异常 GetCurrentIssueNumber-->IsTradeDay  date=%s", issueDate)
			return issueNumber, issueDate, beginDate, err
		}
		if isTradeDay == false {
			if issueDate.After(minDate) {
				issueDate = issueDate.AddDate(0, 0, -1)
			} else {
				//如果当前时间小于等于下周一  则从下下周五开始判断交易日
				issueDate = issueDate.AddDate(0, 0, 11)
			}
		}
		maxCall--
	}
	if isTradeDay {
		issueNumber = issueDate.Format("20060102")

		//获取到下一期的开奖日后 查询这周的第一个交易日作为活动的开始日期
		beginDate = issueDate.AddDate(0, 0, -(int(issueDate.Weekday()) - 1))
		findFirstTrade := false
		for i := 0; i < int(issueDate.Weekday())-1; i++ {
			isTradeDay, err = dataapi.IsTradeDay(beginDate)
			if err != nil {
				shareGuessMarketActivityLogger.ErrorFormat(err, "GetNextIssueNumber 获取当前第一个交易日异常 判断交易日异常-->IsTradeDay  date=%s", issueDate)
				return issueNumber, issueDate, beginDate, err
			}
			if isTradeDay {
				findFirstTrade = true
				break
			}
			beginDate = beginDate.AddDate(0, 0, -1)
		}
		if findFirstTrade == false {
			err = errors.New("获取当前的期号异常 未找到第一个交易日 findFirstTrade=false")
			shareGuessMarketActivityLogger.ErrorFormat(err, err.Error())
			return issueNumber, issueDate, beginDate, err
		}
		return issueNumber, issueDate, beginDate, nil
	} else {
		return issueNumber, issueDate, beginDate, errors.New("未获取到交易日，查询领数字活动期号失败")
	}
}

// GetNumber 获取活动数字
func (srv *GuessMarketActivityService) GetNumber(userInfo userhome_model.UserInfo, activityId int64) (number string, err error) {
	activity, err := srv.guessMarketActivityRepo.GetCurrentBeginingActivity()
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "查询当前活动信息【异常】 GetNumber->GetCurrentBeginingActivity")
		return number, err
	}
	if activity == nil {
		return number, errors.New("当前活动未开始")
	}

	// 查询数字是否已被领取  如果没有领取则写入数据库
	usrInGuessSrv := NewUserInGuessMarketActivityService()
	joinRecord, err := usrInGuessSrv.GetUserInGuessMarketActivity(userInfo.ID, activity.IssueNumber)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "获取用户参与记录异常 userInfoId=%d issueNumber=%s", userInfo.ID, activity.IssueNumber)
		return number, err
	}
	//已经抽过奖 则直接返回抽过的数字
	if joinRecord != nil {
		number = joinRecord.Result
		return number, nil
	}
	//没有抽奖记录 则开始领取数字
	number, err = srv.RedisCache.RPop(EMNET_GuessMarketActivity_NumQueue_CacheKey)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GetNumber 获取活动数字【异常】")
		return number, err
	}
	//记录用户领取的数字
	_, err = usrInGuessSrv.Insert(userInfo.UID, userInfo.ID, userInfo.NickName, activity.IssueNumber, number, activity.GuessMarketActivityId)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GetNumber 写入用户领取记录【异常】")
		return number, err
	}
	//添加用户活动参与记录
	usrInActivitySrv := NewUserInActivityService()
	_, err = usrInActivitySrv.InsertUserInActivity(userInfo.ID, activityId)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "GetNumber 写入用户活动参与记录【异常】")
	}
	return number, err
}

// getLuckNum 获取开奖号码
func (srv *GuessMarketActivityService) getLuckNum(endDate time.Time) (luckNum string, err error) {
	sh000001Price, _, tradeDate, err := dataapi.GetSH000001ClosePrice()
	if err != nil {
		return luckNum, err
	}
	year, month, day := tradeDate.Date()
	year2, month2, day2 := endDate.Date()
	if year != year2 || month != month2 || day != day2 {
		errMsg := fmt.Sprintf("SH000001 未获取到期望的收盘价格  期望的时间为 %s  获取到的收盘价格时间为 %s", endDate, tradeDate)
		err = errors.New(errMsg)
		shareGuessMarketActivityLogger.ErrorFormat(err, errMsg)
		return luckNum, err
	}
	sz399001Price, _, tradeDate, err := dataapi.GetSZ399001ClosePrice()
	if err != nil {
		return luckNum, err
	}
	year, month, day = tradeDate.Date()
	if year != year2 || month != month2 || day != day2 {
		errMsg := fmt.Sprintf("SZ399001 未获取到期望的收盘价格  期望的时间为 %s  获取到的收盘价格时间为 %s", endDate, tradeDate)
		err = errors.New(errMsg)
		shareGuessMarketActivityLogger.ErrorFormat(err, errMsg)
		return luckNum, err
	}
	//上证小数两位
	arr := strings.Split(sh000001Price, ".")
	if len(arr) == 2 {
		luckNum = strings.Repeat("0", 2-len(arr[1])) + arr[1]
	} else {
		luckNum = "00"
	}

	//深证小数两位
	arr = strings.Split(sz399001Price, ".")
	if len(arr) == 2 {
		luckNum = luckNum + strings.Repeat("0", 2-len(arr[1])) + arr[1]
	} else {
		luckNum = luckNum + "00"
	}
	return luckNum, nil
}

// PublishActivityResult 公布开奖号码
func (srv *GuessMarketActivityService) PublishActivityResult(issueNumber string, luckNum string) (n int64, err error) {

	n, err = srv.guessMarketActivityRepo.PublishActivityResult(issueNumber, luckNum)
	if err != nil {
		shareGuessMarketActivityLogger.ErrorFormat(err, "PublishActivityResult 公布开奖号码异常 issueNumber=%s luckNum=%s", issueNumber, luckNum)
	} else {
		//刷新缓存
		go srv.refreshCurrentActivity()
		go srv.refreshGuessMarketActivityList(6)
	}
	return n, err
}
