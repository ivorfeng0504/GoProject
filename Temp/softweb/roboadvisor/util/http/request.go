package _http

import (
	"github.com/devfeel/dotweb"
	"net/http"
)

// GetRequestValue 从请求中获取值
func GetRequestValue(ctx dotweb.Context, key string) (value string) {
	if ctx.Request().Method == http.MethodGet {
		value = ctx.QueryString(key)
	} else {
		value = ctx.PostFormValue(key)
	}
	return value
}
