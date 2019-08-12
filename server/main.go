package main

import (
	"strconv"

	"github.com/train-formula/freyja/pool"
	"github.com/train-formula/october"
	"go.uber.org/zap"

	"runtime"

	"github.com/train-formula/freyja/handler"
	"github.com/valyala/fasthttp"
)

func main() {

	october.MustInitServiceFromEnv()

	configurator := october.NewEnvConfigurator()

	cfg := &Config{}

	configurator.MustDecodeEnv(cfg, "")

	runtime.GOMAXPROCS(runtime.NumCPU())

	pool := pool.NewLilliputPool(runtime.NumCPU(), 8192, int(cfg.MaxContentLength*2))

	fasthandler := &handler.Handler{
		Logger:           zap.L(),
		MaxContentLength: cfg.MaxContentLength,
		S3Host:           "s3-us-west-2.amazonaws.com",
		ValidBuckets:     cfg.MustExtractValidBuckets(),
		DefaultQuality:   85,
		LilliputPool:     pool,
	}

	fastServer := &fasthttp.Server{
		Handler: fasthandler.Handle,
		Name:    "freyja",
		//Handler: s3APiHandler(Bucket),
		GetOnly:                       true,
		DisableHeaderNamesNormalizing: true,
		DisableKeepalive:              true,
		MaxRequestBodySize:            0,
		//ReadBufferSize:                1 << 10,
		//ReadTimeout:                   time.Second * 1,
	}

	if err := fastServer.ListenAndServe(cfg.Host + ":" + strconv.Itoa(cfg.Port)); err != nil {
		panic(err)
	}

}
