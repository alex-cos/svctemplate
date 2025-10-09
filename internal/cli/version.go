package cli

import (
	"fmt"

	"github.com/alex-cos/scvtemplate/version"
	"github.com/spf13/cobra"
)

// nolint: forbidigo
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s (built on %s)\n",
			version.GetVersion(),
			version.GetBuildDate())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
