package topic

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers"
	"git.emoney.cn/softweb/roboadvisor/protected/model/hotspot"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"strings"
	"html/template"
)

//获取runtimcache
var rtcache = handlers.NewWebRuntimeCache("topic_",60*2)

// GetTopicList 获取主题列表
func GetTopicList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()

	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	cachekey := "GetTopicList:" + currpageStr + ":" + pageSizeStr
	cacheKey_total := "GetTopicList_totalCount:" + currpageStr + ":" + pageSizeStr
	obj, exist := rtcache.GetCache(cachekey)
	obj_totalcount, exist := rtcache.GetCache(cacheKey_total)
	if !exist {
		expertTopicService := expertnews.NewExpertNews_TopicService()
		topiclist, totalCount, err := expertTopicService.GetTopicListPage(currpage, pageSize)

		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}

		rtcache.SetCache(cachekey, topiclist)
		rtcache.SetCache(cacheKey_total, totalCount)
		obj = topiclist
		obj_totalcount = totalCount
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = obj
	response.TotalCount = obj_totalcount.(int)
	return ctx.WriteJson(response)
}

// GetTopicInfoByID 获取主题详情
func GetTopicInfoByID(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	TopicIDStr := ctx.QueryString("TopicID")
	TopicID, err := strconv.Atoi(TopicIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "TopicID不正确"
		return ctx.WriteJson(response)
	}

	cachekey := "GetTopicInfoByID:" + TopicIDStr
	obj, exist := rtcache.GetCache(cachekey)
	if !exist {
		expertTopicService := expertnews.NewExpertNews_TopicService()
		topicInfo, err := expertTopicService.GetTopicInfoByID(TopicID)

		if err != nil {
			response.RetCode = -2
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		rtcache.SetCache(cachekey, topicInfo)
		obj = topicInfo
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = obj
	return ctx.WriteJson(response)
}


func GetHotspotById(ctx dotweb.Context) error{

	GroupName := ctx.QueryString("bkname")
	GroupSummary := ctx.QueryString("titlecontent")
	HotspotIdStr := ctx.QueryString("hotspotid")

	HotspotId, err := strconv.Atoi(HotspotIdStr)
	if err != nil {
		return ctx.WriteString("hotspotid 错误")
	}

	var hotspotinfo  = new(model.Hotspot)

	if HotspotId!=0 {
		topicsrv := expertnews.NewExpertNews_TopicService()
		hotspotinfo,err = topicsrv.GetHotspotInfo(HotspotId)
	}

	if hotspotinfo==nil || HotspotId==0 {
		hotspotinfo  = new(model.Hotspot)
		hotspotinfo.TopicSummary = GroupSummary
		hotspotinfo.GroupName = GroupName
	}

	//20181029-新增需求 客户端约定*号替换为换行
	hotspotinfo.TopicSummaryHTML = template.HTML(strings.Replace(hotspotinfo.TopicSummary, "*", "<br/>", -1))

	if len(hotspotinfo.TopicPic)==0{
		hotspotinfo.TopicPic = "http://test.static.emoney.cn:8081/webupload/JRPTSoftManage/2018/7/anser_FJbEJcbC-BeGI-EAEb-JEba-cbCJDEaHBdDI.jpg"
	}

	ctx.ViewData().Set("info",hotspotinfo)

	err = ctx.View("hotspottop.html")
	return err
}


//
func GetHotspotContentById(ctx dotweb.Context) error{
	HotspotIdStr := ctx.QueryString("hotspotid")
	HotspotId, err := strconv.Atoi(HotspotIdStr)
	BkCode := ctx.QueryString("bkcode")


	if err != nil {
		return ctx.WriteString("hotspotid 错误")
	}
	topicsrv := expertnews.NewExpertNews_TopicService()
	hotspotinfo,err := topicsrv.GetHotspotInfo(HotspotId)

	if hotspotinfo ==nil {
		return ctx.WriteHtml(`
		<script>
				location.href = "/myoptional/myoptional/bknews?BlockCode=`+BkCode+`"
		</script>
		`)
	}else{
		return ctx.WriteHtml(`<html><body style="background-color:#fff">`+hotspotinfo.TopicContent+`</body></html>`)
	}

}

//主题-个股排序
type RelationStocks []*expertnews2.Topic_Stock

func (s RelationStocks) Len() int {
	return len(s)
}
func (s RelationStocks) Less(i, j int) bool {
	return s[i].StockSortIndex > s[j].StockSortIndex
}
func (s RelationStocks) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}