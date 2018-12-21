package _error

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"github.com/devfeel/dotweb"
)

// NotFound 404
func NotFound(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "404.html")
}

// ServerError 500
func ServerError(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "500.html")
}

// NotAuth 401
func NotAuth(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "401.html")
}
