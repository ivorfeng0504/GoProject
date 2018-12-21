package myoptional

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	myoptional_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	myoptional_vm "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/crypto"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"sort"
	"strconv"
	"strings"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type StockNewsInformationService struct {
	service.BaseService
	stockNewsInformationRepo *myoptional_repo.StockNewsInformationRepository
}

var (
	shareStockNewsInformationRepo   *myoptional_repo.StockNewsInformationRepository
	shareStockNewsInformationLogger dotlog.Logger
)

const (
	stockNewsInformationServiceName                              = "StockNewsInformationService"
	EMNET_StockNewsInformation_PreCacheKey                       = "StockNewsInformationService:"
	EMNET_StockNewsInformation_GetStockNewsInfo_CacheKey         = EMNET_StockNewsInformation_PreCacheKey + "GetStockNewsInfo:"
	EMNET_StockNewsInformation_SyncStockNewsInformation_CacheKey = "SyncStockNewsInformation:"
	EMNET_StockNewsInformation_CacheSeconds                      = 60 * 6
)

func NewStockNewsInformationService() *StockNewsInformationService {
	service := &StockNewsInformationService{
		stockNewsInformationRepo: shareStockNewsInformationRepo,
	}
	service.RedisCache=_cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return service
}

func init() {
	protected.RegisterServiceLoader(stockNewsInformationServiceName, stockNewsInformationServiceLoader)
}

func stockNewsInformationServiceLoader() {
	shareStockNewsInformationRepo = myoptional_repo.NewStockNewsInformationRepository(protected.DefaultConfig)
	shareStockNewsInformationLogger = dotlog.GetLogger(stockNewsInformationServiceName)
}

// GetStockNewsInfo 获取股票相关资讯
func (srv *StockNewsInformationService) GetStockNewsInfo(startTime string, stockList []string, top int) (newsList []*myoptional_model.StockNewsInformation, err error) {
	if stockList == nil || len(stockList) == 0 {
		return nil, nil
	}
	switch config.CurrentConfig.ReadDB_MyOptional {
	case config.ReadDB_JustDB:
		return srv.getStockNewsInfoDB(startTime, stockList, top)
	case config.ReadDB_CacheOrDB_UpdateCache:
		newsList, err = srv.getStockNewsInfoCache(startTime, stockList, top)
		if err == nil && newsList == nil {
			newsList, err = srv.refreshStockNewsInfo(startTime, stockList, top)
		}
		return newsList, err
	case config.ReadDB_RefreshCache:
		newsList, err = srv.refreshStockNewsInfo(startTime, stockList, top)
		return newsList, err
	default:
		return srv.getStockNewsInfoCache(startTime, stockList, top)
	}
}

// getStockNewsInfoDB 获取股票相关资讯-读取数据库
func (srv *StockNewsInformationService) getStockNewsInfoDB(startTime string, stockList []string, top int) (newsList []*myoptional_model.StockNewsInformation, err error) {
	if stockList == nil || len(stockList) == 0 {
		return nil, nil
	}
	newsList, err = srv.stockNewsInformationRepo.GetStockNewsInfo(startTime, stockList, top)
	if err != nil {
		shareStockNewsInformationLogger.ErrorFormat(err, "GetStockNewsInfo 获取股票相关资讯 异常 startTime=%s stockList=%s top=%d", startTime, _json.GetJsonString(stockList), top)
	}
	return newsList, err
}

// getStockNewsInfoCache 获取股票相关资讯-读取缓存
func (srv *StockNewsInformationService) getStockNewsInfoCache(startTime string, stockList []string, top int) (newsList []*myoptional_model.StockNewsInformation, err error) {
	cacheKey := EMNET_StockNewsInformation_GetStockNewsInfo_CacheKey + startTime + ":" + strconv.Itoa(top) + ":" + calcHash(stockList)
	err = srv.RedisCache.GetJsonObj(cacheKey, &newsList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return newsList, err
}

// refreshStockNewsInfo 获取股票相关资讯-刷新股票
func (srv *StockNewsInformationService) refreshStockNewsInfo(startTime string, stockList []string, top int) (newsList []*myoptional_model.StockNewsInformation, err error) {
	cacheKey := EMNET_StockNewsInformation_GetStockNewsInfo_CacheKey + startTime + ":" + strconv.Itoa(top) + ":" + calcHash(stockList)
	newsList, err = srv.getStockNewsInfoDB(startTime, stockList, top)
	if err != nil {
		return newsList, err
	}
	if newsList != nil {
		srv.RedisCache.Set(cacheKey, _json.GetJsonString(newsList), EMNET_StockNewsInformation_CacheSeconds)
	}
	return newsList, err
}

// InsertStockNewsInfo 插入一个新的股票资讯信息
func (srv *StockNewsInformationService) InsertStockNewsInfo(model myoptional_model.StockNewsInformation) (err error) {
	err = srv.stockNewsInformationRepo.InsertStockNewsInfo(model)
	if err != nil {
		shareStockNewsInformationLogger.ErrorFormat(err, "InsertStockNewsInfo 插入一个新的股票资讯信息 异常 model=%s", _json.GetJsonString(model))
	}
	return err
}

// GetMaxNewsInformationId 获取关系表中最大的NewsInformationId
func (srv *StockNewsInformationService) GetMaxNewsInformationId() (newsInformationId int64, err error) {
	newsInformationId, err = srv.stockNewsInformationRepo.GetMaxNewsInformationId()
	if err != nil {
		shareStockNewsInformationLogger.ErrorFormat(err, "GetMaxNewsInformationId 获取关系表中最大的NewsInformationId 异常")
	}
	return newsInformationId, err
}

// SyncStockNewsInformation 同步最新的资讯信息到关系表及相关缓存
func (srv *StockNewsInformationService) SyncStockNewsInformation(top int) (err error) {
	//获取当前已同步的最大资讯Id
	maxNewsInformationId, err := srv.GetMaxNewsInformationId()
	if err != nil {
		shareStockNewsInformationLogger.ErrorFormat(err, "SyncStockNewsInformation 同步最新的资讯信息到关系表 GetMaxNewsInformationId->异常")
		return err
	}

	//查询出新增的资讯集合
	newsInfoSrv := expertnews.NewNewsInformationService()
	newsInfoList, err := newsInfoSrv.GetStockNewsInfoList(maxNewsInformationId)
	if err != nil {
		shareStockNewsInformationLogger.ErrorFormat(err, "SyncStockNewsInformation 同步最新的资讯信息到关系表 GetNewsInfoList->异常 maxNewsInformationId=%d", maxNewsInformationId)
		return err
	}
	if newsInfoList == nil || len(newsInfoList) == 0 {
		shareStockNewsInformationLogger.Debug("SyncStockNewsInformation 同步最新的资讯信息到关系表 当前没有新增的资讯需要同步")
	} else {
		shareStockNewsInformationLogger.Debug("SyncStockNewsInformation 同步最新的资讯信息到关系表 当开始同步新增的资讯 start")
		//遍历集合，根据股票拆分出数据存储到StockNewsInformation中
		for _, newsInfo := range newsInfoList {
			stockListSrc := strings.Split(newsInfo.SecurityCode, ",")
			if stockListSrc == nil || len(stockListSrc) == 0 {
				continue
			}
			for _, stockCodeSrc := range stockListSrc {
				stockCodeInfo := strings.Split(stockCodeSrc, "|")
				if stockCodeInfo == nil || len(stockCodeInfo) != 2 {
					continue
				}
				stockType := stockCodeInfo[0]
				stockCode := stockCodeInfo[1]
				model := myoptional_model.StockNewsInformation{
					SecurityCode:      stockCode,
					SecurityCodeType:  stockType,
					NewsInformationId: newsInfo.NewsInformationId,
					PublishTime:       newsInfo.PublishTime,
				}
				err = srv.InsertStockNewsInfo(model)
				if err != nil {
					shareStockNewsInformationLogger.ErrorFormat(err, "SyncStockNewsInformation 同步最新的资讯信息到关系表 InsertStockNewsInfo 插入关系表异常 model=%s", _json.GetJsonString(model))
					return err
				}
			}
		}
		shareStockNewsInformationLogger.Debug("SyncStockNewsInformation 同步最新的资讯信息到关系表 当开始同步新增的资讯 end")
	}

	//查询出所有有资讯的股票
	stockCodeList, err := srv.GetStockCodeList()
	if err != nil {
		shareStockNewsInformationLogger.ErrorFormat(err, "SyncStockNewsInformation 同步最新的资讯信息到关系表 GetStockCodeList->异常")
		return err
	}
	if stockCodeList == nil || len(stockCodeList) == 0 {
		shareStockNewsInformationLogger.WarnFormat("SyncStockNewsInformation 同步最新的资讯信息到关系表 没有需要同步的股票")
		return nil
	}
	//遍历所有股票，将每个股票代码最新的N条有效资讯存储到缓存中
	for _, stockCode := range stockCodeList {
		stockNewsList, err := srv.getStockNewsInfoDB("2018-01-01", []string{stockCode}, top)
		if err != nil {
			shareStockNewsInformationLogger.ErrorFormat(err, "SyncStockNewsInformation 同步最新的资讯信息到关系表 遍历所有股票，将每个股票代码最新的N条有效资讯存储到缓存中->异常 stockCode=%s", stockCode)
			return err
		}
		if stockNewsList != nil && len(stockNewsList) > 0 {
			cacheKey := EMNET_StockNewsInformation_SyncStockNewsInformation_CacheKey + strconv.Itoa(top)
			err = srv.RedisCache.HSet(cacheKey, stockCode, _json.GetJsonString(stockNewsList))
		}
	}
	return err
}

// GetStockCodeList 获取包含资讯的所有股票代码
func (srv *StockNewsInformationService) GetStockCodeList() (stockCodeList []string, err error) {
	stockCodeList, err = srv.stockNewsInformationRepo.GetStockCodeList()
	if err != nil {
		shareStockNewsInformationLogger.ErrorFormat(err, "GetStockCodeList 获取包含资讯的所有股票代码 异常")
	}
	return stockCodeList, err
}

// GetStockNewsInfoByHSet 根据股票代码获取相关资讯（直接读HSet缓存） perTop每只股票的最大数目 totalTop 数据集结果总数
func (srv *StockNewsInformationService) GetStockNewsInfoByHSet(stockList []string, perTop int, totalTop int) (newsList []*myoptional_model.StockNewsInformation, err error) {
	var newsListVM myoptional_vm.StockNewsInformationList
	if stockList == nil || len(stockList) == 0 {
		return nil, nil
	}
	cacheKey := EMNET_StockNewsInformation_SyncStockNewsInformation_CacheKey + strconv.Itoa(perTop)
	var fields []interface{}
	stockDistinctDict := make(map[string]bool)
	for _, stock := range stockList {
		if len(stock) == 0 {
			continue
		}
		//如果是7代码去除第一位
		if len(stock) > 6 {
			stock = stock[1:]
		}
		_, isExist := stockDistinctDict[stock]
		if isExist == false {
			fields = append(fields, stock)
			stockDistinctDict[stock] = true
		}
	}

	if len(fields) == 0 {
		return nil, nil
	}

	//如果大于1000个股票 则取前1000个
	if len(fields) > 1000 {
		fields = fields[:1000]
	}

	result, err := srv.RedisCache.HMGet(cacheKey, fields...)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		return nil, nil
	}
	//循环所有结果
	distinctDict := make(map[string]bool)
	for _, jsonStr := range result {
		var tempList []*myoptional_model.StockNewsInformation
		jsonErr := _json.Unmarshal(jsonStr, &tempList)
		if jsonErr == nil && tempList != nil && len(tempList) > 0 {
			//循环每一条资讯 去重
			for _, newsInfo := range tempList {
				_, isExist := distinctDict[newsInfo.Id]
				if isExist == false {
					newsListVM = append(newsListVM, newsInfo)
					distinctDict[newsInfo.Id] = true
				}
			}

		}
	}
	if err != nil {
		return nil, err
	}
	//排序
	sort.Sort(newsListVM)
	//获取指定数量的资讯
	if len(newsListVM) > totalTop {
		newsListVM = newsListVM[:totalTop]
	}
	newsList = newsListVM
	return newsList, err
}

func calcHash(stockList []string) string {
	hash := _crypto.MD5(_json.GetJsonString(stockList))[:8]
	return hash
}
