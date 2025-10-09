package cli

import (
	"fmt"
	"net/http"

	"github.com/alex-cos/scvtemplate/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show application status",
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	cmgr := config.New(nil)
	// load config, if fail then panic
	cfg, err := cmgr.Load(configFile)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:%d/config", cfg.Admin.HTTPPort)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("status failed: %s", resp.Status)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status failed: %s", resp.Status)
	}

	fmt.Println("status is ok") // nolint: forbidigo
	return nil
}
