package strategyapi

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/stockapi"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/cache"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

const (
	strategyApiSeviceName     = "StrategyApiSevice"
	StrategyStockPoolCacheKey = strategyApiSeviceName + "StrategyStockPool:"
)

var (
	shareStrategyApiSeviceLogger dotlog.Logger
	shareStrategyApiSeviceRedis  cache.RedisCache
	NilStockPool                 = errors.New("股池内容为空")
)

func init() {
	protected.RegisterServiceLoader(strategyApiSeviceName, func() {
		shareStrategyApiSeviceLogger = dotlog.GetLogger(strategyApiSeviceName)
		shareStrategyApiSeviceRedis = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	})
}

// GetShuangXiangPaoPool 获取双响炮股池
func GetShuangXiangPaoPool(top int) (stockList []*userhome_model.StockInfo, err error) {
	stockList, err = GetCommonPoolStockList(top, config.CurrentConfig.ShuangXiangPaoStrategyKey, ReadCacheAfter)
	if err != nil {
		shareStrategyApiSeviceLogger.ErrorFormat(err, "获取双响炮股池异常")
	}
	return stockList, err
}

// GetCommonPoolStockList 获取通用股池股票
// top 如果大于0 则取前top个股票，否则取全部
// strategyKey 策略Key
// readCacheMode 0则不读取任何缓存  1 先读取缓存，缓存不存在则读取接口 2 先读取接口，如果接口中无数据，则读取缓存 3只读取缓存
func GetCommonPoolStockList(top int, strategyKey string, readCacheMode int) (stockList []*userhome_model.StockInfo, err error) {
	count := 0
	stockPool, err := GetCommPool(strategyKey, config.CurrentConfig.AppId, readCacheMode)
	if err != nil {
		shareStrategyApiSeviceLogger.ErrorFormat(err, "获取通用炮股池股票异常 strategyKey=%s", strategyKey)
		return stockList, err
	}
	shareStrategyApiSeviceLogger.DebugFormat("GetCommonPoolStockList 获取通用炮股池股票 result=%s", _json.GetJsonString(stockPool))
	if stockPool == nil || len(stockPool) == 0 {
		return stockList, NilStockPool
	}
	stockCodeDict, err := stockapi.GetStockListDictCache()
	if err != nil {
		shareStrategyApiSeviceLogger.ErrorFormat(err, "获取通用炮股池股票异常 获取码表异常 strategyKey=%s", strategyKey)
		return nil, err
	}
	loc, _ := time.LoadLocation("Local")
	for _, pool := range stockPool {
		codeStr := strconv.Itoa(pool.Code)
		if len(codeStr) > 6 {
			codeStr = codeStr[1:]
		}
		if top > 0 && count >= top {
			break
		}
		stock := &userhome_model.StockInfo{
			StockCode: codeStr,
		}
		stock.CreateTime, _ = time.ParseInLocation("200601021504", "20"+strconv.Itoa(pool.Time), loc)

		stock.StockName = stockCodeDict[stock.StockCode]
		if len(stock.StockName) == 0 {
			shareStrategyApiSeviceLogger.WarnFormat("获取通用炮股池股票异常 存在不在码表中的股票 股票代码为%s strategyKey=%s", stock.StockCode, strategyKey)
			continue
		}
		stockList = append(stockList, stock)
		count++
	}
	return stockList, nil
}

const NoCache int = 0
const ReadCacheBefore int = 1
const ReadCacheAfter int = 2
const ReadCacheOnly int = 3

// GetCommPool 获取股池
// 如果接口中读取不到最新的数据 则从本地缓存中获取
// key 策略Key appId 应用Id
// readCacheMode 0则不读取任何缓存  1 先读取缓存，缓存不存在则读取接口 2 先读取接口，如果接口中无数据，则读取缓存  3只读取缓存
func GetCommPool(key string, appId string, readCacheMode int) (stockPool CommonPoolResponseList, err error) {
	cacheKey := StrategyStockPoolCacheKey + key
	if readCacheMode == ReadCacheBefore || readCacheMode == ReadCacheOnly {
		err = shareStrategyApiSeviceRedis.GetJsonObj(cacheKey, &stockPool)
		if err == redis.ErrNil {
			err = nil
		}
		if readCacheMode == ReadCacheOnly {
			return stockPool, err
		}
	}
	if stockPool == nil || len(stockPool) == 0 {

		apiUrl := config.CurrentConfig.CommonPoolApi
		if key[:2] == "10" {
			apiUrl = config.CurrentConfig.OnTradingPoolApi
		}
		if len(apiUrl) == 0 {
			if key[:2] == "10" {
				shareStrategyApiSeviceLogger.ErrorFormat(err, "获取股池接口地址配置不正确 configkey=OnTradingPoolApi")
			} else {
				shareStrategyApiSeviceLogger.ErrorFormat(err, "获取股池接口地址配置不正确 configkey=CommonPoolApi")
			}
			return stockPool, errors.New("获取股池接口地址配置不正确")
		}
		date := time.Now().Format("20060102")
		apiUrl = fmt.Sprintf(apiUrl, appId, key, date)
		body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
		if errReturn != nil {
			shareStrategyApiSeviceLogger.ErrorFormat(err, "GetCommPool-> 异常 URL=【%s】，body=【%s】", apiUrl, body)
			return stockPool, errReturn
		} else {
			shareStrategyApiSeviceLogger.DebugFormat("GetCommPool-> 正常 URL=【%s】，body=【%s】", apiUrl, body)
		}
		_ = contentType
		_ = intervalTime
		apiGatewayResp := contract.ApiGatewayResponse{}
		err = _json.Unmarshal(body, &apiGatewayResp)
		if err != nil {
			return stockPool, err
		}
		if apiGatewayResp.RetCode != 0 {
			return stockPool, errors.New(apiGatewayResp.RetMsg)
		}
		err = _json.Unmarshal(apiGatewayResp.Message, &stockPool)
		if err != nil {
			shareStrategyApiSeviceLogger.ErrorFormat(err, "GetCommPool->Unmarshal API数据反序列化异常 data=%s", apiGatewayResp.Message)
		}
	}
	//如果接口中读取不到最新的数据 则从本地缓存中获取
	if err == nil && len(stockPool) == 0 {
		//如果允许读取缓存 则从redis中尝试读取最近的数据
		if readCacheMode == ReadCacheAfter {
			err = shareStrategyApiSeviceRedis.GetJsonObj(cacheKey, &stockPool)
			if err == redis.ErrNil {
				return stockPool, nil
			}
		}
	} else if err == nil && len(stockPool) > 0 {
		//获取到最新数据 则更新缓存
		shareStrategyApiSeviceRedis.SetJsonObj(cacheKey, stockPool)
	} else {
		shareStrategyApiSeviceLogger.ErrorFormat(err, "GetCommPool 异常")
	}
	//sort.Sort(stockPool)
	return stockPool, err
}

type CommonPoolResponse struct {
	//时间 1809161429
	Time     int    `json:"time"`
	Code     int    `json:"code"`
	Price    int    `json:"price"`
	Strategy string `json:"strategy"`
	Max      int    `json:"max"`
}

type CommonPoolResponseList []*CommonPoolResponse

func (list CommonPoolResponseList) Len() int {
	return len(list)
}

func (list CommonPoolResponseList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list CommonPoolResponseList) Less(i, j int) bool {
	return list[j].Time < list[i].Time
}
