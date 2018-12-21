package repository

import (
	sysSql "database/sql"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"time"
)

type UserReceivedStockRepository struct {
	repository.BaseRepository
}

func NewUserReceivedStockRepository(conf *protected.ServiceConfig) *UserReceivedStockRepository {
	repo := &UserReceivedStockRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertUserStock 用户领取股票
func (repo *UserReceivedStockRepository) InsertUserStock(userStock *userhome_model.UserReceivedStock) (id int64, err error) {
	sql := "INSERT INTO [UserHome_UserReceivedStock] ([ReceiveStockActivityId],[CreateTime],[StockCreateTime],[IssueNumber],[StockList],[UserInfoId],[UID],[ReportUrl])VALUES(?,GETDATE(),?,?,?,?,?,?)"
	id, err = repo.Insert(sql, userStock.ReceiveStockActivityId, time.Time(userStock.StockCreateTime).Format("2006-01-02 15:04:05"), userStock.IssueNumber, userStock.StockList, userStock.UserInfoId, userStock.UID, userStock.ReportUrl)
	return id, err
}

// GetUserStockToday 获取用户当日领取的股票
func (repo *UserReceivedStockRepository) GetUserStockToday(userInfoId int) (userStock *userhome_model.UserReceivedStock, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserReceivedStock] WHERE [UserInfoId]=? AND CONVERT(NVARCHAR(8),GETDATE(),112)=CONVERT(NVARCHAR(8),[CreateTime],112) ORDER BY [CreateTime] DESC"
	userStock = new(userhome_model.UserReceivedStock)
	err = repo.FindOne(userStock, sql, userInfoId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return userStock, err
}

// UpdateUserReceivedStock 更新用户领取的研报
func (repo *UserReceivedStockRepository) UpdateUserReceivedStock(userReceivedStockId int64, reportUrl string) (err error) {
	sql := "UPDATE [UserHome_UserReceivedStock] SET [ReportUrl]=? WHERE [UserReceivedStockId]=?"
	_, err = repo.Update(sql, reportUrl, userReceivedStockId)
	return err
}

// GetUserStockHistory 获取用户领取的股票历史
func (repo *UserReceivedStockRepository) GetUserStockHistory(userInfoId int, top int) (userStockList []*userhome_model.UserReceivedStock, err error) {
	sql := fmt.Sprintf("SELECT TOP %d * FROM [UserHome_UserReceivedStock] WHERE [UserInfoId]=? ORDER BY [CreateTime] DESC", top)
	err = repo.FindList(&userStockList, sql, userInfoId)
	return userStockList, err
}
