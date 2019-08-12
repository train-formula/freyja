package handler

import "github.com/valyala/fasthttp"

func (h *Handler) imageResponseHandler(ctx *fasthttp.RequestCtx, imageBuf []byte, etag []byte) {

	now := ctx.Time().UTC()
	expires := now.Add(expiresDuration)

	ctx.Response.Header.SetContentType("image/jpeg")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetLastModified(now)
	ctx.Response.Header.SetContentLength(len(imageBuf))

	// Cache-Control
	ctx.Response.Header.SetCanonical(cacheControlHeader, okCacheControlValue)

	// Expires
	var expiresDst []byte
	expiresDst = fasthttp.AppendHTTPDate(expiresDst, expires)
	ctx.Response.Header.SetCanonical(expiresHeader, expiresDst)

	if len(etag) > 0 {
		ctx.Response.Header.SetCanonical(etagHeader, etag)
	}

	ctx.SetBody(imageBuf)

	return

}
