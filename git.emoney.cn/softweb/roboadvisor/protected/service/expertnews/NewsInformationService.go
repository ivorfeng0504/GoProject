package expertnews

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	expertnews_model "git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	expertnews_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/mapper"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
)

type NewsInformationService struct {
	service.BaseService
	newsRepo *expertnews_repo.NewsInformationRepository
}

var (
	shareNewsInformationRepo   *expertnews_repo.NewsInformationRepository
	shareNewsInformationLogger dotlog.Logger
)

const (
	ExpertNews_NewsInformationServiceName              = "NewsInformationService"
	EMNET_NewsInformation_CachePreKey                  = "EMNET:NewsInformationService"
	EMNET_NewsInformation_InsertQueue_CacheKey         = EMNET_NewsInformation_CachePreKey + ":InsertQueue"
	EMNET_NewsInformation_InsertErrorQueue_CacheKey    = EMNET_NewsInformation_CachePreKey + ":InsertErrorQueue"
	EMNET_NewsInformation_UpdateQueue_CacheKey         = EMNET_NewsInformation_CachePreKey + ":UpdateQueue"
	EMNET_NewsInformation_NewsTemplateCache_CacheKey   = EMNET_NewsInformation_CachePreKey + ":NewsTemplateCache"
	EMNET_NewsInformation_GetTodayNewsCache_CacheKey   = EMNET_NewsInformation_CachePreKey + ":GetTodayNewsCache:"
	EMNET_NewsInformation_GetTodayNewsCacheV2_CacheKey = EMNET_NewsInformation_CachePreKey + ":GetTodayNewsCacheV2:"
	EMNET_NewsInformation_GetClosingNewsCache_CacheKey = EMNET_NewsInformation_CachePreKey + ":GetClosingNewsCache:"
	EMNET_NewsInformation_GetNewsInfoCache_CacheKey    = EMNET_NewsInformation_CachePreKey + ":GetNewsInfoCache:"
	//当前最新的要闻
	EMNET_NewsInformation_GetNewsInfoCache_Newst_CacheKey   = EMNET_NewsInformation_CachePreKey + ":GetNewsInfoCache:Newst"
	EMNET_NewsInformation_GetHotNewsInfoCache_CacheKey      = EMNET_NewsInformation_CachePreKey + ":GetHotNewsInfoCache:"
	EMNET_NewsInformation_GetTopicNewsInfoCache_CacheKey    = EMNET_NewsInformation_CachePreKey + ":GetTopicNewsInfoCache:"
	EMNET_NewsInformation_GetNewsInfoByIdCache_CacheKey     = EMNET_NewsInformation_CachePreKey + ":GetNewsInfoById:"
	EMNET_NewsInformation_GetNewsInfoTopNCache_CacheKey     = EMNET_NewsInformation_CachePreKey + ":GetNewsInfoTopNCache:"
	EMNET_NewsInformation_GetNewsInfoByIdCache_CacheSeconds = 60 * 60 * 24 * 3
	EMNET_NewsInformation_GetBlockNewsInfomation_CacheKey   = EMNET_NewsInformation_CachePreKey + ":GetBlockNewsInfomation:"
	EMNET_NewsInformation_Detail_CacheKey                   = EMNET_NewsInformation_CachePreKey + ":Detail:"
	//当前板块数据同步处理的最大版本号
	EMNET_NewsInformation_CurrentBlockNewsSyncVersion_CacheKey = EMNET_NewsInformation_CachePreKey + ":CurrentBlockNewsSyncVersion:"

	//最新的N条要闻
	NewsInfoTopN = 30
)

func NewNewsInformationService() *NewsInformationService {
	srv := &NewsInformationService{
		newsRepo: shareNewsInformationRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(ExpertNews_NewsInformationServiceName, newsInformationServiceLoader)
}

func newsInformationServiceLoader() {
	shareNewsInformationRepo = expertnews_repo.NewNewsInformationRepository(protected.DefaultConfig)
	shareNewsInformationLogger = dotlog.GetLogger(ExpertNews_NewsInformationServiceName)
}

// InsertOrUpdate 插入或更新资讯数据
func (srv *NewsInformationService) InsertOrUpdate(id string, params []string, values []string) (err error) {
	err = srv.newsRepo.InsertOrUpdate(id, params, values)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "InsertOrUpdate 插入或更新资讯数据 id=%s params=%s  values=%s", id, _json.GetJsonString(params), _json.GetJsonString(values))
	} else {
		//刷新缓存
	}
	return err
}

// AddInsertQueue 添加到数据库写入队列
func (srv *NewsInformationService) AddInsertQueue(params []string, values []string) (err error) {
	insertData := InsertData{
		Params: params,
		Values: values,
	}
	_, err = srv.RedisCache.LPush(EMNET_NewsInformation_InsertQueue_CacheKey, insertData)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "AddInsertQueue 添加到数据库写入队列异常 insertData=%s", _json.GetJsonString(insertData))
		srv.AddInsertErrorQueue(params, values)
	}
	return err
}

// AddInsertErrorQueue 添加失败写入错误队列中
func (srv *NewsInformationService) AddInsertErrorQueue(params []string, values []string) (err error) {
	insertData := InsertData{
		Params: params,
		Values: values,
	}
	_, err = srv.RedisCache.LPush(EMNET_NewsInformation_InsertErrorQueue_CacheKey, insertData)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "AddInsertErrorQueue 添加失败写入错误队列中异常 insertData=%s", _json.GetJsonString(insertData))
	}
	return err
}

// GetMaxSyncRowVersion 获取指定日期里最新版本号 不传递date则取当前最大版本号
func (srv *NewsInformationService) GetMaxSyncRowVersion(date *time.Time) (version string, err error) {
	return srv.newsRepo.GetMaxSyncRowVersion(date)
}

// AddUpdateQueue 添加到更新队列中 触发资讯更新任务
func (srv *NewsInformationService) AddUpdateQueue(tableName string) (err error) {
	_, err = srv.RedisCache.LPush(EMNET_NewsInformation_UpdateQueue_CacheKey, tableName)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "AddUpdateQueue 添加到更新队列中异常 tableName=%s", tableName)
	}
	return err
}

// PopUpdateQueue 获取资讯更新队列
func (srv *NewsInformationService) PopUpdateQueue() (tableName string, err error) {
	tableName, err = srv.RedisCache.LPop(EMNET_NewsInformation_UpdateQueue_CacheKey)
	if err == redis.ErrNil {
		return "", nil
	}
	if err == nil {
		shareNewsInformationLogger.InfoFormat("PopUpdateQueue 获取队列数据  tableName=%s ", tableName)
	} else {
		shareNewsInformationLogger.ErrorFormat(err, "PopUpdateQueue 获取队列数据异常  tableName=%s ", tableName)
	}
	return tableName, err
}

// ClearUpdateQueue 清除资讯更新队列
func (srv *NewsInformationService) ClearUpdateQueue() (err error) {
	err = srv.RedisCache.Delete(EMNET_NewsInformation_UpdateQueue_CacheKey)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "ClearUpdateQueue 清除资讯更新队列异常  ")
	}
	return err
}

// UpdateClickNum 更新点击数
func (srv *NewsInformationService) UpdateClickNum(newsId int64, clickNum int64) (err error) {
	if newsId <= 0 || clickNum <= 0 {
		err = errors.New("资讯Id或点击数不正确")
	} else {
		err = srv.newsRepo.UpdateClickNum(newsId, clickNum)
	}
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "更新咨询点击数异常 newsId=%d clickNum=%d", newsId, clickNum)
	}
	return err
}

// UpdateStockNewsClickNum 更新点击数
func (srv *NewsInformationService) UpdateStockNewsClickNum(newsId int64, clickNum int64) (err error) {
	if newsId <= 0 || clickNum <= 0 {
		err = errors.New("资讯Id或点击数不正确")
	} else {
		err = srv.newsRepo.UpdateStockNewsClickNum(newsId, clickNum)
	}
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "更新咨询点击数异常 newsId=%d clickNum=%d", newsId, clickNum)
	}
	return err
}

// GetNewsTemplateCache获取资讯模板缓存
func (srv *NewsInformationService) GetNewsTemplateCache() (templates []string, err error) {
	cacheKey := EMNET_NewsInformation_NewsTemplateCache_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &templates)
	if err == redis.ErrNil {
		return nil, nil
	}
	return templates, err
}

// SetNewsTemplateCache 设置资讯模板缓存
func (srv *NewsInformationService) SetNewsTemplateCache(templates []string) (err error) {
	cacheKey := EMNET_NewsInformation_NewsTemplateCache_CacheKey
	_, err = srv.RedisCache.SetJsonObj(cacheKey, templates)
	return err
}

// GetTodayNewsCache 获取今日头条-缓存
func (srv *NewsInformationService) GetTodayNewsCache() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetTodayNewsCache_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// GetTodayNewsCacheV2 获取今日头条-缓存-第2版
func (srv *NewsInformationService) GetTodayNewsCacheV2() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetTodayNewsCacheV2_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// GetClosingNewsCache 获取盘后预测资讯-缓存
func (srv *NewsInformationService) GetClosingNewsCache() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetClosingNewsCache_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// GetNewsInfoCache 获取指定日期的要闻-缓存
func (srv *NewsInformationService) GetNewsInfoCache(date time.Time) (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetNewsInfoCache_CacheKey + date.Format("2006-01-02")
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// GetNewsInfoTopNCache 获取最新的N条要闻-缓存
func (srv *NewsInformationService) GetNewsInfoTopNCache() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetNewsInfoTopNCache_CacheKey + strconv.Itoa(NewsInfoTopN)
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// GetYaoWenHomePageData 获取热读页面数据
func (srv *NewsInformationService) GetYaoWenHomePageData() (pageData *agent.YaoWenHomePageData, err error) {
	pageData = &agent.YaoWenHomePageData{}
	var result []*agent.ExpertNewsInfo
	//最新热读
	newsList, err := srv.GetNewsInfoTopNCache()
	if err == nil && newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(newsList, &result)
		if err == nil {
			//去除最后一天的数据
			length := len(result)
			if result != nil && length > 0 {
				lastDateIndex := -1
				lastDate := ""
				for i := length - 1; i >= 0; i-- {
					publishTime := time.Time(result[i].PublishTime).Format("20060102")
					if lastDate != "" && lastDate != publishTime {
						lastDateIndex = i + 1
						break
					}
					lastDate = publishTime
				}
				if lastDateIndex > 0 {
					result = result[:lastDateIndex]
				}
			}
			pageData.NewsList = result
		}
	}

	//热门热读
	var hotResult []*agent.ExpertNewsInfo
	hotNewsList, err := srv.GetHotNewsInfoCache()
	if err == nil && newsList != nil && len(newsList) > 0 {
		err = mapper.MapperSlice(hotNewsList, &hotResult)
		if err == nil {
			pageData.HotNewsList = hotResult
		}
	}
	return pageData, err
}

// GetNewstNewsInfoCache 返回当前缓存中最新的要闻信息
func (srv *NewsInformationService) GetNewstNewsInfoCache() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetNewsInfoCache_Newst_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// GetHotNewsInfoCache 获取热门资讯-缓存
func (srv *NewsInformationService) GetHotNewsInfoCache() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetHotNewsInfoCache_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// GetTopicNewsInfoCache 获取主题的相关资讯-缓存
func (srv *NewsInformationService) GetTopicNewsInfoCache(topicId int) (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetTopicNewsInfoCache_CacheKey + strconv.Itoa(topicId)
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// RefreshTodayNewsCache 获取今日头条-刷新缓存
func (srv *NewsInformationService) RefreshTodayNewsCache() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetTodayNewsCache_CacheKey
	newsList, err = srv.newsRepo.GetTodayNews()
	if err != nil {
		return nil, err
	}
	if newsList != nil && len(newsList) > 0 {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, newsList)
		//刷新详情
		srv.RefreshNewsDetail(newsList)
	}
	return newsList, err
}

// RefreshTodayNewsCacheV2 获取今日头条-刷新缓存-第2版
func (srv *NewsInformationService) RefreshTodayNewsCacheV2() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetTodayNewsCacheV2_CacheKey
	newsList, err = srv.newsRepo.GetTodayNewsV2()
	if err != nil {
		return nil, err
	}
	if newsList != nil && len(newsList) > 0 {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, newsList)
		//刷新详情
		srv.RefreshNewsDetail(newsList)
	}
	return newsList, err
}

// RefreshClosingNewsCache 获取盘后预测资讯-刷新缓存
func (srv *NewsInformationService) RefreshClosingNewsCache() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetClosingNewsCache_CacheKey
	//templates, err := srv.GetNewsTemplateCache()
	//if err != nil {
	//	return nil, err
	//}
	templates := []string{"涨停预测", "明日猜想"}
	newsList, err = srv.newsRepo.GetClosingNews(templates)
	if err != nil {
		return nil, err
	}
	if newsList != nil && len(newsList) > 0 {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, newsList)
		//刷新详情
		srv.RefreshNewsDetail(newsList)
	}
	return newsList, err
}

// RefreshNewsInfoCache 获取指定日期的要闻-刷新缓存
func (srv *NewsInformationService) RefreshNewsInfoCache(date time.Time) (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetNewsInfoCache_CacheKey + date.Format("2006-01-02")
	templates, err := srv.GetNewsTemplateCache()
	if err != nil {
		return nil, err
	}
	newsList, err = srv.newsRepo.GetNewsInfo(templates, date)
	if err != nil {
		return nil, err
	}
	if newsList != nil && len(newsList) > 0 {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, newsList)
		//刷新详情
		srv.RefreshNewsDetail(newsList)
		//如果是当天的要闻则更新最新的要闻信息
		if time.Now().Format("2006-01-02") == date.Format("2006-01-02") {
			_, err = srv.RedisCache.SetJsonObj(EMNET_NewsInformation_GetNewsInfoCache_Newst_CacheKey, newsList)
		}
	}
	return newsList, err
}

// RefreshNewsInfoCache 获取指定日期的要闻-刷新缓存
func (srv *NewsInformationService) RefreshNewsInfoTopNCache(top int) (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetNewsInfoTopNCache_CacheKey + strconv.Itoa(top)
	templates, err := srv.GetNewsTemplateCache()
	if err != nil {
		return nil, err
	}
	newsList, err = srv.newsRepo.GetNewsInfoTopN(templates, top)
	if err != nil {
		return nil, err
	}
	if newsList != nil && len(newsList) > 0 {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, newsList)
		//刷新详情
		srv.RefreshNewsDetail(newsList)
	}
	return newsList, err
}

// RefreshHotNewsInfoCache 获取热门资讯-刷新缓存
func (srv *NewsInformationService) RefreshHotNewsInfoCache() (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetHotNewsInfoCache_CacheKey
	templates, err := srv.GetNewsTemplateCache()
	if err != nil {
		return nil, err
	}
	newsList, err = srv.newsRepo.GetHotNewsInfo(templates)
	if err != nil {
		return nil, err
	}
	if newsList != nil && len(newsList) > 0 {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, newsList)
		//刷新详情
		srv.RefreshNewsDetail(newsList)
	}
	return newsList, err
}

// RefreshTopicNewsInfoCache 获取主题的相关资讯-刷新缓存
func (srv *NewsInformationService) RefreshTopicNewsInfoCache(topic *expertnews_model.ExpertNews_Topic) (newsList []*expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetTopicNewsInfoCache_CacheKey + strconv.Itoa(topic.ID)
	//如果主题或者主题的板块信息不存在 则清除缓存
	if topic == nil || len(topic.RelatedBKInfo) == 0 {
		srv.RedisCache.Delete(cacheKey)
		return nil, nil
	}
	err = _json.Unmarshal(topic.RelatedBKInfo, &topic.RelatedBKList)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "RefreshTopicNewsInfoCache 获取主题的相关资讯-刷新缓存失败 反序列化板块信息异常 bklist=%s", topic.RelatedBKInfo)
		return nil, err
	}
	newsList, err = srv.newsRepo.GetTopicNewsInfo(topic.RelatedBKList)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "RefreshTopicNewsInfoCache 获取主题的相关资讯-刷新缓存失败 获取主题信息异常 bklist=%s", topic.RelatedBKInfo)
		return nil, err
	}
	if newsList != nil && len(newsList) > 0 {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, newsList)
		//刷新详情
		srv.RefreshNewsDetail(newsList)
	}
	return newsList, err
}

// RefreshCache 刷新所有缓存
func (srv *NewsInformationService) RefreshCache() (err error) {
	_, err = srv.RefreshTodayNewsCacheV2()
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, " 刷新今日头条 执行异常")
		return err
	}
	_, err = srv.RefreshClosingNewsCache()
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, " 刷新盘后预测资讯 执行异常")
		return err
	}

	//刷新最近三天的要闻数据
	//dateList := []time.Time{time.Now(), time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, -2)}
	//for _, date := range dateList {
	//	_, err = srv.RefreshNewsInfoCache(date)
	//	if err != nil {
	//		shareNewsInformationLogger.ErrorFormat(err, " 刷新指定日期的要闻 执行异常 date=%s", date.Format("2006-01-02"))
	//		return err
	//	}
	//}

	//刷新最新的N条要闻
	_, err = srv.RefreshNewsInfoTopNCache(NewsInfoTopN)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, " 刷新最新的N条要闻 执行异常 top=%d", NewsInfoTopN)
		return err
	}

	_, err = srv.RefreshHotNewsInfoCache()
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, " 刷新热门资讯 执行异常")
		return err
	}
	topicSrv := NewExpertNews_TopicService()
	topicList, err := topicSrv.GetTopicList()
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, " 刷新今日头条 执行异常")
		return err
	}
	if topicList == nil || len(topicList) == 0 {
		return nil
	}
	for _, topic := range topicList {
		newsList, err := srv.RefreshTopicNewsInfoCache(topic)
		if err != nil {
			shareNewsInformationLogger.ErrorFormat(err, " 刷新所有主题的相关资讯 执行异常 topicId=%d newsList=【%s】", topic.ID, _json.GetJsonString(newsList))
		} else {
			shareNewsInformationLogger.DebugFormat(" 刷新所有主题的相关资讯 执行完毕 topicId=%d newsList=【%s】", topic.ID, _json.GetJsonString(newsList))
		}
	}
	return err
}

// GetNewsInfoById 根据Id获取资讯信息
func (srv *NewsInformationService) GetNewsInfoById(newsInfoId int64) (newsInfo *expertnews_model.NewsInformation, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.GetNewsInfoByIdDB(newsInfoId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		newsInfo, err = srv.GetNewsInfoByIdCache(newsInfoId)
		if err == nil && newsInfo == nil {
			newsInfo, err = srv.RefreshNewsInfoById(newsInfoId)
		}
		return newsInfo, err
	case config.ReadDB_RefreshCache:
		newsInfo, err = srv.RefreshNewsInfoById(newsInfoId)
		return newsInfo, err
	default:
		return srv.GetNewsInfoByIdCache(newsInfoId)
	}
}

// GetNewsInfoByIdDB 根据Id获取资讯信息-读取数据库
func (srv *NewsInformationService) GetNewsInfoByIdDB(newsInfoId int64) (newsInfo *expertnews_model.NewsInformation, err error) {
	newsInfo, err = srv.newsRepo.GetNewsInfoById(newsInfoId)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "GetNewsInfoById 根据Id获取资讯信息异常 newsInfoId=%d", newsInfoId)
	}
	return newsInfo, err
}

// GetNewsInfoByIdCache 根据Id获取资讯信息-读取缓存
func (srv *NewsInformationService) GetNewsInfoByIdCache(newsInfoId int64) (newsInfo *expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetNewsInfoByIdCache_CacheKey + strconv.FormatInt(newsInfoId, 10)
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsInfo)
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "GetNewsInfoByIdCache 根据Id获取资讯信息异常 newsInfoId=%d", newsInfoId)
	}
	return newsInfo, err
}

// RefreshNewsInfoById 根据Id获取资讯信息-刷新缓存
func (srv *NewsInformationService) RefreshNewsInfoById(newsInfoId int64) (newsInfo *expertnews_model.NewsInformation, err error) {
	cacheKey := EMNET_NewsInformation_GetNewsInfoByIdCache_CacheKey + strconv.FormatInt(newsInfoId, 10)
	newsInfo, err = srv.GetNewsInfoByIdDB(newsInfoId)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "RefreshNewsInfoById->GetNewsInfoByIdDB 根据Id获取资讯信息-刷新缓存异常 newsInfoId=%d", newsInfoId)
		return nil, err
	}
	if newsInfo != nil {
		err = srv.RedisCache.Set(cacheKey, _json.GetJsonString(newsInfo), EMNET_NewsInformation_GetNewsInfoByIdCache_CacheSeconds)
	}
	return newsInfo, err
}

// RefreshNewsDetail 刷新指定的资讯详情 newsList为从数据库查询的最新内容
func (srv *NewsInformationService) RefreshNewsDetail(newsList []*expertnews_model.NewsInformation) {
	if newsList == nil {
		return
	}
	go func() {
		for _, news := range newsList {
			if news != nil && news.NewsInformationId > 0 {
				cacheKey := EMNET_NewsInformation_GetNewsInfoByIdCache_CacheKey + strconv.FormatInt(news.NewsInformationId, 10)
				srv.RedisCache.Set(cacheKey, _json.GetJsonString(news), EMNET_NewsInformation_GetNewsInfoByIdCache_CacheSeconds)
			}
		}
	}()
}

// GetStockNewsInfoList 获取大于指定NewsInformationId的所有股票相关的数据集合
func (srv *NewsInformationService) GetStockNewsInfoList(newsInformationId int64) (newsList []*expertnews_model.NewsInformation, err error) {
	newsList, err = srv.newsRepo.GetStockNewsInfoList(newsInformationId)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "GetNewsInfoList 获取大于指定NewsInformationId的所有数据集合 newsInformationId=%d", newsInformationId)
	}
	return newsList, err
}

// GetBlockNewsInfomation 获取板块相关资讯-缓存
func (srv *NewsInformationService) GetBlockNewsInfomation(blockCode string, top int64) (newsList []*expertnews_model.NewsInformation, err error) {
	if len(blockCode) == 0 {
		return nil, nil
	}
	//获取板块下最新的top条数据
	cacheKey := EMNET_NewsInformation_GetBlockNewsInfomation_CacheKey + blockCode
	result, err := srv.RedisCache.ZRevRange(cacheKey, 0, top-1)
	if err == redis.ErrNil || len(result) == 0 {
		return nil, nil
	}
	var field []interface{}
	for _, item := range result {
		field = append(field, item)
	}
	//获取资讯详情
	jsonArr, err := srv.RedisCache.HMGet(EMNET_NewsInformation_Detail_CacheKey, field...)
	if err == redis.ErrNil || len(jsonArr) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	jsonStr := "[" + strings.Join(jsonArr, ",") + "]"
	err = _json.Unmarshal(jsonStr, &newsList)
	return newsList, err
}

// SyncBlockNewsInformation 同步板块资讯信息
func (srv *NewsInformationService) SyncBlockNewsInformation() (err error) {
	//获取当前已同步的最大版本号
	maxVersion := ""
	err = srv.RedisCache.GetJsonObj(EMNET_NewsInformation_CurrentBlockNewsSyncVersion_CacheKey, &maxVersion)
	if err != nil && err != redis.ErrNil {
		shareNewsInformationLogger.ErrorFormat(err, "SyncBlockNewsInformation 同步板块资讯信息 GetMaxSyncRowVersion->异常")
		return err
	}

	//查询出新增的板块资讯集合
	newsInfoList, err := srv.GetBlockNewsInfoListAfterVersion(maxVersion)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "SyncBlockNewsInformation 同步板块资讯信息 GetBlockNewsInfoListAfterVersion->异常 maxVersion=%s", maxVersion)
		return err
	}
	if newsInfoList == nil || len(newsInfoList) == 0 {
		shareNewsInformationLogger.Debug("SyncBlockNewsInformation 同步板块资讯信息 当前没有新增的资讯需要同步")
	} else {
		shareNewsInformationLogger.Debug("SyncBlockNewsInformation 同步板块资讯信息 当开始同步新增的资讯 start")
		//遍历集合
		for _, newsInfo := range newsInfoList {
			blockListSrc := strings.Split(newsInfo.BlockCode, ",")
			if blockListSrc == nil || len(blockListSrc) == 0 {
				continue
			}
			for _, blockCodeSrc := range blockListSrc {
				if len(blockCodeSrc) < 4 {
					continue
				}
				_, err = srv.RedisCache.ZAdd(EMNET_NewsInformation_GetBlockNewsInfomation_CacheKey+blockCodeSrc, time.Time(newsInfo.PublishTime).UnixNano()/1000000, newsInfo.NewsInformationId)
				if err != nil {
					shareNewsInformationLogger.ErrorFormat(err, "SyncBlockNewsInformation 同步板块资讯信息 ZAdd 异常 newsInfo=%s", _json.GetJsonString(newsInfo))
					return err
				}
				err = srv.RedisCache.HSet(EMNET_NewsInformation_Detail_CacheKey, strconv.FormatInt(newsInfo.NewsInformationId, 10), _json.GetJsonString(newsInfo))
				if err != nil {
					shareNewsInformationLogger.ErrorFormat(err, "SyncBlockNewsInformation 同步板块资讯信息 HSet详情 异常 newsInfo=%s", _json.GetJsonString(newsInfo))
					return err
				}
			}
		}
		//记录已同步的最大版本号
		srv.RedisCache.SetJsonObj(EMNET_NewsInformation_CurrentBlockNewsSyncVersion_CacheKey, newsInfoList[len(newsInfoList)-1].SyncRowVersion)
		shareNewsInformationLogger.Debug("SyncBlockNewsInformation 同步板块资讯信息 当开始同步新增的资讯 end")
	}
	return err
}

// GetBlockNewsInfoListAfterVersion 获取大于指定版本的板块资讯列表
func (srv *NewsInformationService) GetBlockNewsInfoListAfterVersion(version string) (newsList []*expertnews_model.NewsInformation, err error) {
	newsList, err = srv.newsRepo.GetBlockNewsInfoListAfterVersion(version)
	if err != nil {
		shareNewsInformationLogger.ErrorFormat(err, "GetBlockNewsInfoListAfterVersion 获取大于指定版本的板块资讯列表 异常")
	}
	return newsList, err
}

type InsertData struct {
	Params []string
	Values []string
}
