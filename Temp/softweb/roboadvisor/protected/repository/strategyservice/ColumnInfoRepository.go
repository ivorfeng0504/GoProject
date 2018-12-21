package strategyservice

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	strategyservice_model "git.emoney.cn/softweb/roboadvisor/protected/model/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"strings"
	"time"
)

type ColumnInfoRepository struct {
	repository.BaseRepository
}

func NewColumnInfoRepository(conf *protected.ServiceConfig) *ColumnInfoRepository {
	repo := &ColumnInfoRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertColumnInfo 插入一条新的栏目
func (repo *ColumnInfoRepository) InsertColumnInfo(info strategyservice_model.ColumnInfo) (id int, err error) {
	sql := "INSERT INTO [ColumnInfo] ([ColumnName],[ColumnDesc],[AppID],[IsDeleted],[CreateTime]) VALUES(?,?,?,0,GETDATE())"
	id64, err := repo.Insert(sql, info.ColumnName, info.ColumnDesc, info.AppID)
	id = int(id64)
	return id, err
}

// UpdateColumnInfo 更新栏目
func (repo *ColumnInfoRepository) UpdateColumnInfo(columnName string, columnDesc string, id int) (err error) {
	sql := "UPDATE [ColumnInfo] SET [ColumnName]=?,[ColumnDesc]=? WHERE ID=?"
	_, err = repo.Update(sql, columnName, columnDesc, id)
	return err
}

// GetNewestNewsInfoByColumnList 获取栏目中最新的条资讯或视频
func (repo *ColumnInfoRepository) GetNewestNewsInfoByColumnList(columnList []int, newsType int) (newsInfo *model.NewsInfo, err error) {
	newsInfo = new(model.NewsInfo)
	containSql := strings.Repeat("?,", len(columnList))
	containSql = containSql[:len(containSql)-1]
	sql := `  SELECT TOP 1 * FROM [EMoney_RoboAdvisor].[dbo].[Column_News_Relation] AS TB1
  JOIN [EMoney_RoboAdvisor].[dbo].[NewsInfo] AS TB2 ON TB1.[NewsID]=TB2.[ID]
  WHERE TB1.[IsDeleted]=0 AND TB2.[IsDeleted]=0 AND TB1.[ColumnID] IN (` + containSql + `) AND [NewsType]=?
  ORDER BY TB2.[CreateTime] DESC`
	var params []interface{}
	for _, column := range columnList {
		params = append(params, column)
	}
	params = append(params, newsType)
	err = repo.FindOne(newsInfo, sql, params...)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsInfo, err
}

func (repo *ColumnInfoRepository) GetStrategyNewsListByPage_MultiMedia_ByDate(date time.Time) (newsList []*expertnews.ExpertNews_MultiMedia_List, err error) {
	sql := `SELECT a.*,b.AudioURL,b.VideoPlayURL,b.LiveURL,b.LiveVideoURL,b.Live_StartTime,b.Live_EndTime FROM dbo.[NewsInfo] a  WITH(NOLOCK)
LEFT JOIN  dbo.NewsInfo_MultiMedia b  WITH(NOLOCK)
ON a.ID=b.NewsID where a.newstype=3 AND a.isdeleted=0 AND 
(
(b.Live_StartTime IS NULL AND CONVERT(varchar(100),a.createtime, 23)=?)
OR
(b.Live_StartTime IS NOT NULL AND CONVERT(varchar(100),b.Live_StartTime, 23)=?)
)
order by a.CreateTime desc`
	err = repo.FindList(&newsList, sql, date.Format("2006-01-02"), date.Format("2006-01-02"))
	if err == sysSql.ErrNoRows {
		return newsList, nil
	}
	return newsList, err
}