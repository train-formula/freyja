package freyja

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"time"
)

var faviconURI = []byte("/favicon.ico")
var healthCheckURI = []byte("/_health")

var etagEnclosing = "\""
var etagHeaderString = "ETag"
var etagHeader = []byte(etagHeaderString)

var expiresHeader = []byte("Expires")
var expiresDuration = time.Second * 2592000

var okCacheControlValue = []byte("public, must-revalidate, max-age=2592000")

var errorContentType = []byte("text/plain")

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	if bytes.Compare(ctx.RequestURI(), faviconURI) == 0 {
		ctx.SetStatusCode(fasthttp.StatusNotFound)

		return
	}


	if bytes.Compare(ctx.RequestURI(), healthCheckURI) == 0 {

		healthHandler(ctx)

		return
	}
}