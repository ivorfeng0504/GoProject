package myoptional

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	myoptional_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	myoptional_vmmodel "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
	"unicode/utf8"
)

type StockTalkMsgService struct {
	service.BaseService
	stockTalkMsgRepo *myoptional_repo.StockTalkMsgRepository
}

var (
	shareStockTalkMsgRepo   *myoptional_repo.StockTalkMsgRepository
	shareStockTalkMsgLogger dotlog.Logger
)

const (
	stockTalkMsgServiceName                          = "StockTalkMsgService"
	EMNET_StockTalkMsg_PreCacheKey                   = "EMoney:MyOptional:StockTalkMsgBll:"
	EMNET_StockTalk_InsertStockTalkMsgQueue_CacheKey = EMNET_StockTalkMsg_PreCacheKey + "InsertStockTalkMsgQueue:"
)

func NewStockTalkMsgService() *StockTalkMsgService {
	service := &StockTalkMsgService{
		stockTalkMsgRepo: shareStockTalkMsgRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return service
}

func init() {
	protected.RegisterServiceLoader(stockTalkMsgServiceName, stockTalkMsgServiceLoader)
}

func stockTalkMsgServiceLoader() {
	shareStockTalkMsgRepo = myoptional_repo.NewStockTalkMsgRepository(protected.DefaultConfig)
	shareStockTalkMsgLogger = dotlog.GetLogger(stockTalkMsgServiceName)
}

// InsertStockTalkMsgQueue 写入信息到队列中
func (srv *StockTalkMsgService) InsertStockTalkMsgQueue(model *myoptional_vmmodel.StockTalkMsgVM) (err error) {
	if model == nil {
		return errors.New("提交的数据不能为空")
	}
	model.Content = strings.Trim(model.Content, " ")
	model.LiveRoomName = strings.Trim(model.LiveRoomName, " ")
	if len(model.Content) == 0 {
		return errors.New("内容不能为空")
	}
	if len(model.LiveRoomName) == 0 {
		return errors.New("直播间名称不能为空")
	}
	if len(model.StockInfoList) == 0 {
		return errors.New("股票代码不能为空")
	}
	if model.LiveRoomId <= 0 {
		return errors.New("直播间Id不正确")
	}
	//直播间过滤  以下直播间的数据才进行记录
	//拐点形态204、顺势中继205、突破压力206、基本面策略207、抄底三剑客210、长波突进211、中继回踩217
	allowLiveRoom := map[int64]bool{204: true, 205: true, 206: true, 207: true, 210: true, 211: true, 217: true}
	if allowLiveRoom[model.LiveRoomId] == false {
		return nil
	}
	_, err = srv.RedisCache.LPush(EMNET_StockTalk_InsertStockTalkMsgQueue_CacheKey, _json.GetJsonString(model))
	if err != nil {
		shareStockTalkMsgLogger.ErrorFormat(err, "InsertStockTalkMsgQueue 写入信息到队列中 异常 model=%s", _json.GetJsonString(model))
	}
	return err
}

// InsertStockTalkMsg 写入数据到数据库中
func (srv *StockTalkMsgService) InsertStockTalkMsg(model *myoptional_vmmodel.StockTalkMsgVM) (err error) {
	if model == nil {
		return errors.New("提交的数据不能为空")
	}
	model.Content = strings.Trim(model.Content, " ")
	model.LiveRoomName = strings.Trim(model.LiveRoomName, " ")
	if len(model.Content) == 0 {
		return errors.New("内容不能为空")
	}
	if len(model.LiveRoomName) == 0 {
		return errors.New("直播间名称不能为空")
	}
	if len(model.StockInfoList) == 0 {
		return errors.New("股票代码不能为空")
	}
	if model.LiveRoomId <= 0 {
		return errors.New("直播间Id不正确")
	}
	for _, stockInfo := range model.StockInfoList {
		insertModel := myoptional_model.StockTalkMsg{
			LiveRoomId:   model.LiveRoomId,
			LiveRoomName: model.LiveRoomName,
			StockCode:    stockInfo.StockCode,
			StockName:    stockInfo.StockName,
			Content:      model.Content,
			SendTime:     model.SendTime,
			ModifyUser:   "System",
		}
		//如果字数超限  则忽略
		if utf8.RuneCountInString(insertModel.Content) > 1000 {
			continue
		}
		if len(model.ImageList) > 0 {
			insertModel.ImageList = _json.GetJsonString(model.ImageList)
		}
		//如果图片数量超限  则忽略图片
		if utf8.RuneCountInString(insertModel.ImageList) > 2000 {
			insertModel.ImageList = ""
		}
		id, err := srv.stockTalkMsgRepo.InsertStockTalkMsg(&insertModel)
		if err != nil {
			shareStockTalkMsgLogger.ErrorFormat(err, "InsertStockTalkMsg 写入数据到数据库中 异常 model=%s insertModel=%s", _json.GetJsonString(model), _json.GetJsonString(insertModel))
			continue
		} else {
			shareStockTalkMsgLogger.DebugFormat("InsertStockTalkMsg 写入数据到数据库中 成功 主键id为%d model=%s insertModel=%s", id, _json.GetJsonString(model), _json.GetJsonString(insertModel))
		}
	}
	return err
}

// ProcessStockTalkMsgQueue 处理消息队列，将队列中的数据插入数据库
func (srv *StockTalkMsgService) ProcessStockTalkMsgQueue() (err error) {
	for {
		json, err := srv.RedisCache.RPop(EMNET_StockTalk_InsertStockTalkMsgQueue_CacheKey)
		if err == redis.ErrNil {
			return nil
		}
		var result [][]myoptional_vmmodel.StockTalkMsgVM
		err = _json.Unmarshal(json, &result)
		if err != nil {
			shareStockTalkMsgLogger.ErrorFormat(err, "ProcessStockTalkMsgQueue->Unmarshal 反序列异常 json=%s", json)
			continue
		}
		if len(result) == 0 || len(result[0]) == 0 {
			continue
		}
		stockTalkMsgVM := &result[0][0]
		err = srv.InsertStockTalkMsg(stockTalkMsgVM)
		if err != nil {
			shareStockTalkMsgLogger.ErrorFormat(err, "ProcessStockTalkMsgQueue->InsertStockTalkMsg 写入数据库异常 stockTalkMsgVM=%s", _json.GetJsonString(stockTalkMsgVM))
			continue
		}
		time.Sleep(time.Millisecond * 10)
	}

	return err
}
