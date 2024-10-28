package cmd

import (
	"fmt"

	"github.com/gkwa/myher/core/gomod"
	"github.com/spf13/cobra"
)

var (
	concurrent          int
	alternatingComments bool
)

var downgradeCmd = &cobra.Command{
	Use:   "downgrade",
	Short: "Generate downgrade commands for dependencies",
	Long:  `Generates commands to downgrade each direct dependency to its second-latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		svc := gomod.NewService(logger)

		commands, err := svc.GenerateDowngradeCommands(concurrent, alternatingComments)
		if err != nil {
			logger.Error(err, "Failed to generate downgrade commands")
			return
		}

		for _, cmd := range commands {
			fmt.Println(cmd)
		}
	},
}

func init() {
	rootCmd.AddCommand(downgradeCmd)
	downgradeCmd.Flags().IntVarP(&concurrent, "concurrent", "c", 5, "number of concurrent version checks")
	downgradeCmd.Flags().BoolVar(&alternatingComments, "enable-alternating-comments", false, "enable alternating commented commands")
}
