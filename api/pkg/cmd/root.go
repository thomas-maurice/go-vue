package cmd

import "github.com/spf13/cobra"

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts the api server",
	Long:  "",
}

func init() {
	initGenKeyCmd()
	initServerCmd()
	initHashPassCmd()

	rootCmd.AddCommand(genKeyCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(hashPassCmd)
}
