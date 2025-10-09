package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "apppname",
	Short: "A flexible web server template with hot reload and observability",
	Long: `apppname is a Go-based web service template that supports:
- dynamic configuration reloads
- multiple HTTP servers
- observability via Prometheus
- TLS reloads
- structured logging
- graceful shutdown`,
	RunE: runRoot,
}

func init() {
	rootCmd.Flags().StringVarP(&configFile, "configFile", "c", "", "Configuration file to use")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) // nolint: forbidigo
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return serveCmd.RunE(cmd, args)
	}
	return nil
}
