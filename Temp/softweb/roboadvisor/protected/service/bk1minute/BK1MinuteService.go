package bk1minute

import (
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/const"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/protected/model/bk1minute"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

type BK1MinutesService struct {
	service.BaseService
}

var (
	//shareBK1MinutesLogger 共享的Logger实例
	shareBK1MinutesLogger dotlog.Logger
)

const (
	RedisKey_BK1Minutes = _const.RedisKey_NewsPre + "BK1Minutes:GetBK1MinutesInfo"

	BK1MinutesServiceName = "BK1MinutesServiceLogger"
)

func NewBK1MinutesService() *BK1MinutesService {
	bk1minutesService := &BK1MinutesService{
	}
	bk1minutesService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return bk1minutesService
}


// GetBK1MinutesInfoByBKCode 根据指定bkcode获取板块信息
func (service *BK1MinutesService) GetBK1MinutesInfoByBKCode(bkcode string) (*bk1minute.BK1MinutesInfo,error) {
	if bkcode == "" {
		return nil, errors.New("must set bkcode")
	}

	//从redis缓存获取banner
	rediskey := RedisKey_BK1Minutes
	bkinfo := new(bk1minute.BK1MinutesInfo)
	result, err := service.RedisCache.HGet(rediskey, bkcode)
	err = _json.Unmarshal(result, bkinfo)
	if err != nil {
		shareBK1MinutesLogger.ErrorFormat(err, "获取板块1分钟数据失败")
		return nil, err
	}

	return bkinfo, err
}


func init() {
	protected.RegisterServiceLoader(BK1MinutesServiceName, bk1MinutesServiceLoader)
}

func bk1MinutesServiceLoader() {
	shareBK1MinutesLogger = dotlog.GetLogger(BK1MinutesServiceName)
}