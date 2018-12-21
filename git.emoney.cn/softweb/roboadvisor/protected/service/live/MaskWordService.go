package live

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	liverepo "git.emoney.cn/softweb/roboadvisor/protected/repository/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strings"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type MaskWordService struct {
	service.BaseService
	maskWordRepo *liverepo.MaskWordRepository
}

var (
	// shareMaskWordRepo 共享的仓储
	shareMaskWordRepo *liverepo.MaskWordRepository

	// shareMaskWordLogger 共享的Logger实例
	shareMaskWordLogger dotlog.Logger
)

const (
	MaskWordServiceName                  = "MaskWordService"
	EMNET_MaskWord_CachePreKey           = "EMoney:Live:MaskWord"
	EMNET_MaskWord_MaskWordList_CacheKey = EMNET_MaskWord_CachePreKey + ":MaskWordList"
)

func NewMaskWordService() *MaskWordService {
	service := &MaskWordService{
		maskWordRepo: shareMaskWordRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return service
}

// GetMaskWordList 获取所有屏蔽字
func (service *MaskWordService) ProcessMaskWord(content string) (string, error) {
	maskWordList, err := service.GetMaskWordList()
	if err != nil {
		return "", err
	}
	if maskWordList == nil || len(maskWordList) == 0 {
		return content, nil
	}
	for _, mask := range maskWordList {
		if mask.MaskName != "" {
			content = strings.Replace(content, mask.MaskName, "***", -1)
		}
	}
	return content, nil
}

// GetMaskWordList 获取所有屏蔽字
func (service *MaskWordService) GetMaskWordList() (maskWordList []*livemodel.MaskWord, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getMaskWordListDB()
	default:
		return service.getMaskWordListCache()
	}
}

// getMaskWordListDB 获取所有屏蔽字
func (service *MaskWordService) getMaskWordListDB() (maskWordList []*livemodel.MaskWord, err error) {
	maskWordList, err = service.maskWordRepo.GetMaskWordList()
	return maskWordList, err
}

// getMaskWordListCache 获取所有屏蔽字
func (service *MaskWordService) getMaskWordListCache() (maskWordList []*livemodel.MaskWord, err error) {
	cacheKey := EMNET_MaskWord_MaskWordList_CacheKey
	err = service.RedisCache.GetJsonObj(cacheKey, &maskWordList)
	if err == redis.ErrNil {
		return maskWordList, nil
	}
	if err != nil {
		shareMaskWordLogger.Error(err, "读取缓存-获取所有屏蔽字 getMaskWordListCache 异常")
	}
	return maskWordList, err
}

func init() {
	protected.RegisterServiceLoader(MaskWordServiceName, maskWordServiceLoader)
}

func maskWordServiceLoader() {
	shareMaskWordRepo = liverepo.NewMaskWordRepository(protected.DefaultConfig)
	shareMaskWordLogger = dotlog.GetLogger(MaskWordServiceName)
}
