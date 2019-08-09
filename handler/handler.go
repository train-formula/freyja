package handler

import (
	"bytes"
	"context"
	"github.com/discordapp/lilliput"
	"github.com/jolestar/go-commons-pool"
	"github.com/train-formula/freyja/parse"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"time"
)

var okS3Statuses = map[int]struct{}{
	200: struct{}{},
	304: struct{}{},
	201: struct{}{},
	204: struct{}{},
	206: struct{}{},
}


type Handler struct {
	S3Host string
	NoVarnish bool
	VarnishPort string
	ValidBuckets [][]byte
	DefaultQuality uint32
	LilliputPool *pool.ObjectPool
}

func (h *Handler) getLilliputOps(ctx context.Context) (*lilliput.ImageOps, error) {
	obj, err := h.LilliputPool.BorrowObject(ctx)
	if err != nil {
		return nil, err
	}

	return obj.(*lilliput.ImageOps), nil
}

func (h *Handler) releaseLilliputOps(ctx context.Context, ops *lilliput.ImageOps) error {
	return h.LilliputPool.ReturnObject(ctx, ops)
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

	//https://s3-us-west-2.amazonaws.com/formula-tester/Image.jpeg
	s3Uri := []byte("formula-tester/Image.jpeg")
	if h.NoVarnish {
		stdReq, stdReqErr = http.NewRequest("GET", string(append([]byte("http://"+h.S3Host), s3Uri...)), nil)
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

		if stdResp != nil && stdResp.Body != nil {
			stdResp.Body.Close()
		}

		h.errorHandler(ctx, fasthttp.StatusInternalServerError, internalServerError)

		return
	}

	s3Status := stdResp.StatusCode

	if _, statusOK := okS3Statuses[s3Status]; !statusOK {

		h.errorHandler(ctx, s3Status, []byte("Error from S3"))

		return

	}

	respBuf, respBufErr := ioutil.ReadAll(stdResp.Body)
	if respBufErr != nil {
		stdResp.Body.Close()

		h.errorHandler(ctx, fasthttp.StatusInternalServerError, internalServerError)

		return
	}

	stdResp.Body.Close()

	var s3Etag []byte

	if stdResp.Header != nil {
		s3Etag = []byte(stdResp.Header.Get(etagHeaderString))
	}

	decoder, err := lilliput.NewDecoder(respBuf)
	if err != nil {
		h.errorHandler(ctx, fasthttp.StatusInternalServerError, []byte("Error constructing image decoder"))
		return
	}

	ops, err := h.getLilliputOps(ctx)
	if err != nil {
		h.errorHandler(ctx, fasthttp.StatusInternalServerError, []byte("Error retrieving image operation from pool"))
		return
	}
}