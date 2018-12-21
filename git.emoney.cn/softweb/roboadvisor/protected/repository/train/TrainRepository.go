package train

import (
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected/model/train"
	"time"
	"database/sql"
)

type TrainRepository struct {
	repository.BaseRepository
}

var (
	// shareUserLogger 共享的Logger实例
	shareUserLogger dotlog.Logger
)

func NewTrainRepository(conf *protected.ServiceConfig) *TrainRepository {
	repo := &TrainRepository{}
	repo.Init(conf.UserTrainDBConn)
	return repo
}

// GetTrainListByPID 根据PID获取所属培训课程(获取最近一个月)
func (repo *TrainRepository) GetTrainListByPID(pid string) (trainlist []*train.NetworkMeetingInfo,err error) {
	pid = "%" + pid + "%"
	sqlstr := "SELECT * FROM [Traning_NetworkMeeting] WITH(NOLOCK) WHERE chblist_pro like ? and class_date>=DateAdd(MM,-1,GETDATE()) and CONVERT(varchar(100), class_date, 23)<=GETDATE() ORDER BY class_date DESC"
	err = repo.FindList(&trainlist, sqlstr, pid)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return trainlist, err
}

// GetTrainListByDate 根据日期获取智盈培训课程
func (repo *TrainRepository) GetTrainListByDateAndPid(pid string,date *time.Time) (trainlist []*train.NetworkMeetingInfo,err error) {
	dateFormat := date.Format("2006-01-02")
	pid = "%" + pid + "%"
	sql := "SELECT * FROM [Traning_NetworkMeeting] WITH(NOLOCK) WHERE chblist_pro like ? AND CONVERT(varchar(100),class_date,23)=CONVERT(varchar(100),?,23) ORDER BY class_date DESC"
	err = repo.FindList(&trainlist, sql, pid, dateFormat)
	return trainlist, err
}


// GetTrainListByTag 根据标签获取智盈培训课程(获取最近三个月)
func (repo *TrainRepository) GetTrainListByTag(pid string,traintag int) (trainlist []*train.NetworkMeetingInfo,err error) {
	sql := "SELECT * FROM [Traning_NetworkMeeting] WITH(NOLOCK) WHERE chblist_pro like '%?%' AND TrainTag=? ORDER BY class_date DESC"
	err = repo.FindList(&trainlist, sql, pid, traintag)
	return trainlist, err
}


// GetTrainListByArea 根据地区获取智盈培训课程（获取最近三个月）
func (repo *TrainRepository) GetTrainListByArea(pid string,area string) (trainlist []*train.NetworkMeetingInfo,err error) {
	sql := "SELECT * FROM [Traning_NetworkMeeting] WITH(NOLOCK) WHERE chblist_pro like '%?%' AND ddlarea=? ORDER BY class_date DESC"
	err = repo.FindList(&trainlist, sql, pid, area)
	return trainlist, err
}
