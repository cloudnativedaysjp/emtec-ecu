package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/schema"
)

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "about OBS Studios which cnd-operation-server connects to",
}

// trackListCmd represents the "track list" command
var trackListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		logger := getLogger()
		ctx := logr.NewContext(context.Background(), logger)
		//
		// Get Tracks
		//
		if err := createClient(cndOperationServerAddress); err != nil {
			logger.Error(err, "createClient was failed")
			os.Exit(1)
		}
		resp, err := trackClient.ListTrack(ctx, &emptypb.Empty{})
		if err != nil {
			logger.Error(err, "trackClient.ListTrack was failed")
			os.Exit(1)
		}
		//
		// Print
		//
		for _, track := range resp.Tracks {
			fmt.Printf("%s (trackId: %d)\n", track.ObsHost, track.TrackId)
		}
	},
}

var trackEnableCmd = &cobra.Command{
	Use: "enable",
	Run: func(cmd *cobra.Command, args []string) {
		logger := getLogger()
		ctx := logr.NewContext(context.Background(), logger)
		//
		// Enable Track
		//
		if err := createClient(cndOperationServerAddress); err != nil {
			logger.Error(err, "createClient was failed")
			os.Exit(1)
		}
		if _, err := trackClient.EnableAutomation(
			ctx, &pb.SwitchAutomationRequest{TrackId: trackId},
		); err != nil {
			logger.Error(err, "trackClient.EnableAutomation was failed")
			os.Exit(1)
		}
		//
		// Print
		//
		fmt.Println("Done")
	},
}

var trackDisableCmd = &cobra.Command{
	Use: "disable",
	Run: func(cmd *cobra.Command, args []string) {
		logger := getLogger()
		ctx := logr.NewContext(context.Background(), logger)
		//
		// Disable Track
		//
		if err := createClient(cndOperationServerAddress); err != nil {
			logger.Error(err, "createClient was failed")
			os.Exit(1)
		}
		if _, err := trackClient.DisableAutomation(
			ctx, &pb.SwitchAutomationRequest{TrackId: trackId},
		); err != nil {
			logger.Error(err, "trackClient.DisableAutomation was failed")
			os.Exit(1)
		}
		//
		// Print
		//
		fmt.Println("Done")
	},
}

func init() {
	trackEnableCmd.PersistentFlags().Int32VarP(&trackId, "track-id", "t", 0, "Track ID on Dreamkast")
	_ = trackEnableCmd.MarkPersistentFlagRequired("track-id")

	trackCmd.AddCommand(trackListCmd)
	trackCmd.AddCommand(trackEnableCmd)
	trackCmd.AddCommand(trackDisableCmd)
	rootCmd.AddCommand(trackCmd)
}
