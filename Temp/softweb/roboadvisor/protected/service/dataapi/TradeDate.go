package dataapi

import (
	"encoding/json"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	TradeDateCacheKey = "dataapi.TradeDateCacheKey."
	IsTrade           = "1"
	NotTrade          = "0"
)

// IsTradeDay 是否是交易日
func IsTradeDay(date time.Time) (bool, error) {
	cacheProvider := _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	cacheKey := TradeDateCacheKey + date.Format("2006-01-02")
	isTradeStr := ""
	err := cacheProvider.GetJsonObj(cacheKey, &isTradeStr)
	if err == nil {
		return isTradeStr == IsTrade, nil
	}
	apiUrl := `http://dataapi.emoney.cn/platformapi/indicator/execcondition?token=C6E7A4620478B6227B5228C84CC06E9FFC457A3C&condition=vq_dbo_tb_pub_10210(*,"%20f003v_10210=%27012001%27%20AND%20DATEDIFF(dd,[f002d_10210],%27` +
		date.Format("2006-01-02") +
		`%27)=0",null,1)`
	resp, err := http.Get(apiUrl)
	defer func() { resp.Body.Close() }()
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("请求异常")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	//这个接口是utf8带BOM头格式，需要手动去除BOM
	if len(data) > 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	respData := new(DataApiResponse)
	err = json.Unmarshal(data, respData)
	if err != nil {
		return false, err
	}
	if respData.Success && len(respData.Data) > 0 && respData.Data[0].Success && len(respData.Data[0].Result) > 2 {
		cacheProvider.SetJsonObj(cacheKey, IsTrade)
		return true, nil
	}
	cacheProvider.SetJsonObj(cacheKey, NotTrade)
	return false, nil
}

// IsTradeDayToDay 今天是否是交易日
func IsTradeDayToDay() (bool, error) {
	now := time.Now()
	return IsTradeDay(now)
}

// IsTradeTimeNow 当前是否为交易时间
func IsTradeTimeNow() (bool, error) {
	now := time.Now()
	startTime, err := time.Parse("2006-1-2 15:04 MST", now.Format("2006-1-2")+" 08:30 CST")
	if err != nil {
		return false, err
	}
	endTime, err := time.Parse("2006-1-2 15:04 MST", now.Format("2006-1-2")+" 15:30 CST")
	if err != nil {
		return false, err
	}
	if now.Before(startTime) || now.After(endTime) {
		return false, nil
	}
	return IsTradeDayToDay()
}

// GetLastTradeWeekDay 获取指定日期所在周的最后一个交易日
func GetLastTradeWeekDay(date time.Time) (tradeWeek time.Time, err error) {
	isTradeDay := false
	currentWeek := int(date.Weekday())
	if currentWeek == 0 {
		currentWeek = 7
	}
	date = date.AddDate(0, 0, 5-currentWeek)
	//轮询周五到周一 判断哪个是最后的交易日
	for i := 0; i < 5; i++ {
		isTradeDay, err = IsTradeDay(date)
		if err != nil {
			return tradeWeek, err
		}
		//如果取到交易日则直接返回
		if isTradeDay {
			break
		} else {
			//未取到交易日则再次往前推
			date = date.AddDate(0, 0, -1)
		}
	}
	if isTradeDay {
		return date, nil
	} else {
		return date, ErrNoTradeDay
	}
}

// GetNextTradeDay 获取指定日期的下一个交易日
func GetNextTradeDay(date time.Time) (nextTradeDay time.Time, err error) {
	tryCount := 14
	for tryCount > 0 {
		date = date.AddDate(0, 0, 1)
		yes, err := IsTradeDay(date)
		if err != nil {
			return nextTradeDay, err
		}
		if yes {
			nextTradeDay = date
			return nextTradeDay, nil
		}
		tryCount--
	}
	err = errors.New("多次尝试查询下一个交易日，为查询到下一个交易日数据")
	return nextTradeDay, err
}

var ErrNoTradeDay = errors.New("dataapi: 交易日不存在")
