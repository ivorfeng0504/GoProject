package service

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"time"
)

type ReceiveStockActivityService struct {
	service.BaseService
	receiveStockActivityRepo *userhome_repo.ReceiveStockActivityRepository
}

var (
	shareReceiveStockActivityRepo   *userhome_repo.ReceiveStockActivityRepository
	shareReceiveStockActivityLogger dotlog.Logger
)

const (
	ReceiveStockActivityServiceName                       = "ReceiveStockActivityService"
	EMNET_ReceiveStockActivity_CachePreKey                = "EMNET:ReceiveStockActivityService:"
	EMNET_ReceiveStockActivity_GetNewstStockPool_CacheKey = EMNET_ReceiveStockActivity_CachePreKey + "GetNewstStockPool:"
	EMNET_ReceiveStockActivity_CacheSeconds               = 30 * 60
)

func NewReceiveStockActivityService() *ReceiveStockActivityService {
	srv := &ReceiveStockActivityService{
		receiveStockActivityRepo: shareReceiveStockActivityRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(ReceiveStockActivityServiceName, receiveStockActivityServiceLoader)
}

func receiveStockActivityServiceLoader() {
	shareReceiveStockActivityRepo = userhome_repo.NewReceiveStockActivityRepository(protected.DefaultConfig)
	shareReceiveStockActivityLogger = dotlog.GetLogger(ReceiveStockActivityServiceName)
}

// InsertStockPoolByReportUrl 新增活动
func (srv *ReceiveStockActivityService) InsertStockPoolByReportUrl(reportUrl string, issueNumber string) (id int64, err error) {
	id, err = srv.receiveStockActivityRepo.InsertStockPoolByReportUrl(reportUrl, issueNumber)
	if err != nil {
		shareReceiveStockActivityLogger.ErrorFormat(err, "InsertStockPoolByReportUrl 新增活动异常 reportUrl= %s  issueNumber=%s", reportUrl, issueNumber)
	} else {
		//刷新缓存
		go srv.refreshNewstStockPool()
	}
	return id, err
}

// UpdateStockPoolByReportUrl 更新活动
func (srv *ReceiveStockActivityService) UpdateStockPoolByReportUrl(receiveStockActivityId int64, reportUrl string) (err error) {
	err = srv.receiveStockActivityRepo.UpdateStockPoolByReportUrl(receiveStockActivityId, reportUrl)
	if err != nil {
		shareReceiveStockActivityLogger.ErrorFormat(err, "UpdateStockPoolByReportUrl 更新活动 receiveStockActivityId= %d  reportUrl=%s", receiveStockActivityId, reportUrl)
	} else {
		//刷新缓存
		go srv.refreshNewstStockPool()
	}
	return err
}

// InsertStockPool 新增一个股票池
func (srv *ReceiveStockActivityService) InsertStockPool(stocks []*userhome_model.StockInfo, issueNumber string) (id int64, err error) {
	id, err = srv.receiveStockActivityRepo.InsertStockPool(stocks, issueNumber)
	if err != nil {
		shareReceiveStockActivityLogger.ErrorFormat(err, "新增股票池异常 股票信息为 %s", _json.GetJsonString(stocks))
	} else {
		//刷新缓存
		go srv.refreshNewstStockPool()
	}
	return id, err
}

// GetNewstStockPool 获取最新的股票池
func (srv *ReceiveStockActivityService) GetNewstStockPool() (stock *userhome_model.ReceiveStockActivity, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getNewstStockPoolDB()
	case config.ReadDB_CacheOrDB_UpdateCache:
		stock, err = srv.getNewstStockPoolCache()
		if err == nil && stock == nil {
			stock, err = srv.refreshNewstStockPool()
		}
		return stock, err
	case config.ReadDB_RefreshCache:
		stock, err = srv.refreshNewstStockPool()
		return stock, err
	default:
		return srv.getNewstStockPoolCache()
	}
}

// getNewstStockPoolDB 获取最新的股票池-读取数据库
func (srv *ReceiveStockActivityService) getNewstStockPoolDB() (stock *userhome_model.ReceiveStockActivity, err error) {
	stock, err = srv.receiveStockActivityRepo.GetNewstStockPool()
	if err != nil {
		shareReceiveStockActivityLogger.ErrorFormat(err, "获取最新的股票池-读取数据库 异常")
	}
	return stock, err
}

// getNewstStockPoolCache 获取最新的股票池-读取缓存
func (srv *ReceiveStockActivityService) getNewstStockPoolCache() (stock *userhome_model.ReceiveStockActivity, err error) {
	cacheKey := EMNET_ReceiveStockActivity_GetNewstStockPool_CacheKey + time.Now().Format("20060102")
	err = srv.RedisCache.GetJsonObj(cacheKey, &stock)
	if err == redis.ErrNil {
		return nil, nil
	}
	return stock, err
}

// refreshNewstStockPool 获取最新的股票池-刷新缓存
func (srv *ReceiveStockActivityService) refreshNewstStockPool() (stock *userhome_model.ReceiveStockActivity, err error) {
	cacheKey := EMNET_ReceiveStockActivity_GetNewstStockPool_CacheKey + time.Now().Format("20060102")
	stock, err = srv.getNewstStockPoolDB()
	if err != nil {
		return nil, err
	}
	if stock != nil {
		srv.RedisCache.Set(cacheKey, _json.GetJsonString(stock), EMNET_ReceiveStockActivity_CacheSeconds)
	}
	return stock, err
}
