package live

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	liverepo "git.emoney.cn/softweb/roboadvisor/protected/repository/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/time"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type LiveContentService struct {
	service.BaseService
	liveContentRepo *liverepo.LiveContentRepository
}

var (
	// shareLiveContentRepo 共享的仓储
	shareLiveContentRepo *liverepo.LiveContentRepository

	// shareLiveContentLogger 共享的Logger实例
	shareLiveContentLogger dotlog.Logger
)

const (
	liveContentServiceName = "LiveContentService"
	//EMNET开头标识需要与.NET项目中的KEY保持一致
	EMNET_LiveContent_CachePreKey                         = "EMoney:Live:LiveContent"
	EMNET_LiveContent_LiveContentList                     = EMNET_LiveContent_CachePreKey + ":LiveContentList"
	EMNET_LiveContent_GetContentListByLidIncreNumCacheKey = EMNET_LiveContent_CachePreKey + ":LiveContentSortedSet"
	EMNET_LiveContent_GetTopContentCacheKey               = EMNET_LiveContent_CachePreKey + ":TopContentCache:"
)

func NewLiveContentService() *LiveContentService {
	service := &LiveContentService{
		liveContentRepo: shareLiveContentRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return service
}

// GetContentList 根据topicId查询直播内容
func (service *LiveContentService) GetContentList(topicId int) (contentList []*livemodel.LiveContent, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getContentListDB(topicId)
	default:
		return service.getContentListCache(topicId)
	}
}

// getContentListCache 根据topicId查询直播内容
func (service *LiveContentService) getContentListCache(topicId int) (contentList []*livemodel.LiveContent, err error) {
	cacheKey := EMNET_LiveContent_LiveContentList + strconv.Itoa(topicId)
	err = service.RedisCache.GetJsonObj(cacheKey, &contentList)
	if err == redis.ErrNil {
		return contentList, nil
	}
	if err != nil {
		shareLiveContentLogger.Error(err, " getContentListCache 根据topicId查询直播内容 异常")
	}
	return contentList, err
}

// getContentListDB 根据topicId查询直播内容
func (service *LiveContentService) getContentListDB(topicId int) (contentList []*livemodel.LiveContent, err error) {
	contentList, err = service.liveContentRepo.GetContentList(topicId)
	if err != nil {
		shareLiveContentLogger.Error(err, " getContentListDB 根据topicId查询直播内容 异常")
	}
	return contentList, err
}

// GetContentIncre 获取递增的直播内容
// topicId 主题Id
// lastDate上次最后一条消息的时间
func (service *LiveContentService) GetContentIncre(topicId int, lastDate time.Time) ([]*livemodel.LiveContent, error) {
	contentList, err := service.GetContentList(topicId)
	if err != nil {
		return contentList, err
	}
	if contentList == nil {
		return nil, nil
	}

	//筛选新增的消息
	var increContentList []*livemodel.LiveContent
	for _, content := range contentList {
		if time.Time(content.CreateTime).After(lastDate) {
			increContentList = append(increContentList, content)
		}
	}
	return increContentList, nil
}

// GetTopContent 根据topicId查询置顶的直播内容
func (service *LiveContentService) GetTopContent(topicId int) (topContent *livemodel.LiveContent, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getTopContentDB(topicId)
	default:
		return service.getTopContentCache(topicId)
	}
}

// getTopContentCache 根据topicId查询置顶的直播内容
func (service *LiveContentService) getTopContentCache(topicId int) (topContent *livemodel.LiveContent, err error) {
	cacheKey := EMNET_LiveContent_GetTopContentCacheKey + strconv.Itoa(topicId)
	err = service.RedisCache.GetJsonObj(cacheKey, &topContent)
	if err == redis.ErrNil {
		return topContent, nil
	}
	if err != nil {
		shareLiveContentLogger.Error(err, " getTopContentCache 根据topicId查询置顶的直播内容 异常")
	}
	return topContent, err
}

// getTopContentDB 根据topicId查询置顶的直播内容
func (service *LiveContentService) getTopContentDB(topicId int) (topContent *livemodel.LiveContent, err error) {
	topContent, err = service.liveContentRepo.GetTopContent(topicId)
	if err != nil {
		shareLiveContentLogger.Error(err, " getTopContentDB 根据topicId查询置顶的直播内容 异常")
	}
	return topContent, err
}

//GetContentListByLidIncreNum 根据直播室Id获取指定时间的发布内容
//从EMoney.ContentLivePlat.Business.LiveContentBll中迁移过来，逻辑与其保持一致
func (service *LiveContentService) GetContentListByLidIncreNum(topicId int, date time.Time) int {
	start := date.Sub(_time.ToDay()).Seconds()*1000 + 10
	end := time.Now().Sub(_time.ToDay()).Seconds() * 1000
	key := EMNET_LiveContent_GetContentListByLidIncreNumCacheKey + strconv.Itoa(topicId)
	count, err := service.RedisCache.ZCount(key, int64(start), int64(end))
	if err != nil {
		count = 0
	}
	return count
}

func init() {
	protected.RegisterServiceLoader(liveContentServiceName, liveContentServiceLoader)
}

func liveContentServiceLoader() {
	shareLiveContentRepo = liverepo.NewLiveContentRepository(protected.DefaultConfig)
	shareLiveContentLogger = dotlog.GetLogger(liveContentServiceName)
}
