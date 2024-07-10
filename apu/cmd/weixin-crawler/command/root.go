package command

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var log zerolog.Logger

var rootCmd = &cobra.Command{
	Use:   "weixin-crawler",
	Short: "采集微信公众号文章列表、详情及阅读量",
}

func init() {
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
