package stockapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/cache"
	"github.com/devfeel/dotlog"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

const (
	stockInfoApiServiceName              = "StockInfoApiService"
	EMNET_GetStockTableDict_CacheKey     = "StockInfoApiService:GetStockTableDict"
	EMNET_GetStockListA_CacheKey         = "StockInfoApiService:GetStockListA"
	EMNET_GetStockTableDict_CacheSeconds = 60 * 60
)

var (
	shareStockInfoApiServiceLogger dotlog.Logger
	shareStockInfoApiServiceRedis  cache.RedisCache
)

func init() {
	protected.RegisterServiceLoader(stockInfoApiServiceName, func() {
		shareStockInfoApiServiceLogger = dotlog.GetLogger(stockInfoApiServiceName)
		shareStockInfoApiServiceRedis = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	})
}

// GetStockName 根据股票代码获取股票名称
func GetStockName(stockCode string) (stockName string, err error) {
	stockDataMap, err := GetStockListDictCache()
	if err != nil {
		shareStockInfoApiServiceLogger.ErrorFormat(err, "GetStockName 根据股票代码获取股票名称 获取码表异常 股票代码为%s", stockCode)
		return stockName, err
	}
	stockName, success := stockDataMap[stockCode]
	if success {
		return stockName, nil
	} else {
		err = errors.New("获取股票名称失败，该股票代码在码表中不存在")
		shareStockInfoApiServiceLogger.ErrorFormat(err, "GetStockName 获取股票名称失败 该股票代码在码表中不存在 股票代码为%s", stockCode)
		return stockName, err
	}
}

// GetStockListDictCache 获取A股码表缓存
func GetStockListDictCache() (map[string]string, error) {
	stockDataMap := make(map[string]string)
	cacheKey := EMNET_GetStockTableDict_CacheKey
	err := shareStockInfoApiServiceRedis.GetJsonObj(cacheKey, &stockDataMap)
	if err != nil || stockDataMap == nil || len(stockDataMap) == 0 {
		stockTable, innerErr := GetStockListA()
		if innerErr != nil {
			shareStockInfoApiServiceLogger.Error(innerErr, "GetStockListDictCache 获取A股码表缓存 获取码表异常")
			return nil, innerErr
		}
		for _, stock := range stockTable {
			stockDataMap[stock.Code] = stock.Name
		}
		shareStockInfoApiServiceRedis.Set(cacheKey, _json.GetJsonString(stockDataMap), EMNET_GetStockTableDict_CacheSeconds)
		err = innerErr
	}
	return stockDataMap, err
}

func GetStockListACache() (stockList []*StockData, error error) {
	cacheKey := EMNET_GetStockListA_CacheKey
	err := shareStockInfoApiServiceRedis.GetJsonObj(cacheKey, &stockList)

	if err != nil || stockList == nil || len(stockList) == 0 {
		stockTable, err := GetStockListA()
		if err != nil {
			shareStockInfoApiServiceLogger.Error(err, "GetStockListACache 获取A股码表缓存 获取码表异常")
			return nil, err
		}

		stockList = stockTable
		shareStockInfoApiServiceRedis.Set(cacheKey, _json.GetJsonString(stockTable), EMNET_GetStockTableDict_CacheSeconds)
	}

	return stockList, err

}

// GetStockListA 获取A股码表
func GetStockListA() (stockTable []*StockData, err error) {
	apiUrl := config.CurrentConfig.StockTableApi
	if len(apiUrl) == 0 {
		shareStockInfoApiServiceLogger.ErrorFormat(err, "获取A股码表接口地址配置不正确 configkey=StockTableApi")
		return stockTable, errors.New("获取A股码表接口地址配置不正确")
	}
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		return stockTable, errReturn
	}
	_ = contentType
	_ = intervalTime
	shareStockInfoApiServiceLogger.DebugFormat("GetStockListA 查询A股码表 请求地址为：%s  结果为：%s", apiUrl, body)
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		return stockTable, err
	}
	if apiGatewayResp.RetCode != 0 {
		return stockTable, errors.New(apiGatewayResp.RetMsg)
	}
	emcodeResponse := new(EmcodeResponse)
	err = _json.Unmarshal(apiGatewayResp.Message, &emcodeResponse)
	if err != nil {
		return stockTable, err
	}
	if emcodeResponse.ErrNo != 0 {
		err = errors.New(emcodeResponse.ErrMsg)
		return stockTable, err
	}
	if emcodeResponse.Result == nil || len(emcodeResponse.Result) == 0 {
		err = errors.New("码表数据为空")
		return stockTable, err
	}
	stockTable = emcodeResponse.Result
	return stockTable, err
}

// GetQuotes 获取个股或者板块行情
func GetQuotes(emCode string, codetype string) (quotesData *QuotesData, err error) {
	apiUrl := config.CurrentConfig.StockQUOTESUrl
	if len(apiUrl) == 0 {
		shareStockInfoApiServiceLogger.ErrorFormat(err, "获取个股或者板块行情接口地址配置不正确 configkey=StockQUOTESUrl")
		return quotesData, errors.New("获取个股或者板块行情接口地址配置不正确")
	}

	strCode := emCode
	if codetype == "BK" {
		if len(emCode) == 6 {
			strCode = emCode[2:6]
		}
		strCode = "BK" + strCode
	} else {
		firstNum := emCode[0:1]
		if firstNum == "6" {
			strCode = "sh" + emCode
		} else {
			strCode = "sz" + emCode
		}
	}

	apiUrl = fmt.Sprintf(apiUrl, strCode)

	resp, err := http.Get(apiUrl)
	defer func() { resp.Body.Close() }()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("请求异常")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//这个接口是utf8带BOM头格式，需要手动去除BOM
	if len(data) > 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}

	emQuotesResponse := new(EmQuotesResponse)
	err = json.Unmarshal(data, &emQuotesResponse)

	if err != nil {
		shareStockInfoApiServiceLogger.ErrorFormat(err, "获取个股或者板块行情接口反序列化失败")
		return quotesData, err
	}

	if emQuotesResponse.Success {
		var emquotesData []*QuotesData
		err = _json.Unmarshal(emQuotesResponse.Data, &emquotesData)

		if len(emquotesData) > 0 {
			return emquotesData[0], nil
		}
	}

	return nil, nil
}

// GetStockPE 获取个股PE
func GetStockPEList() (peList [][]string, err error) {
	apiUrl := config.CurrentConfig.StockPEUrl
	if len(apiUrl) == 0 {
		shareStockInfoApiServiceLogger.ErrorFormat(err, "获取个股PE接口地址配置不正确 configkey=StockPEUrl")
		return peList, errors.New("获取个股PE接口地址配置不正确")
	}
	nowdate := time.Now()

	apiUrl = fmt.Sprintf(apiUrl, nowdate.Format("2006-01-02"))

	resp, err := http.Get(apiUrl)
	defer func() { resp.Body.Close() }()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("请求异常")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//这个接口是utf8带BOM头格式，需要手动去除BOM
	if len(data) > 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}

	peResponse := new(PEResponse)
	err = json.Unmarshal(data, &peResponse)
	if err != nil {
		shareStockInfoApiServiceLogger.ErrorFormat(err, "获取个股PE接口反序列化失败 data=%s",string(data))
		return nil, err
	}

	if peResponse.Success {
		peData := peResponse.Data
		if len(peData) > 0 {
			peResult := peData[0]
			if peResult.Success {
				peList := peResult.Result

				return peList, nil
			}
		}
	}

	return nil, nil
}

//识别图片中的股票代码
func DiscernStockCode(file multipart.File) (stocklist []string, err error) {

	filebuf := &bytes.Buffer{}
	leng, err := filebuf.ReadFrom(file)

	if leng > 512000 {
		return nil, errors.New("超出500k")
	}

	actoken, _ := GetAccessToken()
	OCRUrl := config.CurrentConfig.Baidu_OCR_Accurate_Url + actoken //"https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic?access_token=" + actoken

	b64str, _ := EncodeImage2B64(filebuf)
	u := url.Values{}
	u.Set("image", string(b64str))

	body, contentType, intervalTime, errReturn := _http.HttpPost(OCRUrl, u.Encode(), "application/x-www-form-urlencoded")

	fmt.Println(body)

	if errReturn != nil {
		return nil, errReturn
	}
	_ = contentType
	_ = intervalTime

	reg := regexp.MustCompile("(((4|8)[\\d]{5})|((00|30|60)[\\d]{4})|(60[\\d]{4})|(01|02|08)[\\d]{4})")
	retarr := reg.FindAllString(body, -1)

	txtbingos := DiscernStockText(body)

	if len(txtbingos) > 0 {
		for _, s := range txtbingos {
			retarr = append(retarr, s.Code)
		}
	}

	sortDupMap := make(map[string]int)
	distStockArr := make([]string, 0)

	for _, v := range retarr {
		if _, exist := sortDupMap[v]; !exist {
			sortDupMap[v] = 1
			distStockArr = append(distStockArr, v)
		}
	}
	return distStockArr, nil

	//sort.Strings(retarr)
	//disarr := make([]string,0)
	//sortsource := Duplicate(retarr)
	//for _,v :=range sortsource{
	//	disarr = append(disarr,v.(string))
	//}
	//return disarr, nil

}

//识别股票名称
func DiscernStockText(src string) (bingoStocks []*StockData) {
	allstock, _ := GetStockListACache()
	bingos := make([]*StockData, 0)

	//bingoDataMap := make(map[string]StockData)
	bingoIndexMap := make(map[int]*StockData)
	indexarr := make([]int, 0)

	for _, v := range allstock {
		if i := strings.Index(src, v.Name); i > 0 {
			//bingos = append(bingos,v)
			//bingoDataMap[v.Name] = *v
			bingoIndexMap[i] = v
			indexarr = append(indexarr, i)
		}
	}

	sort.Ints(indexarr)

	for _, y := range indexarr {
		bingos = append(bingos, bingoIndexMap[y])
	}

	return bingos
}

//图片转Base64
func EncodeImage2B64(buf *bytes.Buffer) (imgb64str string, err error) {
	disbytes := make([]byte, (buf.Len()*15)/10)
	base64.StdEncoding.Encode(disbytes, buf.Bytes())
	truebytes := bytes.TrimRight(disbytes, "\x00")

	b64str := string(truebytes[:])
	return b64str, nil
}

func GetAccessToken() (accesstoken string, err error) {
	// 调用getAccessToken()获取的 access_token建议根据expires_in 时间 设置缓存
	// 返回token示例
	//TOKEN := "24.adda70c11b9786206253ddb70affdc46.2592000.1493524354.282335-1234567"
	// 百度云中开通对应服务应用的 API Key 建议开通应用的时候多选服务
	clientId := config.CurrentConfig.Baidu_OCR_APIKey //"Y7Otl4433eiWk6Rag4Omh9cG"
	// 百度云中开通对应服务应用的 Secret Key
	clientSecret := config.CurrentConfig.Baidu_OCR_SecretKey   //"HUYjrjh1IroWh3AIniTXSWRE74ttzcC6"
	authhost := config.CurrentConfig.Baidu_OCR_AccessToken_Rul //"https://aip.baidubce.com/oauth/2.0/token"

	body, contentType, intervalTime, errReturn := _http.HttpPost(authhost, "grant_type=client_credentials&client_id="+clientId+"&client_secret="+clientSecret, "application/x-www-form-urlencoded")

	if errReturn != nil {
		return "-1", errReturn
	}
	_ = contentType
	_ = intervalTime

	acctoken := new(AccessTokenInfo)
	err = _json.Unmarshal(body, &acctoken)

	if err != nil {
		return "-1", err
	}

	return acctoken.Access_token, nil
}

func Duplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}

type EmcodeResponse struct {
	RequestType string       `json:"RequestType"`
	ErrNo       int          `json:"ErrNo"`
	ErrMsg      string       `json:"ErrMsg"`
	Result      []*StockData `json:"Result"`
}

type StockData struct {
	Code     string   `json:"code"`
	Name     string   `json:"name"`
	Emcode   string   `json:"emcode"`
	Station  string   `json:"station"`
	PinYin   []string `json:"pinyin"`
	Category []int    `json:"category"`
}

type EmQuotesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Data    string `json:"data"`
}
type QuotesData struct {
	//个股名称
	N string
	//个股代码
	C string
	//最新价
	P float32
	//涨跌幅
	F float32
}

type PEResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Token   string            `json:"token"`
	Data    []*PEDataResponse `json:"data"`
}

type PEDataResponse struct {
	Success bool       `json:"success"`
	Errmsg  string     `json:"errmsg"`
	Index   int        `json:"index"`
	Result  [][]string `json:"result"`
}

type AccessTokenInfo struct {
	Access_token  string `json:"access_token"`
	Session_key   string `json:"session_key"`
	Scope         string `json:"scope"`
	Refresh_token string `json:"refresh_token"`
	//Session_secret string `json:"session_secret"`
	//Expires_in int64 `json:"expires_in"`
}
