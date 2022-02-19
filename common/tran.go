package common

import (
	"context"
	"strconv"

	"github.com/hqbobo/frame/common/log"

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
