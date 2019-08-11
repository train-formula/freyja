package handler

import (
	"bytes"
	"fmt"
	"github.com/discordapp/lilliput"
	"github.com/train-formula/freyja/pool"

	"github.com/train-formula/freyja/parse"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
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
	S3Host string
	NoVarnish bool
	VarnishPort string
	ValidBuckets [][]byte
	DefaultQuality uint32
	LilliputPool *pool.LilliputPool
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


	parsed, respStatus := parse.ParseURI(ctx.RequestURI(), h.ValidBuckets, h.DefaultQuality)

	if respStatus != fasthttp.StatusOK {
		h.errorHandler(ctx, respStatus, nil)

		return
	}

	var stdReq *http.Request
	var stdReqErr error

	if h.NoVarnish {
		stdReq, stdReqErr = http.NewRequest("GET", string(append([]byte("http://"+h.S3Host), parsed.S3Uri()...)), nil)
	} else {
		stdReq, stdReqErr = http.NewRequest("GET", string(append([]byte("http://localhost:"+h.VarnishPort), s3Uri...)), nil)
	}

	if stdReqErr != nil {
		h.errorHandler(ctx, fasthttp.StatusInternalServerError, s3InternalServerError)

		return
	}



	stdReq.Header.Set("Host", h.S3Host)
	stdReq.Host = h.S3Host

	stdResp, stdRespErr := http.DefaultClient.Do(stdReq)

	if stdRespErr != nil || stdResp == nil || stdResp.Body == nil {

		fmt.Println(stdRespErr)
		if stdResp != nil && stdResp.Body != nil {
			stdResp.Body.Close()
		}

		h.errorHandler(ctx, fasthttp.StatusInternalServerError, internalServerError)

		return
	}

	fmt.Println("OK 2")

	s3Status := stdResp.StatusCode

	fmt.Println(stdResp.ContentLength)

	if _, statusOK := okS3Statuses[s3Status]; !statusOK {

		if _, statusNotFound := notFoundS3Statuses[s3Status]; statusNotFound {
			h.errorHandler(ctx, s3Status, []byte("Image not found"))
		} else {
			h.errorHandler(ctx, s3Status, []byte("Error from S3"))
		}


		return

	}

	buf := bytebufferpool.Get()
	_, err := io.Copy(buf, stdResp.Body)
	if err != nil {
		stdResp.Body.Close()

		h.errorHandler(ctx, fasthttp.StatusInternalServerError, internalServerError)

		return
	}

	stdResp.Body.Close()

	var s3Etag []byte

	if stdResp.Header != nil {
		s3Etag = []byte(stdResp.Header.Get(etagHeaderString))
	}

	decoder, err := lilliput.NewDecoder(buf.Bytes())
	if err != nil {
		h.errorHandler(ctx, fasthttp.StatusInternalServerError, []byte("Error constructing image decoder"))
		return
	}

	opsPair := h.LilliputPool.Get()

	final, err := opsPair.Ops.Transform(decoder, &lilliput.ImageOptions{
		FileType:".jpeg",
		Width:int(parsed.Opts.Width),
		Height:int(parsed.Opts.Height),
		ResizeMethod:lilliput.ImageOpsFit,
		NormalizeOrientation:false,
		EncodeOptions: map[int]int{
			lilliput.JpegProgressive: 1,
			lilliput.JpegQuality: 100,
		},

	}, opsPair.Buf)

	decoder.Close()
	bytebufferpool.Put(buf)

	if err != nil {
		fmt.Println(err)
		h.errorHandler(ctx, fasthttp.StatusInternalServerError, []byte("Error transforming image"))
		return

	}

	h.imageHandler(ctx, final, s3Etag)

	h.LilliputPool.Put(opsPair)
}