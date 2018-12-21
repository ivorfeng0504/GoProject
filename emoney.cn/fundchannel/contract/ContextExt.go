package contract

import (
	"emoney.cn/fundchannel/global"
	"emoney.cn/fundchannel/config"
	"github.com/devfeel/dotweb"
	"io/ioutil"
	"strings"
)

const (
	//静态资源域名占位符，域名末尾不包含/
	StaticServerHost = "$StaticServerHost$"
	//静态资源环境配置后缀占位符
	StaticResEnv = "$StaticResEnv$"
	//静态资源版本号
	StaticVersion = "$ResourceVersion$"
	//类虚拟目录路径
	ServerVirtualPath = "$ServerVirtualPath$"
	SSOQuery          = "$SSOQuery$"
)

// RenderHtml 直接返回Html页面，不进行服务端渲染，并且将页面缓存在服务器中
// processHtmlFuncs HTML内容处理器 对页面内容进行自定义处理
func RenderHtml(ctx dotweb.Context, filename string, processHtmlFuncs ...func(string) string) error {
	html, exist := global.GlobalItemMap.Get(filename)
	ssoQuery := GetSSOQuery(ctx)
	if exist && config.CurrentConfig.DisabledHtmlCache == false {
		for _, process := range processHtmlFuncs {
			html = process(html.(string))
		}
		html = strings.Replace(html.(string), SSOQuery, ssoQuery, -1)
		return ctx.WriteHtml(html)
	}
	data, err := ioutil.ReadFile(config.CurrentConfig.ResourcePath + "views/" + filename)
	if err != nil {
		//防止死循环
		if filename == "404.html" {
			return ctx.WriteString("页面丢失了！")
		}
		return RenderHtml(ctx, "404.html")
	}
	html = string(data)
	html = strings.Replace(html.(string), StaticServerHost, config.CurrentConfig.StaticServerHost, -1)
	html = strings.Replace(html.(string), StaticResEnv, config.CurrentConfig.StaticResEnv, -1)
	html = strings.Replace(html.(string), StaticVersion, config.CurrentConfig.ResourceVersion, -1)
	html = strings.Replace(html.(string), ServerVirtualPath, config.CurrentConfig.ServerVirtualPath, -1)

	global.GlobalItemMap.Set(filename, html)
	html = strings.Replace(html.(string), SSOQuery, ssoQuery, -1)

	//与用户相关的页面差异不写入缓存
	for _, process := range processHtmlFuncs {
		html = process(html.(string))
	}
	return ctx.WriteHtml(html)
}
