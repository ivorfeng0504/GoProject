package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
)

type NewsInfoRepository struct {
	repository.BaseRepository
}

func NewNewsInfoRepository(conf *protected.ServiceConfig) *NewsInfoRepository {
	repo := &NewsInfoRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetNewsInfoByID 根据newsID获取资讯
func (repo *NewsInfoRepository) GetNewsInfoByID(newsid int) (*model.NewsInfo,error){
	sql := "SELECT * FROM [NewsInfo] WITH(NOLOCK) WHERE ID = ? AND IsDeleted=0"
	newsinfo := new(model.NewsInfo)
	err := repo.FindOne(newsinfo, sql, newsid)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsinfo, err
}

// GetNewsListByColumnID 根据ColumnID查询资讯内容
func (repo *NewsInfoRepository) GetNewsListByColumnID(ColumnID int) (newsList []*model.NewsInfo, err error) {
	sql := "SELECT b.*,a.IsTop FROM Column_News_Relation a LEFT JOIN NewsInfo b ON a.NewsID=b.ID  WHERE b.IsDeleted=0 AND a.IsDeleted=0 AND a.ColumnID=? ORDER BY [CreateTime] DESC"
	err = repo.FindList(&newsList, sql, ColumnID)
	return newsList, err
}

// GetSeriesNewsListByNewsID 根据主题（系列课程）id查询所有资讯课程
func (repo *NewsInfoRepository) GetSeriesNewsListByNewsID(TopicID int) (newsList []*model.NewsInfo,err error) {
	sql := "SELECT b.* FROM Topic_News_Relation a  LEFT JOIN NewsInfo b ON a.NewsID=b.ID  WHERE b.IsDeleted=0 AND a.IsDeleted=0 AND a.TopicID=? ORDER BY [CreateTime] DESC"
	err = repo.FindList(&newsList, sql, TopicID)
	return newsList, err
}

// GetNewsListByStrategyIDAndColID 根据策略id和栏目id查询所有资讯课程
func (repo *NewsInfoRepository) GetNewsListByStrategyIDAndColID(columnID int,strategyID int) (newsList []*model.NewsInfo,err error) {
	sql := `SELECT b.*,a.StrategyInfoID,c.IsTop FROM Strategy_News_Relation a
	LEFT JOIN NewsInfo b ON a.NewsID=b.ID
	LEFT JOIN Column_News_Relation c on b.id=c.NewsID
	WHERE b.IsDeleted=0 AND a.IsDeleted=0 AND c.IsDeleted=0 and a.StrategyInfoID=? and c.ColumnID=?`
	err = repo.FindList(&newsList, sql, strategyID,columnID)
	return newsList, err
}

// UpdateClickNum 更新点击数
func (repo *NewsInfoRepository) UpdateClickNum(newsId int64, clickNum int64) (err error) {
	sql := "UPDATE [NewsInfo] SET [ClickNum]=? WHERE [ID]=? and ([ClickNum] IS NULL OR [ClickNum]<?)"
	_, err = repo.Update(sql, clickNum, newsId, clickNum)
	return err
}

// UpdateVideoPlayClickNum 更新视频播放数
func (repo *NewsInfoRepository) UpdateVideoPlayNum(newsId int64, playNum int64) (err error) {
	sql := "UPDATE [NewsInfo] SET [VideoPlayNum]=? WHERE [ID]=? and ([VideoPlayNum] IS NULL OR [VideoPlayNum]<?)"
	_, err = repo.Update(sql, playNum, newsId, playNum)
	return err
}


// GetNewsListByClicknum 获取最新10条热门文章
func (repo *NewsInfoRepository) GetNewsListByClicknum(ColumnID int) (newsList []*model.NewsInfo, err error) {
	sql := "SELECT top 10 b.*,a.IsTop FROM Column_News_Relation a LEFT JOIN NewsInfo b ON a.NewsID=b.ID  WHERE b.IsDeleted=0 AND a.IsDeleted=0 AND a.ColumnID=? ORDER BY b.ClickNum DESC"
	err = repo.FindList(&newsList, sql, ColumnID)
	return newsList, err
}

// GetNewsListByClicknum 获取最新10条热门文章(策略学习使用),根据文章分类
func (repo *NewsInfoRepository) GetNewsListByClicknum_clxx(ColumnID int,topNum int,newsType int) (newsList []*model.NewsInfo, err error) {
	sql := "SELECT top " + strconv.Itoa(topNum) + " b.ID,b.Title,b.NewsType,a.IsTop,b.CreateTime,b.LastModifyTime,b.ClickNum FROM Column_News_Relation a LEFT JOIN NewsInfo b ON a.NewsID=b.ID  WHERE b.IsDeleted=0 AND a.IsDeleted=0 AND a.ColumnID=? AND b.NewsType=? AND DATEDIFF(dd,b.LastModifyTime,GETDATE())<=30 ORDER BY b.ClickNum DESC"
	err = repo.FindList(&newsList, sql, ColumnID, newsType)
	return newsList, err
}

// 获取最近一个月多媒体课程-不区分策略
func (repo *NewsInfoRepository) GetMultiMediaNewsList_1Month() (newsList []*expertnews.ExpertNews_MultiMedia_List,err error) {
	sqlstr := `SELECT a.*,b.AudioURL,b.VideoPlayURL,b.LiveURL,b.LiveVideoURL,b.Live_StartTime,b.Live_EndTime FROM dbo.[NewsInfo] a
	left join  dbo.NewsInfo_MultiMedia b
	on a.ID=b.NewsID where a.newstype=3 and a.isdeleted=0 and a.createtime>=DateAdd(MM,-1,getdate()) order by a.CreateTime desc`
	err = repo.FindList(&newsList, sqlstr)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsList, err
}





