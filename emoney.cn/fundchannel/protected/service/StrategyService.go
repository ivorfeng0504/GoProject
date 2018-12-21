package service

import (
	"emoney.cn/fundchannel/protected"
	"emoney.cn/fundchannel/util/cache"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"emoney.cn/fundchannel/util/http"
	"emoney.cn/fundchannel/config"
	"strconv"
)

type StrategyService struct {
	BaseService
}

var (
	fundChannelStrategyLogger dotlog.Logger
)

const (
	//开头标识需要与.NET项目中的KEY保持一致
	StrategyListKey     = FundChannelKey + "Strategy.StrategyList"
	StrategyServiceName = FundChannelKey + "StrategyService"
	StrategyInfoHkey    = FundChannelKey + "Strategy.Info.code_"
	//StrategyInfoFkey    = FundChannelKey + "Strategy.Info.date_"
	FundTimeLineKey = FundChannelKey + "Fund.TimeLine.code_"
)

func NewStrategyService() *StrategyService {
	service := &StrategyService{}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return service
}

func init() {
	protected.RegisterServiceLoader(StrategyServiceName, strategyServiceLoader)
}

func strategyServiceLoader() {
	fundChannelStrategyLogger = dotlog.GetLogger("StrategyService")
}

func (service *StrategyService) GetStrategyList() (strategyListJson string, err error) {
	strategyListJson, err = service.RedisCache.GetString(StrategyListKey)
	if err == redis.ErrNil {
		return strategyListJson, nil
	}
	if err != nil {
		fundChannelStrategyLogger.Error(err, " GetStrategyList 获取基金频道配置信息 异常")
	}
	return strategyListJson, err
}

func (service *StrategyService) GetStrategyInfoByCode(code string) (strategyInfoJson string, err error) {
	hKey := StrategyInfoHkey + code
	strategyInfo, err := service.RedisCache.HGetAll(hKey)
	strategyInfoJson = "["
	for _, v := range strategyInfo {
		if len(v) > 0 {
			strategyInfoJson += v + ","
		}
	}
	s := strategyInfoJson[len(strategyInfoJson)-1 : len(strategyInfoJson)]
	if s == "," {
		strategyInfoJson = strategyInfoJson[:len(strategyInfoJson)-1]
	}

	strategyInfoJson += "]"

	if err == redis.ErrNil {
		return strategyInfoJson, nil
	}
	if err != nil {
		fundChannelStrategyLogger.Error(err, " GetStrategyList 获取基金频道配置信息 异常")
	}
	return strategyInfoJson, err
}

func (service *StrategyService) GetFundTimeLineByCode(code string) (jsonData string, err error) {
	jsonData, err = service.RedisCache.GetString(FundTimeLineKey + code)
	if err == redis.ErrNil {
		return jsonData, nil
	}
	if err != nil {
		fundChannelStrategyLogger.Error(err, " GetFundTimeLine 获取基金 实时估值数据 异常")
	}
	return jsonData, err
}

func (service *StrategyService) GetEncryptMobile(userId string) (jsonData string, err error) {
	url := config.CurrentConfig.Boundgroupqrylogin + "&uid=" + userId
	jsonData, _, _, err = _http.HttpGet(url)
	if err != nil {
		fundChannelStrategyLogger.Error(err, " QueryUserRiskInfo 获取用户加密手机号 异常")
	}
	return jsonData, err
}

func (service *StrategyService) QueryUserRiskInfo(mobileNumber string) (jsonData string, err error) {

	postBody := "{\"Mobile\":\"" + mobileNumber + "\"}"
	fundChannelStrategyLogger.Debug(" QueryUserRiskInfo postBody=" + postBody)
	fundChannelStrategyLogger.Debug(" QueryUserRiskInfo url=" + config.CurrentConfig.QueryUserRiskInfo)
	jsonData, _, _, err = _http.HttpPost(config.CurrentConfig.QueryUserRiskInfo, postBody, "application/json")
	fundChannelStrategyLogger.Debug(" QueryUserRiskInfo jsonData=" + jsonData)

	if err != nil {
		fundChannelStrategyLogger.Error(err, " QueryUserRiskInfo 获取用户等级数据 异常")
	}
	return jsonData, err
}

// GetUserInfoByUID 根据UID获取用户信息from redis
func (service *StrategyService) GetUserInfoByUID(UID int64) (userJson string, err error) {
	cacheKey := "UserInfoServiceLogger:GetUserInfoByUID" + strconv.FormatInt(UID, 10)

	//get from redis
	userJson, err = service.RedisCache.GetString(cacheKey)

	if err != redis.ErrNil {
		fundChannelStrategyLogger.Error(err, " GetUserInfoByUID 根据UID获取用户信息 异常")
	}

	return userJson, nil
}
