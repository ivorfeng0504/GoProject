package live

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	liverepo "git.emoney.cn/softweb/roboadvisor/protected/repository/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dbsync"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type LiveQuestionAnswerService struct {
	service.BaseService
	liveQuestionAnswerRepo *liverepo.LiveQuestionAnswerRepository
}

var (
	// shareLiveQuestionAnswerRepo 共享的仓储
	shareLiveQuestionAnswerRepo *liverepo.LiveQuestionAnswerRepository

	// shareLiveQuestionAnswerLogger 共享的Logger实例
	shareLiveQuestionAnswerLogger dotlog.Logger
)

const (
	liveQuestionAnswerServiceName                              = "LiveQuestionAnswerService"
	EMNET_LiveQuestionAnswer_CachePreKey                       = "EMoney:Live:LiveQuestionAnswer"
	EMNET_LiveQuestionAnswer_LiveQuestionAnswerListCacheKey    = EMNET_LiveQuestionAnswer_CachePreKey + ":LiveQuestionAnswerListV1.16"
	EMNET_LiveQuestionAnswer_TodayRepliedQuestionNumCacheKey   = EMNET_LiveQuestionAnswer_CachePreKey + ":TodayRepliedQuestionNum"
	EMNET_LiveQuestionAnswer_TodayALLQuestionNumCacheKey       = EMNET_LiveQuestionAnswer_CachePreKey + ":TodayALLQuestionNum"
	EMNET_LiveQuestionAnswer_TodayNORepliedQuestionNumCacheKey = EMNET_LiveQuestionAnswer_CachePreKey + ":TodayNORepliedQuestionNum"
	EMNET_LiveQuestionAnswer_GetMyAskCount_CacheKey            = EMNET_LiveQuestionAnswer_CachePreKey + ":GetMyAskCount:"
	EMNET_LiveQuestionAnswer_MyAskCacheKey                     = EMNET_LiveQuestionAnswer_CachePreKey + ":LiveMyAsk"
)

func NewLiveQuestionAnswerService() *LiveQuestionAnswerService {
	service := &LiveQuestionAnswerService{
		liveQuestionAnswerRepo: shareLiveQuestionAnswerRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return service
}

// GetLiveQuestionAnswerListByAnswered 获取指定直播间某天已回答的问答列表
func (service *LiveQuestionAnswerService) GetLiveQuestionAnswerListByAnswered(roomId int, date time.Time) (answerList []*livemodel.LiveQuestionAnswer, err error) {
	questionList, err := service.GetLiveQuestionAnswerList(roomId, date, 0)
	if err != nil {
		return nil, err
	}

	//筛选出已经回答并且不是私聊的问答
	for _, question := range questionList {
		if question.Status == 2 && question.IsPrivate == 0 {
			answerList = append(answerList, question)
		}
	}
	return answerList, nil
}

// GetLiveQuestionAnswerList 获取指定直播间某天的问答列表，如果UserId>0则筛选指定用户的提问列表
func (service *LiveQuestionAnswerService) GetLiveQuestionAnswerList(roomId int, date time.Time, userId int) (answerList []*livemodel.LiveQuestionAnswer, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getLiveQuestionAnswerListDB(roomId, date, userId)
	default:
		return service.getLiveQuestionAnswerListCache(roomId, date, userId)
	}
}

// getLiveQuestionAnswerListCache 获取指定直播间某天的问答列表，如果UserId>0则筛选指定用户的提问列表
func (service *LiveQuestionAnswerService) getLiveQuestionAnswerListCache(roomId int, date time.Time, userId int) (answerList []*livemodel.LiveQuestionAnswer, err error) {
	var answerListSource []*livemodel.LiveQuestionAnswer
	dateStr := date.Format("2006/01/02")
	cacheKey := EMNET_LiveQuestionAnswer_LiveQuestionAnswerListCacheKey + dateStr + strconv.Itoa(roomId)
	err = service.RedisCache.GetJsonObj(cacheKey, &answerListSource)
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		shareLiveQuestionAnswerLogger.Error(err, " getLiveQuestionAnswerListCache 获取直播间问答异常，cachekey="+cacheKey)
		return nil, err
	}
	if answerListSource == nil || len(answerListSource) == 0 || userId <= 0 {
		return answerListSource, nil
	}

	//如果UserId>0则筛选指定用户的提问列表
	for _, answer := range answerListSource {
		if answer.AskUserId == userId {
			answerList = append(answerList, answer)
		}
	}
	return answerList, nil
}

// getLiveQuestionAnswerListDB 获取指定直播间某天的问答列表，如果UserId>0则筛选指定用户的提问列表
func (service *LiveQuestionAnswerService) getLiveQuestionAnswerListDB(roomId int, date time.Time, userId int) (answerList []*livemodel.LiveQuestionAnswer, err error) {
	return service.liveQuestionAnswerRepo.GetLiveQuestionAnswerList(roomId, date, userId)
}

// AddLiveQuestionAnswer 用户提问
func (service *LiveQuestionAnswerService) AddLiveQuestionAnswer(askUserId int, askUserName string, askContent string, roomId int, source string) (id int, err error) {
	//屏蔽字过滤
	maskService := NewMaskWordService()
	askContent, err = maskService.ProcessMaskWord(askContent)
	if err != nil {
		return -1, err
	}
	id, err = service.liveQuestionAnswerRepo.AddLiveQuestionAnswer(askUserId, askUserName, askContent, roomId, source)
	if err != nil {
		shareLiveQuestionAnswerLogger.Error(err, "AddLiveQuestionAnswer 用户提问 异常，askUserId="+strconv.Itoa(askUserId))
		return id, err
	}
	//提问成功更新缓存
	go service.freshCacheLiveQuestionAnswer(roomId)
	go service.getRepliedQuestionNum(roomId)
	go service.getTodayAllQuestionNum(roomId)
	go service.getNoReplyQuestionNum(roomId)
	go service.freshCacheLiveMyAsk(askUserId)
	go service.refreshMyAskCount(askUserId)
	//数据库同步
	dbsync.Sync(_const.SyncTable_LiveQuestionAnswer)
	return id, err
}

// getRepliedQuestionNum 获取并刷新今日已回复问题总数-缓存刷新 from .net LiveQuestionAnswerBll
func (service *LiveQuestionAnswerService) getRepliedQuestionNum(roomId int) int64 {
	dateStr := time.Now().Format("2006/01/02")
	result, err := service.liveQuestionAnswerRepo.GetTodayRepliedQuestionNum(roomId)
	if err != nil {
		shareLiveQuestionAnswerLogger.Error(err, "刷新缓存异常 getRepliedQuestionNum roomId="+strconv.Itoa(roomId))
	} else {
		cacheKey := EMNET_LiveQuestionAnswer_TodayRepliedQuestionNumCacheKey + dateStr + strconv.Itoa(roomId)
		service.RedisCache.SetJsonObj(cacheKey, result)
	}
	return result
}

// getTodayAllQuestionNum 获取并刷新今日所有问题总数-缓存刷新 from .net LiveQuestionAnswerBll
func (service *LiveQuestionAnswerService) getTodayAllQuestionNum(roomId int) int64 {
	dateStr := time.Now().Format("2006/01/02")
	result, err := service.liveQuestionAnswerRepo.GetTodayAllQuestionNum(roomId)
	if err != nil {
		shareLiveQuestionAnswerLogger.Error(err, "刷新缓存异常 getTodayAllQuestionNum roomId="+strconv.Itoa(roomId))
	} else {
		cacheKey := EMNET_LiveQuestionAnswer_TodayALLQuestionNumCacheKey + dateStr + strconv.Itoa(roomId)
		service.RedisCache.SetJsonObj(cacheKey, result)
	}
	return result
}

// getNoReplyQuestionNum 获取并刷新今日未回复问题总数-缓存刷新 from .net LiveQuestionAnswerBll
func (service *LiveQuestionAnswerService) getNoReplyQuestionNum(roomId int) int64 {
	dateStr := time.Now().Format("2006/01/02")
	result, err := service.liveQuestionAnswerRepo.GetTodayNoRepliedQuestionNum(roomId)
	if err != nil {
		shareLiveQuestionAnswerLogger.Error(err, "刷新缓存异常 getNoReplyQuestionNum roomId="+strconv.Itoa(roomId))
	} else {
		cacheKey := EMNET_LiveQuestionAnswer_TodayNORepliedQuestionNumCacheKey + dateStr + strconv.Itoa(roomId)
		service.RedisCache.SetJsonObj(cacheKey, result)
	}
	return result
}

// freshCacheLiveMyAsk 持久化问答数据(个人所有的问答)-缓存刷新 from .net LiveQuestionAnswerBll
func (service *LiveQuestionAnswerService) freshCacheLiveMyAsk(askUserId int) (answerList []*livemodel.LiveQuestionAnswer, err error) {
	answerList, err = service.liveQuestionAnswerRepo.GetAllLiveQuestionAnswerListByAskUserId(askUserId)
	if err != nil {
		shareLiveQuestionAnswerLogger.Error(err, "刷新缓存异常 freshCacheLiveMyAsk askUserId="+strconv.Itoa(askUserId))
		return answerList, err
	}
	cacheKey := EMNET_LiveQuestionAnswer_MyAskCacheKey + strconv.Itoa(askUserId)
	_, err = service.RedisCache.SetJsonObj(cacheKey, answerList)
	return answerList, err
}

// freshCacheLiveQuestionAnswer 获取并刷新直播间当天问答列表-缓存刷新 from .net LiveQuestionAnswerBll
func (service *LiveQuestionAnswerService) freshCacheLiveQuestionAnswer(roomId int) (answerList []*livemodel.LiveQuestionAnswer, err error) {
	//提问成功更新缓存
	dateStr := time.Now().Format("2006/01/02")
	cacheKey := EMNET_LiveQuestionAnswer_LiveQuestionAnswerListCacheKey + dateStr + strconv.Itoa(roomId)
	answerList, err = service.getLiveQuestionAnswerListDB(roomId, time.Now(), 0)
	if err != nil {
		shareLiveQuestionAnswerLogger.Error(err, "刷新缓存异常 freshCacheLiveQuestionAnswer roomId="+strconv.Itoa(roomId))
		return nil, err
	}
	_, err = service.RedisCache.SetJsonObj(cacheKey, answerList)
	return answerList, err
}

// GetMyAskCount 获取某个用户的提问总数
func (srv *LiveQuestionAnswerService) GetMyAskCount(askUserId int) (count int64, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return srv.getMyAskCountDB(askUserId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		count, err = srv.getMyAskCountCache(askUserId)
		if err == nil {
			count, err = srv.refreshMyAskCount(askUserId)
		}
		return count, err
	case config.ReadDB_RefreshCache:
		count, err = srv.refreshMyAskCount(askUserId)
		return count, err
	default:
		return srv.getMyAskCountCache(askUserId)
	}
}

// getMyAskCountDB 获取某个用户的提问总数-读取数据库
func (srv *LiveQuestionAnswerService) getMyAskCountDB(askUserId int) (count int64, err error) {
	count, err = srv.liveQuestionAnswerRepo.GetMyAskCount(askUserId)
	return count, err
}

// getMyAskCountCache 获取某个用户的提问总数-读取缓存
func (srv *LiveQuestionAnswerService) getMyAskCountCache(askUserId int) (count int64, err error) {
	cacheKey := EMNET_LiveQuestionAnswer_GetMyAskCount_CacheKey + strconv.Itoa(askUserId)
	err = srv.RedisCache.GetJsonObj(cacheKey, &count)
	if err == redis.ErrNil {
		return 0, nil
	}
	return count, err
}

// refreshMyAskCount 获取某个用户的提问总数-刷新缓存
func (srv *LiveQuestionAnswerService) refreshMyAskCount(askUserId int) (count int64, err error) {
	cacheKey := EMNET_LiveQuestionAnswer_GetMyAskCount_CacheKey + strconv.Itoa(askUserId)
	count, err = srv.getMyAskCountDB(askUserId)
	if err != nil {
		return count, err
	}
	srv.RedisCache.SetJsonObj(cacheKey, count)
	return count, err
}

// CheckPID 检查用户是否可以提问  liveRoomId暂不使用 20181018变更
func (srv *LiveQuestionAnswerService) CheckPID(pid int, userId int, liveRoomId int) (canIAsk bool, tip string, err error) {
	//canIAsk = false
	//tip = "当前用户不能提问哦"
	//体验版（888012000,888012001,888022000） 每日最多提问3次，只能提问3个月
	//体验过期 付费过期 游客 （888012400,888010400,888022400,888020400，888030000,888030400）无权限提问
	//智盈大师正式版（888020000，888000000） 588(888010000)无限制
	unlimitList := map[int]bool{888020000: true, 888000000: true, 888010000: true}
	limitTip := map[int]string{ /*888010000: "抱歉，您今日的问答体验已达上限，欢迎明日互动。", */ 888012000: "抱歉，您今日的问答体验已达上限，欢迎明日互动。", 888012001: "抱歉，您今日的问答体验已达上限，欢迎明日互动。", 888022000: "抱歉，您今日的问答体验已达上限，欢迎明日互动。"}
	noPermitTip := map[int]string{888012400: "抱歉，您的问答权限已经到期，详情请询客服400-602-1015", 888010400: "抱歉，您的问答权限已经到期，详情请询客服400-602-1015", 888022400: "抱歉，您的问答体验期已经结束，详情请询客服400-670-6668", 888020400: "抱歉，您的问答体验期已经结束，详情请询客服400-670-6668", 888030000: "为了保障您互动问答的体验完整，请先登陆软件。", 888030400: "为了保障您互动问答的体验完整，请先登陆软件。"}
	if unlimitList[pid] {
		canIAsk = true
		shareLiveQuestionAnswerLogger.DebugFormat("CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【不受限】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	}
	tip, exist := noPermitTip[pid]
	if exist {
		canIAsk = false
		shareLiveQuestionAnswerLogger.DebugFormat("CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【无权用户】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	}
	//todayAskCacheKey := fmt.Sprintf("%s:%d:%d", time.Now().Format("20060102"), liveRoomId, pid)
	//firstAskCacheKey := fmt.Sprintf("first:%d:%d", liveRoomId, pid)
	todayAskCacheKey := fmt.Sprintf("%s:%d", time.Now().Format("20060102"), pid)
	firstAskCacheKey := fmt.Sprintf("first:%d", pid)
	todayAskCountStr, err := srv.RedisCache.HGet(todayAskCacheKey, strconv.Itoa(userId))
	if err != nil && err != redis.ErrNil {
		shareLiveQuestionAnswerLogger.ErrorFormat(err, "CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【异常】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	}
	todayAskCount := 0
	if err != redis.ErrNil {
		todayAskCount, err = strconv.Atoi(todayAskCountStr)
	}
	if todayAskCount >= 3 {
		canIAsk = false
		tip = "每天只能提问3次！"
		if len(limitTip[pid]) > 0 {
			tip = limitTip[pid]
		}
		shareLiveQuestionAnswerLogger.DebugFormat("CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【超出限制】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	}
	firstAskDateStr, err := srv.RedisCache.HGet(firstAskCacheKey, strconv.Itoa(userId))
	if err != nil && err != redis.ErrNil {
		shareLiveQuestionAnswerLogger.ErrorFormat(err, "CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【异常】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	}
	if err == redis.ErrNil || len(firstAskDateStr) == 0 {
		canIAsk = true
		shareLiveQuestionAnswerLogger.DebugFormat("CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【限制范围内】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	}
	firstAskDate, err := time.Parse("2006-01-02 15:04:05", firstAskDateStr)
	if err != nil {
		shareLiveQuestionAnswerLogger.ErrorFormat(err, "CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【异常】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	}
	if firstAskDate.AddDate(0, 3, 0).After(time.Now()) {
		canIAsk = true
		shareLiveQuestionAnswerLogger.DebugFormat("CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【限制范围内】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	} else {
		canIAsk = false
		tip = "抱歉，您的问答体验期已经结束！"
		shareLiveQuestionAnswerLogger.DebugFormat("CheckPID 检查用户是否可以提问 pid=%d userId=%d liveRoomId=%d 【超出限制】 canIAsk=%s tip=%s", pid, userId, liveRoomId, strconv.FormatBool(canIAsk), tip)
		return canIAsk, tip, err
	}
	return canIAsk, tip, err
}

// IncreAskPermission 递增用户的提问次数 liveRoomId暂不使用 20181018变更
func (srv *LiveQuestionAnswerService) IncreAskPermission(pid int, userId int, liveRoomId int) (err error) {
	//todayAskCacheKey := fmt.Sprintf("%s:%d:%d", time.Now().Format("20060102"), liveRoomId, pid)
	//firstAskCacheKey := fmt.Sprintf("first:%d:%d", liveRoomId, pid)
	todayAskCacheKey := fmt.Sprintf("%s:%d", time.Now().Format("20060102"), pid)
	firstAskCacheKey := fmt.Sprintf("first:%d", pid)
	n, err := srv.RedisCache.HIncrBy(todayAskCacheKey, strconv.Itoa(userId), 1)
	if err != nil {
		shareLiveQuestionAnswerLogger.ErrorFormat(err, "IncreAskPermission 递增用户的提问次数 pid=%d userId=%d liveRoomId=%d 【异常】", pid, userId, liveRoomId)
		return err
	}
	shareLiveQuestionAnswerLogger.DebugFormat("IncreAskPermission 递增用户的提问次数 pid=%d userId=%d liveRoomId=%d 【递增成功】 n=%d", pid, userId, liveRoomId, n)

	firstAskDateStr, err := srv.RedisCache.HGet(firstAskCacheKey, strconv.Itoa(userId))
	if err != nil && err != redis.ErrNil {
		shareLiveQuestionAnswerLogger.ErrorFormat(err, "IncreAskPermission 递增用户的提问次数 pid=%d userId=%d liveRoomId=%d 【异常】", pid, userId, liveRoomId)
		return err
	}
	if err == redis.ErrNil || len(firstAskDateStr) == 0 {
		err = srv.RedisCache.HSet(firstAskCacheKey, strconv.Itoa(userId), time.Now().Format("2006-01-02 15:04:05"))
		shareLiveQuestionAnswerLogger.DebugFormat("IncreAskPermission 递增用户的提问次数 pid=%d userId=%d liveRoomId=%d 【首次提问】", pid, userId, liveRoomId)
		return err
	}
	firstAskDate, err := time.Parse("2006-01-02 15:04:05", firstAskDateStr)
	if err != nil && firstAskDate.Before(time.Now()) {
		err = srv.RedisCache.HSet(firstAskCacheKey, strconv.Itoa(userId), time.Now().Format("2006-01-02 15:04:05"))
		shareLiveQuestionAnswerLogger.DebugFormat("IncreAskPermission 递增用户的提问次数 pid=%d userId=%d liveRoomId=%d 【首次时间无效，时间被重置】 firstAskDateStr=%s", pid, userId, liveRoomId, firstAskDateStr)
		return err
	} else {
		return nil
	}
}

func init() {
	protected.RegisterServiceLoader(liveQuestionAnswerServiceName, liveQuestionAnswerServiceLoader)
}

func liveQuestionAnswerServiceLoader() {
	shareLiveQuestionAnswerRepo = liverepo.NewLiveQuestionAnswerRepository(protected.DefaultConfig)
	shareLiveQuestionAnswerLogger = dotlog.GetLogger(liveQuestionAnswerServiceName)
}
