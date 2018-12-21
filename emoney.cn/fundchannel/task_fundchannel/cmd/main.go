package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"emoney.cn/fundchannel/config"
	"emoney.cn/fundchannel/const"
	"emoney.cn/fundchannel/core/exception"
	"emoney.cn/fundchannel/global"
	"emoney.cn/fundchannel/protected"
	"emoney.cn/fundchannel/task_fundchannel/server"
	"emoney.cn/fundchannel/util/file"
	"time"
)

var (
	configFile string
	configPath string
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			ex := exception.CatchError(_const.Global_ProjectName+":main", err)
			global.InnerLogger.Error(fmt.Errorf("%v", err), ex.GetDefaultLogString())
			os.Stdout.Write([]byte(ex.GetDefaultLogString()))
		}
	}()

	parseFlag()

	//全局初始化
	err := global.Init(configPath)
	if err != nil {
		panic(err)
	}

	//加载全局xml配置文件
	//初始化app.config
	config.InitConfig(configFile)

	//服务初始化工作
	err = protected.Init()
	if err != nil {
		global.InnerLogger.Error(err, "protected.InitConfig失败 "+err.Error())
		fmt.Println("protected.InitConfig失败 " + err.Error())
		return
	}

	//监听系统信号
	//go listenSignal()

	//启动监听服务
	server.StartTaskService(configPath)

	time.Sleep(time.Minute * 60 * 24 * 50)
}

func parseFlag() {
	server.RunEnv = os.Getenv(server.RunEnv_Flag)
	if server.RunEnv == "" {
		server.RunEnv = server.RunEnv_Develop
	}

	configPath = _file.GetCurrentDirectory() + "/conf/" + server.RunEnv
	//load app config
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	if configFile == "" {
		configFile = configPath + "/app.conf"
	}

}

func listenSignal() {
	c := make(chan os.Signal, 1)
	//syscall.SIGSTOP
	signal.Notify(c, syscall.SIGHUP)
	for {
		s := <-c
		global.InnerLogger.Info("signal::ListenSignal [" + s.String() + "]")
		switch s {
		case syscall.SIGHUP: //配置重载
			global.InnerLogger.Info("signal::ListenSignal reload config begin...")
			//重新加载配置文件
			config.InitConfig(configFile)
			global.InnerLogger.Info("signal::ListenSignal reload config end")
		default:
			return
		}
	}
}
