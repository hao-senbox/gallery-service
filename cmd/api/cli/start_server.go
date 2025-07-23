package cli

import (
	"gallery-service/internal/application"
	"github.com/spf13/cobra"
)

const StartServerCommand = "start-server"
const VersionServer = "1.0.0"

var httpServer = &cobra.Command{
	Use:     StartServerCommand,
	Short:   "Start server",
	Version: VersionServer,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		app, appErr := application.New(configPath)
		if appErr != nil {
			return appErr
		}

		return app.Run()
	},
}

func init() {
	cmd.AddCommand(httpServer)
}
