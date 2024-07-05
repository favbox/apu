package command

import (
	"net/http"
	"os"
	"slices"

	"apu/cmd/wechat-proxy/addon"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	addr       string
	log        zerolog.Logger
	allowHosts = []string{
		"mp.weixin.qq.com:443",
		"i.weread.qq.com:443",
	}
)

var rootCmd = &cobra.Command{
	Use:   "微信网络代理",
	Short: "拦截微信/读书会话，以获取文章列表及阅读量请求权",
	Run:   startProxy(),
}

func init() {
	rootCmd.Flags().StringVarP(&addr, "addr", "a", ":9090", "")

	addr = ":9090"
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(func() zerolog.Level {
			if os.Getenv("DEBUG") != "" {
				return zerolog.DebugLevel
			} else {
				return zerolog.InfoLevel
			}
		}()).
		With().
		Timestamp().
		Logger()
}

func startProxy() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		opts := &proxy.Options{Addr: addr}
		p, err := proxy.NewProxy(opts)
		if err != nil {
			log.Fatal().Err(err)
		}

		p.SetShouldInterceptRule(func(req *http.Request) bool {
			return slices.Contains(allowHosts, req.Host)
		})

		p.AddAddon(&addon.WechatAddon{}) // 拦截微信请求
		p.AddAddon(&addon.WereadAddon{}) // 拦截微信读书请求

		err = p.Start()
		if err != nil {
			log.Fatal().Err(err)
		}
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err)
	}
}
