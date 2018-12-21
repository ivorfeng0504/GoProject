package myoptional

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"strconv"
	"time"
)

type StockTalkRepository struct {
	repository.BaseRepository
}

func NewStockTalkRepository(conf *protected.ServiceConfig) *StockTalkRepository {
	repo := &StockTalkRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertStockTalk 插入一个微股吧评论
func (repo *StockTalkRepository) InsertStockTalk(model *myoptional_model.StockTalk) (err error) {
	sql := "INSERT INTO [StockTalk] ([NickName],[Avatar],[StockCode],[StockName],[Content],[ModifyUser],[UID],[IsValid],[IsDeleted],[ModifyTime],[CreateTime],[StockTalkType]) VALUES(?,?,?,?,?,?,?,0,0,GETDATE(),GETDATE(),?)"
	_, err = repo.Insert(sql, model.NickName, model.Avatar, model.StockCode, model.StockName, model.Content, model.NickName, model.UID, model.StockTalkType)
	return err
}

// InsertStockTalkWithAdmin 管理员插入一个微股吧评论 直接审核通过 ModifyUser使用原始值
func (repo *StockTalkRepository) InsertStockTalkWithAdmin(model *myoptional_model.StockTalk) (id int64, err error) {
	sql := "INSERT INTO [StockTalk] ([NickName],[Avatar],[StockCode],[StockName],[Content],[ModifyUser],[UID],[IsValid],[IsDeleted],[ModifyTime],[CreateTime],[TalkLevel],[StockTalkType]) VALUES(?,?,?,?,?,?,?,1,0,GETDATE(),GETDATE(),?,?)"
	id, err = repo.Insert(sql, model.NickName, model.Avatar, model.StockCode, model.StockName, model.Content, model.ModifyUser, model.UID, model.TalkLevel, model.StockTalkType)
	return id, err
}

// GetStockTalkById 通过主键Id获取评论
func (repo *StockTalkRepository) GetStockTalkById(stockTalkId int64) (model *myoptional_model.StockTalk, err error) {
	sql := "SELECT TOP 1 * FROM [StockTalk] WITH(NOLOCK) WHERE [StockTalkId]=?"
	model = new(myoptional_model.StockTalk)
	err = repo.FindOne(model, sql, stockTalkId)
	if err == sysSql.ErrNoRows {
		err = nil
		model = nil
	}
	return model, err
}

// GetStockTalkListByStockCode 通过股票代码获取评论列表 Top20
func (repo *StockTalkRepository) GetStockTalkListByStockCode(stockCode string) (list []*myoptional_model.StockTalk, err error) {
	sql := "SELECT TOP 20 * FROM [StockTalk] WITH(NOLOCK) WHERE [IsDeleted]=0 AND [IsValid]=1 AND StockCode=? ORDER BY [TalkLevel] DESC,[CreateTime] DESC"
	err = repo.FindList(&list, sql, stockCode)
	if err == sysSql.ErrNoRows {
		err = nil
		list = nil
	}
	return list, err
}

// GetStockCodeList 获取微股吧码表
func (repo *StockTalkRepository) GetStockCodeList() (stockCodeList []*model.StockInfo, err error) {
	sql := "SELECT DISTINCT [StockCode],[StockName] FROM [StockTalk]  WHERE [IsValid]=1 AND [IsDeleted]=0 GROUP BY [StockCode],[StockName]"
	err = repo.FindList(&stockCodeList, sql)
	if err == sysSql.ErrNoRows {
		err = nil
		stockCodeList = nil
	}
	return stockCodeList, err
}

// GetRepeatNewsInfo 获取重复的数据集
func (repo *StockTalkRepository) GetRepeatStockTalk() (repeatStockTalkList []*myoptional_model.RepeatStockTalk, err error) {
	sql := `SELECT * FROM (SELECT [StockCode],[StockName],[Content],COUNT(1) AS RepeatCount FROM [StockTalk] WHERE [IsDeleted]=0 AND [IsValid]=1 AND [UID]='' GROUP BY [StockCode],[Content],[StockName])AS TB1 WHERE [RepeatCount]>1`
	err = repo.FindList(&repeatStockTalkList, sql)
	return repeatStockTalkList, err
}

// DeleteRepeatNewsInfo 删除重复的数据-逻辑删除
func (repo *StockTalkRepository) DeleteRepeatStockTalk(repeatStockTalk *myoptional_model.RepeatStockTalk) (n int64, err error) {
	if repeatStockTalk.RepeatCount <= 1 {
		return 0, nil
	}
	repeatStockTalk.RepeatCount--
	sql := `  UPDATE [StockTalk] SET [IsDeleted]=1,[ModifyTime]=GETDATE() WHERE [StockTalkId] IN (
  SELECT TOP  ` + strconv.Itoa(repeatStockTalk.RepeatCount) + `[StockTalkId] FROM [StockTalk] WHERE [StockCode]=? AND [StockName]=? AND [Content]=? AND [UID]='' AND [IsDeleted]=0 AND [IsValid]=1 ORDER BY [CreateTime]
  )`
	n, err = repo.Update(sql, repeatStockTalk.StockCode, repeatStockTalk.StockName, repeatStockTalk.Content)
	return n, err
}

// GetDeletedStockTalkListToday 获取今天删除的数据
func (repo *StockTalkRepository) GetDeletedStockTalkListToday() (list []*myoptional_model.StockTalk, err error) {
	sql := "SELECT * FROM [StockTalk] WITH(NOLOCK) WHERE [IsDeleted]=1 AND [ModifyTime]>?"
	err = repo.FindList(&list, sql, time.Now().Format("2006-01-02"))
	if err == sysSql.ErrNoRows {
		err = nil
		list = nil
	}
	return list, err
}

// DeleteStockTalk 根据股票代码和类型删除数据-逻辑删除
func (repo *StockTalkRepository) DeleteStockTalk(stockCode string, stockTalkType int) (n int64, err error) {
	if len(stockCode) == 0 {
		return 0, nil
	}
	sql := `UPDATE [StockTalk] SET [IsDeleted]=1,[ModifyTime]=GETDATE() WHERE [StockCode]=? AND [StockTalkType]=? AND [UID]='' AND [IsDeleted]=0 AND [IsValid]=1 `
	n, err = repo.Update(sql, stockCode, stockTalkType)
	return n, err
}
