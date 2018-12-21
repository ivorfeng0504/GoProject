package service

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/learnstock"
	"git.emoney.cn/softweb/roboadvisor/protected/service/live"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type UserMedalService struct {
	service.BaseService
	userMedalRepo *userhome_repo.UserMedalRepository
}

var (
	shareUserMedalRepo   *userhome_repo.UserMedalRepository
	shareUserMedalLogger dotlog.Logger
)

const (
	UserMedalServiceName                               = "UserMedalService"
	EMNET_UserMedal_CachePreKey                        = "EMNET:UserMedalService"
	EMNET_UserMedal_GetUserMedal_CacheKey              = EMNET_UserMedal_CachePreKey + ":GetUserMedal:"
	EMNET_UserMedal_GetUserMedalInfoAll_CacheKey       = EMNET_UserMedal_CachePreKey + ":GetUserMedalInfoAll:"
	EMNET_UserMedal_GetQualifiedInvestorLevel_CacheKey = EMNET_UserMedal_CachePreKey + ":GetQualifiedInvestorLevel:"
	EMNET_UserMedal_CacheSeconds                       = 15 * 60
)

func NewUserMedalService() *UserMedalService {
	srv := &UserMedalService{
		userMedalRepo: shareUserMedalRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserMedalServiceName, userMedalServiceLoader)
}

func userMedalServiceLoader() {
	shareUserMedalRepo = userhome_repo.NewUserMedalRepository(protected.DefaultConfig)
	shareUserMedalLogger = dotlog.GetLogger(UserMedalServiceName)
}

// GetMedalName 根据勋章类型获取勋章的名称
func (srv *UserMedalService) GetMedalName(medalType int) string {
	medalMap := map[int]string{
		_const.Medal_WenYuBiaoBing:     "文娱标兵",
		_const.Medal_XueXiWeiYuan:      "学习委员",
		_const.Medal_ZiXunDaKa:         "资讯大咖",
		_const.Medal_GuShiDaHeng:       "股市大亨",
		_const.Medal_HaoWenBoShi:       "好问博士",
		_const.Medal_GaoPengManZuo:     "高朋满座",
		_const.Medal_HuDongDaRen:       "互动达人",
		_const.Medal_SeedUser:          "种子用户",
		_const.Medal_QualifiedInvestor: "合格投资者",
	}
	v, exist := medalMap[medalType]
	if exist == false {
		v = "未知的勋章"
	}
	return v
}

// GetUserMedalInfoAll 获取用户的勋章信息
func (srv *UserMedalService) GetUserMedalInfoAll(user userhome_model.UserInfo) (medalInfo *agent.UserMedalInfo) {
	cacheKey := EMNET_UserMedal_GetUserMedalInfoAll_CacheKey + ":" + strconv.FormatInt(user.UID, 10)
	err := srv.RedisCache.GetJsonObj(cacheKey, &medalInfo)
	if err == nil && medalInfo != nil {
		return medalInfo
	}
	medalInfo = &agent.UserMedalInfo{
		WenYuBiaoBingLevel: srv.GetMedalLevelNoError(user, _const.Medal_WenYuBiaoBing),
		XueXiWeiYuanLevel:  srv.GetMedalLevelNoError(user, _const.Medal_XueXiWeiYuan),
		ZiXunDaKaLevel:     srv.GetMedalLevelNoError(user, _const.Medal_ZiXunDaKa),
		GuShiDaHengLevel:   srv.GetMedalLevelNoError(user, _const.Medal_GuShiDaHeng),
		//该勋章暂时下线
		//HaoWenBoShiLevel:       srv.GetMedalLevelNoError(user, _const.Medal_HaoWenBoShi),
		GaoPengManZuoLevel:     srv.GetMedalLevelNoError(user, _const.Medal_GaoPengManZuo),
		HuDongDaRenLevel:       srv.GetMedalLevelNoError(user, _const.Medal_HuDongDaRen),
		SeedUserLevel:          srv.GetMedalLevelNoError(user, _const.Medal_SeedUser),
		QualifiedInvestorLevel: srv.GetMedalLevelNoError(user, _const.Medal_QualifiedInvestor),
	}
	srv.RedisCache.Set(cacheKey, _json.GetJsonString(medalInfo), EMNET_UserMedal_CacheSeconds)
	return medalInfo
}

// GetMedalLevelNoError 获取好问博士等级 medalType详见MedalType.go
func (srv *UserMedalService) GetMedalLevelNoError(user userhome_model.UserInfo, medalType int) (level int) {
	level, err := srv.GetMedalLevel(user, medalType)
	if err != nil {
		level = 0
	}
	return level
}

// GetMedalLevel 获取勋章等级 medalType详见MedalType.go
func (srv *UserMedalService) GetMedalLevel(user userhome_model.UserInfo, medalType int) (level int, err error) {
	switch medalType {
	case _const.Medal_HaoWenBoShi:
		level, err = srv.GetHaoWenBoShiLevel(user)
		break
	case _const.Medal_WenYuBiaoBing:
		level, err = srv.GetWenYuBiaoBingLevel(user)
		break
	case _const.Medal_HuDongDaRen:
		level, err = srv.GetHuDongDaRenLevel(user)
		break
	case _const.Medal_SeedUser:
		level, err = srv.GetSeedUserLevel(user)
		break
	case _const.Medal_QualifiedInvestor:
		level, err = srv.GetQualifiedInvestorLevel(user)
		break
	}
	//处理勋章新增
	medal, err2 := srv.GetUserMedal(user.UID, medalType)
	if err2 != nil {
		shareUserMedalLogger.ErrorFormat(err, "GetMedalLevel->GetUserMedal 获取用户勋章信息异常 uid=%s medalType=%d", strconv.FormatInt(user.UID, 10), medalType)
	}
	if level > 0 && (medal == nil || medal.MedalLevel < level) {
		medal = &userhome_model.UserMedal{
			UserInfoId: user.ID,
			UID:        user.UID,
			MedalType:  medalType,
			MedalName:  srv.GetMedalName(medalType),
			MedalLevel: level,
		}
		_, err2 = srv.InsertUserMedal(medal)
		if err2 != nil {
			shareUserMedalLogger.ErrorFormat(err, "GetMedalLevel->InsertUserMedal 插入新的勋章异常 勋章信息=%s", jsonutil.GetJsonString(medal))
		}
	}

	return level, err
}

// GetHaoWenBoShiLevel 获取好问博士等级
func (srv *UserMedalService) GetHaoWenBoShiLevel(user userhome_model.UserInfo) (level int, err error) {
	usrSrv := live.NewUserService()
	liveUser, err := usrSrv.GetUserByUID(user.UID)
	if err != nil {
		return level, err
	}
	if liveUser == nil {
		return 0, nil
	}
	questionSrv := live.NewLiveQuestionAnswerService()
	askCount, err := questionSrv.GetMyAskCount(liveUser.UserId)
	if err != nil {
		return level, err
	}
	if askCount >= 100 {
		level = 5
	} else if askCount >= 80 {
		level = 4
	} else if askCount >= 50 {
		level = 3
	} else if askCount >= 20 {
		level = 2
	} else if askCount >= 10 {
		level = 1
	} else {
		level = 0
	}
	return level, err
}

// GetWenYuBiaoBingLevel 获取文娱标兵等级
func (srv *UserMedalService) GetWenYuBiaoBingLevel(user userhome_model.UserInfo) (level int, err error) {
	usrInActivitySrv := NewUserInActivityService()
	joinRecordList, err := usrInActivitySrv.GetUserInActivityList(user.ID)
	if err != nil {
		return level, err
	}
	if joinRecordList == nil {
		return 0, nil
	}
	joinCount := len(joinRecordList)
	if joinCount >= 50 {
		level = 5
	} else if joinCount >= 30 {
		level = 4
	} else if joinCount >= 20 {
		level = 3
	} else if joinCount >= 10 {
		level = 2
	} else if joinCount >= 5 {
		level = 1
	} else {
		level = 0
	}
	return level, err
}

// GetHuDongDaRenLevel 获取互动达人等级
func (srv *UserMedalService) GetHuDongDaRenLevel(user userhome_model.UserInfo) (level int, err error) {
	userInGuessChangeSrv := NewUserInGuessChangeActivityService()
	joinCount, err := userInGuessChangeSrv.GetUserJoinCount(user.ID)
	if joinCount >= 50 {
		level = 5
	} else if joinCount >= 30 {
		level = 4
	} else if joinCount >= 20 {
		level = 3
	} else if joinCount >= 10 {
		level = 2
	} else if joinCount >= 5 {
		level = 1
	} else {
		level = 0
	}
	return level, err
}

// GetSeedUserLevel 获取种子用户等级
func (srv *UserMedalService) GetSeedUserLevel(user userhome_model.UserInfo) (level int, err error) {
	seedUserSrv := NewSeedUserInfoService()
	seedUser, err := seedUserSrv.GetSeedUserInfoByCidCache(strconv.FormatInt(user.UID, 10))
	if err != nil {
		return level, err
	}
	if seedUser != nil && seedUser.IsDeleted == false {
		level = 1
	}
	return level, err
}

// GetQualifiedInvestorLevel 获取合格投资者等级
func (srv *UserMedalService) GetQualifiedInvestorLevel(user userhome_model.UserInfo) (level int, err error) {
	cacheKey := EMNET_UserMedal_GetQualifiedInvestorLevel_CacheKey + strconv.FormatInt(user.UID, 10)
	err = srv.RedisCache.GetJsonObj(cacheKey, &level)
	if err == redis.ErrNil {
		finished, err := learnstock.HasStrategyPowerByCid(user.UID)
		if err != nil {
			shareUserMedalLogger.ErrorFormat(err, "GetQualifiedInvestorLevel 获取合格投资者等级 接口查询异常 UID=%d", user.UID)
			return level, err
		}
		if finished {
			level++
		}
		err = srv.RedisCache.Set(cacheKey, level, EMNET_UserMedal_CacheSeconds)
		return level, err
	}
	return level, err
}

// InsertUserMedal插入一个新的勋章
func (srv *UserMedalService) InsertUserMedal(medal *userhome_model.UserMedal) (id int64, err error) {
	id, err = srv.userMedalRepo.InsertUserMedal(medal)
	if err != nil {
		shareUserMedalLogger.ErrorFormat(err, "插入新的勋章出现异常，勋章信息=%s", jsonutil.GetJsonString(medal))
	} else {
		shareUserMedalLogger.DebugFormat("插入勋章成功 勋章信息=%s", jsonutil.GetJsonString(medal))
		srv.refreshUserMedal(medal.UID, medal.MedalType)
	}
	return id, err
}

// GetUserMedal 查询某个用户的某个勋章
func (srv *UserMedalService) GetUserMedal(uid int64, medalType int) (userMedal *userhome_model.UserMedal, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getUserMedalDB(uid, medalType)
	case config.ReadDB_CacheOrDB_UpdateCache:
		userMedal, err = srv.getUserMedalCache(uid, medalType)
		if err == nil && userMedal == nil {
			userMedal, err = srv.refreshUserMedal(uid, medalType)
		}
		return userMedal, err
	case config.ReadDB_RefreshCache:
		userMedal, err = srv.refreshUserMedal(uid, medalType)
		return userMedal, err
	default:
		return srv.getUserMedalCache(uid, medalType)
	}
}

// getUserMedalDB 查询某个用户的某个勋章
func (srv *UserMedalService) getUserMedalDB(uid int64, medalType int) (userMedal *userhome_model.UserMedal, err error) {
	userMedal, err = srv.userMedalRepo.GetUserMedal(uid, medalType)
	if err != nil {
		shareUserMedalLogger.ErrorFormat(err, "getUserMedalDB 查询某个用户的某个勋章异常 uid=%s medalType=%d", strconv.FormatInt(uid, 10), medalType)
	}
	return userMedal, err
}

// getUserMedalCache 查询某个用户的某个勋章
func (srv *UserMedalService) getUserMedalCache(uid int64, medalType int) (userMedal *userhome_model.UserMedal, err error) {
	cacheKey := fmt.Sprintf("%s%s:%d", EMNET_UserMedal_GetUserMedal_CacheKey, strconv.FormatInt(uid, 10), medalType)
	err = srv.RedisCache.GetJsonObj(cacheKey, &userMedal)
	if err == redis.ErrNil {
		return nil, nil
	}
	return userMedal, err
}

// refreshUserMedal 查询某个用户的某个勋章
func (srv *UserMedalService) refreshUserMedal(uid int64, medalType int) (userMedal *userhome_model.UserMedal, err error) {
	cacheKey := fmt.Sprintf("%s%s:%d", EMNET_UserMedal_GetUserMedal_CacheKey, strconv.FormatInt(uid, 10), medalType)
	userMedal, err = srv.getUserMedalDB(uid, medalType)
	if err != nil {
		return nil, err
	}
	if userMedal != nil {
		srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(userMedal), EMNET_UserMedal_CacheSeconds)
	}
	return userMedal, err
}
