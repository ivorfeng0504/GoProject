package live

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type MaskWordRepository struct {
	repository.BaseRepository
}

func NewMaskWordRepository(conf *protected.ServiceConfig) *MaskWordRepository {
	repo := &MaskWordRepository{}
	repo.Init(conf.ContentLivePlatDBConn)
	return repo
}

// GetMaskWordList 获取所有屏蔽字
func (repo *MaskWordRepository) GetMaskWordList() (maskWordList []*livemodel.MaskWord, err error) {
	sql := "SELECT * FROM [MaskWord] WHERE [IsDelete]=0"
	err = repo.FindList(&maskWordList, sql)
	return maskWordList, err
}
