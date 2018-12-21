package myoptional

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"strconv"
	"strings"
	"time"
)

type StockNewsInformationRepository struct {
	repository.BaseRepository
}

func NewStockNewsInformationRepository(conf *protected.ServiceConfig) *StockNewsInformationRepository {
	repo := &StockNewsInformationRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetStockNewsInfo 获取股票相关资讯
func (repo *StockNewsInformationRepository) GetStockNewsInfo(startTime string, stockList []string, top int) (newsList []*myoptional_model.StockNewsInformation, err error) {
	if stockList == nil || len(stockList) == 0 {
		return nil, nil
	}
	stockListStr := strings.Repeat("?,", len(stockList))
	stockListStr = stockListStr[:len(stockListStr)-1]
	sql := "SELECT TOP " + strconv.Itoa(top) + ` TB1.[StockNewsInformationId],TB1.[NewsInformationId],TB1.[publish_time],TB2.[security_code],TB2.[article_summary],TB2.[article_title],TB2.[data_url],TB2.[NewsTags],TB2.[id] FROM [StockNewsInformation] AS TB1
    JOIN [NewsInformation] AS TB2 ON TB1.[NewsInformationId]=TB2.[NewsInformationId]
    WHERE TB2.[is_checked]=1 AND TB2.[is_valid]=1 AND TB1.[publish_time]>? AND LEN(TB1.[security_code])>0 AND TB1.[security_code] IN (` + stockListStr + `) AND TB1.[security_code_type]=? ORDER BY TB1.[publish_time] DESC`
	var args []interface{}
	args = append(args, startTime)
	for _, stock := range stockList {
		args = append(args, stock)
	}
	args = append(args, _const.StockType_A)
	err = repo.FindList(&newsList, sql, args...)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsList, err
}

// InsertStockNewsInfo 插入一个新的股票资讯信息
func (repo *StockNewsInformationRepository) InsertStockNewsInfo(model myoptional_model.StockNewsInformation) (err error) {
	sql := "INSERT INTO [StockNewsInformation] ([NewsInformationId],[publish_time],[security_code],[security_code_type],[CreateTime]) VALUES(?,?,?,?,GETDATE())"
	_, err = repo.Insert(sql, model.NewsInformationId, time.Time(model.PublishTime), model.SecurityCode, model.SecurityCodeType)
	return err
}

// GetMaxNewsInformationId 获取关系表中最大的NewsInformationId
func (repo *StockNewsInformationRepository) GetMaxNewsInformationId() (newsInformationId int64, err error) {
	sql := "SELECT MAX([NewsInformationId]) FROM [StockNewsInformation]"
	newsInformationIdObj, err := repo.QueryMax(sql)
	if err == sysSql.ErrNoRows || newsInformationIdObj == nil {
		return newsInformationId, nil
	}
	if err == nil {
		newsInformationId = newsInformationIdObj.(int64)
	}
	return newsInformationId, err
}

// GetStockCodeList 获取包含资讯的所有股票代码-A股
func (repo *StockNewsInformationRepository) GetStockCodeList() (stockCodeList []string, err error) {
	sql := "SELECT [security_code] FROM [StockNewsInformation] WHERE [security_code_type]=? GROUP BY [security_code]"
	var newsList []*myoptional_model.StockNewsInformation
	err = repo.FindList(&newsList, sql, _const.StockType_A)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	if newsList != nil {
		for _, news := range newsList {
			stockCodeList = append(stockCodeList, news.SecurityCode)
		}
	}
	return stockCodeList, err
}
