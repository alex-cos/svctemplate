package cli

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/alex-cos/scvtemplate/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reloadCmd)
}

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Trigger a configuration reload via the admin API",
	RunE:  runReload,
}

func runReload(cmd *cobra.Command, args []string) error {
	cmgr := config.New(nil)
	// load config, if fail then panic
	cfg, err := cmgr.Load(configFile)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:%d/admin/reload", cfg.Admin.HTTPPort)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		return fmt.Errorf("failed to call reload endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("reload failed: %s", resp.Status)
	}

	fmt.Println("Configuration reload triggered successfully") // nolint: forbidigo
	return nil
}
