package handler

import "github.com/valyala/fasthttp"

func (h *Handler)  healthHandler(ctx *fasthttp.RequestCtx) {
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
