package server

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/web_api_live/handlers/live"
	"github.com/devfeel/dotweb"
)

func InitRoute(server *dotweb.HttpServer) {
	server.SetVirtualPath(config.CurrentConfig.ServerVirtualPath)
	server.Router().ServerFile("/static/*filepath", config.CurrentConfig.ResourcePath+`static`)
	g := server.Group("/api/live")

	//调用示例
	//$.post("http://127.0.0.1:8092/api/live/istradetime",{},function(data){console.log(JSON.stringify(data))})
	g.POST("/istradetime", live.IsTradeTime)
	g.POST("/addquestion", live.AddQuestion)
	g.POST("/livequestion", live.GetLiveQuestionAnswerList)
	g.POST("/livecontent", live.GetLiveContent)
	g.POST("/hasnewmessage", live.HasNewMessage)
	g.POST("/addroompermit", live.AddLiveUserInRoom)
	g.POST("/removeroompermit", live.RemoveLiveUserInRoom)
	g.POST("/getuserroomlist", live.GetUserRoomList)

	//用户操作
	g.POST("/getuserbyuserid", live.GetUserById)
	g.POST("/getuserbyaccount", live.GetUserByAccount)
	g.POST("/getuserbyuid", live.GetUserByUID)
	g.POST("/adduser", live.AddUser)

	//在websocket网关中注册的地址，websocket服务回调地址  授权校验
	g.GET("/checktoken", live.CheckToken)
	//在websocket网关中注册的地址，long poll服务回调地址  获取消息
	g.GET("/hasnewmessage", live.GetHasNewMessage)
}
