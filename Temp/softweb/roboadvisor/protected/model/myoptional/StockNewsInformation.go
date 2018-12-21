package myoptional

import "github.com/devfeel/mapper"

//股票相关资讯
type StockNewsInformation struct {
	//主键
	StockNewsInformationId int64
	NewsInformationId      int64
	PublishTime            mapper.JSONTime `mapper:"publish_time"`
	SecurityCodeType       string          `mapper:"security_code_type"`
	SecurityCode           string          `mapper:"security_code"`
	//CreateTime             mapper.JSONTime

	//非数据库字段
	ArticleSummary string `mapper:"article_summary"`
	ArticleTitle   string `mapper:"article_title"`
	DataUrl        string `mapper:"data_url"`
	NewsTags       string `mapper:"NewsTags"`
	Id             string `mapper:"id"`
}
