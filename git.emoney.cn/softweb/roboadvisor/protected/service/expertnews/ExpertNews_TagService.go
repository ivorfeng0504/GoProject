package expertnews

import (
	"encoding/json"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"strconv"
	"strings"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected/service/resapi"
)

type ExpertNews_TagService struct {
	service.BaseService
}

var (
	// shareExpertNewsLogger 共享的Logger实例
	shareExpertNewsTagLogger dotlog.Logger
)

const (
	CacheKey_NewsListByColIDAndSIDAndTagID_zset = _const.RedisKey_NewsPre + "ExpertNews_TagRelation.GetNewsListByColIDAndSIDAndTagID_zset:" //资讯列表zset（根据栏目和标签）
	CacheKey_NewsListByColIDAndSIDAndTagID_hset = _const.RedisKey_NewsPre + "ExpertNews_TagRelation.GetNewsListByColIDAndSIDAndTagID_hset:" //资讯列表hset（根据栏目和标签）

	CacheKey_NewsListByColIDAndClientSID_zset = _const.RedisKey_NewsPre + "GetNewsListByColumnIDAndClientStrategyID_zset:"//重要提示非置顶文章zset
	CacheKey_NewsListByColIDAndClientSID_hset = _const.RedisKey_NewsPre + "GetNewsListByColumnIDAndClientStrategyID_hset:"//重要提示非置顶文章hset

	CacheKey_TopNewsListByColIDAndClientSID_zset = _const.RedisKey_NewsPre + "GetTopNewsListByColumnIDAndClientStrategyID_zset:"//重要提示置顶文章zset
	CacheKey_TopNewsListByColIDAndClientSID_hset = _const.RedisKey_NewsPre + "GetTopNewsListByColumnIDAndClientStrategyID_hset:"//重要提示置顶文章hset

	expertNewsTagInfoServiceName = "ExpertNews_TagServiceServiceLogger"
)

func NewExpertNews_TagService() *ExpertNews_TagService {
	expertNewsService := &ExpertNews_TagService{
	}
	expertNewsService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return expertNewsService
}

/* GetTopNewsList_ImportantTips 分页获取策略重要提示
columnID:专家资讯栏目id
StrategyID：客户端策略ID(
*/
func (service *ExpertNews_TagService) GetTopNewsList_ImportantTips(StrategyID int) ([]*model.NewsInfo, error) {
	var newslist []*model.NewsInfo
	var err error
	var redisKey_newslist_zset string
	var redisKey_newslist_hset string

	columnID := config.CurrentConfig.ColumnIDImportantTips

	redisKey_newslist_zset = CacheKey_TopNewsListByColIDAndClientSID_zset + columnID + "_" + strconv.Itoa(StrategyID)
	redisKey_newslist_hset = CacheKey_TopNewsListByColIDAndClientSID_hset + columnID + "_" + strconv.Itoa(StrategyID)

	// redis 数据分页获取(只取1条）
	var retStrs []string
	retStrs, err = service.RedisCache.ZRevRange(redisKey_newslist_zset, 0, 0)
	if err != nil {
		shareExpertNewsTagLogger.ErrorFormat(err, "GetTopNewsList_ImportantTips redis获取策略资讯列表(重要提示置顶一条)失败ZRevRange key=%s", redisKey_newslist_zset)
		return nil, err
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
		shareExpertNewsTagLogger.ErrorFormat(err, "GetTopNewsList_ImportantTips 获取策略资讯列表(重要提示置顶一条)反序列化失败 newslist=%s", stringByte)
		return nil, err
	}

	return newslist, err
}

/* GetNewsListPage_ImportantTips 分页获取策略重要提示
columnID:专家资讯栏目id
TagID：专家策略ID(0:只根据栏目查询)
pageIndex：当前页面
PageSize：一页显示条数
*/
func (service *ExpertNews_TagService) GetNewsListPage_ImportantTips(StrategyID int, pageIndex int, PageSize int) ([]*model.NewsInfo, int, error) {
	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--

	var newslist []*model.NewsInfo
	var err error
	var redisKey_newslist_zset string
	var redisKey_newslist_hset string

	columnID := config.CurrentConfig.ColumnIDImportantTips

	redisKey_newslist_zset = CacheKey_NewsListByColIDAndClientSID_zset + columnID + "_" + strconv.Itoa(StrategyID)
	redisKey_newslist_hset = CacheKey_NewsListByColIDAndClientSID_hset + columnID + "_" + strconv.Itoa(StrategyID)

	// redis 数据分页获取
	var retStrs []string
	retStrs, err = service.RedisCache.ZRevRange(redisKey_newslist_zset, int64(startnum), int64(endnum))
	if err != nil {
		shareExpertNewsTagLogger.ErrorFormat(err, "GetNewsListPage_ImportantTips redis分页获取策略资讯列表(重要提示)失败ZRevRange key=%s begin=%d end=%d", redisKey_newslist_zset, startnum, endnum)
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
		shareExpertNewsTagLogger.ErrorFormat(err, "GetNewsListPage_ImportantTips 获取策略资讯列表(重要提示)反序列化失败 newslist=%s", stringByte)
		return nil, 0, err
	}

	return newslist, 0, err
}

/* GetTagNewsListPage 分页获取策略资讯列表
columnID:专家资讯栏目id
TagID：专家策略ID(0:只根据栏目查询)
pageIndex：当前页面
PageSize：一页显示条数
*/
func (service *ExpertNews_TagService) GetStrategyTagNewsListPage(columnID int,StrategyID int, TagID int, pageIndex int64, PageSize int64) ([]*model.NewsInfo, int, error) {
	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--

	return service.GetStrategyTagNewsListByCachePage(columnID, StrategyID, TagID, startnum, endnum)
}

// GetTagNewsListByCachePage redis分页获取策略资讯列表(区分标签)
func (service *ExpertNews_TagService) GetStrategyTagNewsListByCachePage(columnID int,StrategyID int, TagID int, begin int64, end int64) ([]*model.NewsInfo, int, error) {
	var newslist []*model.NewsInfo
	var err error
	var totalnum int
	var redisKey_newslist_zset string
	var redisKey_newslist_hset string

	redisKey_newslist_zset = CacheKey_NewsListByColIDAndSIDAndTagID_zset + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID) + "_" + strconv.Itoa(TagID)
	redisKey_newslist_hset = CacheKey_NewsListByColIDAndSIDAndTagID_hset + strconv.Itoa(columnID) + "_" + strconv.Itoa(StrategyID) + "_" + strconv.Itoa(TagID)

	// redis 数据分页获取
	var retStrs []string
	totalnum, err = service.RedisCache.ZCard(redisKey_newslist_zset)
	if err != nil {
		shareExpertNewsTagLogger.ErrorFormat(err, "GetStrategyTagNewsListByCachePage redis分页获取策略资讯列表(区分标签)失败ZCard key=%s", redisKey_newslist_zset)
		return nil, 0, err
	}
	retStrs, err = service.RedisCache.ZRevRange(redisKey_newslist_zset, begin, end)
	if err != nil {
		shareExpertNewsTagLogger.ErrorFormat(err, "GetStrategyTagNewsListByCachePage redis分页获取策略资讯列表(区分标签)失败ZRevRange key=%s begin=%d end=%d", redisKey_newslist_zset, begin, end)
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
		shareExpertNewsTagLogger.ErrorFormat(err, "GetStrategyTagNewsListByCachePage 获取策略资讯列表(区分标签)反序列化失败 newslist=%s", stringByte)
		return nil, 0, err
	}

	return newslist, totalnum, err
}

// 获取专家策略ID 根据客户端策略ID
func(service *ExpertNews_TagService) GetStrategyIDByClientStrategyID(clientStrategyID int) (int,error) {
	//confRelation := config.CurrentConfig.ClientExpertStrategyRelation
	confRelation, _ := resapi.GetClientExpertStrategyConfig()

	relationList := strings.Split(confRelation, "|")
	for i, _ := range relationList {
		relation_str := relationList[i]
		relation := strings.Split(relation_str, ",")

		if len(relation) > 0 {
			confClientStrategyID, err := strconv.Atoi(relation[0])
			if err != nil {
				shareExpertNewsTagLogger.ErrorFormat(err, "GetStrategyIDByClientStrategyID 客户端策略ID和专家策略ID对应关系转换 relation=%s", relation)
				return 0, nil
			}
			if clientStrategyID == confClientStrategyID {
				expertStrategID, err := strconv.Atoi(relation[1])
				if err != nil {
					shareExpertNewsTagLogger.ErrorFormat(err, "GetStrategyIDByClientStrategyID 客户端策略ID和专家策略ID对应关系转换 relation=%s", relation)
				}
				return expertStrategID, nil
			}

		}
	}
	return 0, nil
}

func init() {
	protected.RegisterServiceLoader(expertNewsTagInfoServiceName, expertNews_TagServiceLoader)
}

func expertNews_TagServiceLoader() {
	shareExpertNewsTagLogger = dotlog.GetLogger(expertNewsTagInfoServiceName)
}
