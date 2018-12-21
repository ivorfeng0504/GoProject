package service

import (
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/protected/repository/panwai"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	"git.emoney.cn/softweb/roboadvisor/const"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type BannerInfoService struct {
	service.BaseService
	bannerRepo *repository.BannerInfoRepository
}

var (
	//shareBannerRepo 共享的BannerInfoRepository实例
	shareBannerRepo *repository.BannerInfoRepository

	//shareBannerLogger 共享的Logger实例
	shareBannerLogger dotlog.Logger
)

const (
	RedisKey_GetBannerInfoByID = _const.RedisKey_NewsPre + "NewsInfo:GetBannerInfoByID:"
	RedisKey_GetBannerList = _const.RedisKey_NewsPre + "BannerInfo:GetBannerList"

	bannerInfoServiceName = "BannerInfoServiceLogger"
)

func NewBannerInfoService() *BannerInfoService {
	bannerService := &BannerInfoService{
		bannerRepo: shareBannerRepo,
	}
	bannerService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return bannerService
}


// GetBannerInfoList 获取banner列表
func (service *BannerInfoService) GetBannerInfoList() ([]*model.BannerInfo,error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.GetBannerInfoListDB()
	default:
		return service.GetBannerInfoListCache()
	}
}
func (service *BannerInfoService) GetBannerInfoListDB() ([]*model.BannerInfo,error) {
	results, err := service.bannerRepo.GetBannerInfoList()
	if err == nil {
		if len(results) <= 0 {
			results = nil
			err = errors.New("not exists this new info")
		}
	}
	return results, err
}
func (service *BannerInfoService) GetBannerInfoListCache() ([]*model.BannerInfo,error) {
	var results []*model.BannerInfo
	var err error
	redisKey := RedisKey_GetBannerList
	//get from redis
	err = service.RedisCache.GetJsonObj(redisKey, &results)
	if err == nil {
		return results, nil
	}
	return results, err
}


// GetBannerInfoByID 根据指定bannerID获取BannerInfo
func (service *BannerInfoService) GetBannerInfoByID(bannerID int) (banner *model.BannerInfo,err error) {
	if bannerID <= 0 {
		return nil, errors.New("must set bannerID")
	}

	//从redis缓存获取banner
	rediskey := RedisKey_GetBannerInfoByID + string(bannerID)

	bannerResult := new(model.BannerInfo)

	jsonerr := service.RedisCache.GetJsonObj(rediskey, bannerResult)
	if jsonerr == nil {
		return bannerResult, jsonerr
	}

	return banner, err
}


func init() {
	protected.RegisterServiceLoader(bannerInfoServiceName, bannerServiceLoader)
}

func bannerServiceLoader() {
	shareBannerRepo = repository.NewBannerInfoRepository(protected.DefaultConfig)
	shareBannerLogger = dotlog.GetLogger(bannerInfoServiceName)
}