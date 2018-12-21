package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/mapper"
	"time"
)

// GetTodayNews 获取今日头条
func GetTodayNews() (newsList []*ExpertNewsInfo, err error) {
	req := contract.NewApiRequest()
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/gettodaynews", req)
	if err != nil {
		return newsList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return newsList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &newsList)
	return newsList, nil
}

//GetClosingNews 获取盘后预测资讯
func GetClosingNews() (newsList []*ExpertNewsInfo, err error) {
	req := contract.NewApiRequest()
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/getclosingnews", req)
	if err != nil {
		return newsList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return newsList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &newsList)
	return newsList, nil
}

//GetNewsInfo 获取指定日期的要闻
func GetNewsInfo(date *time.Time) (newsList []*ExpertNewsInfo, err error) {
	now := time.Now()
	if date == nil {
		date = &now
	}
	req := contract.NewApiRequest()
	req.RequestData = date
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/getnews", req)
	if err != nil {
		return newsList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return newsList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &newsList)
	return newsList, nil
}

// GetNewsInfoTopN 获取最新的N条要闻
func GetNewsInfoTopN() (newsList []*ExpertNewsInfo, err error) {
	req := contract.NewApiRequest()
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/getnewstopn", req)
	if err != nil {
		return newsList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return newsList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &newsList)
	if err == nil {
		//去除最后一天的数据
		length := len(newsList)
		if newsList != nil && length > 0 {
			lastDateIndex := -1
			lastDate := ""
			for i := length - 1; i >= 0; i-- {
				publishTime := time.Time(newsList[i].PublishTime).Format("20060102")
				if lastDate != "" && lastDate != publishTime {
					lastDateIndex = i + 1
					break
				}
				lastDate = publishTime
			}
			if lastDateIndex > 0 {
				newsList = newsList[:lastDateIndex]
			}
		}
	}
	return newsList, nil
}

// GetHotNewsInfo 获取热门资讯
func GetHotNewsInfo() (newsList []*ExpertNewsInfo, err error) {
	req := contract.NewApiRequest()
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/gethotnews", req)
	if err != nil {
		return newsList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return newsList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &newsList)
	return newsList, nil
}

// GetYaoWenHomePageData 获取热读页面数据
func GetYaoWenHomePageData() (pageData *YaoWenHomePageData, err error) {
	req := contract.NewApiRequest()
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/getyaowenpagedata", req)
	if err != nil {
		return pageData, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return pageData, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &pageData)
	return pageData, nil
}

// GetTopicNewsInfo 获取主题的相关资讯
func GetTopicNewsInfo(topicId int) (newsList []*ExpertNewsInfo, err error) {
	req := contract.NewApiRequest()
	req.RequestData = topicId
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/gettopicnews", req)
	if err != nil {
		return newsList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return newsList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &newsList)
	return newsList, nil
}

// GetNewsInfoDetail 获取资讯详情
func GetNewsInfoDetail(newsInfoId int64) (newsInfo *ExpertNewsInfo, err error) {
	req := contract.NewApiRequest()
	req.RequestData = newsInfoId
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/expertnews/getnewsdetail", req)
	if err != nil {
		return nil, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return nil, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &newsInfo)
	return newsInfo, nil
}

// NotifyUpdate 通知更新资讯数据库
func NotifyUpdate() (err error) {
	req := contract.NewApiRequest()
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/expertnews/notifyupdate?tableName="+_const.NewsInformationTable, req)
	if err != nil {
		return errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return errors.New(response.RetMsg)
	}
	return nil
}

// 资讯要闻信息
type ExpertNewsInfo struct {
	//新闻主键Id
	NewsInformationId int64
	//文章标题
	ArticleTitle string `mapper:"article_title"`
	//文章短标题
	ArticleShortTitle string `mapper:"article_short_title"`
	//发布时间
	PublishTime        mapper.JSONTime `mapper:"publish_time"`
	CreateTime         mapper.JSONTime `mapper:"create_time"`
	PublistTime_format string
	//内容详情
	DataUrl string `mapper:"data_url"`
	//内容摘要
	ArticleSummary string `mapper:"article_summary"`
	//资讯模板
	NewsTags string `mapper:"NewsTags"`
}

//资讯详情
type ExpertNewsInfoDetail struct {
	PreviousNews *ExpertNewsInfo
	CurrentNews  *ExpertNewsInfo
	NextNews     *ExpertNewsInfo
}

type YaoWenHomePageData struct {
	//最新要闻信息
	NewsList []*ExpertNewsInfo

	//热门资讯
	HotNewsList []*ExpertNewsInfo
}
