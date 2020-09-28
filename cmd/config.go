package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(showConfigCmd)
}

var showConfigCmd = &cobra.Command{
	Use:   "show-config [group]",
	Short: "Show config of a group",
	Long:  `Show config of a group`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Config for %s: %v\n", args[0], viper.Get(args[0]))
	},
}
