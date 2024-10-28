package cmd

import (
	"github.com/gkwa/myher/core/gomod"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse go.mod dependencies",
	Long:  `Parse direct dependencies from go.mod file`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		svc := gomod.NewService(logger)

		info, err := svc.GetModuleInfo()
		if err != nil {
			logger.Error(err, "Failed to get module info")
			return
		}

		svc.PrettyPrint(info)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
