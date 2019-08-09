package main

import (
	"context"
	"fmt"
	"github.com/discordapp/lilliput"
	"github.com/jolestar/go-commons-pool"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var RootCmd = &cobra.Command{
	Use: "freyja",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var objectPoolConfig = &pool.ObjectPoolConfig{
	LIFO:                     true,
	MaxTotal:                 -1,
	MaxIdle:                  100,
	MinIdle:                  10,
	MinEvictableIdleTime:     -1,
	SoftMinEvictableIdleTime: time.Minute*30,
	NumTestsPerEvictionRun:   100,
	EvictionPolicyName:       pool.DefaultEvictionPolicyName,
	EvitionContext:           context.Background(),
	TestOnCreate:             false,
	TestOnBorrow:             false,
	TestOnReturn:             false,
	TestWhileIdle:            false,
	TimeBetweenEvictionRuns:  time.Second*10,
	BlockWhenExhausted:       true}

var ServerCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {


		factory := pool.NewPooledObjectFactorySimple(
			func(context.Context) (interface{}, error) {
				return lilliput.NewImageOps(8192),
					nil
			})

		objpool := pool.NewObjectPool(context.Background(), factory, objectPoolConfig)

		fastServer := &fasthttp.Server{
			Handler: fastHTTPHandler,
			Name:    "freyja",
			//Handler: s3APiHandler(Bucket),
			GetOnly: true,
			DisableHeaderNamesNormalizing: true,
			DisableKeepalive:              true,
			MaxRequestBodySize:            0,
			//ReadBufferSize:                1 << 10,
			//ReadTimeout:                   time.Second * 1,
		}

		fmt.Println(fastServer)

		return nil
	},
}

func init() {
	ServerCmd.Flags().Bool("noVarnish", false, "Skip varnish cache? (DEVELOPMENT ONLY)")
	ServerCmd.Flags().IntP("varnishPort", "v", 80, "Local varnish port so varnish can cache outbound s3 requests")
	ServerCmd.Flags().StringP("host", "H", "0.0.0.0", "HTTP host to use (if not using unix socket)")
	ServerCmd.Flags().IntP("port", "p", 8081, "HTTP port to expose (if not using unix socket)")
	ServerCmd.Flags().StringP("unix", "u", "", "File path for unix socket (if not using HTTP port)")

	RootCmd.AddCommand(ServerCmd)
}


func main() {

	fmt.Println(runtime.NumCPU())
	fmt.Println(os.LookupEnv("http_proxy"))
	fmt.Println(os.LookupEnv("HTTP_PROXY"))
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}