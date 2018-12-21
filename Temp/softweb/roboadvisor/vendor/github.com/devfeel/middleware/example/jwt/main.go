package main

import (
	"errors"
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/middleware/jwt"
	"strconv"
	"time"
)

const JwtContextKey = "jwtuser"

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	//如果不设置，默认不启用，且默认为当前目录
	app.SetEnabledLog(true)

	//开启development模式
	app.SetDevelopmentMode()

	//设置路由
	InitRoute(app.HttpServer)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	payload, exists := ctx.Items().Get(JwtContextKey)
	return ctx.WriteString("custom jwt context => ", payload, exists)
}

func Login(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	config := parseJwtConfig(ctx.AppItems().Get("CustomJwtConfig"))
	if config == nil {
		return ctx.WriteString("custom login failed, token config not exists")
	}
	m := make(map[string]interface{})
	m["userid"] = "loginuser"
	m["userip"] = ctx.RemoteIP()
	token, err := jwt.GeneratorToken(config, m)
	if err != nil || token == "" {
		return ctx.WriteString("custom login failed, token create failed, ", err.Error())
	}

	ctx.SetCookieValue(config.Name, token, 0)
	return ctx.WriteString("custom login is ok, token => ", token)
}

func Logout(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	config := parseJwtConfig(ctx.AppItems().Get("CustomJwtConfig"))
	if config == nil {
		return ctx.WriteString("logout failed, token config not exists")
	}
	ctx.RemoveCookie(config.Name)
	return ctx.WriteString("logout is ok")
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index).Use(NewCustomJwt(server.DotApp))
	server.Router().GET("/Login", Login)
	server.Router().GET("/Logout", Logout)
}

func NewSimpleJwt(app *dotweb.DotWeb) dotweb.Middleware {
	option := &jwt.Config{
		SigningKey: []byte("devfeel/dotweb"), //must input
		//use cookie
		Extractor: jwt.ExtractorFromCookie,
	}
	app.Items.Set("SimpleJwtConfig", option)
	return jwt.Middleware(option)
}

func NewCustomJwt(app *dotweb.DotWeb) dotweb.Middleware {
	option := &jwt.Config{
		TTL:           time.Second * 20,         //default is 24 hour
		ContextKey:    JwtContextKey,            //default is dotjwt-user
		SigningKey:    []byte("devfeel/dotweb"), //must input
		SigningMethod: jwt.SigningMethodHS256,   //default is SigningMethodHS256
		ExceptionHandler: func(ctx dotweb.Context, err error) {
			//TODO:log err info
			ctx.WriteString("no authorization, please login first")
		},
		AddonValidator: func(config *jwt.Config, ctx dotweb.Context) error {
			//example: check user ip
			jwtobj, exists := ctx.Items().Get(JwtContextKey)
			if !exists {
				return errors.New("no token exists")
			}
			jwtmap := jwtobj.(map[string]interface{})

			jwtUserIp := jwtmap["userip"].(string)
			requestIp := ctx.RemoteIP()
			fmt.Println("jwtUserIp", jwtUserIp, " requestIp:", requestIp)
			if jwtUserIp != requestIp {
				return errors.New("ip is not match")
			}
			return nil
		},
		//use cookie
		Extractor: jwt.ExtractorFromCookie,
	}

	app.Items.Set("CustomJwtConfig", option)

	return jwt.Middleware(option)
}

func parseJwtConfig(c interface{}, exists bool) (config *jwt.Config) {
	if c == nil || !exists {
		return nil
	}
	config = c.(*jwt.Config)
	return config
}
