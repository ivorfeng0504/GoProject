package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/repository/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"strconv"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"strings"
	"encoding/json"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	"sort"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"errors"
)

type ExpertNews_FocusStrategyService struct {
	service.BaseService
	focusStrategyRepo *expertnews2.FocusStrategyRepository
}

var (
	shareExpertNews_FocusStrategyRepo *expertnews2.FocusStrategyRepository
	// shareExpertNewsLogger 共享的Logger实例
	shareExpertNewsFocusStrategyLogger dotlog.Logger
)

const (
	expertNewsFocusStrategyServiceName = "expertNewsFocusStrategyServiceLogger"

	//CacheKey_FocusStrategyListByUID = _const.RedisKey_NewsPre+"FocusStrategyByUID:"
	CacheKey_FocusStrategyListByUID = _const.RedisKey_NewsPre+"zset_FocusStrategyByUID:"
	CacheKey_FocusStrategyListByStrategyID = _const.RedisKey_NewsPre+"zset_FocusStrategyByStrategyID:"

)

func NewExpertNews_FocusStrategyService() *ExpertNews_FocusStrategyService {
	expertNewsFocusStrategyService := &ExpertNews_FocusStrategyService{
		focusStrategyRepo: shareExpertNews_FocusStrategyRepo,
	}
	expertNewsFocusStrategyService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return expertNewsFocusStrategyService
}

// GetFocusStrategyCount 获取策略被关注的总量
func (service *ExpertNews_FocusStrategyService) GetFocusStrategyCount(strategyId int) (count int,err error) {
	//根据策略存入关注的用户uid
	rediskey_sid := CacheKey_FocusStrategyListByStrategyID + strconv.Itoa(strategyId)
	count, err = service.RedisCache.ZCard(rediskey_sid)
	if err != nil {
		shareExpertNewsFocusStrategyLogger.Error(err, fmt.Sprintf("获取关注策略数量失败 sid=&d", strategyId))
		return 0, err
	}
	return count, err
}

// AddFocusStrategy 关注策略信息记录
func (service *ExpertNews_FocusStrategyService) AddFocusStrategy(uid int64,strategyId int) (err error) {
	err = service.focusStrategyRepo.AddFocusStrategy(uid, strategyId)

	if err != nil {
		shareExpertNewsFocusStrategyLogger.Error(err, fmt.Sprintf("%s关注策略%s失败", uid, strategyId))
		return err
	}

	//存入redis
	rediskey := CacheKey_FocusStrategyListByUID + strconv.FormatInt(uid, 10)

	score := int64(strategyId)
	service.RedisCache.ZRem(rediskey, strategyId)
	service.RedisCache.ZAdd(rediskey, score, strategyId)

	//根据策略存入关注的用户uid
	rediskey_sid := CacheKey_FocusStrategyListByStrategyID + strconv.Itoa(strategyId)
	service.RedisCache.ZRem(rediskey_sid, uid)
	service.RedisCache.ZAdd(rediskey_sid, uid, uid)

	return nil
}

// RemoveFocusStrategy 取消关注信息记录
func (service *ExpertNews_FocusStrategyService) RemoveFocusStrategy(uid int64,strategyId int) (err error) {
	err = service.focusStrategyRepo.RemoveFocusStrategy(uid, strategyId)

	if err != nil {
		shareExpertNewsFocusStrategyLogger.Error(err, fmt.Sprintf("%s取消关注策略%s失败", uid, strategyId))
		return err
	}

	//从已关注缓存列表中删除
	rediskey := CacheKey_FocusStrategyListByUID + strconv.FormatInt(uid, 10)
	service.RedisCache.ZRem(rediskey, strategyId)

	rediskey_sid := CacheKey_FocusStrategyListByStrategyID + strconv.Itoa(strategyId)
	service.RedisCache.ZRem(rediskey_sid, uid)
	return nil
}

// HasFocusStrategy 是否关注过该策略
func (service *ExpertNews_FocusStrategyService) HasFocusStrategy(uid int64,strategyId int) (bool,error) {
	rediskey := CacheKey_FocusStrategyListByUID + strconv.FormatInt(uid, 10)
	retmsg, err := service.RedisCache.ZRevRange(rediskey, 0, -1)
	if err != nil {
		shareExpertNewsFocusStrategyLogger.Error(err, fmt.Sprintf("查询%s是否关注过策略%s失败", uid, strategyId))
		return false, err
	}
	for i, _ := range retmsg {
		if retmsg[i] == strconv.Itoa(strategyId) {
			return true, nil
		}
	}
	//retmsg, err := service.RedisCache.HGet(rediskey, strconv.Itoa(strategyId))

	return false, nil
}

// GetFocusStrategyByUID 根据UID获取关注的策略信息（包含最新一条直播和最新一条资讯）
func (service *ExpertNews_FocusStrategyService) GetFocusStrategyByUID(uid int64) ([]*expertnews.FocusStrategyInfo,bool,error) {
	var focusStrategyList []*expertnews.FocusStrategyInfo
	rediskey_focusSIDS := CacheKey_FocusStrategyListByUID + strconv.FormatInt(uid, 10)
	rediskey_allStrategyList_live := CacheKey_ReceiveUpdateMsgToStrategy + ":1"
	rediskey_allStrategyList_news := CacheKey_ReceiveUpdateMsgToStrategy + ":2"
	hasFocus := true

	retStrs, err := service.RedisCache.ZRevRange(rediskey_focusSIDS, 0, 100)
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetFocusStrategyByUID redis获取关注策略ID失败 key=%s", rediskey_focusSIDS)
		return nil, false, err
	}

	var args_live []interface{}
	var args_news []interface{}
	for _, v := range retStrs {
		//关注策略的直播是否有数据
		hasval, _ := service.RedisCache.HGet(rediskey_allStrategyList_live, v)
		if len(hasval) > 0 {
			args_live = append(args_live, v)
		}
		//关注策略的资讯是否有数据
		hasval, _ = service.RedisCache.HGet(rediskey_allStrategyList_news, v)
		if len(hasval) > 0 {
			args_news = append(args_news, v)
		}
	}

	//如果用户没有关注策略，则显示热门观点策略
	var args_all_live []interface{}
	var args_all_news []interface{}
	if len(retStrs) <= 0 {
		hasFocus = false
		var newslist []*model.NewsInfo
		_ = service.RedisCache.GetJsonObj(CacheKey_ExpertNewsList_Top6, &newslist)
		for i, _ := range newslist {
			var sid = strconv.Itoa(newslist[i].ExpertStrategyID)
			hasval, _ := service.RedisCache.HGet(rediskey_allStrategyList_live, sid)
			if len(hasval) > 0 {
				args_all_live = append(args_all_live, sid)
			}
			hasval, _ = service.RedisCache.HGet(rediskey_allStrategyList_news, sid)
			if len(hasval) > 0 {
				args_all_news = append(args_all_news, sid)
			}
		}
		args_live = args_all_live
		args_news = args_all_news
	}

	//直播
	stringResult_live, err := service.RedisCache.HMGet(rediskey_allStrategyList_live, args_live...)
	//资讯
	stringResult_news, err := service.RedisCache.HMGet(rediskey_allStrategyList_news, args_news...)

	stringResult := append(stringResult_live, stringResult_news...)

	stringByte := "[" + strings.Join(stringResult, ",") + "]"

	err = json.Unmarshal([]byte(stringByte), &focusStrategyList)

	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetFocusStrategyByUID 获取策略资讯列表反序列化失败 newslist=%s", stringByte)
		return nil, hasFocus, err
	}

	sort.Sort(FocusStrategyDatas(focusStrategyList))

	return focusStrategyList, hasFocus, err
}

//主题-个股排序
type FocusStrategyDatas []*expertnews.FocusStrategyInfo

func (s FocusStrategyDatas) Len() int {
	return len(s)
}
func (s FocusStrategyDatas) Less(i, j int) bool {
	return s[i].SendTime.After(s[j].SendTime)
}
func (s FocusStrategyDatas) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}

// FocusLive 关注/取消关注益圈圈直播
func (service *ExpertNews_FocusStrategyService) FocusLive(account string,ssourl string,optype int,liveid int) (error) {
	url := config.CurrentConfig.FocusLiveUrl

	json := fmt.Sprintf("Lid=%d&Username=%s&sso=%s&optype=%d", liveid, account, ssourl, optype)

	shareExpertNewsFocusStrategyLogger.InfoFormat("关注/取消关注益圈圈直播 url:%s 入参:%s", url, json)
	body, contentType, intervalTime, errReturn := _http.HttpPost(url, json, "application/x-www-form-urlencoded")
	if errReturn != nil {
		shareExpertNewsFocusStrategyLogger.WarnFormat("关注/取消关注益圈圈直播异常:%s 参数:%s", errReturn.Error(), json)
		return errReturn
	}
	_ = contentType
	_ = intervalTime
	//解析网关响应结果
	apiGatewayResp := contract.ApiGatewayResponse{}
	err := _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		shareExpertNewsFocusStrategyLogger.WarnFormat("关注/取消关注益圈圈直播响应结果解析异常:%s 参数:%s", err.Error(), json)
		return err
	}
	if apiGatewayResp.RetCode != 0 {
		shareExpertNewsFocusStrategyLogger.WarnFormat("关注/取消关注益圈圈直播响应异常 %s", apiGatewayResp.RetMsg)
		return errors.New(apiGatewayResp.RetMsg)
	}

	//解析业务结果
	result := new(FocusLiveResult)
	err = _json.Unmarshal(apiGatewayResp.Message, result)
	if err != nil {
		shareExpertNewsFocusStrategyLogger.ErrorFormat(err, "关注/取消关注益圈圈直播异常 反序列化异常 地址=%s 入参:%s 结果：%s", url, json, body)
		return err
	}

	if result == nil || result.RetCode != "0" {
		shareExpertNewsFocusStrategyLogger.ErrorFormat(err, "关注/取消关注益圈圈直播异常 地址=%s 入参:%s 结果：%s", url, json, body)
		return errors.New("关注/取消关注失败")
	}
	shareExpertNewsFocusStrategyLogger.DebugFormat("关注/取消关注成功 地址=%s  入参:%s  结果：%s", url, json, body)
	return nil
}

type FocusLiveModel struct {
	Lid int
	Optype int // 操作类型，0关注 1取消关注
	Username string
	Sso string //ssoURL
}

type FocusLiveResult struct {
	RetCode string
	RetMsg   string
	Data interface{}
}

func init() {
	protected.RegisterServiceLoader(expertNewsFocusStrategyServiceName, expertNews_FocusStrategyLoader)
}

func expertNews_FocusStrategyLoader() {
	shareExpertNewsFocusStrategyLogger = dotlog.GetLogger(expertNewsFocusStrategyServiceName)
	shareExpertNews_FocusStrategyRepo = expertnews2.NewFocusStrategyRepository(protected.DefaultConfig)
}
