package contract_expert

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"strings"
)

const (
	//SSO串
	SSOQueryString = "$SSOQueryString$"
	ExpertPageData = "$PageData$"
)

func isSSOString(str string) bool {
	return strings.Contains(str, "uid") &&
		strings.Contains(str, "pid") &&
		strings.Contains(str, "token") &&
		strings.Contains(str, "tid") &&
		strings.Contains(str, "sid")
}

func getSSOQueryString(ctx dotweb.Context) string {
	var url string
	referer := ctx.Request().Referer()
	if len(referer) == 0 || isSSOString(referer) == false {
		referer = ctx.Request().Url()
	}

	queryStr := ""
	questionMarkIndex := strings.Index(referer, "?")
	if questionMarkIndex >= 0 {
		queryStr += referer[questionMarkIndex+1:]
	}
	if strings.Index(url, "?") >= 0 {
		url = url + "&" + queryStr
	} else {
		url = url + "?" + queryStr
	}
	//去除打开注册弹窗的标识
	url = strings.Replace(url, "&openreg=1", "", -1)
	//去除打开指定活动的标识 使用无效的参数替换
	url = strings.Replace(url, "showActivity=", "ignore=", -1)
	return url
}

func RenderExpertHtml(ctx dotweb.Context, filename string) error {
	ssoStr := getSSOQueryString(ctx)
	processHtml := func(html string) string {
		html = strings.Replace(html, SSOQueryString, ssoStr, -1)
		return html
	}
	return contract.RenderHtml(ctx, filename, processHtml)
}

func RenderExpertHtmlWithPageData(ctx dotweb.Context, filename string, pageData interface{}) error {
	ssoStr := getSSOQueryString(ctx)
	processHtml := func(html string) string {
		html = strings.Replace(html, SSOQueryString, ssoStr, -1)
		if pageData != nil {
			html = strings.Replace(html, ExpertPageData, _json.GetJsonString(pageData), -1)
		}
		return html
	}
	return contract.RenderHtml(ctx, filename, processHtml)
}
