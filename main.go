package main

import (
	"fmt"
	"github.com/train-formula/freyja/pool"

	"github.com/train-formula/freyja/handler"
	"github.com/valyala/fasthttp"
	"os"
	"runtime"
)

func main() {

	fmt.Println(runtime.NumCPU())
	fmt.Println(os.LookupEnv("http_proxy"))
	fmt.Println(os.LookupEnv("HTTP_PROXY"))
	runtime.GOMAXPROCS(runtime.NumCPU())

	pool := pool.NewLilliputPool(runtime.NumCPU(), 8192, 10*1024*1024)

	fasthandler := &handler.Handler{
		S3Host:"s3-us-west-2.amazonaws.com",
		NoVarnish:true,
		VarnishPort:"1234",
		ValidBuckets: [][]byte{
			[]byte("formula-tester"),
		},
		DefaultQuality:85,
		LilliputPool:pool,
	}
	fastServer := &fasthttp.Server{
		Handler: fasthandler.Handle,
		Name:    "freyja",
		//Handler: s3APiHandler(Bucket),
		GetOnly: true,
		DisableHeaderNamesNormalizing: true,
		DisableKeepalive:              true,
		MaxRequestBodySize:            0,
		//ReadBufferSize:                1 << 10,
		//ReadTimeout:                   time.Second * 1,
	}


	if err := fastServer.ListenAndServe("0.0.0.0:8081"); err != nil {
		panic(err)
	}


}