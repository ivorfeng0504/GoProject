package protected

type ServiceConfig struct {
	//默认数据库连接
	DefaultDBConn string
	//图文直播数据库连接
	ContentLivePlatDBConn string
	//智投数据库连接
	EMoney_RoboAdvisorDBConn string
	//智投fee账号信息相关数据库连接
	FeedbDBConn string
	//智投emoney账号信息相关数据库连接
	EmoneyDBConn string
	//用户培训数据连接
	UserTrainDBConn string
}