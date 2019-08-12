package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/discordapp/lilliput"
	"github.com/train-formula/freyja/parse"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (h *Handler) imageHandler(ctx *fasthttp.RequestCtx) {

	parsed, respStatus := parse.ParseURI(ctx.RequestURI(), h.ValidBuckets, h.DefaultQuality)

	if respStatus != fasthttp.StatusOK {
		h.errorHandler(ctx, respStatus, nil)

		return
	}

	stdReq, err := http.NewRequest("GET", string(append([]byte("http://"+h.S3Host), parsed.S3Uri()...)), nil)

	if err != nil {
		h.Logger.Error("Error creating http request", zap.Error(err))
		h.errorHandler(ctx, fasthttp.StatusInternalServerError, s3InternalServerError)

		return
	}

	stdReq.Header.Set("Host", h.S3Host)
	stdReq.Host = h.S3Host

	downloadStart := time.Now()

	stdResp, stdRespErr := http.DefaultClient.Do(stdReq)

	if stdRespErr != nil || stdResp == nil || stdResp.Body == nil {

		if stdResp != nil && stdResp.Body != nil {
			stdResp.Body.Close()
		}

		h.errorHandler(ctx, fasthttp.StatusInternalServerError, internalServerError)

		return
	}

	s3Status := stdResp.StatusCode

	if stdResp.ContentLength > h.MaxContentLength {
		h.errorHandler(ctx, fasthttp.StatusBadRequest, []byte("Image too large"))
		return
	}

	if _, statusOK := okS3Statuses[s3Status]; !statusOK {

		if _, statusNotFound := notFoundS3Statuses[s3Status]; statusNotFound {
			h.errorHandler(ctx, s3Status, []byte("Image not found"))
		} else {
			h.errorHandler(ctx, s3Status, []byte("Error from S3"))
		}

		return

	}

	buf := bytebufferpool.Get()
	_, err = io.Copy(buf, stdResp.Body)
	if err != nil {
		stdResp.Body.Close()
		h.Logger.Error("Error reading image", zap.Error(err))

		h.errorHandler(ctx, fasthttp.StatusInternalServerError, internalServerError)

		return
	}

	downloadElapsed := time.Since(downloadStart)

	stdResp.Body.Close()

	if int64(buf.Len()) > h.MaxContentLength {
		h.errorHandler(ctx, fasthttp.StatusBadRequest, []byte("Image too large"))
		return
	}

	var s3Etag []byte

	if stdResp.Header != nil {
		s3Etag = []byte(stdResp.Header.Get(etagHeaderString))
	}

	processingStart := time.Now()

	decoder, err := lilliput.NewDecoder(buf.Bytes())
	if err != nil {
		h.Logger.Error("Error decoding image", zap.Error(err))
		h.errorHandler(ctx, fasthttp.StatusInternalServerError, []byte("Error constructing image decoder"))
		return
	}

	opsPair := h.LilliputPool.Get()

	final, err := opsPair.Ops.Transform(decoder, &lilliput.ImageOptions{
		FileType:             ".jpeg",
		Width:                int(parsed.Opts.Width),
		Height:               int(parsed.Opts.Height),
		ResizeMethod:         lilliput.ImageOpsFit,
		NormalizeOrientation: false,
		EncodeOptions: map[int]int{
			lilliput.JpegProgressive: 1,
			lilliput.JpegQuality:     90,
		},
	}, opsPair.Buf)

	decoder.Close()
	bytebufferpool.Put(buf)

	if err != nil {
		h.Logger.Error("Error transforming image", zap.Error(err))
		h.errorHandler(ctx, fasthttp.StatusInternalServerError, []byte("Error transforming image"))
		return

	}

	processingElapsed := time.Since(processingStart)

	h.Logger.Debug("Image processing performance", zap.Duration("download", downloadElapsed), zap.Duration("processing", processingElapsed))

	h.imageResponseHandler(ctx, final, s3Etag)

	h.LilliputPool.Put(opsPair)
}
