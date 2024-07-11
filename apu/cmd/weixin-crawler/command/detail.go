package command

import "github.com/spf13/cobra"

var detailCmd = &cobra.Command{
	Use:   "detail",
	Short: "抓取文章详情",
}
