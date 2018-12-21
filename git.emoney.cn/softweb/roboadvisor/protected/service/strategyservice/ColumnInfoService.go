package strategyservice

import (
	"encoding/json"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	strategyservice_model "git.emoney.cn/softweb/roboadvisor/protected/model/strategyservice"
	strategyservice_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	strategyservice_vmmodel "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/mapper"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
)

type ColumnInfoService struct {
	service.BaseService
	columnInfoRepo *strategyservice_repo.ColumnInfoRepository
}

var (
	shareColumnInfoRepo   *strategyservice_repo.ColumnInfoRepository
	shareColumnInfoLogger dotlog.Logger
)

const (
	columnInfoServiceServiceName = "ColumnInfoService"
	EMNET_ColumnInfo_PreCacheKey = "ColumnInfoService:"

	CacheKey_Getstrategynewsdetaillist             = _const.RedisKey_NewsPre + "clxx_GetStrategyNewsDetailList"
	CacheKey_Getstrategynewslist                   = _const.RedisKey_NewsPre + "clxx_GetNewsInfoByColIDAndNewsType:"
	EMNET_ColumnInfo_GetIndexStrategyNews_CacheKey = EMNET_ColumnInfo_PreCacheKey + "GetIndexStrategyNews:"

	//用户培训-策略实战（最近一个月）
	RedisKey_1MonthMultiMediaList = _const.RedisKey_NewsPre + "GetMultiMediaNewsList_1Month"
)

func NewColumnInfoService() *ColumnInfoService {
	service := &ColumnInfoService{
		columnInfoRepo: shareColumnInfoRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return service
}

// GetStrategyNewsList 根据栏目id和资讯类型分页获取资讯列表
/*
newsIDStr：批量资讯ID
colID：栏目ID
newsType=0:文章
newsType=1:视频
*/
func (service *ColumnInfoService) GetStrategyNewsList(colID int, newsType int, pageIndex int64, PageSize int64) ([]*model.NewsInfo, error) {
	rediskey := fmt.Sprintf(CacheKey_Getstrategynewslist+"%d_%d", colID, newsType)
	return service.GetStrategyNewsListByPage(rediskey, pageIndex, PageSize)
}

// GetIndexStrategyNews 获取首页最新的一条策略资讯与策略视频
func (srv *ColumnInfoService) GetIndexStrategyNewsCache(clientStrategyGroupId int) (result *strategyservice_vmmodel.IndexStrategyNewsVM, err error) {
	cacheKey := EMNET_ColumnInfo_GetIndexStrategyNews_CacheKey
	jsonStr, err := srv.RedisCache.HGet(cacheKey, strconv.Itoa(clientStrategyGroupId))
	if err == redis.ErrNil {
		return nil, nil
	}
	err = _json.Unmarshal(jsonStr, &result)
	return result, err
}

// RefreshIndexStrategyNewsCache 刷新首页最新的一条策略资讯与策略视频 clientStrategyGroupId 组策略Id  columnList栏目Id集合
func (srv *ColumnInfoService) RefreshIndexStrategyNewsCache(clientStrategyGroupId int, columnList []int) (result *strategyservice_vmmodel.IndexStrategyNewsVM, err error) {
	cacheKey := EMNET_ColumnInfo_GetIndexStrategyNews_CacheKey
	result = new(strategyservice_vmmodel.IndexStrategyNewsVM)
	result.StrategyNews, err = srv.columnInfoRepo.GetNewestNewsInfoByColumnList(columnList, _const.NewsType_News)
	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "RefreshIndexStrategyNewsCache 刷新首页最新的一条策略资讯与策略视频 -->资讯 异常 clientStrategyGroupId=%d columnList=%s", clientStrategyGroupId, _json.GetJsonString(columnList))
		return result, err
	}
	result.StrategyVideoNews, err = srv.columnInfoRepo.GetNewestNewsInfoByColumnList(columnList, _const.NewsType_Video)
	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "RefreshIndexStrategyNewsCache 刷新首页最新的一条策略资讯与策略视频 -->视频 异常 clientStrategyGroupId=%d columnList=%s", clientStrategyGroupId, _json.GetJsonString(columnList))
		return result, err
	}
	err = srv.RedisCache.HSet(cacheKey, strconv.Itoa(clientStrategyGroupId), _json.GetJsonString(result))
	return result, err
}

// 分页获取资讯列表
func (service *ColumnInfoService) GetStrategyNewsListByPage(rediskey string, pageIndex int64, PageSize int64) ([]*model.NewsInfo, error) {
	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--

	var result []*model.NewsInfo
	var newsList []*model.NewsInfo
	newsIDStr, err := service.RedisCache.ZRevRange(rediskey, startnum, endnum)

	//取出资讯ID放入string数组，HMGet使用
	args := make([]interface{}, len(newsIDStr))
	for i, v := range newsIDStr {
		args[i] = v
	}

	newsDetailskey := CacheKey_Getstrategynewsdetaillist
	stringResult, err := service.RedisCache.HMGet(newsDetailskey, args...)

	stringByte := "[" + strings.Join(stringResult, ",") + "]"

	err = json.Unmarshal([]byte(stringByte), &newsList)
	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "GetStrategyNewsDetailList 策略学习资讯列表反序列化失败 newslist=%s", stringByte)
		return nil, err
	}

	for _, v := range newsList {
		newsinfo := v
		v.Summary = ""
		v.NewsContent = ""

		result = append(result, newsinfo)
	}
	return newsList, err
}

// GetStrategyNewsDetailList 根据栏目id和资讯类型批量获取资讯列表
func (service *ColumnInfoService) GetStrategyNewsDetailList(newsIDStr []string) ([]*model.NewsInfo, error) {
	var newsList []*model.NewsInfo
	rediskey := CacheKey_Getstrategynewsdetaillist

	//取出资讯ID放入string数组，HMGet使用
	args := make([]interface{}, len(newsIDStr))
	for i, v := range newsIDStr {
		args[i] = v
	}
	stringResultSrc, err := service.RedisCache.HMGet(rediskey, args...)

	//去除无效的空数据 防止序列化报错
	var stringResult []string
	for _, item := range stringResultSrc {
		if len(item) > 0 {
			stringResult = append(stringResult, item)
		}
	}
	stringByte := "[" + strings.Join(stringResult, ",") + "]"
	err = json.Unmarshal([]byte(stringByte), &newsList)

	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "GetStrategyNewsDetailList 策略学习资讯列表反序列化失败 newslist=%s", stringByte)
		return nil, err
	}

	return newsList, err
}

// redis读取其他相关资讯
func (service *ColumnInfoService) GetNewstStrategyNews(clientStrategyIdList []string, newsType int) (newsList []*model.NewsInfo, err error) {
	var resultall []*model.NewsInfo
	colDict, err := NewClientStrategyInfoService().GetClientStrategyInfoDict()
	if err != nil {
		return newsList, err
	}
	//去重字典
	distinctDict := make(map[int]bool)
	for _, v := range clientStrategyIdList {
		clientStrategyId, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		colID, success := colDict[clientStrategyId]
		if success == false {
			continue
		}
		newslist, err := service.GetStrategyNewsList(colID, newsType, 1, 5)
		if err == nil && newslist != nil {
			for _, item := range newslist {
				_, exist := distinctDict[item.ID]
				if exist == false {
					resultall = append(resultall, item)
					distinctDict[item.ID] = true
				}
			}
		}
	}
	return resultall, err
}

// GetStrategyNewsList_MultiMedia_ByDate 根据栏目id和资讯类型获取某天的资讯列表
func (service *ColumnInfoService) GetStrategyNewsList_MultiMedia_ByDate(colID int, newsType int, date *time.Time) ([]*strategyservice_vmmodel.ExpertNews_MultiMedia_List, error) {
	allList, err := service.GetStrategyNewsListByPage_MultiMedia_ByDate_All(date)
	if err != nil {
		return nil, err
	}
	var filterList []*strategyservice_vmmodel.ExpertNews_MultiMedia_List
	showPos := `"ShowPosition":"` + strconv.Itoa(colID) + `"`
	for _, item := range allList {
		if strings.Contains(item.ShowPosition, showPos) {
			filterList = append(filterList, &strategyservice_vmmodel.ExpertNews_MultiMedia_List{
				ID:             item.ID,
				NewsType:       item.NewsType,
				Title:          item.Title,
				Summary:        item.Summary,
				CoverImg:       item.CoverImg,
				TagInfo:        item.TagInfo,
				ClickNum:       item.ClickNum,
				AudioURL:       item.AudioURL,
				VideoPlayURL:   item.VideoPlayURL,
				LiveURL:        item.LiveURL,
				LiveVideoURL:   item.LiveVideoURL,
				Live_StartTime: mapper.JSONTime(item.Live_StartTime),
				Live_EndTime:   mapper.JSONTime(item.Live_EndTime),
				CreateTime:     mapper.JSONTime(item.CreateTime),
				LastModifyTime: mapper.JSONTime(item.LastModifyTime),
				ShowPosition:   item.ShowPosition,
				From:           item.From,
			})
		}
	}
	return filterList, nil
}

// GetStrategyNewsList_MultiMedia 根据栏目id和资讯类型分页获取资讯列表
/*
colID：栏目ID
newsType=3:多媒体资讯类型
*/
func (service *ColumnInfoService) GetStrategyNewsList_MultiMedia(colID int, newsType int) ([]*strategyservice_vmmodel.ExpertNews_MultiMedia_List, error) {
	rediskey := fmt.Sprintf(CacheKey_Getstrategynewslist+"%d_%d", colID, newsType)

	//获取最近一个月资讯列表
	nTime := time.Now().AddDate(0, 0, 2)
	before1MonthTime := nTime.AddDate(0, -1, -2)

	minscore := strconv.FormatInt(before1MonthTime.Unix()*1000, 10)
	maxscore := strconv.FormatInt(nTime.Unix()*1000, 10)

	var newsList []*strategyservice_vmmodel.ExpertNews_MultiMedia_List
	newsIDStr, err := service.RedisCache.ZREVRangeByScore(rediskey, maxscore, minscore, false)

	shareColumnInfoLogger.InfoFormat("GetStrategyNewsList_MultiMedia 实战培训最近一个月数据时间戳 rediskey:%s begin:%d end:%d newsIDStr:%s", rediskey, maxscore, minscore, newsIDStr)

	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "GetStrategyNewsListByPage_MultiMedia 获取策略学习多媒体资讯列表失败 rediskey=%s maxscore=%s", rediskey, maxscore)
		return nil, err
	}

	//取出资讯ID放入string数组，HMGet使用
	args := make([]interface{}, len(newsIDStr))
	for i, v := range newsIDStr {
		args[i] = v
	}

	newsDetailskey := CacheKey_Getstrategynewsdetaillist
	stringResult, err := service.RedisCache.HMGet(newsDetailskey, args...)

	stringByte := "[" + strings.Join(stringResult, ",") + "]"

	err = json.Unmarshal([]byte(stringByte), &newsList)
	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "GetStrategyNewsListByPage_MultiMedia 策略学习多媒体资讯列表反序列化失败 newslist=%s", stringByte)
		return nil, err
	}

	return newsList, err
}

// GetStrategyNewsListByPage_MultiMedia 分页获取资讯列表-多媒体资讯
func (service *ColumnInfoService) GetStrategyNewsListByPage_MultiMedia(rediskey string, pageIndex int64, PageSize int64) ([]*strategyservice_vmmodel.ExpertNews_MultiMedia_List, error) {
	startnum := pageIndex
	endnum := pageIndex * PageSize

	if pageIndex == 1 {
		startnum = 0
	} else {
		startnum = PageSize * (pageIndex - 1)
	}
	endnum--

	var newsList []*strategyservice_vmmodel.ExpertNews_MultiMedia_List
	newsIDStr, err := service.RedisCache.ZRevRange(rediskey, startnum, endnum)

	//取出资讯ID放入string数组，HMGet使用
	args := make([]interface{}, len(newsIDStr))
	for i, v := range newsIDStr {
		args[i] = v
	}

	newsDetailskey := CacheKey_Getstrategynewsdetaillist
	stringResult, err := service.RedisCache.HMGet(newsDetailskey, args...)

	stringByte := "[" + strings.Join(stringResult, ",") + "]"

	err = json.Unmarshal([]byte(stringByte), &newsList)
	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "GetStrategyNewsListByPage_MultiMedia 策略学习多媒体资讯列表反序列化失败 newslist=%s", stringByte)
		return nil, err
	}

	return newsList, err
}

// GetStrategyNewsDetailList_MultiMedia 根据栏目id和资讯类型批量获取多媒体资讯列表
func (service *ColumnInfoService) GetStrategyNewsDetailList_MultiMedia(newsIDStr []string) ([]*strategyservice_vmmodel.ExpertNews_MultiMedia_List, error) {
	var newsList []*strategyservice_vmmodel.ExpertNews_MultiMedia_List
	rediskey := CacheKey_Getstrategynewsdetaillist

	//取出资讯ID放入string数组，HMGet使用
	args := make([]interface{}, len(newsIDStr))
	for i, v := range newsIDStr {
		args[i] = v
	}
	stringResultSrc, err := service.RedisCache.HMGet(rediskey, args...)

	//去除无效的空数据 防止序列化报错
	var stringResult []string
	for _, item := range stringResultSrc {
		if len(item) > 0 {
			stringResult = append(stringResult, item)
		}
	}
	stringByte := "[" + strings.Join(stringResult, ",") + "]"
	err = json.Unmarshal([]byte(stringByte), &newsList)

	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "GetStrategyNewsDetailList_MultiMedia 策略学习多媒体资讯详情反序列化失败 newslist=%s", stringByte)
		return nil, err
	}

	return newsList, err
}

// GetStrategyNewsListByPage_MultiMedia_ByDate_All 分页获取指定日期的实战培训课程
func (service *ColumnInfoService) GetStrategyNewsListByPage_MultiMedia_ByDate_All(date *time.Time) ([]*expertnews.ExpertNews_MultiMedia_List, error) {
	var allList []*expertnews.ExpertNews_MultiMedia_List
	cacheKey := fmt.Sprintf("%s:GetStrategyNewsListByPage_MultiMedia_ByDate:%s", EMNET_ColumnInfo_PreCacheKey, date.Format("2006-01-02"))
	err := service.RedisCache.GetJsonObj(cacheKey, &allList)
	if err != nil || len(allList) == 0 {
		//读取数据库
		allList, err = service.columnInfoRepo.GetStrategyNewsListByPage_MultiMedia_ByDate(*date)
		if err != nil || len(allList) == 0 {
			return allList, err
		}
		service.RedisCache.Set(cacheKey, _json.GetJsonString(allList), 60*20)
	}
	return allList, nil
}

// GetStrategyNewsListByPage_MultiMedia_ByDate_Page 分页获取指定日期的实战培训课程
func (service *ColumnInfoService) GetStrategyNewsListByPage_MultiMedia_ByDate_Page(currpage int64, pageSize int64, date *time.Time) ([]*expertnews.ExpertNews_MultiMedia_List, int, error) {
	allList, err := service.GetStrategyNewsListByPage_MultiMedia_ByDate_All(date)
	if err != nil {
		return allList, 0, err
	}
	//分页计算
	total := len(allList)
	if total == 0 {
		return nil, total, nil
	}
	pageList, total := service.GetPage(allList, int(pageSize), int(currpage))
	return pageList, total, nil
}
func (service *ColumnInfoService) GetPage(allList []*expertnews.ExpertNews_MultiMedia_List, pageSize int, currPage int) ([]*expertnews.ExpertNews_MultiMedia_List, int) {
	var retList []*expertnews.ExpertNews_MultiMedia_List
	//分页
	totalCount := len(allList)

	if currPage <= 0 {
		currPage = 1
	}
	if pageSize <= 0 || totalCount <= 0 {
		return nil, totalCount
	}

	if pageSize > totalCount {
		pageSize = totalCount
	}

	if currPage > (totalCount/pageSize)+1 {
		return nil, totalCount
	}

	beginnum := (currPage - 1) * pageSize
	endnum := pageSize

	if currPage == 1 {
		retList = allList[:endnum]
	} else {
		endnum = pageSize * currPage
		//数量不够，endnum=开始num+剩余数量
		if (totalCount - beginnum) < pageSize {
			endnum = beginnum + (totalCount - beginnum)
		}

		if beginnum >= totalCount {
			return nil, totalCount
		}
		retList = allList[beginnum:endnum]
	}
	return retList, totalCount
}

// GetClSZList_1Month 分页获取最近一个月实战培训课程
func (service *ColumnInfoService) GetStrategyNewsListByPage_MultiMedia1Month(currpage int64, pageSize int64) ([]*expertnews.ExpertNews_MultiMedia_List, int, error) {
	var retList []*expertnews.ExpertNews_MultiMedia_List
	var newsList []*expertnews.ExpertNews_MultiMedia_List

	//获取最近一个月实战培训课程数据
	redisKey := RedisKey_1MonthMultiMediaList
	err := service.RedisCache.GetJsonObj(redisKey, &newsList)

	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "GetStrategyNewsListByPage_MultiMedia1Month 获取最近一个月实战培训课程列表失败")
		return nil, 0, err
	}

	//分页
	totalCount := int64(len(newsList))
	intTotalCount := len(newsList)

	if currpage > (totalCount/pageSize)+1 {
		return nil, intTotalCount, err
	}

	if pageSize > totalCount {
		pageSize = totalCount
	}

	beginnum := (currpage - 1) * pageSize
	endnum := pageSize
	if endnum > int64(len(newsList)) {
		endnum = int64(len(newsList))
	}

	if currpage == 1 {
		retList = newsList[:endnum]
	} else {
		endnum = pageSize * currpage
		//数量不够，endnum=开始num+剩余数量
		if (totalCount - beginnum) < pageSize {
			endnum = beginnum + (totalCount - beginnum)
		}
		retList = newsList[beginnum:endnum]
	}

	return retList, intTotalCount, err
}

func init() {
	protected.RegisterServiceLoader(columnInfoServiceServiceName, columnInfoServiceLoader)
}

func columnInfoServiceLoader() {
	shareColumnInfoRepo = strategyservice_repo.NewColumnInfoRepository(protected.DefaultConfig)
	shareColumnInfoLogger = dotlog.GetLogger(columnInfoServiceServiceName)
}

// InsertColumnInfo 插入一条新的栏目
func (srv *ColumnInfoService) InsertColumnInfo(info strategyservice_model.ColumnInfo) (id int, err error) {
	id, err = srv.columnInfoRepo.InsertColumnInfo(info)
	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "InsertColumnInfo 插入一条新的栏目 异常 info=%s", _json.GetJsonString(info))
	}
	return id, err
}

// UpdateColumnInfo 更新栏目
func (srv *ColumnInfoService) UpdateColumnInfo(columnName string, columnDesc string, id int) (err error) {
	err = srv.columnInfoRepo.UpdateColumnInfo(columnName, columnDesc, id)
	if err != nil {
		shareColumnInfoLogger.ErrorFormat(err, "UpdateColumnInfo 更新栏目 异常 columnName=%s columnDesc=%s id=%d", columnName, columnDesc, id)
	}
	return err
}
