package dataapi

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	expertnews_srv "git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/cache"
	"github.com/devfeel/dotlog"
	"strings"
	"time"
)

const (
	newsApiServiceName               = "NewsApiService"
	EMNET_NewsApiService_PreCahceKey = "NewsApiService:"
	//EMNET_NewsApiService_SyncRowVersion_Day_CacheKey = EMNET_NewsApiService_PreCahceKey + "SyncRowVersion:Day:"
	//EMNET_NewsApiService_SyncRowVersion_CacheKey = EMNET_NewsApiService_PreCahceKey + "SyncRowVersion:"
)

var (
	shareNewsApiServiceLogger dotlog.Logger
	shareNewsApiServiceRedis  cache.RedisCache
)

func init() {
	protected.RegisterServiceLoader(newsApiServiceName, func() {
		shareNewsApiServiceLogger = dotlog.GetLogger(newsApiServiceName)
		shareNewsApiServiceRedis = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	})
}

var (
	TableFields []string
)

func init() {
	//TableFields = []string{"id", "image_id", "article_title", "article_short_title", "article_regular_title", "article_from_url", "article_from_website_id", "article_from_website_name", "article_source", "article_source_channel", "article_author", "article_channel_id", "article_channel_name", "article_time", "article_tags_id", "article_tags_name", "article_type_id", "article_type_name", "article_category_id", "article_category_name", "article_keywords", "article_summary", "article_property", "article_comments", "article_rating_id", "article_rating_name", "importance_ratings", "cyclical_id", "cyclical_name", "data_source_id", "data_source_name", "block_code", "block_name", "block_type", "security_code", "security_name", "theme_id", "theme_name", "path_locator", "priority", "source_news_id", "publish_time", "status", "is_deliver", "is_checked", "is_valid", "is_broadcast", "is_sync_file", "creator", "editor", "author", "comments", "create_time", "modify_time", "SyncRowVersion", "is_top", "bespeak_status", "bespeak_time", "upload_cloud_status", "NewsTags", "NewsTagsID","data_url"}
	TableFields = []string{"id", "image_id", "article_title", "article_short_title", "article_regular_title", "article_from_url", "article_from_website_id", "article_from_website_name", "article_source", "article_source_channel", "article_author", "article_channel_id", "article_channel_name", "article_time", "article_tags_id", "article_tags_name", "article_type_id", "article_type_name", "article_category_id", "article_category_name", "article_keywords", "article_summary", "article_property", "article_comments", "importance_ratings", "block_code", "block_name", "block_type", "security_code", "security_name", "theme_id", "theme_name", "source_news_id", "publish_time", "status", "is_deliver", "is_checked", "is_valid", "author", "comments", "create_time", "modify_time", "SyncRowVersion", "is_top", "NewsTags", "NewsTagsID", "data_url", "head_lines"}
}

// UpdateNewsInformation 更新咨询信息
func UpdateNewsInformation() (newsList []*expertnews.NewsInformation, err error) {
	return UpdateNewsInformationByDate(nil)
}

// UpdateNewsInformation 更新指定日期的咨询数据,如果date为nil，则不筛选日期
func UpdateNewsInformationByDate(date *time.Time) (newsList []*expertnews.NewsInformation, err error) {
	filter := ""
	newsSrv := expertnews_srv.NewNewsInformationService()
	if date != nil {
		//根据日期筛选
		filter = fmt.Sprintf("%s>%s", _const.NewsInfo_Field_Publish_Time, "%27"+date.Format("2006-01-02")+"%27")
		var syncRowVersion string
		syncRowVersion, versionErr := newsSrv.GetMaxSyncRowVersion(date)
		//如果查询版本号出错 则直接返回，防止过滤条件为空，导致接口数据查询过大
		if versionErr != nil {
			return newsList, versionErr
		}
		if versionErr == nil && len(syncRowVersion) > 0 {
			filter = fmt.Sprintf("%s%s%s>%s", filter, "%20AND%20", _const.NewsInfo_Field_SyncRowVersion, "%27"+syncRowVersion+"%27")
		}
	} else {
		//不根据日期筛选
		var syncRowVersion string
		syncRowVersion, versionErr := newsSrv.GetMaxSyncRowVersion(date)
		//如果查询版本号出错 则直接返回，防止过滤条件为空，导致接口数据查询过大
		if versionErr != nil {
			return newsList, versionErr
		}
		if versionErr == nil && len(syncRowVersion) > 0 {
			filter = fmt.Sprintf("%s>%s", _const.NewsInfo_Field_SyncRowVersion, "%27"+syncRowVersion+"%27")
		}
	}

	apiUrl := config.CurrentConfig.NewsInformationApi
	if len(apiUrl) == 0 {
		shareNewsApiServiceLogger.ErrorFormat(err, "获取资讯信息接口地址配置不正确 configkey=NewsInformationApi")
		return newsList, errors.New("获取资讯信息接口地址配置不正确")
	}
	apiUrl = fmt.Sprintf(apiUrl, filter, strings.Join(TableFields, ","))
	result, err := GetDataApi(apiUrl)
	if err != nil {
		shareNewsApiServiceLogger.ErrorFormat(err, "UpdateNewsInformationByDate 获取资讯信息异常 请求地址为%s", apiUrl)
		return nil, err
	}
	if result == nil || len(result) <= 2 {
		return nil, nil
	}
	//所有行数据
	rows := result[2:]
	for _, row := range rows {
		//写入数据库
		err = newsSrv.InsertOrUpdate(row[0], TableFields, row)
		if err != nil {
			//写入日志
			shareNewsApiServiceLogger.ErrorFormat(err, "InsertOrUpdate 写入或更新数据库异常  params=%s  values=%s", _json.GetJsonString(TableFields), _json.GetJsonString(row))
			newsSrv.AddInsertErrorQueue(TableFields, row)
		}
	}
	return newsList, err
}

// GetNewsTemplate 获取资讯模板信息
func GetNewsTemplate() (templates []string, err error) {
	apiUrl := config.CurrentConfig.NewsInfomationTemplateApi
	if len(apiUrl) == 0 {
		shareNewsApiServiceLogger.ErrorFormat(err, "获取资讯模板信息接口地址配置不正确 configkey=NewsInfomationTemplateApi")
		return templates, errors.New("获取资讯模板信息接口地址配置不正确")
	}
	result, err := GetDataApi(apiUrl)
	if err != nil {
		shareNewsApiServiceLogger.ErrorFormat(err, "GetNewsTemplate 获取资讯模板信息-> GetDataApi 异常 请求地址为%s", apiUrl)
		return templates, err
	}
	if result == nil || len(result) <= 2 {
		return templates, nil
	}
	//所有行数据
	rows := result[2:]
	for _, row := range rows {
		templates = append(templates, row[0])
	}
	return templates, nil
}
