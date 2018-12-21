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
)

type SeedUserInfoService struct {
	service.BaseService
	seedUserInfoRepo *userhome_repo.SeedUserInfoRepository
}

var (
	shareSeedUserInfoRepo   *userhome_repo.SeedUserInfoRepository
	shareSeedUserInfoLogger dotlog.Logger
)

const (
	UserHomeSeedUserInfoServiceName                   = "UserHomeSeedUserInfoService"
	EMNET_SeedUserInfo_CachePreKey                    = "EMNET:UserHome_SeedUserInfoBLL:"
	EMNET_SeedUserInfo_GetSeedUserInfoByCid_CachePKey = EMNET_SeedUserInfo_CachePreKey + "GetSeedUserInfoByCid:"

	EMNET_SeedUserInfo_CacheSeconds = 30 * 60
)

func NewSeedUserInfoService() *SeedUserInfoService {
	srv := &SeedUserInfoService{
		seedUserInfoRepo: shareSeedUserInfoRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeSeedUserInfoServiceName, userHomeSeedUserInfoServiceLoader)
}

func userHomeSeedUserInfoServiceLoader() {
	shareSeedUserInfoRepo = userhome_repo.NewSeedUserInfoRepository(protected.DefaultConfig)
	shareSeedUserInfoLogger = dotlog.GetLogger(UserHomeSeedUserInfoServiceName)
}

// GetSeedUserInfoByCid 根据CID获取种子用户信息
func (srv *SeedUserInfoService) GetSeedUserInfoByCid(cid string) (seedUser *userhome_model.SeedUserInfo, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getSeedUserInfoByCidDB(cid)
	case config.ReadDB_CacheOrDB_UpdateCache:
		seedUser, err = srv.GetSeedUserInfoByCidCache(cid)
		if err == nil && seedUser == nil {
			seedUser, err = srv.refreshSeedUserInfoByCid(cid)
		}
		return seedUser, err
	case config.ReadDB_RefreshCache:
		seedUser, err = srv.refreshSeedUserInfoByCid(cid)
		return seedUser, err
	default:
		return srv.GetSeedUserInfoByCidCache(cid)
	}
}

// getSeedUserInfoByCidDB 根据CID获取种子用户信息-读取数据库
func (srv *SeedUserInfoService) getSeedUserInfoByCidDB(cid string) (seedUser *userhome_model.SeedUserInfo, err error) {
	if len(cid) == 0 {
		return nil, nil
	}
	seedUser, err = srv.seedUserInfoRepo.GetSeedUserInfoByCid(cid)
	if err != nil {
		shareSeedUserInfoLogger.ErrorFormat(err, "GetSeedUserInfoByCid 根据CID获取种子用户信息 异常 cid=%s", cid)
	}
	return seedUser, err
}

// GetSeedUserInfoByCidCache 根据CID获取种子用户信息-读取缓存
func (srv *SeedUserInfoService) GetSeedUserInfoByCidCache(cid string) (seedUser *userhome_model.SeedUserInfo, err error) {
	if len(cid) == 0 {
		return nil, nil
	}
	cacheKey := EMNET_SeedUserInfo_GetSeedUserInfoByCid_CachePKey
	json, err := srv.RedisCache.HGet(cacheKey, cid)
	if err == redis.ErrNil {
		return nil, nil
	}
	err = _json.Unmarshal(json, &seedUser)
	return seedUser, err
}

// refreshSeedUserInfoByCid 根据CID获取种子用户信息-刷新缓存
func (srv *SeedUserInfoService) refreshSeedUserInfoByCid(cid string) (seedUser *userhome_model.SeedUserInfo, err error) {
	if len(cid) == 0 {
		return nil, nil
	}
	seedUser, err = srv.getSeedUserInfoByCidDB(cid)
	if err != nil {
		shareSeedUserInfoLogger.ErrorFormat(err, "refreshSeedUserInfoByCid 根据CID获取种子用户信息-刷新缓存 异常 cid=%s", cid)
		return nil, err
	}
	if seedUser != nil {
		cacheKey := EMNET_SeedUserInfo_GetSeedUserInfoByCid_CachePKey
		err = srv.RedisCache.HSet(cacheKey, cid, _json.GetJsonString(seedUser))
	}
	return seedUser, err
}

// RefreshSeedUserInfoList 刷新所有种子用户信息
func (srv *SeedUserInfoService) RefreshSeedUserInfoList() (total int, err error) {
	seedUserList, err := srv.seedUserInfoRepo.GetSeedUserInfoList()
	if err != nil {
		shareSeedUserInfoLogger.ErrorFormat(err, "RefreshSeedUserInfoList 刷新所有种子用户信息 异常")
		return total, err
	}
	if seedUserList == nil || len(seedUserList) == 0 {
		shareSeedUserInfoLogger.DebugFormat("RefreshSeedUserInfoList 刷新所有种子用户信息-->没有种子用户数据")
		return total, nil
	}
	cacheKey := EMNET_SeedUserInfo_GetSeedUserInfoByCid_CachePKey
	for _, seedUser := range seedUserList {
		if len(seedUser.Cid) == 0 {
			continue
		}
		if seedUser.IsDeleted {
			_, err = srv.RedisCache.HDel(cacheKey, seedUser.Cid)
		} else {
			err = srv.RedisCache.HSet(cacheKey, seedUser.Cid, _json.GetJsonString(seedUser))
			if err == nil {
				total++
			}
		}
	}
	return total, err
}
