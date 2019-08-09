package handler

import "github.com/valyala/fasthttp"

var errorContentType = []byte("text/plain")

func (h *Handler) errorHandler(ctx *fasthttp.RequestCtx, statusCode int, body []byte) {
	ctx.Response.Header.SetContentTypeBytes(errorContentType)
	ctx.Response.Header.SetStatusCode(statusCode)
	ctx.Response.Header.SetContentLength(0)
	ctx.Response.Header.SetLastModified(ctx.Time().UTC())

	// Cache-Control
	ctx.Response.Header.SetCanonical(cacheControlHeader, noCacheControlValue)
	// Pragma
	ctx.Response.Header.SetCanonical(pragmaHeader, noCachePragmaValue)

	if len(body) > 0 {
		ctx.SetBody(body)
	}
}