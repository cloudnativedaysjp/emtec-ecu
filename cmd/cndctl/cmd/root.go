package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cndctl",
	Short: "CLI tool for control OBS using WebSocket",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	cndOperationServerAddress string
	directly                  bool
	obsHost                   string
	obsPassword               string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cndOperationServerAddress, "emtec-ecu-address", "d", "127.0.0.1:20080", `Address of emtec-ecu (format: "<host>:<port>")`)
	rootCmd.PersistentFlags().BoolVarP(&directly, "directly", "", false, `If this flag is true, CLI connect to OBS by WebSocket instead of connecting to emtec-ecu`)
	rootCmd.PersistentFlags().StringVarP(&obsHost, "obs-host", "", "", ``)
	rootCmd.PersistentFlags().StringVarP(&obsPassword, "obs-password", "", "", ``)

	rootCmd.MarkFlagsRequiredTogether("directly", "obs-host", "obs-password")
}
