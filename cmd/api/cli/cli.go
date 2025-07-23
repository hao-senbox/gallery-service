package cli

import (
	"github.com/spf13/cobra"
	"sync"
)

var (
	configPath string
	verbose    bool
)

var cmd = &cobra.Command{
	Use:   "cmd",
	Short: "ShortDescription ",
	Long:  `Long Description`,
}

func Execute() error {
	return cmd.Execute()
}

var once sync.Once

func init() {
	once.Do(func() {
		// Here you will define your flags and configuration settings.
		// Cobra supports persistent flags, which, if defined here,
		// will be global for your application.
		cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yml", "path to config file")
		cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")

		// Cobra also supports local flags, which will only run
		// when this action is called directly.
		cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	})
}
