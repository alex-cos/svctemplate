package cli

import (
	"log/slog"
	"os"

	"github.com/alex-cos/scvtemplate/config"
	"github.com/alex-cos/scvtemplate/internal/app"
	"github.com/alex-cos/scvtemplate/pkg/dynamicLevel"
	"github.com/alex-cos/scvtemplate/pkg/logx"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	RunE:  runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	// logger with dynamic level
	dynLevel := &dynamicLevel.DynamicLevel{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       dynLevel,
		ReplaceAttr: nil,
	}))
	slog.SetDefault(logger)
	logx.Set(logger)

	cmgr := config.New(logger)
	defer cmgr.Close()

	// load config, if fail then panic
	cfg, err := cmgr.Load(configFile)
	if err != nil {
		return err
	}
	dynLevel.SetLevel(dynamicLevel.ParseLogLevel(cfg.LogLevel))
	logx.L().Info("initial configuration loaded", slog.Any("config", cfg))

	return app.Execute(cmgr, dynLevel)
}
