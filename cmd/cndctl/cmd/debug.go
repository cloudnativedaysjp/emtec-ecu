package cmd

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

// debugCmd represents the track command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "debug for cnd-operation-server",
}

// debugListLocalmemCmd represents the "track list" command
var debugListLocalmemCmd = &cobra.Command{
	Use: "list-localmem",
	Run: func(cmd *cobra.Command, args []string) {
		logger := getLogger()
		ctx := logr.NewContext(context.Background(), logger)
		//
		// List Sharedmem
		//
		if err := createClient(cndOperationServerAddress); err != nil {
			logger.Error(err, "createClient was failed")
			os.Exit(1)
		}
		resp, err := debugClient.ListSharedmem(ctx, &emptypb.Empty{})
		if err != nil {
			logger.Error(err, "debugClient.ListSharedmem was failed")
			os.Exit(1)
		}
		//
		// Print
		//
		pp.Print(resp) // TODO
	},
}

func init() {
	debugCmd.AddCommand(debugListLocalmemCmd)
	rootCmd.AddCommand(debugCmd)
}
