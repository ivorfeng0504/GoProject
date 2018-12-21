package service

import (
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/const"
	"fmt"
	"encoding/json"
	"github.com/axgle/mahonia"
	"strings"
	"time"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"github.com/devfeel/cache"
	"errors"
)

//定义
type MarketService struct {
	service.BaseService
	innerOutRedis cache.RedisCache
	onTradingRedis cache.RedisCache
	afterTradingRedis cache.RedisCache
}
//日志
var (
	// strategyLogger 共享的Logger实例
	marketApiLogger dotlog.Logger
)
//常量
const (
	RedisKey_Atmosphere = _const.RedisKey_NewsPre + "Market:Atmosphere:"
	InitServiceName = "MarketService"

)
//初始化
func init() {
	protected.RegisterServiceLoader(InitServiceName, initServiceLoader)
}
//初始化载入动作
func initServiceLoader() {
	marketApiLogger = dotlog.GetLogger(InitServiceName)
}
//Constructor
func MarketConService() *MarketService {
	newsService := &MarketService{
	}

	newsService.RedisCache = _cache.GetRedisCacheProvider(protected.MarketRedisConfig)
	newsService.innerOutRedis = _cache.GetRedisCacheProvider(protected.StrategyInnerOutRedisConfig)
	newsService.onTradingRedis = _cache.GetRedisCacheProvider(protected.OnTradingMonitorPoolRedisConfig)
	newsService.afterTradingRedis = _cache.GetRedisCacheProvider(protected.AfterTradingStrategyPoolRedisConfig)
	return newsService
}


//获取市场氛围 大盘情况
func (service *MarketService) GetAtmosphere() (objstr string, err error ){

	var objarr []string
	redisResult , err := service.RedisCache.EVAL("local lRet = redis.call('ZRANGE','MarketAtmosphere_10013', 0, -1);local lBuf = lRet[1];local lVer = string.byte(lBuf, 1, 1);local lNumber = string.byte(lBuf, 2, 3);local lBegin = 4;local lTable = {};for i=1, lNumber do local lType  = string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lTitleLength = string.byte(lBuf, lBegin, lBegin + 2 - 1);lBegin = lBegin + 2;local lEnd = lBegin + lTitleLength;local lTitle = string.sub(lBuf, lBegin, lEnd - 1);lBegin = lEnd;local lCntLen = string.byte(lBuf, lBegin, lBegin + 2 -1);lBegin = lBegin + 2;lEnd = lBegin + lCntLen;local lText = string.sub(lBuf, lBegin, lEnd - 1);local lPer = {type=lType, title_lenght=lTitleLength, title=lTitle, text_size=lCntLen, text=lText};lBegin = lEnd;table.insert(lTable, cjson.encode(lPer));end return {lNumber, lTable};",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		if len(results) < 2{
			fmt.Println("result err, not 2 items", redisResult)
		}else{
			fmt.Println(results[0])
			if datas, isOk := results[1].([]interface{});!isOk {
				fmt.Println("results[1] to []interface{} err", err)
			}else{

				for _, v:=range datas{
					gbkstr := UInt8ToString(v.([]uint8))
					utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
					objarr = append(objarr,utfstr)
				}

				objstr ="["+ strings.Join(objarr,",") +"]"

				fmt.Println(objstr)
			}
		}

	}

	return objstr,err

}

//GetStrategyPool 获取单个策略的股票池
//daycount 当前时间取历史几天的数据
func (service *MarketService) GetStrategyPool(StrategyID string,DayCount int) (objstr string, err error ){
	datestr := time.Now().Format("20060102")
	if (DayCount==0||DayCount<0) {DayCount=1}
	//if (DayCount>90) {DayCount=90}

	var keys []interface{}

	if DayCount==1 {
		keys = append(keys, StrategyID +"_"+datestr)
	}else{
		keys = append(keys,StrategyID )
		for i:=0;i<DayCount ;i++  {
			d, _ := time.ParseDuration("-"+strconv.Itoa(24*i)+"h")
			datestr = time.Now().Add(d).Format("20060102")
			keys = append(keys, datestr)
		}
	}

	var objarr []string
	var luascript string

	luascript ="local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2)*256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lTable = {}; for i=1, #ARGV do local lRet = redis.call('ZRANGE',KEYS [1]..'_'..ARGV[i], 0, -1);for item in pairs(lRet) do local lBuf = lRet[item];local lBegin = 1;local lStatus = string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf, lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1);if(lPrevStr ~=  'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus == 2)) then local lID = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4;  local lType = string.byte(lBuf, lBegin, lBegin+1-1); lBegin = lBegin + 1; local lPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1);  lBegin = lBegin + 4; local lIClose = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lInputHsl = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4;local lBigOrder = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4 -1), 1);lBegin = lBegin + 4;local lStrategyType = string.byte(lBuf, lBegin, lBegin+1-1); local lPer = {code=lID, time=lTime, price = lPrice, max=lCalcPrice, strategy=lStrategyType};table.insert( lTable,  cjson.encode(lPer));end end end return lTable;"

	//策略股池的移动接口的定制化逻辑 20180504 详询邬逸文 严磊
	if (StrategyID=="60002"||StrategyID=="60003") {
		luascript = "local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte (lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2) *256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lTable = {}; local lRet = redis.call('ZRANGE',KEYS[1]..'_10016', 0, -1); for item in pairs(lRet) do local lBuf = lRet[item];local lBegin = 1;local lStatus =  string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf, lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1);if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus ==  2)) then local lID = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lType = string.byte(lBuf, lBegin, lBegin +1-1); lBegin = lBegin + 1; local lPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lIClose =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lInputHsl = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4;local lBigOrder = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4 -1),  1);lBegin = lBegin + 4;local lStrategyType = string.byte(lBuf, lBegin, lBegin+1-1); local lPer = {code=lID, time=lTime, price = lPrice, max=lCalcPrice, strategy=lStrategyType};table.insert( lTable, cjson.encode(lPer));end end return  lTable;"
	}

	redisResult , err := service.RedisCache.EVAL(luascript,1,keys...)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			if StrategyID=="60001" {
				utfstr = TransferStr(utfstr)
			}else{
				utfstr = TransferStrPre(utfstr)
			}
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)
	}
	return objstr,err
}

//Market.BigOrderStrategyHero 英雄榜
//key:60001Hero_10012（策略ID）
func (service *MarketService) BigOrderStrategyHero(StrategyID string) (objstr string, err error ){

	var objarr []string
	var luascript string
	luascript = "local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2)*256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE', '"+StrategyID+"Hero_10012', 0, -1);local lTable = {};for item in pairs (lRet) do local lBuf = lRet [item]; local lBegin = 1;local lStatus = string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf, lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1); if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus == 2)) then local lID = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1),nil); lBegin = lBegin + 4; local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lType = string.byte(lBuf, lBegin, lBegin+1-1); lBegin = lBegin + 1; local lPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lIClose = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lInputHsl = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lBigOrder = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4;local lStrategyType = string.byte(lBuf, lBegin, lBegin+1-1);local lPer = {code=lID, time=lTime, price = lPrice,max=lCalcPrice, strategy=lStrategyType}; table.insert( lTable, cjson.encode(lPer)); end end return lTable;"

	//策略股池的移动接口的定制化逻辑 20180504 详询邬逸文 严磊
	if (StrategyID=="60001") {
		luascript = "local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin +  1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2)*256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE',  '"+StrategyID+"Hero_10012', 0, -1);local lTable = {};for item in pairs (lRet) do local lBuf = lRet [item]; local  lPrevStr = string.sub(lBuf, 1, 12 - 1); if (lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{') then local lBegin = 1; local lID = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lTime = GetDwordVal (string.sub(lBuf, lBegin,  lBegin+4- 1), nil); lBegin = lBegin + 4; local lType = string.byte(lBuf, lBegin, lBegin+1-1); lBegin = lBegin + 1; local lPrice = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lIClose = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lInputHsl = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin =  lBegin + 4; local lBigOrder = GetDwordVal(string.sub (lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4;local lStrategyType = string.byte(lBuf, lBegin, lBegin+1-1);local lPer = {code=lID, time=lTime, price = lPrice,max=lCalcPrice,  strategy=lStrategyType}; table.insert( lTable, cjson.encode(lPer)); end end return lTable;"
	}

	redisResult , err := service.RedisCache.EVAL(luascript,0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			if StrategyID=="60001" {
				utfstr = TransferStr(utfstr)
			}else{
				utfstr = TransferStrPre(utfstr)
			}

			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)

	}
	return objstr,err
}

//回踩策略英雄榜
//Market.BigOrderStrategyHeroBack
//func (service *MarketService) BigOrderStrategyHeroBack(StrategyID string) (objstr string, err error ){
//
//	var objarr []string
//	redisResult , err := service.RedisCache.EVAL("local function GetDwordVal(lBuf) local lBegin = 1; local lType1  = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1;  local  lType2 = string.byte(lBuf, lBegin,  lBegin);lBegin = lBegin +  1; local lType3 = string.byte(lBuf, lBegin,lBegin);lBegin = lBegin + 1;local lType4  =  string.byte(lBuf, lBegin,   lBegin); return lType1  +lType2*16*16 +lType3*16*16*16*16+lType4*16*16*16*16*16*16; end local lRet = redis.call('ZRANGE', '"+StrategyID+"Hero_10012', 0, -1);local lTable = {};for item in pairs (lRet) do   local lBuf = lRet [item]; local lBegin = 1;local lStatus = string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf, lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1); if(lPrevStr ~= 'delete_prev')  then   local lID = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1));  lBegin = lBegin + 4;   local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1));  lBegin = lBegin + 4;   local lType = string.byte(lBuf, lBegin, lBegin+1-1);  lBegin = lBegin + 1;   local lPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1));  lBegin = lBegin + 4;   local lCalcPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1));  lBegin = lBegin + 4;   local lIClose = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1));  lBegin = lBegin + 4;   local lInputHsl = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1));  lBegin = lBegin + 4;   local lBigOrder = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1)); lBegin = lBegin + 4;local lStrategyType = string.byte(lBuf, lBegin, lBegin+1-1);local lPer = {code=lID, time=lTime, price =  lPrice,max=lCalcPrice, strategy=lStrategyType};  table.insert(  lTable, cjson.encode(lPer)); end end return lTable;",0)
//
//	if results, isOk := redisResult.([]interface{});!isOk{
//		fmt.Println("interface{} to []interface{} err", err)
//	}else{
//		for _, v:=range results{
//			gbkstr := UInt8ToString(v.([]uint8))
//			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
//			utfstr = TransferStr(utfstr)
//			objarr = append(objarr,utfstr)
//		}
//		objstr ="["+ strings.Join(objarr,",") +"]"
//		fmt.Println(objstr)
//
//	}
//	return objstr,err
//}

// 行情指数标签
// 获得最新交易日上证指数走势上的所有板块标注
func (service *MarketService) UnusualGroup(DateStr string) (objstr string, err error ){

	var objarr []string
	var redisResult interface{}
	var results []interface{}
	var isOk bool

	//判断时间是否为空
	if DateStr== ""{
		//循环往前10天
		breaked := false
		for i:=0;i<10 ;i++  {

			if(!breaked){
				d, _ := time.ParseDuration("-"+strconv.Itoa(24*i)+"h")
				DateStr = time.Now().Add(d).Format("20060102")

				redisResult , err = service.RedisCache.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte (lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2) *256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE','BkMark_"+DateStr+"', 0, -1);local lTable = {};for item in pairs(lRet) do local lBuf = lRet[item];local lBegin = 1;local lStatus =  string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte  (lBuf, lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1);if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus ==  2)) then  local lVer = string.byte(lBuf, lBegin, lBegin+1-1);lBegin = lBegin + 1;local lType = string.byte(lBuf, lBegin, lBegin+1-1);lBegin = lBegin + 1;if (lType > 0) then local lCommodityID = GetDwordVal(string.sub(lBuf, lBegin,  lBegin+3), 1); lBegin = lBegin + 4;local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+3),nil); lBegin = lBegin + 4;local lRange = string.byte(lBuf, lBegin, lBegin+1);lBegin = lBegin+1;if(lRange == 2 or lRange == 3) then local lPer  = {time=lTime, code=lCommodityID, name='', stock={}}; local lStrTmp = cjson.encode(lPer);lStrTmp= string.gsub(lStrTmp, '{}' , '[]');table.insert(lTable, lStrTmp);end elseif(lType == 0) then local lNameLength = string.byte(lBuf, lBegin,  lBegin + 2 - 1); lBegin = lBegin + 2; local lEnd = lBegin + lNameLength;  local lName = string.sub(lBuf, lBegin, lEnd - 1); lBegin = lEnd; local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+3), 1); lBegin = lBegin + 4; local  lCount = string.byte(lBuf, lBegin,  lBegin+1); lBegin = lBegin + 2; local lStockTable = {}; for i=1, lCount do local lGoodsID = GetDwordVal(string.sub(lBuf, lBegin, lBegin+3), 1);lBegin = lBegin + 4;table.insert (lStockTable,  lGoodsID);end  local lRange = string.byte(lBuf, lBegin, lBegin+1);lBegin = lBegin+1;if(lRange == 2 or lRange == 3) then local lPer = {time=lTime, code=0, name=lName, stock=lStockTable};table.insert(lTable, cjson.encode(lPer)); end end  end end return  lTable;",0)
				if results, isOk = redisResult.([]interface{});isOk{
					breaked=true

					for _, v:=range results{
						gbkstr := UInt8ToString(v.([]uint8))
						utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
						objarr = append(objarr,utfstr)
					}
					objstr ="{\"date\":"+DateStr+",\"group\":["+ strings.Join(objarr,",") +"]}"

				}
			}
		}

	}else{
		redisResult , err := service.RedisCache.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte (lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2) *256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE','BkMark_"+DateStr+"', 0, -1);local lTable = {};for item in pairs(lRet) do local lBuf = lRet[item];local lBegin = 1;local lStatus =  string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte  (lBuf, lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1);if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus ==  2)) then  local lVer = string.byte(lBuf, lBegin, lBegin+1-1);lBegin = lBegin + 1;local lType = string.byte(lBuf, lBegin, lBegin+1-1);lBegin = lBegin + 1;if (lType > 0) then local lCommodityID = GetDwordVal(string.sub(lBuf, lBegin,  lBegin+3), 1); lBegin = lBegin + 4;local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+3),nil); lBegin = lBegin + 4;local lRange = string.byte(lBuf, lBegin, lBegin+1);lBegin = lBegin+1;if(lRange == 2 or lRange == 3) then local lPer  = {time=lTime, code=lCommodityID, name='', stock={}}; local lStrTmp = cjson.encode(lPer);lStrTmp= string.gsub(lStrTmp, '{}' , '[]');table.insert(lTable, lStrTmp);end elseif(lType == 0) then local lNameLength = string.byte(lBuf, lBegin,  lBegin + 2 - 1); lBegin = lBegin + 2; local lEnd = lBegin + lNameLength;  local lName = string.sub(lBuf, lBegin, lEnd - 1); lBegin = lEnd; local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+3), 1); lBegin = lBegin + 4; local  lCount = string.byte(lBuf, lBegin,  lBegin+1); lBegin = lBegin + 2; local lStockTable = {}; for i=1, lCount do local lGoodsID = GetDwordVal(string.sub(lBuf, lBegin, lBegin+3), 1);lBegin = lBegin + 4;table.insert (lStockTable,  lGoodsID);end  local lRange = string.byte(lBuf, lBegin, lBegin+1);lBegin = lBegin+1;if(lRange == 2 or lRange == 3) then local lPer = {time=lTime, code=0, name=lName, stock=lStockTable};table.insert(lTable, cjson.encode(lPer)); end end  end end return  lTable;",0)

		if results, isOk := redisResult.([]interface{});!isOk{
			fmt.Println("interface{} to []interface{} err", err)
		}else{

			for _, v:=range results{
				gbkstr := UInt8ToString(v.([]uint8))
				utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
				objarr = append(objarr,utfstr)
			}

			objstr ="{\"date\":"+DateStr+",\"group\":["+ strings.Join(objarr,",") +"]}"

			fmt.Println(objstr)
		}
	}

	return objstr,err
}

//今日策略
func (service *MarketService) GetTodayStrategy(StrategyID string) (objstr string, err error ){
	datestr := time.Now().Format("20060102")
	var objarr []string
	redisResult , err := service.RedisCache.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2)*256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lRet =  redis.call ('ZRANGE','StrategyMgr_10007', 0, -1);local lTable = {};for item in pairs(lRet) do local lBuf = lRet[item];local lBegin = 1;local lStatus = string.byte(lBuf, lBegin,    lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf,   lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1);if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus == 2)) then local lstrID = string.sub(lBuf, lBegin, lBegin + 4 - 1);lBegin = lBegin + 4;local lEnd =  lBegin + 24;local lName =  string.sub(lBuf, lBegin, lEnd - 1);lName = string.format('%s',lName); lBegin = lEnd; local lPower = string.sub(lBuf, lBegin, lBegin + 4 - 1);local lPer = {id=GetDwordVal(lstrID, nil), name=lName, power=GetDwordVal(lPower, 1)};table.insert( lTable, cjson.encode(lPer));end end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			objarr = append(objarr,utfstr)
		}
		objstr ="{\"date\":"+datestr+",\"strategy\":["+ strings.Join(objarr,",") +"]}"
		fmt.Println(objstr)
		fmt.Println(datestr)
	}
	return objstr,err
}

////全部策略
func (service *MarketService) GetStrategyList(StrategyID string) (objstr string, err error ){

	var objarr []string
	redisResult , err := service.RedisCache.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2)*256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lRet =  redis.call ('ZRANGE','StrategyMgr_10008', 0, -1);local lTable = {};for item in pairs(lRet) do local lBuf = lRet[item];local lBegin = 1;local lStatus = string.byte(lBuf, lBegin,    lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf,   lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1);if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus == 2)) then local lstrID = string.sub(lBuf, lBegin, lBegin + 4 - 1);lBegin = lBegin + 4;local lEnd =  lBegin + 24;local lName =  string.sub(lBuf, lBegin, lEnd - 1);lName = string.format('%s',lName); lBegin = lEnd; local lPower = string.sub(lBuf, lBegin, lBegin + 4 - 1);local lPer = {id=GetDwordVal(lstrID, nil), name=lName, power=GetDwordVal(lPower, 1)};table.insert( lTable, cjson.encode(lPer));end end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)
	}
	return objstr,err
}

//回踩预选池
func (service *MarketService) ActiveGoodsPrimary(StrategyID string) (objstr string,err error){
	datestr := time.Now().Format("20060102")
	var objarr []string
	redisResult , err := service.RedisCache.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte (lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then lTotal = -(1 + (255-lType1)+(255-lType2) *256+(255-lType3)*65536+(255-lType4)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE', '"+StrategyID+"Primary_10015', 0, -1);local lTable = {};for item in pairs (lRet) do local lBuf = lRet [item]; local lBegin = 1;local  lStatus = string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf, lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1); if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or  lStatus == 2)) then local lID = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lType = string.byte(lBuf,  lBegin, lBegin+1-1); lBegin = lBegin + 1; local lPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local  lIClose = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lInputHsl = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lBigOrder = GetDwordVal(string.sub(lBuf, lBegin,  lBegin+4-1), 1); lBegin = lBegin + 4;local lStrategyType = string.byte(lBuf, lBegin, lBegin+1-1);local lPer = {code=lID, time=lTime, price = lPrice,max=lCalcPrice, strategy=lStrategyType}; table.insert( lTable, cjson.encode(lPer)); end  end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			utfstr = TransferStrPre(utfstr)
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)
		fmt.Println(datestr)
	}
	return objstr,err

}


//通用策略股池 传入RedisKey
func (service *MarketService) CommonPool(Key string) (objstr string,err error){

	if(Key==""){Key="StrategyMaster_80001"}

	var objarr []string
	redisResult , err := service.RedisCache.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin);   lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin   + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647)   then lTotal = -(1 + (255-lType4)+(255-lType3)*256+(255-lType2)*65536+(255-lType1)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE','"+Key+"', 0, -1);local lTable = {};for item in pairs (lRet) do local lBuf = lRet [item]; local lBegin = 1;local lPrevStr = string.sub  (lBuf, 1, 12 - 1); if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{') then local lID =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1),  nil); lBegin = lBegin + 4; local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil);  lBegin = lBegin + 4; local lType = string.byte(lBuf,  lBegin, lBegin+1-1); lBegin = lBegin + 1; local lPrice = GetDwordVal(string.sub(lBuf, lBegin,  lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lIClose =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin  = lBegin + 4; local lInputHsl = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1);  lBegin = lBegin + 4; local lBigOrder = GetDwordVal(string.sub (lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4;local lStrategyType = string.byte (lBuf, lBegin, lBegin+1-1);local lPer = {code=lID, time=lTime,  price = lPrice,max=lCalcPrice, strategy=lStrategyType}; table.insert( lTable,  cjson.encode(lPer)); end end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			utfstr = TransferStrPre(utfstr)
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)
	}
	return objstr,err

}

//盘中预警股池 传入RedisKey + 日期
func (service *MarketService) OnTradingMonitorPool(Key string,Date string) (objstr string,err error){

	if(Key==""){return "-1",errors.New("need key params")}
	if(Date==""){return "-1",errors.New("need date params")}

	RedisKey := Key +"_"+Date
	var objarr []string
	redisResult , err := service.onTradingRedis.EVAL("local function GetDwordVal(lBuf, bInt) local  lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local  lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte (lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin);  local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal >  2147483647) then lTotal = -(1 + (255-lType4)+(255-lType3)*256+(255-lType2)*65536+(255- lType1)*16777216); end return lTotal; end local lRet = redis.call ('ZRANGE','"+RedisKey+"', 0, -1);local lTable = {};for item in pairs(lRet) do local lBuf  = lRet[item];local lBegin = 1;local lStatus = string.byte(lBuf, lBegin, lBegin);lBegin =  lBegin + 1;local lStatusID = string.byte  (lBuf, lBegin, lBegin+4-1);lBegin = lBegin +  4;local lPrevStr = string.sub(lBuf, 1, 12 - 1);if(lPrevStr ~= 'delete_prev' and string.sub (lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus == 2)) then local lID = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lTime =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lType =  string.byte(lBuf, lBegin, lBegin+1-1);lBegin = lBegin + 1; local lPrice = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lIClose =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lInputHsl  = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4;local  lBigOrder = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4 -1), 1);lBegin = lBegin + 4;local  lStrategyType = string.byte(lBuf, lBegin, lBegin+1-1); local lPer = {code=lID, time=lTime,  price = lPrice, max=lCalcPrice, strategy=lStrategyType};table.insert( lTable, cjson.encode (lPer));end end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			utfstr = TransferStrPre(utfstr)
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)
	}
	return objstr,err
}


//盘后策略选股股池 传入RedisKey + 日期
func (service *MarketService) AfterTradingStrategyPool(Key string,Date string) (objstr string,err error){
	if(Key==""){return "-1",errors.New("need key params")}
	if(Date==""){return "-1",errors.New("need date params")}

	RedisKey := Key +"_"+Date
	var objarr []string
	redisResult , err := service.afterTradingRedis.EVAL("local function GetDwordVal(lBuf, bInt) local  lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin); lBegin = lBegin + 1; local  lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte (lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType4 = string.byte(lBuf,lBegin,lBegin);  local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal >  2147483647) then lTotal = -(1 + (255-lType4)+(255-lType3)*256+(255-lType2)*65536+(255- lType1)*16777216); end return lTotal; end local lRet = redis.call ('ZRANGE','"+RedisKey+"', 0, -1);local lTable = {};for item in pairs(lRet) do local lBuf  = lRet[item];local lBegin = 1;local lStatus = string.byte(lBuf, lBegin, lBegin);lBegin =  lBegin + 1;local lStatusID = string.byte  (lBuf, lBegin, lBegin+4-1);lBegin = lBegin +  4;local lPrevStr = string.sub(lBuf, 1, 12 - 1);if(lPrevStr ~= 'delete_prev' and string.sub (lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus == 2)) then local lID = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lTime =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lType =  string.byte(lBuf, lBegin, lBegin+1-1);lBegin = lBegin + 1; local lPrice = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lIClose =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lInputHsl  = GetDwordVal (string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4;local  lBigOrder = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4 -1), 1);lBegin = lBegin + 4;local  lStrategyType = string.byte(lBuf, lBegin, lBegin+1-1); local lPer = {code=lID, time=lTime,  price = lPrice, max=lCalcPrice, strategy=lStrategyType};table.insert( lTable, cjson.encode (lPer));end end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			utfstr = TransferStrPre(utfstr)
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)
	}
	return objstr,err

}


//内部用通用策略股池 传入RedisKey
//短线宝
func (service *MarketService) InnerOutPool(Key string) (objstr string,err error){

	if(Key==""){Key="StrategyMaster_80001"}

	var objarr []string
	redisResult , err := service.innerOutRedis.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin);   lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin + 1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin   + 1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647)   then lTotal = -(1 + (255-lType4)+(255-lType3)*256+(255-lType2)*65536+(255-lType1)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE','"+Key+"', 0, -1);local lTable = {};for item in pairs (lRet) do local lBuf = lRet [item]; local lBegin = 1;local lPrevStr = string.sub  (lBuf, 1, 12 - 1); if(lPrevStr ~= 'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del{') then local lID =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1),  nil); lBegin = lBegin + 4; local lTime = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), nil);  lBegin = lBegin + 4; local lType = string.byte(lBuf,  lBegin, lBegin+1-1); lBegin = lBegin + 1; local lPrice = GetDwordVal(string.sub(lBuf, lBegin,  lBegin+4-1), 1); lBegin = lBegin + 4; local lCalcPrice =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4; local lIClose =  GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1); lBegin  = lBegin + 4; local lInputHsl = GetDwordVal(string.sub(lBuf, lBegin, lBegin+4-1), 1);  lBegin = lBegin + 4; local lBigOrder = GetDwordVal(string.sub (lBuf, lBegin, lBegin+4-1), 1); lBegin = lBegin + 4;local lStrategyType = string.byte (lBuf, lBegin, lBegin+1-1);local lPer = {code=lID, time=lTime,  price = lPrice,max=lCalcPrice, strategy=lStrategyType}; table.insert( lTable,  cjson.encode(lPer)); end end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			utfstr = TransferStrPre(utfstr)
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)
	}
	return objstr,err

}


//获取策略组列表 2018年7月20日 策略副窗需求
func (service *MarketService) GetStrategyGroup()(objstr string, err error ) {

	var objarr []string
	redisResult , err := service.RedisCache.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin);  lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin +   1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin +  1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then  lTotal = -(1 +   (255-lType4)+(255-lType3)*256+(255-lType2)*65536+(255-lType1)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE', 'StrategyGroup_10018', 0, -1);local lTable = {};for item in pairs (lRet)  do local lBuf = lRet [item];   local lBegin = 1;local lStatus = string.byte(lBuf, lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf,  lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 - 1); if(lPrevStr ~=   'delete_prev' and string.sub(lPrevStr, 1,4) ~= 'del {' and (lStatus == 0 or lStatus == 2)) then local lVer = string.byte(lBuf, lBegin, lBegin+1-1); lBegin = lBegin + 1;local lID = GetDwordVal(string.sub(lBuf,  lBegin, lBegin+4-1), nil);   lBegin = lBegin + 4; local lSize = string.byte(lBuf,  lBegin, lBegin+2-1); lBegin = lBegin + 2; local lName = string.sub(lBuf, lBegin,  lBegin+lSize-1);lBegin = lBegin + lSize;local lGroup = string.byte(lBuf,  lBegin, lBegin+2-1); lBegin   = lBegin + 2; local lTableGroup = {};for  groupItem=1, lGroup, 1 do local lGroupID = GetDwordVal(string.sub(lBuf,  lBegin, lBegin+4-1), nil); lBegin = lBegin + 4; local lPerGroup = {id=lGroupID}; table.insert(lTableGroup,  lPerGroup); end local lPer  = {ID=lID,name=lName,StrategyGroup=lTableGroup}; table.insert( lTable, cjson.encode(lPer)); end end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"
		fmt.Println(objstr)
	}
	return objstr,err

}

//获取专家策略全部策略列表 2018年7月21日 策略副窗需求
func (service *MarketService) GetExpertStrategyList()(objstr string, err error ) {

	var objarr []string
	redisResult , err := service.RedisCache.EVAL("local function GetDwordVal(lBuf, bInt) local lBegin = 1; local lType1 = string.byte(lBuf,lBegin,lBegin);  lBegin = lBegin + 1; local lType2 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin +  1; local lType3 = string.byte(lBuf,lBegin,lBegin);lBegin = lBegin +  1; local lType4 = string.byte(lBuf,lBegin,lBegin); local lTotal = lType1 +lType2*256+lType3*65536+lType4*16777216;if(bInt and lTotal > 2147483647) then  lTotal = -(1 +  (255-lType4)+(255-lType3)*256+(255-lType2)*65536+(255-lType1)*16777216); end return lTotal; end local lRet = redis.call('ZRANGE',  'StrategyMgr_10019', 0, -1);local lTable = {};for item in pairs (lRet) do local lBuf = lRet [item]; local  lBegin = 1;local lStatus = string.byte(lBuf,  lBegin, lBegin);lBegin = lBegin + 1;local lStatusID = string.byte(lBuf, lBegin, lBegin+4-1);lBegin = lBegin + 4;local lPrevStr = string.sub(lBuf, 1, 12 -  1); if(lPrevStr ~= 'delete_prev' and  string.sub(lPrevStr, 1,4) ~= 'del{' and (lStatus == 0 or lStatus == 2)) then local lID = GetDwordVal(string.sub(lBuf,  lBegin, lBegin+4-1), nil); lBegin = lBegin + 4;local lName = string.sub(lBuf,  lBegin, lBegin+24-1);local lPer =  {ID=lID,name=lName}; table.insert( lTable, cjson.encode(lPer));end end return lTable;",0)

	if results, isOk := redisResult.([]interface{});!isOk{
		fmt.Println("interface{} to []interface{} err", err)
	}else{
		for _, v:=range results{
			gbkstr := UInt8ToString(v.([]uint8))
			utfstr := ConvertToString(gbkstr, "gbk", "utf-8")
			objarr = append(objarr,utfstr)
		}
		objstr ="["+ strings.Join(objarr,",") +"]"

	}
	return objstr,err

}




//大单天眼
func TransferStr(sour string) string{
	//移动接口定制化逻辑 20180504 详询邬逸文 严磊
	//items := [][]string{{"0","继续观望"},{"1","适当止盈"},{"2","建议止盈"},{"3","适当止损"},{"4","策略结束"}}
	items := [][]string{{"0","操作期"},{"1","操作期"},{"2","成功"},{"3","操作期"},{"4","失败"},{"5","退出"}}
	for _,v:=range items{
		sour = strings.Replace(sour,("\"strategy\":"+v[0]),("\"strategy\":\""+v[1]+"\""),-1)
	}

	return sour
}


//回踩策略池
func TransferStrPre(sour string) string{
	//移动接口定制化逻辑 20180504 详询邬逸文 严磊
	//items := [][]string{{"0","继续观望"},{"1","适当止盈"},{"2","建议止盈"},{"3","适当止损"},{"4","策略结束"}}
	items := [][]string{{"0","关注"},{"1","成功"},{"2","失败"},{"3","退出"},{"4","退出"},{"5","操作期"}}
	for _,v:=range items{
		sour = strings.Replace(sour,("\"strategy\":"+v[0]),("\"strategy\":\""+v[1]+"\""),-1)
	}

	return sour
}


func UInt8ToString(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func FromJson(str string ,stru interface{}) error{
	err := json.Unmarshal([]byte(str),&stru)
	if err !=nil {
		return err
	}
	return nil
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func ClearUnicStr(src string) string{
	textQuoted := strconv.QuoteToASCII(src)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	fmt.Println(textUnquoted)

	sUnicodev := strings.Split(textUnquoted, "\\u")
	var context string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		context += fmt.Sprintf("%c", temp)
	}
	fmt.Println(context)
	return context
}