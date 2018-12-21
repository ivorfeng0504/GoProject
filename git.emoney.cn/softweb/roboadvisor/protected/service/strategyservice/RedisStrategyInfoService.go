package strategyservice

import (
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	marketservice "git.emoney.cn/softweb/roboadvisor/protected/service/market"
	"git.emoney.cn/softweb/roboadvisor/protected/model/strategyservice"
	"encoding/json"
	"strings"
)


type RedisStrategyInfoService struct{
	service.BaseService
}

var (
	shareRedisStrategyInfoLogger dotlog.Logger
)

func NewRedisStrategyInfoService() *RedisStrategyInfoService {
	service := &RedisStrategyInfoService{	}
	return service
}


//获取策略组列表
func (*RedisStrategyInfoService) GetStrategyGroup()(RedisStrateInfo []strategyservice.RedisStrategyInfo){
	svr := marketservice.MarketConService()
	groupstr,_ := svr.GetStrategyGroup()

	groupstr = strings.Replace(groupstr,"\\u0000","",-1)
	var garr []strategyservice.RedisStrategyInfo
	err := json.Unmarshal([]byte(groupstr),&garr)

	if err!=nil {
		return nil
	}

	return  garr
}


//获取专家策略的策略列表
func (*RedisStrategyInfoService) GetExpertStrategyList()(StrategyMap map[int]strategyservice.Strategy){
	svr := marketservice.MarketConService()
	strategyliststr,_ := svr.GetExpertStrategyList()

	strategyliststr = strings.Replace(strategyliststr,"\\u0000","",-1)


	var sarr *[]strategyservice.Strategy
	err := json.Unmarshal([]byte(strategyliststr),&sarr)

	if err!=nil {
		return nil
	}

	strategyMP := make(map[int]strategyservice.Strategy)

	for _,v := range *sarr {
		strategyMP[v.StrategyId] = v
	}

	return strategyMP

}


func (m *RedisStrategyInfoService) GetStrategyData()(RedisStrategy *[]strategyservice.RedisStrategyInfo){

	strategyGroups := m.GetStrategyGroup()
	strategyMap := m.GetExpertStrategyList()

	retGroups := make([]strategyservice.RedisStrategyInfo,0)

	for _,v :=range strategyGroups{
		//排除基本面看盘 70005 和 70011两个组
		if v.StrategyGroupId!=70005 && v.StrategyGroupId!=70011 {

			if len(v.StrategyGroup)>0{

				var strategylist = make([]strategyservice.Strategy,0)
				for _,x := range v.StrategyGroup {
					strategylist = append(strategylist,strategyMap[x.StrategyId])
				}
				//strategyGroups[i].StrategyList = strategylist
				v.StrategyList = strategylist
			}
			retGroups = append(retGroups,v)
		}

	}

	return &retGroups

}

