package repository

import (
	sysSql "database/sql"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

type ReceiveStockActivityRepository struct {
	repository.BaseRepository
}

func NewReceiveStockActivityRepository(conf *protected.ServiceConfig) *ReceiveStockActivityRepository {
	repo := &ReceiveStockActivityRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertStockPool 新增一个股票池
func (repo *ReceiveStockActivityRepository) InsertStockPool(stocks []*userhome_model.StockInfo, issueNumber string) (id int64, err error) {
	if stocks == nil || len(stocks) == 0 {
		err = errors.New("股票信息不能为空！")
		return id, err
	}
	stockList := ""
	if stocks != nil && len(stocks) > 0 {
		stockList = _json.GetJsonString(stocks)
	}
	sql := "INSERT INTO [UserHome_ReceiveStockActivity] ([StockList],[CreateTime],[IssueNumber])VALUES(?,GETDATE(),?)"
	id, err = repo.Insert(sql, stockList, issueNumber)
	return id, err
}

// InsertStockPoolByReportUrl 新增活动
func (repo *ReceiveStockActivityRepository) InsertStockPoolByReportUrl(reportUrl string, issueNumber string) (id int64, err error) {
	if len(reportUrl) == 0 {
		err = errors.New("reportUrl不能为空！")
		return id, err
	}
	sql := "INSERT INTO [UserHome_ReceiveStockActivity] ([ReportUrl],[CreateTime],[IssueNumber])VALUES(?,GETDATE(),?)"
	id, err = repo.Insert(sql, reportUrl, issueNumber)
	return id, err
}

// UpdateStockPoolByReportUrl 更新活动
func (repo *ReceiveStockActivityRepository) UpdateStockPoolByReportUrl(receiveStockActivityId int64, reportUrl string) (err error) {
	if len(reportUrl) == 0 {
		err = errors.New("reportUrl不能为空！")
		return err
	}
	sql := "UPDATE [UserHome_ReceiveStockActivity] SET [ReportUrl]=? WHERE [ReceiveStockActivityId]=?"
	_, err = repo.Update(sql, reportUrl, receiveStockActivityId)
	return err
}

// GetNewstStockPool 获取最新的股票池
func (repo *ReceiveStockActivityRepository) GetNewstStockPool() (stock *userhome_model.ReceiveStockActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_ReceiveStockActivity] ORDER BY [IssueNumber] DESC"
	stock = new(userhome_model.ReceiveStockActivity)
	err = repo.FindOne(stock, sql)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return stock, err
}
