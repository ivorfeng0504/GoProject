package dataapi

import (
	"encoding/json"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"github.com/devfeel/dotlog"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	marketPriceServiceName = "MarketPriceService"
	StockCodeReplace       = "$StockCode$"
)

var (
	shareMarketPriceServiceLogger dotlog.Logger
)

func init() {
	protected.RegisterServiceLoader(marketPriceServiceName, func() {
		shareMarketPriceServiceLogger = dotlog.GetLogger(marketPriceServiceName)
	})
}

// GetSH000001ClosePrice 获取上证指数当前价格
func GetSH000001ClosePrice() (price string, change string, tradeDate time.Time, err error) {
	price, change, tradeDateStr, err := GetStockPrice("000001")
	if err != nil {
		return price, change, tradeDate, err
	}
	tradeDate, err = time.Parse("2006/1/2 0:00:00", tradeDateStr)
	return price, change, tradeDate, err
}

// GetSZ399001ClosePrice 获取深证指数当前价格
func GetSZ399001ClosePrice() (price string, change string, tradeDate time.Time, err error) {
	price, change, tradeDateStr, err := GetStockPrice("399001")
	if err != nil {
		return price, change, tradeDate, err
	}
	tradeDate, err = time.Parse("2006/1/2 0:00:00", tradeDateStr)
	return price, change, tradeDate, err
}

// GetStockPrice 获取股票当前价格以及涨跌幅
func GetStockPrice(stockCode string) (price string, change string, tradeDate string, err error) {
	apiUrl := config.CurrentConfig.MarketPriceApi
	if len(apiUrl) == 0 {
		shareMarketPriceServiceLogger.ErrorFormat(err, "行情接口地址配置不正确 configkey=MarketPriceApi")
		return price, change, tradeDate, errors.New("行情接口地址配置不正确")
	}
	apiUrl = strings.Replace(apiUrl, StockCodeReplace, stockCode, -1)
	resp, err := http.Get(apiUrl)
	defer func() { resp.Body.Close() }()
	if err != nil {
		return price, change, tradeDate, err
	}
	if resp.StatusCode != http.StatusOK {
		return price, change, tradeDate, errors.New("请求异常")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return price, change, tradeDate, err
	}
	//这个接口是utf8带BOM头格式，需要手动去除BOM
	if len(data) > 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	respData := new(DataApiResponse)
	err = json.Unmarshal(data, respData)
	if err != nil {
		return price, change, tradeDate, err
	}
	if respData.Success && len(respData.Data) > 0 && respData.Data[0].Success && len(respData.Data[0].Result) > 2 {
		price = respData.Data[0].Result[2][0]
		tradeDate = respData.Data[0].Result[2][1]
		change = respData.Data[0].Result[2][2]
		if len(price) == 0 {
			return price, change, tradeDate, errors.New("行情解析不正确")
		}
		digArr := strings.Split(price, ".")
		if len(digArr) == 2 {
			num := digArr[0]
			dec := digArr[1]
			if len(dec) > 2 {
				decNum, err := strconv.Atoi(dec[:2])
				if err != nil {
					return price, change, tradeDate, err
				}
				decNum3, err := strconv.Atoi(dec[2:3])
				if err != nil {
					return price, change, tradeDate, err
				}
				if decNum3 > 4 {
					decNum = decNum + 1
				}
				dec = strconv.Itoa(decNum)
			}
			price = num + "." + dec
		}
		return price, change, tradeDate, err
	}
	return price, change, tradeDate, nil
}
