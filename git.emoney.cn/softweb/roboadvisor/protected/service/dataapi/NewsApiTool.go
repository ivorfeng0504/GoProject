package dataapi

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"strings"
)

func main() {
	BuildNewsStruct()
}

func BuildNewsStruct() {
	typeMaps := map[string]string{"STRING": "string", "DATETIME": "mapper.JSONTime", "BOOL": "bool", "INT64": "int64", "INT32": "int32", "DOUBLE": "string", "UINT8": "int"}
	json := `["id","image_id","article_title","article_short_title","article_regular_title","article_from_url","article_from_website_id","article_from_website_name","article_source","article_source_channel","article_author","article_channel_id","article_channel_name","article_time","article_tags_id","article_tags_name","article_type_id","article_type_name","article_category_id","article_category_name","article_keywords","article_summary","article_property","article_comments","article_rating_id","article_rating_name","importance_ratings","cyclical_id","cyclical_name","data_source_id","data_source_name","block_code","block_name","block_type","security_code","security_name","theme_id","theme_name","path_locator","priority","source_news_id","publish_time","status","is_deliver","is_checked","is_valid","is_broadcast","is_sync_file","creator","editor","author","comments","create_time","modify_time","SyncRowVersion","is_top","bespeak_status","bespeak_time","upload_cloud_status","NewsTags","NewsTagsID"]`
	typeJson := `["DOUBLE","DOUBLE","STRING","STRING","STRING","STRING","DOUBLE","STRING","STRING","STRING","STRING","DOUBLE","STRING","DATETIME","DOUBLE","STRING","DOUBLE","STRING","DOUBLE","STRING","STRING","STRING","STRING","STRING","DOUBLE","STRING","UINT8","DOUBLE","STRING","DOUBLE","STRING","STRING","STRING","STRING","STRING","STRING","STRING","STRING","STRING","STRING","STRING","DATETIME","UINT8","BOOL","UINT8","BOOL","BOOL","BOOL","STRING","STRING","STRING","STRING","DATETIME","DATETIME","INT64","BOOL","INT32","DATETIME","BOOL","STRING","DOUBLE"]`
	var fields []string
	var types []string
	_json.Unmarshal(json, &fields)
	_json.Unmarshal(typeJson, &types)
	goFile := ""
	for index, field := range fields {
		items := strings.Split(field, "_")
		if len(items) == 0 {
			continue
		}
		varName := ""
		for _, item := range items {
			varName += strings.ToUpper(item[0:1]) + item[1:]
		}
		typeName := typeMaps[types[index]]
		goFile += fmt.Sprintf("%s %s `mapper:\"%s\"` \r\n", varName, typeName, field)
	}
	fmt.Println(goFile)
}

func BuildNewsConst() {
	json := `["id","image_id","article_title","article_short_title","article_regular_title","article_from_url","article_from_website_id","article_from_website_name","article_source","article_source_channel","article_author","article_channel_id","article_channel_name","article_time","article_tags_id","article_tags_name","article_type_id","article_type_name","article_category_id","article_category_name","article_keywords","article_summary","article_property","article_comments","article_rating_id","article_rating_name","importance_ratings","cyclical_id","cyclical_name","data_source_id","data_source_name","block_code","block_name","block_type","security_code","security_name","theme_id","theme_name","path_locator","priority","source_news_id","publish_time","status","is_deliver","is_checked","is_valid","is_broadcast","is_sync_file","creator","editor","author","comments","create_time","modify_time","SyncRowVersion","is_top","bespeak_status","bespeak_time","upload_cloud_status","NewsTags","NewsTagsID"]`
	var fields []string
	_json.Unmarshal(json, &fields)
	goFile := ""
	for _, field := range fields {
		items := strings.Split(field, "_")
		if len(items) == 0 {
			continue
		}
		varName := "NewsInfo_Field"
		for _, item := range items {
			varName += "_" + strings.ToUpper(item[0:1]) + item[1:]
		}
		varValue := field
		goFile += fmt.Sprintf(`%s="%s"`+"\r\n", varName, varValue)
	}
	fmt.Println(goFile)
}
