package handler

import (
	"bytes"

	"go.uber.org/zap"

	"github.com/train-formula/freyja/pool"

	"github.com/valyala/fasthttp"
)

var okS3Statuses = map[int]struct{}{
	200: struct{}{},
	304: struct{}{},
	201: struct{}{},
	204: struct{}{},
	206: struct{}{},
}

var notFoundS3Statuses = map[int]struct{}{
	// Can indicates not found if the bucket itself is private
	403: struct{}{},
	404: struct{}{},
}

type Handler struct {
	Logger           *zap.Logger
	MaxContentLength int64
	S3Host           string
	ValidBuckets     [][]byte
	DefaultQuality   uint32
	LilliputPool     *pool.LilliputPool
}

func (h *Handler) Handle(ctx *fasthttp.RequestCtx) {
	if bytes.Compare(ctx.RequestURI(), faviconURI) == 0 {
		ctx.SetStatusCode(fasthttp.StatusNotFound)

		return
	}

	if bytes.Compare(ctx.RequestURI(), healthCheckURI) == 0 {

		h.healthHandler(ctx)

		return
	}

	h.imageHandler(ctx)

}
