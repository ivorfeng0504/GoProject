package server

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/middleware/userhome"
	"git.emoney.cn/softweb/roboadvisor/web_front_userhome/handlers/activity"
	"git.emoney.cn/softweb/roboadvisor/web_front_userhome/handlers/error"
	"git.emoney.cn/softweb/roboadvisor/web_front_userhome/handlers/index"
	"git.emoney.cn/softweb/roboadvisor/web_front_userhome/handlers/mine"
	"git.emoney.cn/softweb/roboadvisor/web_front_userhome/handlers/validatecode"
	"github.com/devfeel/dotweb"
)

func InitRoute(server *dotweb.HttpServer) {
	server.SetVirtualPath(config.CurrentConfig.ServerVirtualPath)
	server.Router().ServerFile("/static/*filepath", config.CurrentConfig.ResourcePath+`static`)
	g := server.Group("/homepage")

	//非置顶通知列表
	//调用示例
	//$.get("http://127.0.0.1:8084/homepage/newslist",{ColumnID:1,StrategyID:0},function(data){console.log(JSON.stringify(data))})
	g.GET("/newslist", index.GetNewsListByColumnID).Use(middleware.NewUserHomeAuthMiddleware())

	//置顶通知信息
	//$.get("http://127.0.0.1:8084/homepage/topnewslist",{ColumnID:1,StrategyID:0},function(data){console.log(JSON.stringify(data))})
	g.GET("/topnewslist", index.GetTopNewsInfoByColumnID).Use(middleware.NewUserHomeAuthMiddleware())

	//单条通知信息
	//$.get("http://127.0.0.1:8084/homepage/newsinfo",{NewsId:1},function(data){console.log(JSON.stringify(data))})
	g.GET("/newsinfo", index.GetNewsInfoByID).Use(middleware.NewUserHomeAuthMiddleware())

	//获取用户信息
	//$.get("http://127.0.0.1:8084/homepage/userinfo",{},function(data){console.log(JSON.stringify(data))})
	g.GET("/userinfo", index.GetUserInfo).Use(middleware.NewUserHomeAuthMiddleware())

	//首页
	g.GET("/index", index.Index).Use(middleware.NewUserHomeSSOMiddleware())

	g.GET("/rsatest", index.Rsatest)

	//活动相关
	g = server.Group("/activity")
	//活动首页
	g.GET("/index", activity.Index).Use(middleware.NewUserHomeAuthMiddleware())
	//领数字活动首页
	g.GET("/guessmarket", activity.GuessMarket).Use(middleware.NewUserHomeAuthMiddleware())
	//双响炮活动首页
	g.GET("/receivestock", activity.ReceiveStock).Use(middleware.NewUserHomeAuthMiddleware())
	//猜涨跌活动首页
	g.GET("/guesschange", activity.GuessChange).Use(middleware.NewUserHomeAuthMiddleware())
	//获取活动列表
	g.POST("/getactivitylist", activity.GetActivityList).Use(middleware.NewUserHomeAuthMiddleware())
	//获取当前用户的连登奖励
	g.POST("/getserialloginrule", activity.GetSerialLoginRule).Use(middleware.NewUserHomeAuthMiddleware())
	//获取领数字活动历史信息
	//调用示例： $.post("http://127.0.0.1:8084/activity/getguessmarkethistoryinfo",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getguessmarkethistoryinfo", activity.GetGuessMarketHistoryInfo).Use(middleware.NewUserHomeAuthMiddleware())
	//获取领数字活动奖品信息
	//调用示例： $.post("http://127.0.0.1:8084/activity/getguessmarketaward",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getguessmarketaward", activity.GetAwardInfo).Use(middleware.NewUserHomeAuthMiddleware())
	//用户领取数字
	//调用示例： $.post("http://127.0.0.1:8084/activity/getguessmarketnumber",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getguessmarketnumber", activity.GetGuessMarketNumber).Use(middleware.NewUserHomeAuthMiddleware())
	//用户领取数字并返回当前活动的详细信息
	//调用示例： $.post("http://127.0.0.1:8084/activity/getguessmarketnumberwithdetail",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getguessmarketnumberwithdetail", activity.GetGuessMarketNumberWithDetail).Use(middleware.NewUserHomeAuthMiddleware())
	//获取当前领数字活动与用户参与信息
	//调用示例： $.post("http://127.0.0.1:8084/activity/getcurrentguessinfo",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getcurrentguessinfo", activity.GetCurrentGuessInfo).Use(middleware.NewUserHomeAuthMiddleware())
	//获取用户勋章
	//调用示例： $.post("http://127.0.0.1:8084/activity/getusermedal",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getusermedal", activity.GetUserMedal).Use(middleware.NewUserHomeAuthMiddleware())
	//获取用户连登天数，如果已经当前已经领了股票 则直接返回股票信息
	//调用示例： $.post("http://127.0.0.1:8084/activity/getreceivestat",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getreceivestat", activity.GetUserSerialLoginDay).Use(middleware.NewUserHomeAuthMiddleware())
	//解密今日股票
	//调用示例： $.post("http://127.0.0.1:8084/activity/getstocklist",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getstocklist", activity.GetStockList).Use(middleware.NewUserHomeAuthMiddleware())
	//获取用户领取的股票历史
	//调用示例： $.post("http://127.0.0.1:8084/activity/getuserstockhistory",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getuserstockhistory", activity.GetUserStockHistory).Use(middleware.NewUserHomeAuthMiddleware())

	//获取当期猜涨跌活动情况
	//调用示例： $.post("http://127.0.0.1:8084/activity/currentguesschange",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/currentguesschange", activity.GetCurrentGuessChange).Use(middleware.NewUserHomeAuthMiddleware())
	//提交猜涨跌竞猜结果
	//调用示例： $.post("http://127.0.0.1:8084/activity/guesschangesubmit",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/guesschangesubmit", activity.GuessChangeSubmit).Use(middleware.NewUserHomeAuthMiddleware())
	//获取我的竞猜和我的奖品信息
	//调用示例： $.post("http://127.0.0.1:8084/activity/myguesschangejoininfo",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/myguesschangejoininfo", activity.GetMyGuessChangeJoinInfo).Use(middleware.NewUserHomeAuthMiddleware())
	//获取我的竞猜-更多信息
	//调用示例： $.post("http://127.0.0.1:8084/activity/myguesschangeinfonewst",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/myguesschangeinfonewst", activity.GetMyGuessChangeInfoNewst).Use(middleware.NewUserHomeAuthMiddleware())
	//获取猜涨跌奖品信息
	//调用示例： $.post("http://127.0.0.1:8084/activity/getguesschangeaward",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/getguesschangeaward", activity.GetGuessChangeAwardInfo).Use(middleware.NewUserHomeAuthMiddleware())

	//我的
	g = server.Group("/mine")
	g.GET("/index", mine.Index).Use(middleware.NewUserHomeAuthMiddleware())
	//我的订单
	g.GET("/myorder", mine.MyOrder).Use(middleware.NewUserHomeAuthMiddleware())
	//订单详情
	g.GET("/myorderdetail", mine.OrderDetail).Use(middleware.NewUserHomeAuthMiddleware())
	//个人信息
	g.GET("/myprofile", mine.MyProfile).Use(middleware.NewUserHomeAuthMiddleware())
	//QQ登录页
	g.GET("/qqlogin", mine.Login_QQ).Use(middleware.NewUserHomeAuthMiddleware())
	//微信登录页
	g.GET("/wechatlogin", mine.Login_WeChat).Use(middleware.NewUserHomeAuthMiddleware())
	//QQ绑定确认页
	g.GET("/qqbind", mine.Bind_QQ).Use(middleware.NewUserHomeAuthMiddleware())
	//微信绑定确认页
	g.GET("/wechatbind", mine.Bind_WeChat).Use(middleware.NewUserHomeAuthMiddleware())
	// 获取用户奖品列表
	g.POST("/getuserawardlist", mine.GetUserAwardList).Use(middleware.NewUserHomeAuthMiddleware())
	// 获取用户产品列表
	g.POST("/getuserproductlist", mine.GetUserProductList).Use(middleware.NewUserHomeAuthMiddleware())
	// 查询订单
	//调用示例： $.post("http://127.0.0.1:8084/mine/queryorderlist",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/queryorderlist", mine.QueryOrderList).Use(middleware.NewUserHomeAuthMiddleware())
	// 验证是否支持快速退货退款
	//调用示例： $.post("http://127.0.0.1:8084/mine/validateorder",{OrderId:'SH1130807194553514'},function(data){console.log(JSON.stringify(data))})
	g.POST("/validateorder", mine.ValidateOrder).Use(middleware.NewUserHomeAuthMiddleware())
	//用户退款提交
	//调用示例： $.post("http://127.0.0.1:8084/mine/refundsubmit",{
	//    "OrderId": "SH1130812165527803",
	//    "Reason": "没有理由",
	//    "Name": "张三",
	//    "BankAccount": "10010011212555",
	//    "BankValue": "民生银行",
	//    "BankDetail": "北京西单民生银行总行",
	//	  "RefundMode":1
	//}
	// ,function(data){console.log(JSON.stringify(data))})
	g.POST("/refundsubmit", mine.RefundSubmit).Use(middleware.NewUserHomeAuthMiddleware())
	//获取订单信息
	//调用示例： $.post("http://127.0.0.1:8084/mine/getorderinfo",{OrderId:'SH1180716152253509_test'},function(data){console.log(JSON.stringify(data))})
	g.POST("/getorderinfo", mine.GetOrderInfo).Use(middleware.NewUserHomeAuthMiddleware())

	//获取验证码
	//$.get("http://127.0.0.1:8084/homepage/getValidateCode",{Mobile:""},function(data){console.log(JSON.stringify(data))})
	g.GET("/getValidateCode", mine.GetValidateCode).Use(middleware.NewUserHomeAuthMiddleware())

	//注册
	//$.get("http://127.0.0.1:8084/homepage/regTyAccount",{Mobile:"",code:"",password:"",confirmpassword:""},function(data){console.log(JSON.stringify(data))})
	g.GET("/regTyAccount", mine.Reg_TyAccount).Use(middleware.NewUserHomeAuthMiddleware())

	//修改密码
	//$.get("http://127.0.0.1:8084/homepage/modifypassword",{Mobile:"",code:"",password:"",confirmpassword:""},function(data){console.log(JSON.stringify(data))})
	g.GET("/modifypassword", mine.ModifyPassword).Use(middleware.NewUserHomeAuthMiddleware())

	//绑定-手机账号
	//$.get("http://127.0.0.1:8084/homepage/bindAccount",{Mobile:"",code:""},function(data){console.log(JSON.stringify(data))})
	g.GET("/bindaccountphone", mine.BindAccountPhone).Use(middleware.NewUserHomeAuthMiddleware())

	//绑定-EM账号
	g.GET("/bindaccountem", mine.BindAccountEM).Use(middleware.NewUserHomeAuthMiddleware())

	//绑定-微信账号
	g.GET("/bindaccountwechat", mine.BindAccountWeChat).Use(middleware.NewUserHomeAuthMiddleware())

	//绑定-QQ账号
	g.GET("/bindaccountqq", mine.BindAccountQQ).Use(middleware.NewUserHomeAuthMiddleware())

	//修改头像
	//$.get("http://127.0.0.1:8084/homepage/modifyheadportrait",{headimg:""},function(data){console.log(JSON.stringify(data))})
	g.GET("/modifyheadportrait", mine.ModifyHeadportrait).Use(middleware.NewUserHomeAuthMiddleware())

	//修改昵称
	//$.get("http://127.0.0.1:8084/homepage/modifynickname",{nickname:""},function(data){console.log(JSON.stringify(data))})
	g.GET("/modifynickname", mine.ModifyNickName).Use(middleware.NewUserHomeAuthMiddleware())

	//获取手机密码
	g.GET("/getpwdbymobile", mine.GetPwdByMobile).Use(middleware.NewUserHomeAuthMiddleware())

	// 解除绑定
	//$.get("http://127.0.0.1:8084/mine/removebind",{cid:""},function(data){console.log(JSON.stringify(data))})
	g.GET("/removebind", mine.RemoveBind).Use(middleware.NewUserHomeAuthMiddleware())

	//获取个人资料、绑定信息等
	g.GET("/getprofile", mine.GetProfile).Use(middleware.NewUserHomeAuthMiddleware())

	//获取已绑定手机号
	g.GET("/getbindmobilebyuid", mine.GetBindMobileByUID).Use(middleware.NewUserHomeAuthMiddleware())

	g = server.Group("/captcha")
	g.GET("/page", validatecode.ShowCaptchaPage)
	g.GET("/image", validatecode.BuffImage)
	g.GET("/reloadimage", validatecode.BuffNewImage)
	g.POST("/verify", validatecode.VerifyCaptcha)
	g.GET("/fetchid", validatecode.GetCaptchaId)
	//错误页面
	g = server.Group("/error")
	//访问地址不存在
	g.GET("/404", _error.NotFound)
	//服务器出错
	g.GET("/500", _error.ServerError)
	//身份认证出错
	g.GET("/401", _error.NotAuth)

}
