package cmd

import (
	"RssReader/depend"
	"RssReader/pkg/bootstrap"
	"context"
	"github.com/spf13/cobra"
)

var feedSpiderCommand = &cobra.Command{
	Use:   "spider",
	Short: "start feed spider",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		i := bootstrap.Helper{}
		i.AddComponent(
			&depend.Redis{},
			&depend.MongoDB{},
			&depend.Spider{},
		)
		err := i.Init(ctx)
		if err != nil {
			panic(err)
			return
		}

		forever := make(chan struct{})
		<-forever
	},
}
