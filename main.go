package freyja

import (
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

var RootCmd = &cobra.Command{
	Use: "freyja",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var ServerCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fastServer := &fasthttp.Server{
			Handler: fastHTTPHandler,
			Name:    "apollo",
			//Handler: s3APiHandler(bucket),
			GetOnly: true,
			DisableHeaderNamesNormalizing: true,
			DisableKeepalive:              true,
			MaxRequestBodySize:            0,
			//ReadBufferSize:                1 << 10,
			//ReadTimeout:                   time.Second * 1,
		}
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
