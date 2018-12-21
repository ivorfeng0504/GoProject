package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/repository/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/stockapi"
	"github.com/devfeel/dotlog"
	"strconv"
	"sort"
	"git.emoney.cn/softweb/roboadvisor/protected/model/hotspot"
	"github.com/garyburd/redigo/redis"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type ExpertNews_TopicService struct {
	service.BaseService
	topicRepo *expertnews2.TopicRepository
}

var (
	shareExpertNews_TopicRepo *expertnews2.TopicRepository
	// shareExpertNewsLogger 共享的Logger实例
	shareExpertNewsTopicLogger dotlog.Logger
)

const (
	expertNewsTopicServiceName = "expertNewsInfoServiceLogger"

	CacheKey_TopicList_zset = _const.RedisKey_NewsPre + "ExpertNews.TopicList_zset"
	CacheKey_TopicList_hset = _const.RedisKey_NewsPre + "ExpertNews.TopicList_hset"
	CacheKey_TopicInfo      = _const.RedisKey_NewsPre + "ExpertNews.GetTopicByID:"
	CacheKey_TopicInfoInnerQuotes = _const.RedisKey_NewsPre+"ExpertNews.TopicInfoInnerQuotes:"
)

func NewExpertNews_TopicService() *ExpertNews_TopicService {
	expertNewsTopicService := &ExpertNews_TopicService{
		topicRepo: shareExpertNews_TopicRepo,
	}
	expertNewsTopicService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return expertNewsTopicService
}

// GetTopicListPage 分页获取主题列表
func (service *ExpertNews_TopicService) GetTopicListPage(pageIndex int, PageSize int) ([]*expertnews.ExpertNews_Topic, int, error) {
	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--
	return service.GetTopicListByCachePage(int64(startnum), int64(endnum))
}

// GetTopicListByCachePage redis分页获取主题列表
func (service *ExpertNews_TopicService) GetTopicListByCachePage(begin int64, end int64) ([]*expertnews.ExpertNews_Topic, int, error) {
	var topicList []*expertnews.ExpertNews_Topic
	var err error
	var totalnum int

	// redis 数据分页获取
	var retStrs []string
	totalnum, err = service.RedisCache.ZCard(CacheKey_TopicList_zset)
	retStrs, err = service.RedisCache.ZRevRange(CacheKey_TopicList_zset, begin, end)

	for i, _ := range retStrs {
		topicID,_:=strconv.Atoi(retStrs[i])
		expertNews_Topic := new(expertnews.ExpertNews_Topic)
		expertNews_Topic, _ = service.GetTopicInfoByID(topicID)
		topicList = append(topicList, expertNews_Topic)
	}

	if err == nil {
		return topicList, totalnum, nil
	}
	return topicList, totalnum, err
}

// GetTopicInfoByID 获取主题详情
func (service *ExpertNews_TopicService) GetTopicInfoByID(topicID int) (*expertnews.ExpertNews_Topic, error) {
	expertNews_Topic := new(expertnews.ExpertNews_Topic)
	redisKey_TopicInfo := CacheKey_TopicInfoInnerQuotes + strconv.Itoa(topicID)

	err := service.RedisCache.GetJsonObj(redisKey_TopicInfo, expertNews_Topic)
	if err == redis.ErrNil {
		err = nil
	}
	return expertNews_Topic, err
}

// task任务使用 定时关联行情数据后存入redis
func (service *ExpertNews_TopicService) GetTopicInfoInnerQuotes(expertNews_Topic *expertnews.ExpertNews_Topic) (topicinfo *expertnews.ExpertNews_Topic, err error) {
	redisKey_TopicInfo := CacheKey_TopicInfoInnerQuotes + strconv.Itoa(expertNews_Topic.ID)

	if expertNews_Topic != nil {

		//获取板块涨跌幅
		var newBKList []*expertnews.Topic_BK
		oldBKList := expertNews_Topic.RelatedBKList

		for i, _ := range oldBKList {
			var bkinfo = new(expertnews.Topic_BK)
			bkinfo = oldBKList[i]

			var quotesData = new(stockapi.QuotesData)
			quotesData, err = stockapi.GetQuotes(bkinfo.BKCode, "BK")

			if err != nil {
				shareExpertNewsTopicLogger.Error(err, "stockapi.GetQuotes获取行情接口失败")
			}

			if quotesData != nil {
				bkinfo.BKF = quotesData.F
			}

			newBKList = append(newBKList, bkinfo)
		}

		//获取所有个股PE列表
		peList, _ := stockapi.GetStockPEList()

		//获取个股最新价和PE
		var newStockList []*expertnews.Topic_Stock
		oldStockList := expertNews_Topic.RelatedStockList
		for i, _ := range oldStockList {
			var stockinfo = new(expertnews.Topic_Stock)
			stockinfo = oldStockList[i]

			var quotesData = new(stockapi.QuotesData)
			quotesData, err = stockapi.GetQuotes(stockinfo.StockCode, "stock")

			if err != nil {
				shareExpertNewsTopicLogger.Error(err, "stockapi.GetQuotes获取行情接口失败")
			}

			if quotesData != nil {
				stockinfo.StockF = quotesData.F     //涨跌幅
				stockinfo.StockPrice = quotesData.P //最新价
			}

			stockinfo.StockPE = GetStockPE(stockinfo.StockCode, peList) //PE

			newStockList = append(newStockList, stockinfo)
		}

		expertNews_Topic.RelatedBKList = newBKList
		expertNews_Topic.RelatedStockList = newStockList
	}

	//关联行情数据后存入redis
	if expertNews_Topic != nil {
		service.RedisCache.SetJsonObj(redisKey_TopicInfo, expertNews_Topic)
	}

	return expertNews_Topic, nil
}

// GetStockPE 获取某个股pe
func GetStockPE(stockCode string,peList [][]string) (pe string) {
	for i, _ := range peList {
		if peList[i][0] == stockCode {
			return peList[i][1]
		}
	}
	return ""
}

// GetTopicList db获取主题列表
func (service *ExpertNews_TopicService) GetTopicList() (topicList []*expertnews.ExpertNews_Topic, err error) {
	topicList, err = service.topicRepo.GetTopicList()
	if err != nil {
		shareExpertNewsTopicLogger.ErrorFormat(err, "获取主题列表不正确。ExpertNews_TopicService.GetTopicList()")
		return nil, err
	}
	return topicList, err
}

func init() {
	protected.RegisterServiceLoader(expertNewsTopicServiceName, expertNews_TopicServiceLoader)
}

func expertNews_TopicServiceLoader() {
	shareExpertNewsTopicLogger = dotlog.GetLogger(expertNewsTopicServiceName)
	shareExpertNews_TopicRepo = expertnews2.NewTopicRepository(protected.DefaultConfig)
}



//统计主题关注的股票的统计排行
func (service *ExpertNews_TopicService) StatTopicFocusStock()(stocks *[]StockFocusTimes, err error){
	topicList, err := service.topicRepo.GetTopicList()

	statmp  := make(map[string]int)
	stocknamemp := make(map[string]string)

	for i, _ := range topicList {
		for _,v :=range topicList[i].RelatedStockList{
			oldint :=  statmp[v.StockCode]
			statmp[v.StockCode] = oldint+1;
			stocknamemp[v.StockCode] = v.StockName
		}
	}
	var newMp = make([]int, 0)
	var newMpKey = make([]string, 0)
	for oldk, v := range statmp {
		newMp = append(newMp, v)
		newMpKey = append(newMpKey, oldk)
	}
	sort.Ints(newMp)

	var focuslist = make([]StockFocusTimes,0)
	var x int ;x=1
	for k, v := range newMp {
		//fmt.Println("根据value排序后的新集合》》  key:", newMpKey[k], "    value:", v)
		t := StockFocusTimes{StockCode:newMpKey[k],StockName:stocknamemp[newMpKey[k]],Times:v}
		focuslist = append(focuslist,t)
		x++
		if x>10 {
			break
		}
	}

	return &focuslist,err
}


type StockFocusTimes struct {
	StockCode string
	StockName string
	Times int
}



//获取天眼策略的题材专题信息
func (srv *ExpertNews_TopicService) GetHotspotInfo(hotspotid int)(hotspotinfo *model.Hotspot,error error){
	rediskey := "RoboAdvisor:Hotspot:Cache:GetHotspotById:"+ strconv.Itoa(hotspotid)
	err := srv.RedisCache.GetJsonObj(rediskey, &hotspotinfo)
	if err == redis.ErrNil {
		return nil, nil
	}
	return hotspotinfo, err

}