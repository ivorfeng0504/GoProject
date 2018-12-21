package service

import (
	"encoding/json"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	"git.emoney.cn/softweb/roboadvisor/protected/repository/panwai"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"strconv"
	"strings"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
)

type NewsInfoService struct {
	service.BaseService
	newsRepo *repository.NewsInfoRepository
}

var (
	// shareNewsRepo NewsInfoRepository
	shareNewsRepo *repository.NewsInfoRepository

	// shareNewsLogger 共享的Logger实例
	shareNewsLogger dotlog.Logger
)

const (
	//根据栏目获取资讯列表 key
	RedisKey_NewsListByColumnID_zset = _const.RedisKey_NewsPre + "NewsInfo_sortset:GetNewsListByColumnID:"
	RedisKey_NewsListByColumnID_hset = _const.RedisKey_NewsPre + "NewsInfo_hashset:GetNewsListByColumnID:"

	//根据栏目和策略获取资讯列表 key
	RedisKey_NewsListByColumnIDAndStrategyID_zset = _const.RedisKey_NewsPre + "NewsInfo_sortset:GetNewsListByColumnIDAndStrategyID:"
	RedisKey_NewsListByColumnIDAndStrategyID_hset = _const.RedisKey_NewsPre + "NewsInfo_hashset:GetNewsListByColumnIDAndStrategyID:"

	//根据系列课程id获取资讯列表 key
	RedisKey_SeriesNewsListByNewsID_zset = _const.RedisKey_NewsPre + "NewsInfo_sortset:GetSeriesNewsListByNewsID:"
	RedisKey_SeriesNewsListByNewsID_hset = _const.RedisKey_NewsPre + "NewsInfo_hashset:GetSeriesNewsListByNewsID:"

	//获取单条资讯信息 key
	RedisKey_NewsInfoByID = _const.RedisKey_NewsPre + "NewsInfo:GetNewsInfoByID:"

	//根据栏目获取置顶资讯信息 key
	RedisKey_TopNewsInfo                        = _const.RedisKey_NewsPre + "NewsInfo_top:GetTopNewsInfoByColumnID:"
	RedisKey_TopNewsInfoByColumnIDAndStrategyID = _const.RedisKey_NewsPre + "NewsInfo_top:GetTopNewsInfoByColumnIDAndStrategyID:"

	//根据栏目获取所有置顶资讯信息 key（用户中心-通知模块使用）
	RedisKey_AllTopNewsInfo = _const.RedisKey_NewsPre + "NewsInfo_alltop:GetTopNewsInfoByColumnID:"

	RedisKey_NewsListByColumnID              = _const.RedisKey_NewsPre + "NewsInfo_zset:GetNewsListByColumnID:"
	RedisKey_NewsListByColumnIDAndStrategyID = _const.RedisKey_NewsPre + "NewsInfo:GetNewsListByColumnIDAndStrategyID:"
	RedisKey_SeriesNewsListByNewsID          = _const.RedisKey_NewsPre + "NewsInfo_zset:GetSeriesNewsListByNewsID:"

	//专家资讯-热门文章主动缓存key
	RedisKey_HotNewsListByClickNum = _const.RedisKey_NewsPre + "ExpertNews.GetHotNewsListByClickNum"

	//策略学习-热门文章
	RedisKey_HotArticleList_CLXX = _const.RedisKey_NewsPre + "CLXX.GetHotArticleList"
	//策略学习-热门视频
	RedisKey_HotVideoList_CLXX = _const.RedisKey_NewsPre + "CLXX.GetHotVideoList"
	//策略学习-热门实战培训
	RedisKey_HotMultiMediaList_CLXX = _const.RedisKey_NewsPre + "CLXX.GetHotMultiMediaList"

	//用户培训-策略实战（最近一个月）
	RedisKey_1MonthMultiMediaList = _const.RedisKey_NewsPre + "GetMultiMediaNewsList_1Month"

	//策略学习-热门资讯显示条数
	Clxx_HotNewsList_Num = 10

	newsInfoServiceName = "NewsInfoServiceLogger"
)

func NewNewsInfoService() *NewsInfoService {
	newsService := &NewsInfoService{
		newsRepo: shareNewsRepo,
	}
	newsService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return newsService
}

// GetAllTopNewsInfoByColumnID 根据栏目获取所有置顶资讯
func (service *NewsInfoService) GetTopNewsInfoByColumnID_userHome(columnID int, pid string) ([]*model.NewsInfo, error) {
	var results []*model.NewsInfo
	var err error
	var redisKey string

	redisKey = RedisKey_AllTopNewsInfo + strconv.Itoa(columnID)

	//get from redis
	err = service.RedisCache.GetJsonObj(redisKey, &results)

	if err != nil {
		return results, nil
	}

	//过滤pid-查询所属置顶资讯
	var newslistByPid []*model.NewsInfo
	for i := range results {
		newsinfo := results[i]
		if strings.Contains(newsinfo.Pid, pid) {
			newslistByPid = append(newslistByPid, newsinfo)
		}
	}
	results = newslistByPid
	if len(results) > 0 {
		results = results[:1]
	}

	return results, err
}

// GetTopNewsInfoByColumnID 根据栏目获取最新一条置顶资讯
func (service *NewsInfoService) GetTopNewsInfoByColumnID(columnID int, StrategyID int) (*model.NewsInfo, error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.GetTopNewsInfoByColumnIDDB(columnID, StrategyID)
	default:
		return service.GetTopNewsInfoByColumnIDCache(columnID, StrategyID)
	}
}
func (service *NewsInfoService) GetTopNewsInfoByColumnIDDB(columnID int, StrategyID int) (*model.NewsInfo, error) {
	var results *model.NewsInfo
	var err error
	return results, err
}
func (service *NewsInfoService) GetTopNewsInfoByColumnIDCache(columnID int, StrategyID int) (*model.NewsInfo, error) {
	result := new(model.NewsInfo)
	var err error
	var redisKey string

	if StrategyID == 0 {
		// 策略ID=0 不区分策略 直接返回置顶资讯
		redisKey = RedisKey_TopNewsInfo + strconv.Itoa(columnID)
	} else {
		// 策略ID>0 区分策略 返回置顶资讯
		redisKey = RedisKey_TopNewsInfoByColumnIDAndStrategyID + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID)
	}

	//get from redis
	err = service.RedisCache.GetJsonObj(redisKey, result)

	if err == nil {
		return result, nil
	}
	return result, err
}

// GetNewsListByColumnID 根据columnID获取资讯课程列表
func (service *NewsInfoService) GetNewsListByColumnID(columnID int, StrategyID int) ([]*model.NewsInfo, error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.GetNewsListByColumnIDDB(columnID, StrategyID)
	default:
		return service.GetNewsListByColumnIDCache(columnID, StrategyID)
	}
}

// GetNewsListByColumnIDPage 分页获取资讯列表(根据栏目、策略)
/*startnum endnum
0 9
10 19
20 29
*/
func (service *NewsInfoService) GetNewsListByColumnIDPage(columnID int, StrategyID int, pageIndex int, PageSize int) ([]*model.NewsInfo, int, error) {
	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--
	return service.GetNewsListByColumnIDCachePage(columnID, StrategyID, int64(startnum), int64(endnum))
}

// GetNewsListByColumnIDPage_userhome 用户中心通知消息分页
// 1、根据栏目获取资讯500条
// 2、根据栏目和pid获取最新一条置顶资讯
// 3、根据pid过滤所属资讯 并排除置顶资讯
// 4、对过滤后资讯分页
func (service *NewsInfoService) GetNewsListByColumnIDPage_userhome(columnID int, pid string, currpage int, pageSize int) ([]*model.NewsInfo, int, error) {
	newslist, totalCount, err := service.GetNewsListByColumnIDPage(columnID, 0, 1, 500)
	topnewslist, err := service.GetTopNewsInfoByColumnID_userHome(columnID, pid)

	//过滤pid-查询所属通知列表
	var newslistByPid []*model.NewsInfo
	for i := range newslist {
		newsinfo := newslist[i]
		if strings.Contains(newsinfo.Pid, pid) {

			//排除置顶一条资讯
			if len(topnewslist) > 0 {
				if newsinfo.ID != topnewslist[0].ID {
					newslistByPid = append(newslistByPid, newsinfo)
				}
			} else {
				newslistByPid = append(newslistByPid, newsinfo)
			}
		}
	}
	totalCount = len(newslistByPid)

	if pageSize > totalCount {
		pageSize = totalCount
	}
	beginnum := (currpage - 1) * pageSize
	endnum := pageSize

	if currpage == 1 {
		newslist = newslistByPid[:endnum]
	} else {
		endnum = pageSize * currpage
		//数量不够，endnum=开始num+剩余数量
		if (totalCount - beginnum) < pageSize {
			endnum = beginnum + (totalCount - beginnum)
		}
		newslist = newslistByPid[beginnum:endnum]
	}

	return newslist, totalCount, err
}

// GetNewsListByColumnIDDB 根据columnID获取资讯课程列表
func (service *NewsInfoService) GetNewsListByColumnIDDB(columnID int, StrategyID int) ([]*model.NewsInfo, error) {
	if columnID <= 0 {
		return nil, errors.New("must set columnID")
	}

	var results []*model.NewsInfo
	var err error
	if StrategyID == 0 {
		// 策略ID=0 不区分策略 直接返回栏目下所有资讯
		results, err = service.newsRepo.GetNewsListByColumnID(columnID)
	} else {
		// 策略ID>0 区分策略 获取该策略下所有资讯
		results, err = service.newsRepo.GetNewsListByStrategyIDAndColID(columnID, StrategyID)
	}

	if err == nil {
		if len(results) <= 0 {
			results = nil
			err = errors.New("not exists this new info")
		}
	}
	return results, err
}

// GetNewsListByColumnIDCache 根据columnID获取资讯课程列表
func (service *NewsInfoService) GetNewsListByColumnIDCache(columnID int, StrategyID int) ([]*model.NewsInfo, error) {
	if columnID <= 0 {
		return nil, errors.New("must set columnID")
	}

	var results []*model.NewsInfo
	var err error
	var redisKey string

	if StrategyID == 0 {
		// 策略ID=0 不区分策略 直接返回栏目下所有资讯
		redisKey = RedisKey_NewsListByColumnID + strconv.Itoa(columnID)
	} else {
		// 策略ID>0 区分策略 获取该策略和栏目下所有资讯
		redisKey = RedisKey_NewsListByColumnIDAndStrategyID + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID)
	}
	//get from redis
	err = service.RedisCache.GetJsonObj(redisKey, &results)
	if err == nil {
		return results, nil
	}
	return results, err
}

func (service *NewsInfoService) GetNewsListByColumnIDCachePage(columnID int, StrategyID int, start int64, end int64) ([]*model.NewsInfo, int, error) {
	var results []*model.NewsInfo
	var err error
	var redisKey_zset string
	var redisKey_hset string
	var totalnum int

	if StrategyID == 0 {
		// 策略ID=0 不区分策略 直接返回栏目下所有资讯
		redisKey_zset = RedisKey_NewsListByColumnID_zset + strconv.Itoa(columnID)
		redisKey_hset = RedisKey_NewsListByColumnID_hset + strconv.Itoa(columnID)
	} else {
		// 策略ID>0 区分策略 获取该策略和栏目下所有资讯
		redisKey_zset = RedisKey_NewsListByColumnIDAndStrategyID_zset + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID)
		redisKey_hset = RedisKey_NewsListByColumnIDAndStrategyID_hset + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID)
	}

	// redis 数据分页获取
	var retStrs []string
	totalnum, err = service.RedisCache.ZCard(redisKey_zset)
	retStrs, err = service.RedisCache.ZRevRange(redisKey_zset, start, end)
	args := make([]interface{}, len(retStrs))
	for i, v := range retStrs {
		args[i] = v
	}
	stringResult, err := service.RedisCache.HMGet(redisKey_hset, args...)
	//fmt.Println(stringResult)
	stringByte := "[" + strings.Join(stringResult, ",") + "]"
	err = json.Unmarshal([]byte(stringByte), &results)

	if err == nil {
		return results, totalnum, nil
	}
	return results, totalnum, err
}

// GetNewsInfoByID 根据newsID获取NewsInfo
func (service *NewsInfoService) GetNewsInfoByID(newsID int) (*model.NewsInfo, error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.GetNewsInfoByIDDB(newsID)
	default:
		return service.GetNewsInfoByIDCache(newsID)
	}
}

// GetNewsInfoByIDDB 根据newsID获取NewsInfo
func (service *NewsInfoService) GetNewsInfoByIDDB(newsID int) (*model.NewsInfo, error) {
	if newsID <= 0 {
		return nil, errors.New("must set newsID")
	}

	result, err := service.newsRepo.GetNewsInfoByID(newsID)
	if err == nil {
		if result == nil {
			err = errors.New("not exists this news info")
		}
	}
	return result, err
}

// GetNewsInfoByIDCache 根据newsID获取NewsInfo
func (service *NewsInfoService) GetNewsInfoByIDCache(newsID int) (*model.NewsInfo, error) {
	if newsID <= 0 {
		return nil, errors.New("must set newsID")
	}

	result := new(model.NewsInfo)
	var err error
	redisKey := RedisKey_NewsInfoByID + strconv.Itoa(newsID)
	//get from redis
	err = service.RedisCache.GetJsonObj(redisKey, result)
	if err == nil {
		return result, nil
	}
	return result, err
}

// GetSeriesNewsListByNewsID 根据topicID获取所有系列课程
func (service *NewsInfoService) GetSeriesNewsListByNewsID(newsID int) ([]*model.NewsInfo, error) {
	if newsID <= 0 {
		return nil, errors.New("must set newsID")
	}

	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.GetSeriesNewsListByNewsIDDB(newsID)
	default:
		return service.GetSeriesNewsListByNewsIDCache(newsID)
	}
}

// GetSeriesNewsListByNewsIDPage 分页获取系列课程列表(根据资讯id)
func (service *NewsInfoService) GetSeriesNewsListByNewsIDPage(newsID int, pageIndex int, PageSize int) ([]*model.NewsInfo, int, error) {
	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--
	return service.GetSeriesNewsListByNewsIDCachePage(newsID, int64(startnum), int64(endnum))
}

// GetSeriesNewsListByNewsIDDB 根据topicID获取所有系列课程
func (service *NewsInfoService) GetSeriesNewsListByNewsIDDB(newsID int) ([]*model.NewsInfo, error) {
	if newsID <= 0 {
		return nil, errors.New("must set newsID")
	}

	result, err := service.newsRepo.GetSeriesNewsListByNewsID(newsID)
	if err == nil {
		if result == nil {
			err = errors.New("not exists this news info")
		}
	}
	return result, err
}

// GetSeriesNewsListByNewsIDCache 根据topicID获取所有系列课程
func (service *NewsInfoService) GetSeriesNewsListByNewsIDCache(newsID int) ([]*model.NewsInfo, error) {
	if newsID <= 0 {
		return nil, errors.New("must set newsID")
	}

	var result []*model.NewsInfo
	var err error
	redisKey := RedisKey_SeriesNewsListByNewsID + strconv.Itoa(newsID)
	//get from redis
	err = service.RedisCache.GetJsonObj(redisKey, &result)
	if err == nil {
		return result, nil
	}

	return result, err
}

// GetSeriesNewsListByNewsIDCache 根据topicID获取所有系列课程
func (service *NewsInfoService) GetSeriesNewsListByNewsIDCachePage(newsID int, start int64, end int64) ([]*model.NewsInfo, int, error) {
	var results []*model.NewsInfo
	var err error
	var totalnum int

	redisKey_zset := RedisKey_SeriesNewsListByNewsID_zset + strconv.Itoa(newsID)
	redisKey_hset := RedisKey_SeriesNewsListByNewsID_hset + strconv.Itoa(newsID)

	// redis 数据分页获取
	var retStrs []string
	totalnum, err = service.RedisCache.ZCard(redisKey_zset)
	retStrs, err = service.RedisCache.ZRevRange(redisKey_zset, start, end)
	args := make([]interface{}, len(retStrs))
	for i, v := range retStrs {
		args[i] = v
	}
	stringResult, err := service.RedisCache.HMGet(redisKey_hset, args...)
	//fmt.Println(stringResult)
	stringByte := "[" + strings.Join(stringResult, ",") + "]"
	err = json.Unmarshal([]byte(stringByte), &results)

	if err == nil {
		return results, totalnum, nil
	}
	return results, totalnum, err
}

// UpdateClickNum 更新点击数
func (srv *NewsInfoService) UpdateClickNum(newsId int64, clickNum int64) (err error) {
	if newsId <= 0 || clickNum <= 0 {
		err = errors.New("咨询Id或点击数不正确")
	} else {
		err = srv.newsRepo.UpdateClickNum(newsId, clickNum)
	}
	if err != nil {
		shareNewsLogger.ErrorFormat(err, "更新咨询点击数异常 newsId=%d clickNum=%d", newsId, clickNum)
	}
	return err
}

// UpdateClickNum 更新点击数
func (srv *NewsInfoService) UpdateVideoPlayNum(newsId int64, playNum int64) (err error) {
	if newsId <= 0 || playNum <= 0 {
		err = errors.New("咨询Id或点击数不正确")
	} else {
		err = srv.newsRepo.UpdateVideoPlayNum(newsId, playNum)
	}
	if err != nil {
		shareNewsLogger.ErrorFormat(err, "更新资讯视频播放数异常 newsId=%d playNum=%d", newsId, playNum)
	}
	return err
}

// db获取10条热门文章存入redis （根据点击量）
func (srv *NewsInfoService) GetNewsListByClicknum(columnID int) (newsList []*model.NewsInfo, err error) {
	redisKey := RedisKey_HotNewsListByClickNum + ":" + strconv.Itoa(columnID)
	if columnID <= 0 {
		err = errors.New("栏目id不正确")
	} else {
		newsList, err = srv.newsRepo.GetNewsListByClicknum(columnID)

		if err != nil {
			shareNewsLogger.ErrorFormat(err, "更新热门文章失败 newsId=%d", columnID)
		} else {
			//db读取后存入redis
			srv.RedisCache.Delete(redisKey)
			srv.RedisCache.SetJsonObj(redisKey, newsList)
		}
	}

	if err != nil {
		shareNewsLogger.ErrorFormat(err, "更新热门文章失败 newsId=%d", columnID)
	}
	return newsList, err
}

// redis获取10条热门文章
func (service *NewsInfoService) GetNewsListByClicknumFromRedis(columnID int) (newsList []*model.NewsInfo, err error) {
	var result []*model.NewsInfo
	redisKey := RedisKey_HotNewsListByClickNum + ":" + strconv.Itoa(columnID)

	if columnID <= 0 {
		err = errors.New("栏目id不正确")
	} else {
		//get from redis
		err = service.RedisCache.GetJsonObj(redisKey, &result)
		if err == nil {
			return result, nil
		}

		return result, err
	}
	return result, err
}

// db获取10条热门文章资讯存入redis-策略学习使用
func (srv *NewsInfoService) GetHotArticleList(columnID int) (newsList []*model.NewsInfo, err error) {
	redisKey := RedisKey_HotArticleList_CLXX
	if columnID <= 0 {
		err = errors.New("栏目id不正确")
	} else {
		newsList, err = srv.newsRepo.GetNewsListByClicknum_clxx(columnID, Clxx_HotNewsList_Num, 0)

		if err != nil {
			shareNewsLogger.ErrorFormat(err, "策略学习-更新热门文章失败 newsId=%d", columnID)
		} else {
			//db读取后存入redis
			newsListJson, err := _json.Marshal(newsList)
			if err == nil {
				srv.RedisCache.HDel(redisKey, strconv.Itoa(columnID))
				srv.RedisCache.HSet(redisKey, strconv.Itoa(columnID), newsListJson)
			}
		}
	}

	if err != nil {
		shareNewsLogger.ErrorFormat(err, "策略学习-更新热门文章失败 newsId=%d", columnID)
	}
	return newsList, err
}

// redis读取策略学习热门文章
func (srv *NewsInfoService) GetHotArticleListFromRedis(clientStrategyIdList []string) (newsList []*model.NewsInfo, err error) {
	var newslist []*model.NewsInfo
	redisKey := RedisKey_HotArticleList_CLXX
	var args []interface{}
	colDict, err := strategyservice.NewClientStrategyInfoService().GetClientStrategyInfoDict()
	if err != nil {
		return newsList, err
	}
	for _, v := range clientStrategyIdList {
		clientStrategyId, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		colID, success := colDict[clientStrategyId]
		if success == false {
			continue
		}
		args = append(args, colID)
	}

	//没有匹配的栏目Id则直接返回
	if len(args) == 0 {
		return newsList, nil
	}

	stringResult, err := srv.RedisCache.HMGet(redisKey, args...)
	if err != nil {
		return newslist, err
	}
	if stringResult == nil || len(stringResult) == 0 {
		return nil, nil
	}
	//循环所有结果
	distinctDict := make(map[int]bool)
	for _, jsonStr := range stringResult {
		var tempList []*model.NewsInfo
		jsonErr := _json.Unmarshal(jsonStr, &tempList)
		if jsonErr == nil && tempList != nil && len(tempList) > 0 {
			//循环每一条资讯 去重
			for _, newsInfo := range tempList {
				_, isExist := distinctDict[newsInfo.ID]
				if isExist == false {
					newslist = append(newslist, newsInfo)
					distinctDict[newsInfo.ID] = true
				}
			}

		}
	}
	return newslist, err
}

// db获取10条热门视频资讯存入redis-策略学习使用
func (srv *NewsInfoService) GetHotVideoList(columnID int) (newsList []*model.NewsInfo, err error) {
	redisKey := RedisKey_HotVideoList_CLXX
	if columnID <= 0 {
		err = errors.New("栏目id不正确")
	} else {
		newsList, err = srv.newsRepo.GetNewsListByClicknum_clxx(columnID, Clxx_HotNewsList_Num, 1)

		if err != nil {
			shareNewsLogger.ErrorFormat(err, "策略学习-更新热门视频资讯失败 newsId=%d", columnID)
		} else {
			newsListJson, err := _json.Marshal(newsList)
			if err == nil {
				srv.RedisCache.HDel(redisKey, strconv.Itoa(columnID))
				srv.RedisCache.HSet(redisKey, strconv.Itoa(columnID), newsListJson)
			}
		}
	}

	if err != nil {
		shareNewsLogger.ErrorFormat(err, "策略学习-更新热门视频资讯失败 newsId=%d", columnID)
	}
	return newsList, err
}

// redis读取策略学习热门视频
func (srv *NewsInfoService) GetHotVideoListFromRedis(clientStrategyIdList []string) (newsList []*model.NewsInfo, err error) {
	var newslist []*model.NewsInfo
	redisKey := RedisKey_HotVideoList_CLXX
	var args []interface{}
	colDict, err := strategyservice.NewClientStrategyInfoService().GetClientStrategyInfoDict()
	for _, v := range clientStrategyIdList {
		clientStrategyId, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		colID, success := colDict[clientStrategyId]
		if success == false {
			continue
		}
		args = append(args, colID)
	}

	stringResult, err := srv.RedisCache.HMGet(redisKey, args...)
	if err != nil {
		return newslist, err
	}
	if stringResult == nil || len(stringResult) == 0 {
		return nil, nil
	}
	//循环所有结果
	distinctDict := make(map[int]bool)
	for _, jsonStr := range stringResult {
		var tempList []*model.NewsInfo
		jsonErr := _json.Unmarshal(jsonStr, &tempList)
		if jsonErr == nil && tempList != nil && len(tempList) > 0 {
			//循环每一条资讯 去重
			for _, newsInfo := range tempList {
				_, isExist := distinctDict[newsInfo.ID]
				if isExist == false {
					newslist = append(newslist, newsInfo)
					distinctDict[newsInfo.ID] = true
				}
			}

		}
	}

	return newslist, err
}

// db获取10条热门多媒体资讯存入redis-策略学习使用
func (srv *NewsInfoService) GetHotMultiMediaList(columnID int) (newsList []*model.NewsInfo, err error) {
	redisKey := RedisKey_HotMultiMediaList_CLXX
	if columnID <= 0 {
		err = errors.New("栏目id不正确")
	} else {
		newsList, err = srv.newsRepo.GetNewsListByClicknum_clxx(columnID, Clxx_HotNewsList_Num, 3)

		if err != nil {
			shareNewsLogger.ErrorFormat(err, "策略学习-更新热门培训资讯失败 newsId=%d", columnID)
		} else {
			newsListJson, err := _json.Marshal(newsList)
			if err == nil {
				srv.RedisCache.HDel(redisKey, strconv.Itoa(columnID))
				srv.RedisCache.HSet(redisKey, strconv.Itoa(columnID), newsListJson)
			}
		}
	}

	if err != nil {
		shareNewsLogger.ErrorFormat(err, "策略学习-更新热门培训资讯失败 newsId=%d", columnID)
	}
	return newsList, err
}

// redis读取策略学习热门培训
func (srv *NewsInfoService) GetHotMultiMediaListFromRedis(clientStrategyIdList []string) (newsList []*model.NewsInfo, err error) {
	var newslist []*model.NewsInfo
	redisKey := RedisKey_HotMultiMediaList_CLXX
	var args []interface{}
	colDict, err := strategyservice.NewClientStrategyInfoService().GetClientStrategyInfoDict()
	for _, v := range clientStrategyIdList {
		clientStrategyId, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		colID, success := colDict[clientStrategyId]
		if success == false {
			continue
		}
		args = append(args, colID)
	}

	stringResult, err := srv.RedisCache.HMGet(redisKey, args...)
	if err != nil {
		return newslist, err
	}
	if stringResult == nil || len(stringResult) == 0 {
		return nil, nil
	}
	//循环所有结果
	distinctDict := make(map[int]bool)
	for _, jsonStr := range stringResult {
		var tempList []*model.NewsInfo
		jsonErr := _json.Unmarshal(jsonStr, &tempList)
		if jsonErr == nil && tempList != nil && len(tempList) > 0 {
			//循环每一条资讯 去重
			for _, newsInfo := range tempList {
				_, isExist := distinctDict[newsInfo.ID]
				if isExist == false {
					newslist = append(newslist, newsInfo)
					distinctDict[newsInfo.ID] = true
				}
			}
		}
	}

	return newslist, err
}

// db获取最近一个月多媒体课程存入redis
func (srv *NewsInfoService) RefreshMultiMediaNewsList_1MonthToRedis() (newsList []*expertnews.ExpertNews_MultiMedia_List,err error) {
	redisKey := RedisKey_1MonthMultiMediaList
	newsList, err = srv.newsRepo.GetMultiMediaNewsList_1Month()

	if err != nil {
		shareNewsLogger.ErrorFormat(err, " 用户培训使用-db获取最近一个月多媒体课程失败 GetMultiMediaNewsList_1Month")
	} else {
		srv.RedisCache.SetJsonObj(redisKey, newsList)
	}

	return newsList, err
}


func init() {
	protected.RegisterServiceLoader(newsInfoServiceName, newsServiceLoader)
}

func newsServiceLoader() {
	shareNewsRepo = repository.NewNewsInfoRepository(protected.DefaultConfig)
	shareNewsLogger = dotlog.GetLogger(newsInfoServiceName)
}
