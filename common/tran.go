package common

import (
	"context"
	"github.com/hqbobo/frame/common/log"
	"strconv"

	"github.com/gin-gonic/gin"
	goTracing "github.com/opentracing/opentracing-go"
)

const (
	commonTransit = "Transid"
	commonSpan    = "TRACESPAN"
)

// GetSpanFromGin 从gin里面获取Span
func GetSpanFromGin(c *gin.Context) context.Context {
	span, _ := c.Get(commonSpan)
	return goTracing.ContextWithSpan(c, span.(goTracing.Span))
}

// SetSpanForGin 设置span到gin
func SetSpanForGin(c *gin.Context) goTracing.Span {
	span := goTracing.GlobalTracer().StartSpan(c.Request.URL.Path)
	c.Set(commonSpan, span)
	return span
}

// GetUserInfo 获取用户新
func GetUserInfo(c *gin.Context) (interface{}, bool) {
	return c.Get("UserInfo")
}

// SetUserInfo 设置用户信息
func SetUserInfo(c *gin.Context, userinfo interface{}) {
	c.Set("UserInfo", userinfo)
}

// GetAdminInfo 获取admin
func GetAdminInfo(c *gin.Context) (interface{}, bool) {
	return c.Get("AdminInfo")
}

// SetAdminInfo 设置admin信息
func SetAdminInfo(c *gin.Context, adminInfo interface{}) {
	c.Set("AdminInfo", adminInfo)
}

func QueryPaging(ctx *gin.Context) (page, limit int64, e error) {
	var err error
	i := ctx.Query("page")
	page, e = strconv.ParseInt(i, 10, 64)
	if err != nil {
		log.Warnln(e)
	}
	i = ctx.Query("limit")
	limit, e = strconv.ParseInt(i, 10, 64)
	if err != nil {
		log.Warnln(e)
	}
	return page, limit, e
}
