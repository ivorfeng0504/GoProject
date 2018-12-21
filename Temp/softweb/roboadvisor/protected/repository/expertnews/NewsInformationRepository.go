package expertnews

import (
	sysSql "database/sql"
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected"
	expertnews_model "git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"strings"
	"time"
)

type NewsInformationRepository struct {
	repository.BaseRepository
}

func NewNewsInformationRepository(conf *protected.ServiceConfig) *NewsInformationRepository {
	repo := &NewsInformationRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// Exist 指定的咨询是否存在
func (repo *NewsInformationRepository) Exist(id string) (isExist bool, err error) {
	sql := "SELECT COUNT(1) FROM [NewsInformation] WHERE [id]=?"
	count, err := repo.Count(sql, id)
	return count > 0, err
}

// InsertOrUpdate 插入或更新资讯数据
func (repo *NewsInformationRepository) InsertOrUpdate(id string, params []string, values []string) (err error) {
	if params == nil || values == nil {
		err = errors.New("参数不能为空！")
		return err
	}
	if len(params) != len(values) {
		err = errors.New("参数的个数与值的个数不相符！")
		return err
	}

	isExist, err := repo.Exist(id)
	if err != nil {
		return err
	}
	//带引号的value
	quoteValues := repo.processValues(values)

	//插入更新时间
	params = append(params, _const.NewsInfo_Field_SyncModifyTime)
	values = append(values, "GETDATE()")
	quoteValues = append(quoteValues, "GETDATE()")

	if isExist == false {
		//新增
		//插入创建时间
		params = append(params, _const.NewsInfo_Field_SyncCreateTime)
		values = append(values, "GETDATE()")
		quoteValues = append(quoteValues, "GETDATE()")

		sql := "INSERT INTO [NewsInformation] (%s) VALUES(%s)"
		sql = fmt.Sprintf(sql, strings.Join(params, ","), strings.Join(quoteValues, ","))
		_, err = repo.Insert(sql)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "InsertOrUpdate 数据库新增异常 sql=%s", sql)
		}
	} else {
		//更新
		sql := "UPDATE [NewsInformation] SET %s WHERE [id]=? AND [SyncRowVersion]<?"
		var updateArr []string
		var rowVersion string
		for index, field := range params {
			if field != _const.NewsInfo_Field_Id {
				updateArr = append(updateArr, fmt.Sprintf("%s=%s", field, quoteValues[index]))
			}
			if field == _const.NewsInfo_Field_SyncRowVersion {
				rowVersion = values[index]
			}
		}
		sql = fmt.Sprintf(sql, strings.Join(updateArr, ","))
		_, err = repo.Update(sql, id, rowVersion)
		if err != nil {
			global.TaskLogger.ErrorFormat(err, "InsertOrUpdate 数据库更新异常 sql=%s id=%s rowVersion=%s", sql, id, rowVersion)
		}
	}
	return err
}

func (repo *NewsInformationRepository) processValues(values []string) (result []string) {
	if values == nil {
		return values
	}
	for _, item := range values {
		//将单引号替换为两个单引号 防止数据插入错误
		item = strings.Replace(item, "'", "''", -1)
		result = append(result, fmt.Sprintf("'%s'", item))
	}
	return result
}

// GetMaxSyncRowVersion 获取指定日期里最新版本号 不传递date则取当前最大版本号
func (repo *NewsInformationRepository) GetMaxSyncRowVersion(date *time.Time) (version string, err error) {
	sql := "SELECT MAX([SyncRowVersion]) FROM [NewsInformation]"
	if date != nil {
		sql = fmt.Sprintf("%s WHERE CONVERT(varchar(100),[publish_time], 23)='%s'", sql, date.Format("2006-01-02"))
	}
	versionObj, err := repo.QueryMax(sql)
	if err == sysSql.ErrNoRows || versionObj == nil {
		return "", nil
	}
	if err == nil {
		version = versionObj.(string)
	}
	return version, err
}

// UpdateClickNum 更新点击数
// 如果传递的clickNum小于数据库中的值则不更新 防止错误数据覆盖数据的值
func (repo *NewsInformationRepository) UpdateClickNum(newsId int64, clickNum int64) (err error) {
	sql := "UPDATE [NewsInformation] SET [ClickNum]=? WHERE [NewsInformationId]=? AND [ClickNum]<?"
	_, err = repo.Update(sql, clickNum, newsId, clickNum)
	return err
}

// UpdateStockNewsClickNum 更新点击数
// 如果传递的clickNum小于数据库中的值则不更新 防止错误数据覆盖数据的值
func (repo *NewsInformationRepository) UpdateStockNewsClickNum(newsId int64, clickNum int64) (err error) {
	sql := "UPDATE [NewsInformation] SET [StockNewsClickNum]=? WHERE [NewsInformationId]=? AND [StockNewsClickNum]<?"
	_, err = repo.Update(sql, clickNum, newsId, clickNum)
	return err
}

// GetHasNewsDate 获取指定日期之前有资讯的最近一个日期
func (repo *NewsInformationRepository) GetLastHasNewsDate(date time.Time) (lastDate time.Time, err error) {
	sql := "SELECT TOP 1 * FROM [NewsInformation] WHERE  CONVERT(varchar(100),[publish_time], 23)<? ORDER BY [publish_time] DESC"
	record := new(expertnews_model.NewsInformation)
	err = repo.FindOne(record, sql, date.Format("2006-01-02"))
	if err == nil {
		lastDate = time.Time(record.PublishTime)
	}
	return lastDate, err
}

// GetTodayNews 获取今日头条
// 规则：不需要根据标签筛选
// 第1条，尝试取重要度5的资讯，如果不存在则取is_top=1 is_deliver=1 重要度为4的一条
// 		----如果没有取到第1条数据，则第1,2,3条取最新的 is_deliver=1 重要度为4的资讯
//		----如果取到了第1条，则第2条取is_top=1  is_deliver=1 重要度为4的一条资讯，第3条取最新的一条 is_deliver=1 重要度为4的资讯
//					----如果第2条没有取到，则第2,3条取最新的 is_deliver=1 重要度为4的资讯
//							----如果上述规则下取到的数据不足3条，则递归取前一天的数据，在相同的规则下补充不足的数据
//							----如果仍不满足则再次递归，直到取到3条为止（限制递归最大深度为10，防止无限递归）
func (repo *NewsInformationRepository) GetTodayNews() (newsList []*expertnews_model.NewsInformation, err error) {
	date := time.Now()
	//保存文章Id，防止返回的文章重复
	newidMap := make(map[string]bool)
	//最大尝试10次
	tryCountMax := 10
	for tryCountMax > 0 {
		//取出指定日期的头条数据
		newsListTmp, err := repo.GetTodayNewsByDate(date)
		if err != nil {
			return newsList, err
		}
		//如果头条数据不为空 则附加到结果集中
		if newsListTmp != nil && len(newsListTmp) > 0 {
			for _, newsInfo := range newsListTmp {
				_, exist := newidMap[newsInfo.Id]
				if exist == false {
					newsList = append(newsList, newsInfo)
					newidMap[newsInfo.Id] = true
				}
			}
		}
		//如果结果集中的数据大于等于3条 则直接取前三条返回
		//否则递归取前一天的数据进行补充
		if len(newsList) >= 3 {
			//只取前3条
			newsList = newsList[:3]
			return newsList, err
		}
		tryCountMax--
		//获取上一个有资讯的日期
		date, err = repo.GetLastHasNewsDate(date)
		if err != nil {
			return newsList, err
		}
	}
	return newsList, err
}

// GetTodayNewsByDate 获取指定日期的头条
// 规则：不需要根据标签筛选
// 第1条，尝试取重要度5的资讯，如果不存在则取is_top=1 is_deliver=1 重要度为4的一条
// 		----如果没有取到第1条数据，则第1,2,3条取最新的 is_deliver=1 重要度为4的资讯
//		----如果取到了第1条，则第2条取is_top=1  is_deliver=1 重要度为4的一条资讯，第3条取最新的一条 is_deliver=1 重要度为4的资讯
//					----如果第2条没有取到，则第2,3条取最新的 is_deliver=1 重要度为4的资讯
func (repo *NewsInformationRepository) GetTodayNewsByDate(date time.Time) (newsList []*expertnews_model.NewsInformation, err error) {
	timeStr := date.Format("2006-01-02")
	// 第1条，尝试取重要度5的资讯，如果不存在则取is_top=1 is_deliver=1 重要度为4的一条
	firstLineSql := "SELECT TOP 1 * FROM [NewsInformation] WHERE CONVERT(varchar(100),[publish_time], 23)=? AND [is_checked]=1 AND [is_valid]=1 AND [article_type_id]='100000110' AND [article_category_id] IN ('100120059','101008060') AND (([importance_ratings]=4 AND [is_top]=1 AND [is_deliver]=1) OR ([importance_ratings]=5)) ORDER BY [importance_ratings] DESC,[publish_time] DESC"
	firstLine := new(expertnews_model.NewsInformation)
	//查询第一条数据
	err = repo.FindOne(firstLine, firstLineSql, timeStr)
	if err == sysSql.ErrNoRows {
		firstLine = nil
		err = nil
	}
	if err != nil {
		return nil, err
	}

	//与第一条去重
	filterFirstLine := ""
	if firstLine != nil {
		filterFirstLine = fmt.Sprintf(" AND [id]!= '%s' ", firstLine.Id)
		//将第一条添加到结果中
		newsList = append(newsList, firstLine)
	}

	otherLineSql := ""
	var otherLine []*expertnews_model.NewsInformation
	if firstLine == nil {
		//如果没有取到第1条数据，则第1,2,3条取最新的 is_deliver=1 重要度为4的资讯
		otherLineSql = "SELECT TOP 3 * FROM [NewsInformation] WHERE CONVERT(varchar(100),[publish_time], 23)=? AND [is_deliver]=1 AND [is_checked]=1 AND [is_valid]=1 AND [article_type_id]='100000110' AND [article_category_id] IN ('100120059','101008060') AND [importance_ratings]=4 ORDER BY [publish_time] DESC"
		err = repo.FindList(&otherLine, otherLineSql, timeStr)
	} else {
		//如果取到了第1条，则第2条取is_top=1  is_deliver=1 重要度为4的一条资讯，第3条取最新的一条 is_deliver=1 重要度为4的资讯
		//如果第2条没有取到，则第2,3条取最新的 is_deliver=1 重要度为4的资讯
		otherLineSql = "(SELECT *,'1' AS [InnerLevel] FROM (SELECT TOP 1 * FROM [NewsInformation] WHERE CONVERT(varchar(100),[publish_time], 23)=? AND [is_top]=1 AND [is_deliver]=1 AND [is_checked]=1 AND [is_valid]=1 AND [article_type_id]='100000110' AND [article_category_id] IN ('100120059','101008060') AND [importance_ratings]=4 %s ORDER BY [publish_time] DESC) AS TB1"
		otherLineSql = fmt.Sprintf(otherLineSql, filterFirstLine)
		otherLineSql += " UNION "
		otherLineSql += "SELECT *,'0' AS [InnerLevel] FROM (SELECT TOP 2 * FROM [NewsInformation] WHERE CONVERT(varchar(100),[publish_time], 23)=? AND [is_deliver]=1 AND [is_checked]=1 AND [is_valid]=1 AND [article_type_id]='100000110' AND [article_category_id] IN ('100120059','101008060') AND [importance_ratings]=4 %s ORDER BY [publish_time] DESC) AS TB2)"
		otherLineSql = fmt.Sprintf(otherLineSql, filterFirstLine)
		otherLineSql += " ORDER BY [InnerLevel] DESC,[publish_time] DESC"
		err = repo.FindList(&otherLine, otherLineSql, timeStr, timeStr)
	}

	//如果出错直接返回
	if err != nil && err != sysSql.ErrNoRows {
		return nil, err
	}

	//如果otherLine第一条和第二条重复 则删除第一条
	if len(otherLine) >= 3 && otherLine[0].Id == otherLine[1].Id {
		otherLine = otherLine[1:]
	}

	newsList = append(newsList, otherLine...)

	//如果获取到的数据满足3条 则直接返回
	if len(newsList) >= 3 {
		newsList = newsList[:3]
		return newsList, nil
	}
	return newsList, nil
}

// GetTodayNewsV2 获取今日头条
func (repo *NewsInformationRepository) GetTodayNewsV2() (newsList []*expertnews_model.NewsInformation, err error) {
	firstSql := "SELECT TOP 1 * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 AND [head_lines]='头条' ORDER BY [publish_time] DESC"
	secondSql := "SELECT TOP 1 * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 AND [head_lines]='次头条' ORDER BY [publish_time] DESC"
	thirdSql := "SELECT TOP 1 * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 AND [head_lines]='第三条' ORDER BY [publish_time] DESC"
	firstModel := new(expertnews_model.NewsInformation)
	secondModel := new(expertnews_model.NewsInformation)
	thirdModel := new(expertnews_model.NewsInformation)
	err = repo.FindOne(firstModel, firstSql)
	if err != nil && err != sysSql.ErrNoRows {
		return nil, err
	}
	err = repo.FindOne(secondModel, secondSql)
	if err != nil && err != sysSql.ErrNoRows {
		return nil, err
	}
	err = repo.FindOne(thirdModel, thirdSql)
	if err != nil && err != sysSql.ErrNoRows {
		return nil, err
	}
	newsList = append(newsList, firstModel)
	newsList = append(newsList, secondModel)
	newsList = append(newsList, thirdModel)
	return newsList, err
}

// GetClosingNews 获取盘后预测资讯
// 规则：包含资讯模板标签的 置顶的最新15条数据
func (repo *NewsInformationRepository) GetClosingNews(templates []string) (newsList []*expertnews_model.NewsInformation, err error) {
	sql := "SELECT TOP 15 * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 AND [article_tags_name]='推荐' AND  [is_top]=1 %s ORDER BY [publish_time] DESC"
	templateSql := repo.GetTemplateSql(templates)
	sql = fmt.Sprintf(sql, templateSql)
	err = repo.FindList(&newsList, sql)
	return newsList, err
}

// GetNewsInfo 获取指定日期的要闻
// 规则：包含资讯模板标签的 指定日期的最新15条数据
func (repo *NewsInformationRepository) GetNewsInfo(templates []string, date time.Time) (newsList []*expertnews_model.NewsInformation, err error) {
	sql := "SELECT TOP 15 * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 AND [article_tags_name]='推荐' AND  CONVERT(varchar(100),[publish_time], 23)=? %s ORDER BY [publish_time] DESC"
	templateSql := repo.GetTemplateSql(templates)
	sql = fmt.Sprintf(sql, templateSql)
	err = repo.FindList(&newsList, sql, date.Format("2006-01-02"))
	return newsList, err
}

// GetNewsInfoTopN 获取指定Top N的要闻
// 规则：包含资讯模板标签的 最新N条数据
func (repo *NewsInformationRepository) GetNewsInfoTopN(templates []string, top int) (newsList []*expertnews_model.NewsInformation, err error) {
	sql := "SELECT TOP %d * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 AND [article_tags_name]='推荐' %s ORDER BY [publish_time] DESC"
	templateSql := repo.GetTemplateSql(templates)
	sql = fmt.Sprintf(sql, top, templateSql)
	err = repo.FindList(&newsList, sql)
	return newsList, err
}

// GetHotNewsInfo 获取热门资讯
// 规则：包含资讯模板标签的 近一周热度最高的15条数据
func (repo *NewsInformationRepository) GetHotNewsInfo(templates []string) (newsList []*expertnews_model.NewsInformation, err error) {
	sql := "SELECT * FROM (SELECT TOP 15 * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 AND [article_tags_name]='推荐' AND  [publish_time]>? %s ORDER BY [ClickNum] DESC ,[publish_time] DESC ) AS TB1 ORDER BY [publish_time] DESC"
	date := time.Now().AddDate(0, 0, -7)
	templateSql := repo.GetTemplateSql(templates)
	sql = fmt.Sprintf(sql, templateSql)
	err = repo.FindList(&newsList, sql, date.Format("2006-01-02"))
	return newsList, err
}

// GetTopicNewsInfo 获取主题的相关资讯
// 规则：与指定板块关联的的15条数据
func (repo *NewsInformationRepository) GetTopicNewsInfo(bkList []*expertnews_model.Topic_BK) (newsList []*expertnews_model.NewsInformation, err error) {
	sql := "SELECT TOP 15 * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 %s ORDER BY [publish_time] DESC"
	bkSQL := repo.GetBKSql(bkList)
	sql = fmt.Sprintf(sql, bkSQL)
	err = repo.FindList(&newsList, sql)
	return newsList, err
}

// GetTemplateSql 获取资讯模板过滤语句
func (repo *NewsInformationRepository) GetTemplateSql(templates []string) string {
	if templates == nil {
		return ""
	}
	sql := " AND [NewsTags] IN(%s) "
	var quoteTemplateList []string
	for _, template := range templates {
		quoteTemplateList = append(quoteTemplateList, fmt.Sprintf("'%s'", template))
	}
	if len(quoteTemplateList) == 0 {
		return ""
	}
	sql = fmt.Sprintf(sql, strings.Join(quoteTemplateList, ","))
	return sql
}

// GetBKSql 获取板块过滤语句
func (repo *NewsInformationRepository) GetBKSql(bkList []*expertnews_model.Topic_BK) string {
	if bkList == nil {
		return ""
	}
	sql := " AND (%s)"
	var bkSqlList []string
	for _, bk := range bkList {
		bkCode := bk.BKCode
		if len(bkCode) > 4 {
			bkCode = bkCode[len(bkCode)-4:]
		}
		bkSqlList = append(bkSqlList, fmt.Sprintf(" [block_code] LIKE '%%%s%%' ", bkCode))
	}
	if len(bkSqlList) == 0 {
		return ""
	}
	sql = fmt.Sprintf(sql, strings.Join(bkSqlList, " OR "))
	return sql
}

// GetNewsInfoById 根据Id获取指定的资讯
func (repo *NewsInformationRepository) GetNewsInfoById(newsInfoId int64) (newsInfo *expertnews_model.NewsInformation, err error) {
	sql := "SELECT TOP 1 * FROM [NewsInformation] WHERE [is_checked]=1 AND [is_valid]=1 AND [NewsInformationId]=?"
	newsInfo = new(expertnews_model.NewsInformation)
	err = repo.FindOne(newsInfo, sql, newsInfoId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsInfo, err
}

// GetStockNewsInfoList 获取大于指定NewsInformationId的所有股票相关的数据集合
func (repo *NewsInformationRepository) GetStockNewsInfoList(newsInformationId int64) (newsList []*expertnews_model.NewsInformation, err error) {
	sql := "SELECT * FROM [NewsInformation] WHERE [NewsInformationId]>? AND LEN([security_code])>0"
	err = repo.FindList(&newsList, sql, newsInformationId)
	return newsList, err
}

// GetBlockNewsInfoListAfterVersion 获取大于指定版本的板块资讯列表
func (repo *NewsInformationRepository) GetBlockNewsInfoListAfterVersion(version string) (newsList []*expertnews_model.NewsInformation, err error) {
	sql := "SELECT * FROM [NewsInformation] WHERE [SyncRowVersion]>? AND LEN([block_code])>0 AND [is_checked]=1 AND [is_valid]=1 ORDER BY [SyncRowVersion] ASC"
	err = repo.FindList(&newsList, sql, version)
	return newsList, err
}
