package cmd

import (
	"RssReader/depend"
	"RssReader/pkg/bootstrap"
	"context"
	"github.com/spf13/cobra"
)

var httpServerCommand = &cobra.Command{
	Use:   "api",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		i := bootstrap.Helper{}
		i.AddComponent(
			&depend.MongoDB{},
			&depend.Api{},
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
