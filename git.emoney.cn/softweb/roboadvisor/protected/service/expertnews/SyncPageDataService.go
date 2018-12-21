package expertnews

import (
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"encoding/json"
	"git.emoney.cn/softweb/roboadvisor/protected/model/yqq"
	"sort"
	"github.com/devfeel/mapper"
)

type SyncPageDataService struct {
	service.BaseService
}

var (
	ExpertSyncPageDataLogger dotlog.Logger
)

const (
	expertSyncPageDataServiceName = "expertSyncPageDataServiceLogger"

	//首页数据redis key
	CacheKey_SyncPageData_Index = _const.RedisKey_NewsPre + "ExpertNews.SyncPageData_Index"
)

func NewSyncPageDataService() *SyncPageDataService {
	syncPageDataService := &SyncPageDataService{
	}
	syncPageDataService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return syncPageDataService
}

// SyncIndexData task同步首页数据入redis（今日头条、策略看盘、主题、专家直播、盘后预测）
func (service *SyncPageDataService) SyncIndexData() (error) {
	IndexDataModel := new(expertnews.IndexDataInfo)

	//今日头条
	var todayNewsList []*agent.ExpertNewsInfo
	newsSrv := NewNewsInformationService()
	newsList, err := newsSrv.GetTodayNewsCacheV2()
	if err != nil {
		todayNewsList = nil
	}
	err = mapper.MapperSlice(newsList, &todayNewsList)
	if err != nil {
		todayNewsList = nil
	}
	IndexDataModel.JrttList = todayNewsList

	//策略看盘
	expertNewsService := NewExpertNews_StrategyService()
	strategyNewslist, err := expertNewsService.GetStrategyNewsList_Top6()
	if err != nil {
		strategyNewslist = nil
	}
	IndexDataModel.ClkpList = strategyNewslist

	//主题
	expertTopicService := NewExpertNews_TopicService()
	topiclist, _, err := expertTopicService.GetTopicListPage(1, 3)
	if err != nil {
		topiclist = nil
	}
	IndexDataModel.TopicList = topiclist

	//专家直播
	YqqLiveService := ConExpertNews_YqqService()
	RoomListStr, err := YqqLiveService.GetLiveRoomList()
	var YqqLive yqq.YqqRet
	err = json.Unmarshal([]byte(RoomListStr), &YqqLive)
	if err != nil {
		topiclist = nil
	}
	yqqroomdata := YqqLive.Data
	if len(yqqroomdata) > 2 {
		sort.Sort(YqqRoomDatas(yqqroomdata))
		yqqroomdata = yqqroomdata[:2]
	}
	IndexDataModel.ZjzbList = yqqroomdata

	//盘后预测
	var closingNewsList []*agent.ExpertNewsInfo
	ph_newsList, err := newsSrv.GetClosingNewsCache()
	if err != nil {
		closingNewsList = nil
	} else {
		err = mapper.MapperSlice(ph_newsList, &closingNewsList)
		if len(closingNewsList) > 2 {
			closingNewsList = closingNewsList[:2]
		}
	}
	IndexDataModel.PhycList = closingNewsList

	_, err = service.RedisCache.SetJsonObj(CacheKey_SyncPageData_Index, IndexDataModel)
	if err != nil {
		ExpertSyncPageDataLogger.ErrorFormat(err, "SyncIndexData 专家资讯首页内容存入redis失败")
	}
	return err
}

// GetSyncIndexData 获取服务自动更新后的首页数据（专家资讯）
func (service *SyncPageDataService) GetSyncIndexData() (*expertnews.IndexDataInfo,error) {
	IndexDataModel := new(expertnews.IndexDataInfo)
	err := service.RedisCache.GetJsonObj(CacheKey_SyncPageData_Index, IndexDataModel)

	if err != nil {
		ExpertSyncPageDataLogger.ErrorFormat(err, "GetSyncIndexData 获取专家资讯首页内容失败")
		return nil, err
	}

	return IndexDataModel, nil
}

//专家直播条件排序
type YqqRoomDatas []yqq.YqqRoom

//Len()
func (s YqqRoomDatas) Len() int {
	return len(s)
}
func (s YqqRoomDatas) Less(i, j int) bool {
	return s[i].FansNum > s[j].FansNum
}
func (s YqqRoomDatas) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}


func init() {
	protected.RegisterServiceLoader(expertSyncPageDataServiceName, SyncPageDataServiceLoader)
}

func SyncPageDataServiceLoader() {
	ExpertSyncPageDataLogger = dotlog.GetLogger(expertSyncPageDataServiceName)
}

