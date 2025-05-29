package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thomas-maurice/api/go-vue/pkg/api"
)

var (
	flagConfigFile string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the server",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := api.NewAPI(flagConfigFile)
		if err != nil {
			return err
		}

		return a.Run()
	},
}

func initServerCmd() {
	serverCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", "config.yaml", "Path to the configuration file")
}
