package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

// debugCmd represents the track command.
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "debug for emtec-ecu",
}

// debugListLocalmemCmd represents the "track list" command.
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
		for trackId, talks := range resp.TalksMap {
			fmt.Printf("--- # track %d: ", trackId)
			if disabled, ok := resp.DisabledMap[trackId]; !ok || !disabled {
				fmt.Println("Enabled")
			} else {
				fmt.Println("Disabled")
			}
			for _, talk := range talks.Talks {
				fmt.Printf("- Id: %v\n", talk.Id)
				fmt.Printf("  TalkName: %v\n", talk.TalkName)
				fmt.Printf("  TrackId: %v\n", talk.TrackId)
				fmt.Printf("  TrackName: %v\n", talk.TrackName)
				fmt.Printf("  EventAbbr: %v\n", talk.EventAbbr)
				fmt.Printf("  SpeakerNames: %v\n", talk.SpeakerNames)
				fmt.Printf("  Type: %v\n", talk.Type)
				fmt.Printf("  StartAt: %v\n", talk.StartAt)
				fmt.Printf("  EndAt: %v\n", talk.EndAt)
			}
		}
	},
}

func init() {
	debugCmd.AddCommand(debugListLocalmemCmd)
	rootCmd.AddCommand(debugCmd)
}
