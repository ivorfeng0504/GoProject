package bk1minute

import (
	"git.emoney.cn/softweb/roboadvisor/web_front_expert/handlers"
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/protected/service/bk1minute"
	"strings"
	bk1minute2 "git.emoney.cn/softweb/roboadvisor/protected/model/bk1minute"
	"git.emoney.cn/softweb/roboadvisor/config"
	"html/template"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/protected/model/hotspot"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
)

//获取runtimcache
var rtcache = handlers.NewWebRuntimeCache("bk1minute_",60)


// 板块一分钟-输出数据
func IndexView(ctx dotweb.Context) error {
	ctx.ViewData().Set("BKInfo", "")
	ctx.ViewData().Set("StaticServerHost", config.CurrentConfig.StaticServerHost)
	ctx.ViewData().Set("ResourceVersion", config.CurrentConfig.ResourceVersion)

	Stockcode := ctx.QueryString("Stockcode")
	if Stockcode == "" {
		return ctx.View("block1Minute.html")
	}

	rs := []rune(Stockcode)
	if len(rs) > 6 {
		Stockcode = string(rs[1:7])
	}
	service := bk1minute.NewBK1MinutesService()
	bkinfo, err := service.GetBK1MinutesInfoByBKCode(Stockcode)
	if err != nil {
		return ctx.View("block1Minute.html")
	}

	if bkinfo != nil {
		upperbk := bkinfo.UpperBK
		lowerbk := bkinfo.LowerBK
		industrycontent := bkinfo.IndustryContent

		//格式化上游板块
		upperbks := strings.Split(upperbk, "，")
		var upperbkList []*bk1minute2.BKInfo
		for i, _ := range upperbks {
			if upperbks[i] != "" {
				upperbkinfo := new(bk1minute2.BKInfo)
				upperbkinfo.BKName = upperbks[i]
				upperbkList = append(upperbkList, upperbkinfo)
			}
		}
		bkinfo.UpperBKList = upperbkList

		//格式化下游板块
		lowerbks := strings.Split(lowerbk, "，")
		var lowerbkList []*bk1minute2.BKInfo
		for i, _ := range lowerbks {
			if lowerbks[i] != "" {
				lowerbkinfo := new(bk1minute2.BKInfo)
				lowerbkinfo.BKName = lowerbks[i]
				lowerbkList = append(lowerbkList, lowerbkinfo)
			}
		}
		bkinfo.LowerBKList = lowerbkList

		//格式化行业内容
		industrycontents := strings.Split(industrycontent, "）、")
		var industryList []*bk1minute2.IndustryInfo
		for i, _ := range industrycontents {
			contents := strings.Split(industrycontents[i], "（")

			if len(contents) > 1 {
				industryInfo := new(bk1minute2.IndustryInfo)
				industryInfo.StockDesc = contents[0]
				industryInfo.StockName = contents[1]

				industryList = append(industryList, industryInfo)
			}
		}
		bkinfo.IndustryContentList = industryList

		keylogic := strings.Replace(bkinfo.BK_KeyLogic, "\n", "<br/>", -1)
		impetus := strings.Replace(bkinfo.BK_Impetus, "\n", "<br/>", -1)

		bkinfo.BK_KeyLogicHTML = template.HTML(keylogic)
		bkinfo.BK_ImpetusHTML = template.HTML(impetus)

	}

	ctx.ViewData().Set("BKInfo", bkinfo)
	return ctx.View("block1Minute.html")
}

// 板块一分钟-上半部分
func BkMinute_UpView(ctx dotweb.Context) error {
	ctx.ViewData().Set("BKInfo", "")
	ctx.ViewData().Set("StaticServerHost", config.CurrentConfig.StaticServerHost)
	ctx.ViewData().Set("ResourceVersion", config.CurrentConfig.ResourceVersion)
	ctx.ViewData().Set("StaticResEnv", config.CurrentConfig.StaticResEnv)

	Stockcode := ctx.QueryString("bkcode")
	bkname := ctx.QueryString("bkname")
	GroupSummary := ctx.QueryString("titlecontent")
	HotspotIdStr := ctx.QueryString("hotspotid")

	GroupSummaryHMTL := strings.Replace(GroupSummary, "*", "<br/>", -1)

	ctx.ViewData().Set("BKName", bkname)
	ctx.ViewData().Set("TopicSummary",template.HTML(GroupSummaryHMTL))
	if Stockcode == "" {
		return ctx.View("bkMinute_up.html")
	}

	rs := []rune(Stockcode)
	if len(rs) > 6 {
		Stockcode = string(rs[1:7])
	}
	cachekey := "BkMinuteByCode:" + Stockcode
	retObj, exist := rtcache.GetCache(cachekey)

	if !exist {
		service := bk1minute.NewBK1MinutesService()
		bkinfo, err := service.GetBK1MinutesInfoByBKCode(Stockcode)
		if err != nil {
			return ctx.View("bkMinute_up.html")
		}

		if bkinfo != nil {
			keylogic := strings.Replace(bkinfo.BK_KeyLogic, "\n", "<br/>", -1)
			impetus := strings.Replace(bkinfo.BK_Impetus, "\n", "<br/>", -1)
			bkinfo.BK_KeyLogicHTML = template.HTML(keylogic)
			bkinfo.BK_ImpetusHTML = template.HTML(impetus)
		}

		rtcache.SetCache(cachekey, bkinfo)
		retObj = bkinfo
	}

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
		hotspotinfo.GroupName = bkname
	}

	GroupSummaryHMTL = strings.Replace(hotspotinfo.TopicSummary, "*", "<br/>", -1)

	ctx.ViewData().Set("TopicSummary",template.HTML(GroupSummaryHMTL))
	ctx.ViewData().Set("BKInfo", retObj)
	return ctx.View("bkMinute_up.html")
}

// 板块一分钟-下半部分
func BkMinute_DownView(ctx dotweb.Context) error{

	ctx.ViewData().Set("BKInfo", "")
	ctx.ViewData().Set("StaticServerHost", config.CurrentConfig.StaticServerHost)
	ctx.ViewData().Set("ResourceVersion", config.CurrentConfig.ResourceVersion)
	ctx.ViewData().Set("StaticResEnv", config.CurrentConfig.StaticResEnv)

	Stockcode := ctx.QueryString("BlockCode")
	if Stockcode == "" {
		return ctx.View("bkMinute_down.html")
	}

	rs := []rune(Stockcode)
	if len(rs) > 6 {
		Stockcode = string(rs[1:7])
	}
	cachekey := "BkMinute_downByCode:" + Stockcode
	retObj, exist := rtcache.GetCache(cachekey)
	if !exist {
		service := bk1minute.NewBK1MinutesService()
		bkinfo, err := service.GetBK1MinutesInfoByBKCode(Stockcode)
		if err != nil {
			return ctx.View("bkMinute_down.html")
		}

		if bkinfo != nil {
			upperbk := bkinfo.UpperBK
			lowerbk := bkinfo.LowerBK
			industrycontent := bkinfo.IndustryContent

			//格式化上游板块
			upperbks := strings.Split(upperbk, "，")
			var upperbkList []*bk1minute2.BKInfo
			for i, _ := range upperbks {
				if upperbks[i] != "" {
					upperbkinfo := new(bk1minute2.BKInfo)
					upperbkinfo.BKName = upperbks[i]
					upperbkList = append(upperbkList, upperbkinfo)
				}
			}
			bkinfo.UpperBKList = upperbkList

			//格式化下游板块
			lowerbks := strings.Split(lowerbk, "，")
			var lowerbkList []*bk1minute2.BKInfo
			for i, _ := range lowerbks {
				if lowerbks[i] != "" {
					lowerbkinfo := new(bk1minute2.BKInfo)
					lowerbkinfo.BKName = lowerbks[i]
					lowerbkList = append(lowerbkList, lowerbkinfo)
				}
			}
			bkinfo.LowerBKList = lowerbkList

			//格式化行业内容
			industrycontents := strings.Split(industrycontent, "）、")
			var industryList []*bk1minute2.IndustryInfo
			for i, _ := range industrycontents {
				contents := strings.Split(industrycontents[i], "（")

				if len(contents) > 1 {
					industryInfo := new(bk1minute2.IndustryInfo)
					industryInfo.StockDesc = contents[0]
					industryInfo.StockName = contents[1]

					industryList = append(industryList, industryInfo)
				}
			}
			bkinfo.IndustryContentList = industryList

			keylogic := strings.Replace(bkinfo.BK_KeyLogic, "\n", "<br/>", -1)
			impetus := strings.Replace(bkinfo.BK_Impetus, "\n", "<br/>", -1)

			bkinfo.BK_KeyLogicHTML = template.HTML(keylogic)
			bkinfo.BK_ImpetusHTML = template.HTML(impetus)

		}

		rtcache.SetCache(cachekey, bkinfo)
		retObj = bkinfo
	}

	ctx.ViewData().Set("BKInfo", retObj)
	return ctx.View("bkMinute_down.html")
}