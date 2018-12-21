package yqq

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"encoding/json"
	"git.emoney.cn/softweb/roboadvisor/protected/model/yqq"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers"
	"strconv"
	"sort"
	"strings"
)

//获取runtimcache
var rtcache = handlers.NewWebRuntimeCache("yqq_",60*1)

//获取所有直播间列表
func GetLiveRoomList(ctx dotweb.Context) error {
	cachekey := "GetLiveRoomList"
	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		return ctx.WriteJson(obj)
	}else{

		YqqLiveService := expertnews.ConExpertNews_YqqService()
		RoomListStr,err:=YqqLiveService.GetLiveRoomList()

		if err != nil {
			return ctx.WriteJson("")
		}

		var ret yqq.YqqRetSimple
		err =json.Unmarshal([]byte(RoomListStr),&ret)
		if err != nil {
			return ctx.WriteJson("")
		}

		back_obj := struct{RetCode string ;RetMsg string; Data []yqq.YqqRoomSimple}{
			RetCode:"0",RetMsg:"OK",Data:ret.Data}

		rtcache.SetCache(cachekey,back_obj)
		return ctx.WriteJson(back_obj)

	}

	return ctx.WriteJson("")

}

///返回专家资讯的
func GetHotTop4LiveRoom(ctx dotweb.Context) error{

	cachekey := "GetLiveRoomTop4List"
	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		return ctx.WriteJson(obj)
	}else{

		YqqLiveService := expertnews.ConExpertNews_YqqService()
		RoomListStr,err:=YqqLiveService.GetLiveRoomList()

		if err != nil {
			return ctx.WriteJson("")
		}

		var ret yqq.YqqRetSimple
		err =json.Unmarshal([]byte(RoomListStr),&ret)
		if err != nil {
			return ctx.WriteJson("")
		}
		if ret.Data == nil {
			return ctx.WriteJson("")
		}
		//sort热度
		indexref := make(map[int]yqq.YqqRoomSimple,0)
		var sortarr []int
		for _,v := range ret.Data{
			indexref[v.VisitAddrNum]  = v
			sortarr = append(sortarr,v.VisitAddrNum)
		}

		//sort.Ints(sortarr)
		sort.Sort(sort.Reverse(sort.IntSlice(sortarr)))

		var retlist []yqq.YqqRoomSimple
		for _,a := range sortarr{
			retlist = append(retlist,indexref[a])
		}

		back_obj := struct{RetCode string ;RetMsg string; Data []yqq.YqqRoomSimple}{
			RetCode:"0",RetMsg:"OK",Data:retlist[0:4]}

		rtcache.SetCache(cachekey,back_obj)
		return ctx.WriteJson(back_obj)

	}

	return ctx.WriteJson("")

}




//获取热门和VIP直播间
func GetHotVipLiveRooms(ctx dotweb.Context) error{
	cachekey := "GetHotVipLiveRooms"
	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		ctx.WriteString(obj.(string))
	}else {
		YqqLiveService := expertnews.ConExpertNews_YqqService()
		HotVipLiveRooms,err:=YqqLiveService.GetHotVipLiveRooms()

		if err != nil {
			return ctx.WriteString("")
		}

		rtcache.SetCache(cachekey,HotVipLiveRooms)

		return ctx.WriteString(HotVipLiveRooms)
	}

	return ctx.WriteString("")
}


//获取热门直播间
func GetSkilbag(ctx dotweb.Context) error{

	LidStr := ctx.QueryString("lid")

	LidInt,_ := strconv.Atoi(LidStr)

	cachekey := "GetSkilbag_"+ strconv.Itoa(LidInt)

	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		ctx.WriteString(obj.(string))
	}else {

		YqqLiveService := expertnews.ConExpertNews_YqqService()

		HotVipLiveRooms,err:=YqqLiveService.GetSkilbag(LidStr)

		if err != nil {
			return ctx.WriteString("")
		}

		rtcache.SetCache(cachekey,HotVipLiveRooms)

		return ctx.WriteString(HotVipLiveRooms)
	}

	return ctx.WriteString("")

}

//获取益圈圈直播统计信息
func GetExpertLiveData(ctx dotweb.Context) error{

	cachekey := "GetExpertLiveData"
	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		ctx.WriteString(obj.(string))
	}else {
		YqqLiveService := expertnews.ConExpertNews_YqqService()

		Jsonstr,err:=YqqLiveService.GetExpertLiveData()

		if err != nil {
			return ctx.WriteString("")
		}

		rtcache.SetCache(cachekey,Jsonstr)

		return ctx.WriteString(Jsonstr)
	}

	return ctx.WriteString("")
}

//按标签查询直播列表
func GetTagLiveRoomInfo(ctx dotweb.Context) error{

	TagsId := ctx.QueryString("tags")
	cachekey := "GetTagLiveRoomInfo_" +TagsId
	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		ctx.WriteString(obj.(string))
	}else {
		YqqLiveService := expertnews.ConExpertNews_YqqService()

		Jsonstr,err:=YqqLiveService.GetTagLiveRoomInfo(TagsId)

		if err != nil {
			return ctx.WriteString("")
		}

		rtcache.SetCache(cachekey,Jsonstr)

		return ctx.WriteString(Jsonstr)
	}

	return ctx.WriteString("")
}


//获取最新直播消息的直播间列表
func GetLatestRoomInfo(ctx dotweb.Context) error{
	cachekey := "GetLatestRoomInfo"
	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		return ctx.WriteJson(obj)
	}else {
		YqqLiveService := expertnews.ConExpertNews_YqqService()

		Jsonstr,err:=YqqLiveService.GetLatestRoomInfo()

		if err != nil {
			return ctx.WriteString("")
		}

		var listobj yqq.ScrollLiveRoomRet

		err = json.Unmarshal([]byte(Jsonstr),&listobj)
		_ = err

		if err!= nil {
			return ctx.WriteJson("")
		}

		rtcache.SetCache(cachekey,listobj)

		return ctx.WriteJson(listobj)
		//return ctx.WriteString(Jsonstr)
	}

}



//返回专家策略的策略页面的分类列表
func GetExpertLiveIndexRooms(ctx dotweb.Context) error{
	cachekey := "GetExpertLiveIndexRooms"
	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		return ctx.WriteJson(obj)
	}else {
		back_obj := struct{RetCode string ;RetMsg string; Data interface{}}{
			RetCode:"0",RetMsg:"OK",Data:sub_Expert_StrategyChanenl_TypeGroup()}

		rtcache.SetCache(cachekey,back_obj)

		return ctx.WriteJson(back_obj)
	}

	return ctx.WriteJson("")

}


//获取专家资讯整个页面的单Key缓存的接口
func Expert_EntireYqq(ctx dotweb.Context) error{
	cachekey := "Expert_EntireYqq"
	obj, exist := rtcache.GetCache(cachekey)

	if exist {
		return ctx.WriteJson(obj)
	}else {
		YqqLiveService := expertnews.ConExpertNews_YqqService()
		entire_obj,err := YqqLiveService.GetExpertYqqEntireData_RedisCache()

		if err!=nil {
			return ctx.WriteString("Expert_EntireYqq() Get RedisData Failed")
		}

		back_obj := struct{RetCode string ;RetMsg string; Data interface{}}{
			RetCode:"0",RetMsg:"OK",Data:entire_obj}

		rtcache.SetCache(cachekey,back_obj)

		return ctx.WriteJson(back_obj)
	}

	return ctx.WriteJson(nil)

}


//返回专家资讯的策略页面定制数据
func Expert_StrategyChannel(ctx dotweb.Context) error{
	cachekey := "Expert_StrategyChannel"
	obj, exist := rtcache.GetCache(cachekey)

	if exist {
		return ctx.WriteJson(obj)
	}else {

		entire_obj := struct {
			TypeGroup interface{};LatestLive interface{};YqqStat interface{}
		}{
			TypeGroup:sub_Expert_StrategyChanenl_TypeGroup(),LatestLive:sub_Expert_StrategyChannel_LatestLive(),YqqStat:sub_Expert_StrategyChannel_YqqStat(),
		}

		back_obj := struct{RetCode string ;RetMsg string; Data interface{}}{
			RetCode:"0",RetMsg:"OK",Data:entire_obj}

		rtcache.SetCache(cachekey,back_obj)
		return ctx.WriteJson(back_obj)
	}

	return ctx.WriteJson(nil)
}


//sub
func sub_Expert_StrategyChanenl_TypeGroup()interface{}{

	YqqLiveService := expertnews.ConExpertNews_YqqService()
	var array_len = 6

	back_data  := make(map[string][]yqq.YqqRoomSimple)

	//1技术交易型 2价值投资型 3波段操作型 4短线激进型
	var simret_1,simret_2,simret_3,simret_4,simretall  yqq.YqqRetSimple

	Jsonstr,err:=YqqLiveService.GetAllTagLiveRoomInfo()
	if err!=nil {
		return nil
	}

	err = json.Unmarshal([]byte(Jsonstr),&simretall)

	for _,v := range simretall.Data {
		if strings.Contains(v.TagNameStr,"技术交易型") {
			simret_1.Data = append(simret_1.Data,v)
		}
		if strings.Contains(v.TagNameStr,"价值投资型") {
			simret_2.Data = append(simret_2.Data,v)
		}
		if strings.Contains(v.TagNameStr,"波段操作型") {
			simret_3.Data = append(simret_3.Data,v)
		}
		if strings.Contains(v.TagNameStr,"短线激进型") {
			simret_4.Data = append(simret_4.Data,v)
		}
	}

	back_data["Technology"] = simret_1.Data[0:array_len]
	back_data["Value"] = simret_2.Data[0:array_len]
	back_data["Operate"] = simret_3.Data[0:array_len]
	back_data["ShortTerm"] =simret_4.Data[0:array_len]

	return back_data

}

//sub
func sub_Expert_StrategyChannel_LatestLive() interface{}{
	YqqLiveService := expertnews.ConExpertNews_YqqService()

	Jsonstr,err:=YqqLiveService.GetLatestRoomInfo()

	var listobj yqq.ScrollLiveRoomRet

	err = json.Unmarshal([]byte(Jsonstr),&listobj)
	_ = err

	return listobj.Data
}

//sub
func sub_Expert_StrategyChannel_YqqStat()interface{}{
	YqqLiveService := expertnews.ConExpertNews_YqqService()
	Jsonstr,err:=YqqLiveService.GetExpertLiveData()
	_ = err
	var statobj  yqq.YqqStatRet
	err = json.Unmarshal([]byte(Jsonstr),&statobj)

	return statobj.Data

}


func GetYqqHomeData(ctx dotweb.Context) error{

	cachekey := "GetYqqHomeData"
	obj, exist := rtcache.GetCache(cachekey)
	if exist {
		return ctx.WriteJson(obj)
	}else {
		YqqLiveService := expertnews.ConExpertNews_YqqService()

		Jsonstr,err:=YqqLiveService.GetYQQHomeData()

		if err != nil {
			return ctx.WriteJson("")
		}

		var yqqhomedata yqq.YqqHomeData
		err = json.Unmarshal([]byte(Jsonstr),&yqqhomedata)

		rtcache.SetCache(cachekey,yqqhomedata)

		return ctx.WriteJson(yqqhomedata)
	}

}