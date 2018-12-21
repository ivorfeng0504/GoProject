package myoptional

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"github.com/devfeel/dotweb"
)

// Index 微股吧和策略首页
func Index(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "zixuan.html")
}

// SaySomething 个股微股吧评论
func SaySomething(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "saysomething.html")
}


// IndexStatic 静态页面临时跳转
func IndexStatic(ctx dotweb.Context) error {
	return ctx.Redirect(302, "http://test.roboadvisor.emoney.cn/experttest/zixuan.html")
}

// RelatedArticle 相关资讯详情
func RelatedArticle(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "relatedArticle.html")
}

// BKNews 板块相关资讯
func BKNews(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "godeyesbknews.html")
}