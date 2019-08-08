package main

import "github.com/valyala/fasthttp"

var healthResp = []byte("OK")
var cacheControlHeader = []byte("Cache-Control")
var noCacheControlValue = []byte("private, no-cache, no-store, must-revalidate, max-age=0")
var pragmaHeader = []byte("Pragma")
var noCachePragmaValue = []byte("no-cache")

func healthHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentLength(len(healthResp))

	ctx.Response.Header.SetLastModified(ctx.Time().UTC())

	// We never want to cache health checks
	// Cache-Control
	ctx.Response.Header.SetCanonical(cacheControlHeader, noCacheControlValue)
	// Pragma
	ctx.Response.Header.SetCanonical(pragmaHeader, noCachePragmaValue)

	ctx.SetBody(healthResp)

	return
}
