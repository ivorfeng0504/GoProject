package expertnews

import "github.com/devfeel/mapper"

type NewsInformation struct {
	Id                     string          `mapper:"id"`
	ImageId                string          `mapper:"image_id"`
	ArticleTitle           string          `mapper:"article_title"`
	ArticleShortTitle      string          `mapper:"article_short_title"`
	ArticleRegularTitle    string          `mapper:"article_regular_title"`
	ArticleFromUrl         string          `mapper:"article_from_url"`
	ArticleFromWebsiteId   string          `mapper:"article_from_website_id"`
	ArticleFromWebsiteName string          `mapper:"article_from_website_name"`
	ArticleSource          string          `mapper:"article_source"`
	ArticleSourceChannel   string          `mapper:"article_source_channel"`
	ArticleAuthor          string          `mapper:"article_author"`
	ArticleChannelId       string          `mapper:"article_channel_id"`
	ArticleChannelName     string          `mapper:"article_channel_name"`
	ArticleTime            mapper.JSONTime `mapper:"article_time"`
	ArticleTagsId          string          `mapper:"article_tags_id"`
	ArticleTagsName        string          `mapper:"article_tags_name"`
	ArticleTypeId          string          `mapper:"article_type_id"`
	ArticleTypeName        string          `mapper:"article_type_name"`
	ArticleCategoryId      string          `mapper:"article_category_id"`
	ArticleCategoryName    string          `mapper:"article_category_name"`
	ArticleKeywords        string          `mapper:"article_keywords"`
	ArticleSummary         string          `mapper:"article_summary"`
	ArticleProperty        string          `mapper:"article_property"`
	ArticleComments        string          `mapper:"article_comments"`
	ArticleRatingId        string          `mapper:"article_rating_id"`
	ArticleRatingName      string          `mapper:"article_rating_name"`
	ImportanceRatings      int             `mapper:"importance_ratings"`
	CyclicalId             string          `mapper:"cyclical_id"`
	CyclicalName           string          `mapper:"cyclical_name"`
	DataSourceId           string          `mapper:"data_source_id"`
	DataSourceName         string          `mapper:"data_source_name"`
	BlockCode              string          `mapper:"block_code"`
	BlockName              string          `mapper:"block_name"`
	BlockType              string          `mapper:"block_type"`
	SecurityCode           string          `mapper:"security_code"`
	SecurityName           string          `mapper:"security_name"`
	ThemeId                string          `mapper:"theme_id"`
	ThemeName              string          `mapper:"theme_name"`
	PathLocator            string          `mapper:"path_locator"`
	Priority               string          `mapper:"priority"`
	SourceNewsId           string          `mapper:"source_news_id"`
	PublishTime            mapper.JSONTime `mapper:"publish_time"`
	Status                 int             `mapper:"status"`
	IsDeliver              bool            `mapper:"is_deliver"`
	IsChecked              int             `mapper:"is_checked"`
	IsValid                bool            `mapper:"is_valid"`
	IsBroadcast            bool            `mapper:"is_broadcast"`
	IsSyncFile             bool            `mapper:"is_sync_file"`
	Creator                string          `mapper:"creator"`
	Editor                 string          `mapper:"editor"`
	Author                 string          `mapper:"author"`
	Comments               string          `mapper:"comments"`
	CreateTime             mapper.JSONTime `mapper:"create_time"`
	ModifyTime             mapper.JSONTime `mapper:"modify_time"`
	SyncRowVersion         string          `mapper:"SyncRowVersion"`
	IsTop                  bool            `mapper:"is_top"`
	BespeakStatus          int32           `mapper:"bespeak_status"`
	BespeakTime            mapper.JSONTime `mapper:"bespeak_time"`
	UploadCloudStatus      bool            `mapper:"upload_cloud_status"`
	NewsTags               string          `mapper:"NewsTags"`
	NewsTagsID             string          `mapper:"NewsTagsID"`
	DataUrl                string          `mapper:"data_url"`
	HeadLines              string          `mapper:"head_lines"`
	//点击数
	ClickNum int64
	//主键Id
	NewsInformationId int64

	//数据库同步创建时间
	SyncCreateTime mapper.JSONTime `mapper:"Sync_CreateTime"`
	//数据库同步最后更新时间
	SyncModifyTime mapper.JSONTime `mapper:"Sync_ModifyTime"`
}
