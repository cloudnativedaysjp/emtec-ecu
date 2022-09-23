package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/obsws"
	pb "github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/schema"
)

// sceneCmd represents the "scene" command
var sceneCmd = &cobra.Command{
	Use:   "scene",
	Short: "operate OBS scene",
}

// sceneListCmd represents the "scene list" command
var sceneListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		logger := getLogger()
		ctx := logr.NewContext(context.Background(), logger)
		//
		// Get Scenes
		//
		var scenes []*pb.Scene
		if directly {
			//
			// directly
			//
			obswsClient, err := obsws.NewObsWebSocketClient(obsHost, obsPassword)
			if err != nil {
				logger.Error(err, "NewObsWebSocketClient was failed")
				os.Exit(1)
			}
			ss, err := obswsClient.ListScenes(ctx)
			if err != nil {
				logger.Error(err, "obswsClient.ListScenes was failed")
				os.Exit(1)
			}
			for _, s := range ss {
				scenes = append(scenes, &pb.Scene{
					Name:             s.Name,
					SceneIndex:       int32(s.SceneIndex),
					IsCurrentProgram: s.IsCurrentProgram,
				})
			}
		} else {
			//
			// via cnd-operation-server
			//
			if err := createClient(cndOperationServerAddress); err != nil {
				logger.Error(err, "createClient was failed")
				os.Exit(1)
			}
			resp, err := sceneClient.ListScene(ctx, &pb.ListSceneRequest{TrackId: trackId})
			if err != nil {
				logger.Error(err, "sceneClient.ListScene was failed")
				os.Exit(1)
			}
			scenes = resp.Scene
		}
		//
		// Print
		//
		for _, scene := range scenes {
			if scene.IsCurrentProgram {
				fmt.Printf("> [%2d] %s\n", scene.SceneIndex, scene.Name)
			} else {
				fmt.Printf("  [%2d] %s\n", scene.SceneIndex, scene.Name)
			}
		}
	},
}

// sceneListCmd represents the "scene list" command
var sceneNextCmd = &cobra.Command{
	Use: "next",
	Run: func(cmd *cobra.Command, args []string) {
		logger := getLogger()
		ctx := logr.NewContext(context.Background(), logger)
		//
		// Update Scene
		//
		if directly {
			//
			// directly
			//
			obswsClient, err := obsws.NewObsWebSocketClient(obsHost, obsPassword)
			if err != nil {
				logger.Error(err, "NewObsWebSocketClient was failed")
				os.Exit(1)
			}
			if err := obswsClient.MoveSceneToNext(ctx); err != nil {
				logger.Error(err, "obswsClient.MoveSceneToNext was failed")
				os.Exit(1)
			}
		} else {
			//
			// via cnd-operation-server
			//
			if err := createClient(cndOperationServerAddress); err != nil {
				logger.Error(err, "createClient was failed")
				os.Exit(1)
			}
			if _, err := sceneClient.MoveSceneToNext(ctx, &pb.MoveSceneToNextRequest{TrackId: trackId}); err != nil {
				logger.Error(err, "sceneClient.MoveSceneToNext was failed")
				os.Exit(1)
			}
		}
		//
		// Print
		//
		fmt.Println("Done")
	},
}

func init() {
	sceneCmd.PersistentFlags().Int32VarP(&trackId, "track-id", "t", 0, "Track ID on Dreamkast")
	_ = sceneCmd.MarkPersistentFlagRequired("track-id")

	sceneCmd.AddCommand(sceneListCmd)
	sceneCmd.AddCommand(sceneNextCmd)
	rootCmd.AddCommand(sceneCmd)
}
