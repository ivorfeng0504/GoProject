package expertnews

import (
	"encoding/json"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/repository/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/click"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"strconv"
	"strings"
	"time"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"github.com/garyburd/redigo/redis"
)

type ExpertNews_StrategyService struct {
	service.BaseService
	expertNewsRepo *expertnews2.StrategyInfoRepository
}

var (
	shareExpertNews_StrategyRepo *expertnews2.StrategyInfoRepository
	// shareExpertNewsLogger 共享的Logger实例
	shareExpertNewsLogger dotlog.Logger
)

const (
	CacheKey_NewsListByColIDAndSID_zset = _const.RedisKey_NewsPre + "ExpertNews.GetNewsListByColumnIDAndStrategyID_zset:" //资讯列表zset（根据栏目和策略）
	CacheKey_NewsListByColIDAndSID_hset = _const.RedisKey_NewsPre + "ExpertNews.GetNewsListByColumnIDAndStrategyID_hset:" //资讯列表hset（根据栏目和策略）


	CacheKey_NewsListByColIDAndSIDAndType_zset = _const.RedisKey_NewsPre + "ExpertNews.GetNewsListByColumnIDAndStrategyIDAndType_zset:" //资讯列表zset（根据栏目和策略）
	CacheKey_NewsListByColIDAndSIDAndType_hset = _const.RedisKey_NewsPre + "ExpertNews.GetNewsListByColumnIDAndStrategyIDAndType_hset:" //资讯列表hset（根据栏目和策略）

	//根据栏目获取资讯列表 key
	CacheKey_NewsListByColumnID_zset = _const.RedisKey_NewsPre + "NewsInfo_sortset:GetNewsListByColumnID:"
	CacheKey_NewsListByColumnID_hset = _const.RedisKey_NewsPre + "NewsInfo_hashset:GetNewsListByColumnID:"

	//获取单条资讯信息 key
	CacheKey_NewsInfoByID = _const.RedisKey_NewsPre + "NewsInfo:GetNewsInfoByID:" //获取单条资讯详情

	CacheKey_ExpertStrategyInfo  = _const.RedisKey_NewsPre + "ExpertNews.GetStrategyInfoByID:" //获取专家策略详情key
	CacheKey_ExpertNewsList_Top6 = _const.RedisKey_NewsPre + "ExpertNews.ExpertNewsList_Top6"  //首页热门观点-策略分组后 最新6条资讯
	CacheKey_ExpertNewsList_Top1_video = _const.RedisKey_NewsPre + "ExpertNews.ExpertNewsList_Top1_video"  //首页热门观点-策略分组后 最新1条视频资讯组合列表

	CacheKey_ReceiveUpdateMsgToStrategy = _const.RedisKey_NewsPre + "ReceiveUpdateMsgToStrategy" //接收直播和资讯信息存入redis

	CacheKey_ExpertStrategyList_hset = _const.RedisKey_NewsPre + "ExpertNews.GetExpertStrategyList_hset" //专家策略列表

	expertNewsInfoServiceName = "ExpertNews_StrategyService"
)

func NewExpertNews_StrategyService() *ExpertNews_StrategyService {
	expertNewsService := &ExpertNews_StrategyService{
		expertNewsRepo: shareExpertNews_StrategyRepo,
	}
	expertNewsService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return expertNewsService
}

/* GetStrategyNewsListPage 分页获取策略资讯列表
columnID:专家资讯栏目id
StrategyID：专家策略ID(0:只根据栏目查询)
pageIndex：当前页面
PageSize：一页显示条数
*/
func (service *ExpertNews_StrategyService) GetStrategyNewsListPage(columnID int, StrategyID int, pageIndex int, PageSize int) ([]*expertnews.ExpertNews_StrategyNewsInfo, int, error) {
	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--

	return service.GetStrategyNewsListByCachePage(columnID, StrategyID, int64(startnum), int64(endnum))
}

// GetStrategyNewsListByCachePage redis分页获取策略资讯列表
func (service *ExpertNews_StrategyService) GetStrategyNewsListByCachePage(columnID int, StrategyID int, begin int64, end int64) ([]*expertnews.ExpertNews_StrategyNewsInfo, int, error) {
	var newslist []*model.NewsInfo
	var err error
	var totalnum int
	var redisKey_newslist_zset string
	var redisKey_newslist_hset string

	//获取资讯缓存数据,策略id=0 只根据栏目ID资讯列表
	if StrategyID == 0 {
		redisKey_newslist_zset = CacheKey_NewsListByColumnID_zset + strconv.Itoa(columnID)
		redisKey_newslist_hset = CacheKey_NewsListByColumnID_hset + strconv.Itoa(columnID)
	} else {
		redisKey_newslist_zset = CacheKey_NewsListByColIDAndSID_zset + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID)
		redisKey_newslist_hset = CacheKey_NewsListByColIDAndSID_hset + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID)
	}

	// redis 数据分页获取
	var retStrs []string
	totalnum, err = service.RedisCache.ZCard(redisKey_newslist_zset)
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetStrategyNewsListByCachePage redis分页获取策略资讯列表失败ZCard key=%s", redisKey_newslist_zset)
		return nil, 0, err
	}
	retStrs, err = service.RedisCache.ZRevRange(redisKey_newslist_zset, begin, end)
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetStrategyNewsListByCachePage redis分页获取策略资讯列表失败ZRevRange key=%s begin=%d end=%d", redisKey_newslist_zset, begin, end)
		return nil, 0, err
	}

	//取出资讯ID放入string数组，HMGet使用
	args := make([]interface{}, len(retStrs))
	for i, v := range retStrs {
		args[i] = v
	}
	stringResult, err := service.RedisCache.HMGet(redisKey_newslist_hset, args...)
	stringByte := "[" + strings.Join(stringResult, ",") + "]"
	err = json.Unmarshal([]byte(stringByte), &newslist)

	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetStrategyNewsListByCachePage 获取策略资讯列表反序列化失败 newslist=%s", stringByte)
		return nil, 0, err
	}

	//资讯关联策略信息 拼装后返回
	var strategyNewInfoList []*expertnews.ExpertNews_StrategyNewsInfo
	for i, _ := range newslist {
		newsInfo := newslist[i]

		if newsInfo != nil {
			strategyNewInfo := new(expertnews.ExpertNews_StrategyNewsInfo)                              //资讯
			clicknum, _ := click.GetClick(_const.ClickType_StrategyNewsInfo, strconv.Itoa(newsInfo.ID)) //获取资讯点击量
			newsInfo.ClickNum = clicknum

			strategyInfo := new(expertnews.ExpertNews_StrategyInfo) //策略
			StrategyID := newsInfo.ExpertStrategyID
			strategyInfo, _ = service.GetStrategyInfoByStrategyID(StrategyID)

			clicknum, _ = click.GetClick(_const.ClickType_StrategyInfo, strconv.Itoa(StrategyID)) //获取策略点击量
			strategyInfo.ClickNum = clicknum

			strategyNewInfo.NewsInfo = newsInfo
			strategyNewInfo.StrategyInfo = strategyInfo
			strategyNewInfoList = append(strategyNewInfoList, strategyNewInfo)
		}
	}

	return strategyNewInfoList, totalnum, err
}

// GetStrategyInfoByID 获取策略详情
func (service *ExpertNews_StrategyService) GetStrategyInfoByStrategyID(StrategyID int) (*expertnews.ExpertNews_StrategyInfo, error) {
	strategyInfo := new(expertnews.ExpertNews_StrategyInfo)
	redisKey_StrategyInfo := CacheKey_ExpertStrategyInfo + strconv.Itoa(StrategyID)

	err := service.RedisCache.GetJsonObj(redisKey_StrategyInfo, strategyInfo)

	if err == redis.ErrNil {
		err = nil
	}
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetStrategyInfoByStrategyID redis获取策略详情失败 key=%s", redisKey_StrategyInfo)
		return nil, err
	}

	if strategyInfo != nil {
		clicknum, _ := click.GetClick(_const.ClickType_StrategyInfo, strconv.Itoa(strategyInfo.ID)) //获取策略阅读总量
		strategyInfo.ClickNum = clicknum

		playnum, _ := click.GetClick(_const.ClickType_StrategyVideoInfo, strconv.Itoa(strategyInfo.ID)) //获取策略点击量
		strategyInfo.VideoPlayNum = playnum

		fansnum, _ := NewExpertNews_FocusStrategyService().GetFocusStrategyCount(StrategyID) //获取策略被关注的总数
		strategyInfo.FansNum = int64(fansnum)
	}

	return strategyInfo, err
}

// GetStrategyNewsInfoByNewsID 获取策略咨询详情
func (service *ExpertNews_StrategyService) GetStrategyNewsInfoByNewsID(NewsID int) (*expertnews.ExpertNews_StrategyNewsInfo, error) {
	strategyNewsInfo := new(expertnews.ExpertNews_StrategyNewsInfo)
	newsInfo := new(model.NewsInfo)
	redisKey_NewsInfo := CacheKey_NewsInfoByID + strconv.Itoa(NewsID)

	err := service.RedisCache.GetJsonObj(redisKey_NewsInfo, newsInfo)
	if err == redis.ErrNil {
		err = nil
		return nil, err
	}
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetStrategyNewsInfoByNewsID redis策略咨询详情失败 key=%s", redisKey_NewsInfo)
		return nil, err
	}

	if newsInfo != nil {
		clicknum, _ := click.GetClick(_const.ClickType_StrategyNewsInfo, strconv.Itoa(newsInfo.ID)) //获取策略点击量
		newsInfo.ClickNum = clicknum
	}

	strategyNewsInfo.NewsInfo = newsInfo
	StrategyInfo, _ := service.GetStrategyInfoByStrategyID(newsInfo.ExpertStrategyID)

	strategyNewsInfo.StrategyInfo = StrategyInfo

	return strategyNewsInfo, nil
}

// GetStrategyNewsList_Top6  首页-热门观点（分组策略后取最新1条文章资讯组合）
func (service *ExpertNews_StrategyService) GetStrategyNewsList_Top6() ([]*expertnews.ExpertNews_StrategyNewsInfo, error) {
	var newslist []*model.NewsInfo
	var strategyNewInfoList []*expertnews.ExpertNews_StrategyNewsInfo
	err := service.RedisCache.GetJsonObj(CacheKey_ExpertNewsList_Top6, &newslist)
	if err == redis.ErrNil {
		err = nil
	}
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetStrategyNewsList_Top6 redis获取首页-热门观点失败 key=%s", CacheKey_ExpertNewsList_Top6)
		return nil, err
	}
	if len(newslist) > 4 {
		newslist = newslist[:4]
	}
	for i, _ := range newslist {
		strategyNewInfo := new(expertnews.ExpertNews_StrategyNewsInfo)
		newsInfo := newslist[i] //资讯

		if newsInfo != nil {
			strategyInfo := new(expertnews.ExpertNews_StrategyInfo) //策略
			StrategyID := newsInfo.ExpertStrategyID
			strategyInfo, _ = service.GetStrategyInfoByStrategyID(StrategyID)

			strategyNewInfo.NewsInfo = newsInfo
			if strategyInfo != nil {
				strategyNewInfo.StrategyInfo = strategyInfo
			}

			strategyNewInfoList = append(strategyNewInfoList, strategyNewInfo)
		}
	}

	return strategyNewInfoList, nil
}

// GetStrategyNewsList_RmgdByPage 首页-热门观点（分组策略后取最新一条资讯组合）
func (service *ExpertNews_StrategyService) GetStrategyNewsList_RmgdByPage(currpage int, pageSize int) ([]*expertnews.ExpertNews_StrategyNewsInfo, int, error) {
	var newslist []*model.NewsInfo
	var newslist_page []*model.NewsInfo
	var strategyNewInfoList []*expertnews.ExpertNews_StrategyNewsInfo
	err := service.RedisCache.GetJsonObj(CacheKey_ExpertNewsList_Top6, &newslist)
	if err == redis.ErrNil {
		err = nil
	}
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetStrategyNewsList_RmgdByPage redis获取文章资讯列表失败 key=%s", CacheKey_ExpertNewsList_Top6)
		return nil, 0, err
	}
	totalCount := len(newslist)

	if pageSize > totalCount {
		pageSize = totalCount
	}
	beginnum := (currpage - 1) * pageSize
	//防止开始数大于总数
	if beginnum > totalCount {
		beginnum = totalCount
	}
	endnum := pageSize

	if currpage == 1 {
		newslist_page = newslist[:endnum]
	} else {
		endnum = pageSize * currpage
		//数量不够，endnum=开始num+剩余数量
		if (totalCount - beginnum) < pageSize {
			endnum = beginnum + (totalCount - beginnum)
		}
		newslist_page = newslist[beginnum:endnum]
	}

	for i, _ := range newslist_page {
		strategyNewInfo := new(expertnews.ExpertNews_StrategyNewsInfo)
		newsInfo := newslist_page[i] //资讯

		if newsInfo!=nil {
			clicknum, _ := click.GetClick(_const.ClickType_StrategyNewsInfo, strconv.Itoa(newsInfo.ID)) //获取资讯点击量
			newsInfo.ClickNum = clicknum

			strategyInfo := new(expertnews.ExpertNews_StrategyInfo) //策略
			StrategyID := newsInfo.ExpertStrategyID
			strategyInfo, _ = service.GetStrategyInfoByStrategyID(StrategyID)
			clicknum, _ = click.GetClick(_const.ClickType_StrategyInfo, strconv.Itoa(StrategyID)) //获取策略点击量
			strategyInfo.ClickNum = clicknum

			strategyNewInfo.NewsInfo = newsInfo
			strategyNewInfo.StrategyInfo = strategyInfo

			strategyNewInfoList = append(strategyNewInfoList, strategyNewInfo)
		}
	}

	return strategyNewInfoList, totalCount, nil
}

// GetNextStrategyNewsListByCurrentNews 获取当前资讯的上一篇和下一篇
func (service *ExpertNews_StrategyService) GetNextStrategyNewsListByCurrNews(newsId int, columnID int, StrategyID int, pageIndex int, pageSize int) ([]*expertnews.ExpertNews_StrategyNewsInfo, error) {
	strategyNewsList, _, err := service.GetStrategyNewsListPage(columnID, StrategyID, pageIndex, 10000)

	var result []*expertnews.ExpertNews_StrategyNewsInfo
	preNews := new(expertnews.ExpertNews_StrategyNewsInfo)  //上一篇
	nextNews := new(expertnews.ExpertNews_StrategyNewsInfo) //下一篇
	currNews := new(expertnews.ExpertNews_StrategyNewsInfo) //当前资讯

	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "根据当前资讯获取上下一条资讯失败 newsId=%d", newsId)
		return nil, err
	}

	if len(strategyNewsList) > 0 {
		//获取当前资讯所在列表的位置
		currIndex := 0
		for i, _ := range strategyNewsList {
			if strategyNewsList[i].NewsInfo != nil && strategyNewsList[i].NewsInfo.ID == newsId {

				currIndex = i
				//记录资讯点击量
				click.AddClick(_const.ClickType_StrategyNewsInfo, strconv.Itoa(newsId))
				//记录资讯对应策略点击量
				click.AddClick(_const.ClickType_StrategyInfo, strconv.Itoa(strategyNewsList[i].NewsInfo.ExpertStrategyID))

				break
			}
		}

		if currIndex == 0 { //当前资讯是第一条，上一篇资讯则为空
			currNews = strategyNewsList[currIndex]

			if len(strategyNewsList) > 1 {
				nextNews = strategyNewsList[currIndex+1]
			}

		} else if currIndex+1 < len(strategyNewsList) { //当前资讯存在当前列表中且不是第一条和最后一条

			preNews = strategyNewsList[currIndex-1]
			currNews = strategyNewsList[currIndex]
			nextNews = strategyNewsList[currIndex+1]

		} else { //当前资讯已是分页列表最后一条

			preNews = strategyNewsList[currIndex-1]
			currNews = strategyNewsList[currIndex]
		}
	}

	result = append(result, preNews)
	result = append(result, currNews)
	result = append(result, nextNews)

	return result, err
}

// UpdateClickNum 更新点击数
func (srv *ExpertNews_StrategyService) UpdateClickNum(newsId int64, clickNum int64) (err error) {
	if newsId <= 0 || clickNum <= 0 {
		err = errors.New("策略Id或点击数不正确")
	} else {
		err = srv.expertNewsRepo.UpdateClickNum(newsId, clickNum)
	}
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "更新策略点击数异常 newsId=%d clickNum=%d", newsId, clickNum)
	}
	return err
}

// UpdateVideoPlayNum 更新点击数
func (srv *ExpertNews_StrategyService) UpdateVideoPlayNum(newsId int64, playNum int64) (err error) {
	if newsId <= 0 || playNum <= 0 {
		err = errors.New("策略Id或点击数不正确")
	} else {
		err = srv.expertNewsRepo.UpdateVideoPlayNum(newsId, playNum)
	}
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "更新策略总点播数异常 newsId=%d playnum=%d", newsId, playNum)
	}
	return err
}

// 接收直播、资讯推送消息-加工后存入redis
/*
type=1: 最新直播消息
type=2: 最新资讯消息
*/
func (service *ExpertNews_StrategyService) ReceiveUpdateMsgToStrategy(msgType int, newsInfo *model.NewsInfo, liveInfo *expertnews.ReciveLive) (err error) {
	rediskey := CacheKey_ReceiveUpdateMsgToStrategy + ":" + strconv.Itoa(msgType)
	var strategyList []*expertnews.ExpertNews_StrategyInfo
	var focusStrategy = new(expertnews.FocusStrategyInfo)
	var sid = 0
	if newsInfo != nil {
		shareExpertNewsLogger.InfoFormat("接收资讯推送消息记录 newsinfo:%s", _json.GetJsonString(newsInfo))
		sid = newsInfo.ExpertStrategyID
	}
	if liveInfo != nil {
		shareExpertNewsLogger.InfoFormat("接收直播推送消息记录 live:%s", _json.GetJsonString(liveInfo))

		strategyList, err = service.expertNewsRepo.GetStrategyListByLiveID(liveInfo.LId)

		if err != nil {
			shareExpertNewsLogger.ErrorFormat(err, "根据直播ID获取策略信息异常 liveID=%d", liveInfo.LId)
		}
		if len(strategyList) > 0 {
			strategyInfo := strategyList[0]
			sid = strategyInfo.ID
		}
	}
	expertNewsStrategy, err := service.GetStrategyInfoByStrategyID(sid)
	if err != redis.ErrNil && err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "根据策略ID获取策略信息异常 strategyID=%d", sid)
		return nil
	}

	if expertNewsStrategy != nil {
		focusStrategy.StrategyID = sid
		focusStrategy.StrategyName = expertNewsStrategy.StrategyName
		focusStrategy.StrategyImg = expertNewsStrategy.StrategyImg
		focusStrategy.LiveId = expertNewsStrategy.LiveID
		focusStrategy.NewsType = msgType
		if msgType == 1 {
			focusStrategy.Title = liveInfo.LiveName
			focusStrategy.NewsID = liveInfo.LId
			focusStrategy.Summary = liveInfo.LiveContent
			focusStrategy.SendTime = time.Time(liveInfo.SendTime)
		}
		if msgType == 2 {
			focusStrategy.Title = newsInfo.Title
			focusStrategy.NewsID = newsInfo.ID
			focusStrategy.Summary = newsInfo.Title
			focusStrategy.SendTime = time.Time(newsInfo.CreateTime)

			//文章类型资讯
			if newsInfo.NewsType == _const.NewsType_News {
				focusStrategy.NewsType = 2
			}
			//视频类型资讯
			if newsInfo.NewsType == _const.NewsType_Video {
				focusStrategy.NewsType = 3
			}
		}
	}
	retObj, err := _json.Marshal(focusStrategy)
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "关注策略主动缓存接口序列化失败")
		return nil
	}
	service.RedisCache.HSet(rediskey, strconv.Itoa(sid), retObj)
	return nil

	//最新直播消息入redis
	//if msgType == 1 && liveInfo != nil {
	//
	//	for i, _ := range strategyList {
	//		strategyInfo := strategyList[i]
	//		strategyID := strconv.Itoa(strategyInfo.ID)
	//
	//		/*
	//			1、redis获取该策略信息（直播、资讯）
	//			2、覆盖最新一条直播信息
	//			3、重新存入redis
	//		*/
	//		retObj, _ := service.RedisCache.HGet(redisKey, strategyID)
	//		if len(retObj) > 0 {
	//			err = _json.Unmarshal(retObj, focusStrategy)
	//			if err != nil {
	//				shareExpertNewsLogger.ErrorFormat(err, "关注策略主动缓存接口反序列化失败")
	//			}
	//		} else {
	//			//获取该策略最新一条资讯
	//			newsList, err := service.expertNewsRepo.GetLatestNewsByExpertStrategyID(strategyInfo.ID)
	//			if err != nil {
	//				shareExpertNewsLogger.ErrorFormat(err, "数据库查询专家策略最新一条资讯失败")
	//			}
	//			if len(newsList) > 0 {
	//				focusStrategy.LatestNews = newsList[0]
	//			}
	//		}
	//		focusStrategy.LatestLive = liveInfo
	//
	//		retObj, err = _json.Marshal(focusStrategy)
	//		if err != nil {
	//			shareExpertNewsLogger.ErrorFormat(err, "关注策略主动缓存接口序列化失败")
	//		}
	//		service.RedisCache.HSet(redisKey, strategyID, retObj)
	//	}
	//}
	//
	////最新资讯消息入redis
	//if msgType == 2 && newsInfo != nil {
	//	strategyID := strconv.Itoa(newsInfo.ExpertStrategyID)
	//	/*
	//		1、redis获取该策略信息（直播、资讯）
	//		2、覆盖最新一条资讯信息
	//		3、重新存入redis
	//	*/
	//	retObj, _ := service.RedisCache.HGet(redisKey, strategyID)
	//	if len(retObj) > 0 {
	//		err = _json.Unmarshal(retObj, focusStrategy)
	//		if err != nil {
	//			shareExpertNewsLogger.ErrorFormat(err, "关注策略主动缓存接口反序列化失败")
	//		}
	//	}
	//	focusStrategy.LatestNews = newsInfo
	//	retObj, err = _json.Marshal(focusStrategy)
	//	if err != nil {
	//		shareExpertNewsLogger.ErrorFormat(err, "关注策略主动缓存接口序列化失败")
	//	}
	//
	//	service.RedisCache.HSet(redisKey, strategyID, retObj)
	//}

}



/* GetNewsListByStrategyIDAndTypePage 分页获取策略资讯列表
columnID:专家资讯栏目id
StrategyID：专家策略ID(0:只根据栏目查询)
newsType:资讯类型 （0：文章  1：视频）
pageIndex：当前页面
PageSize：一页显示条数
*/
func (service *ExpertNews_StrategyService) GetNewsListByStrategyIDAndTypePage(columnID int, StrategyID int, newsType int, pageIndex int64, PageSize int64) ([]*expertnews.ExpertNews_StrategyNewsInfo, int, error) {
	var newslist []*model.NewsInfo
	var err error
	var totalnum int

	redisKey_newslist_zset := CacheKey_NewsListByColIDAndSIDAndType_zset + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID) + "_" + strconv.Itoa(newsType)
	redisKey_newslist_hset := CacheKey_NewsListByColIDAndSIDAndType_hset + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID) + "_" + strconv.Itoa(newsType)

	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--

	// redis 数据分页获取
	var retStrs []string
	totalnum, err = service.RedisCache.ZCard(redisKey_newslist_zset)
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetNewsListByStrategyIDAndTypePage redis分页获取策略资讯列表失败ZCard key=%s", redisKey_newslist_zset)
		return nil, 0, err
	}
	retStrs, err = service.RedisCache.ZRevRange(redisKey_newslist_zset, startnum, endnum)
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetNewsListByStrategyIDAndTypePage redis分页获取策略资讯列表失败ZRevRange key=%s begin=%d end=%d", redisKey_newslist_zset, startnum, endnum)
		return nil, 0, err
	}

	//取出资讯ID放入string数组，HMGet使用
	args := make([]interface{}, len(retStrs))
	for i, v := range retStrs {
		args[i] = v
	}
	stringResult, err := service.RedisCache.HMGet(redisKey_newslist_hset, args...)
	stringByte := "[" + strings.Join(stringResult, ",") + "]"
	err = json.Unmarshal([]byte(stringByte), &newslist)

	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetNewsListByStrategyIDAndTypePage 获取策略资讯列表反序列化失败 newslist=%s", stringByte)
		return nil, 0, err
	}

	//资讯关联策略信息 拼装后返回
	var strategyNewInfoList []*expertnews.ExpertNews_StrategyNewsInfo
	for i, _ := range newslist {
		newsInfo := newslist[i]

		if newsInfo != nil {
			strategyNewInfo := new(expertnews.ExpertNews_StrategyNewsInfo)                              //资讯
			clicknum, _ := click.GetClick(_const.ClickType_StrategyNewsInfo, strconv.Itoa(newsInfo.ID)) //获取资讯点击量
			newsInfo.ClickNum = clicknum
			playnum, _ := click.GetClick(_const.ClickType_StrategyVideoNewsInfo, strconv.Itoa(newsInfo.ID)) //获取资讯视频播放量
			newsInfo.VideoPlayNum = playnum

			strategyInfo := new(expertnews.ExpertNews_StrategyInfo) //策略
			StrategyID := newsInfo.ExpertStrategyID
			strategyInfo, _ = service.GetStrategyInfoByStrategyID(StrategyID)

			clicknum, _ = click.GetClick(_const.ClickType_StrategyInfo, strconv.Itoa(StrategyID)) //获取策略点击量
			if strategyInfo!=nil {
				strategyInfo.ClickNum = clicknum
			}

			strategyNewInfo.NewsInfo = newsInfo
			strategyNewInfo.StrategyInfo = strategyInfo
			strategyNewInfoList = append(strategyNewInfoList, strategyNewInfo)
		}
	}

	return strategyNewInfoList, totalnum, err
}

// 获取专家策略列表
func (service *ExpertNews_StrategyService) GetExpertStrategyList() ([]*expertnews.ExpertNews_StrategyInfo,error) {
	var strategyList []*expertnews.ExpertNews_StrategyInfo

	allfelids, err := service.RedisCache.HGetAll(CacheKey_ExpertStrategyList_hset)
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetExpertStrategyList 获取策略资讯列表失败 key=%s", CacheKey_ExpertStrategyList_hset)
		return nil, err
	}

	for key, _ := range allfelids {
		strjson, err := service.RedisCache.HGet(CacheKey_ExpertStrategyList_hset, key)

		if err != nil {
			shareExpertNewsLogger.ErrorFormat(err, "GetExpertStrategyList 获取策略资讯反序列化失败 key=%s field=%s", CacheKey_ExpertStrategyList_hset, key)
		} else {
			strategyInfo := new(expertnews.ExpertNews_StrategyInfo)
			err = _json.Unmarshal(strjson, strategyInfo)

			strategyList = append(strategyList, strategyInfo)
		}
	}

	return strategyList, nil
}

// GetStrategyNewsList_Top1_video 首页-热门观点（分组策略后取最新1条视频资讯组合）
func (service *ExpertNews_StrategyService) GetStrategyNewsList_Top1_video() ([]*expertnews.ExpertNews_StrategyNewsInfo, error) {
	var newslist []*model.NewsInfo
	var strategyNewInfoList []*expertnews.ExpertNews_StrategyNewsInfo
	err := service.RedisCache.GetJsonObj(CacheKey_ExpertNewsList_Top1_video, &newslist)
	if err == redis.ErrNil {
		err = nil
	}
	if err != nil {
		shareExpertNewsLogger.ErrorFormat(err, "GetStrategyNewsList_Top1_video redis获取首页-热门视频失败 key=%s", CacheKey_ExpertNewsList_Top6)
		return nil, err
	}
	for i, _ := range newslist {
		strategyNewInfo := new(expertnews.ExpertNews_StrategyNewsInfo)
		newsInfo := newslist[i] //资讯

		if newsInfo != nil {
			strategyInfo := new(expertnews.ExpertNews_StrategyInfo) //策略
			StrategyID := newsInfo.ExpertStrategyID
			strategyInfo, _ = service.GetStrategyInfoByStrategyID(StrategyID)

			strategyNewInfo.NewsInfo = newsInfo
			if strategyInfo != nil {
				strategyNewInfo.StrategyInfo = strategyInfo
			}

			strategyNewInfoList = append(strategyNewInfoList, strategyNewInfo)
		}
	}

	return strategyNewInfoList, nil
}

func init() {
	protected.RegisterServiceLoader(expertNewsInfoServiceName, expertNews_StrategyServiceLoader)
}

func expertNews_StrategyServiceLoader() {
	shareExpertNewsLogger = dotlog.GetLogger(expertNewsInfoServiceName)
	shareExpertNews_StrategyRepo = expertnews2.NewStrategyInfoRepository(protected.DefaultConfig)
}
