package _error

import (
	"github.com/devfeel/dotweb"
	"emoney.cn/fundchannel/contract"
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
