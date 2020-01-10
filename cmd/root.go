package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	dynomiteHost string
	dynomitePort int16
)

var rootCmd = &cobra.Command{
	Use:   "powder-monkey",
	Short: "powder-monkey manages your dynomite instances",
	Long: `powder-monkey manages your dynomite instances
including the backend database.
Supports warmups and backups.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dynomiteHost, "dynomite-host", "d", "localhost", "dynomite host address")
	rootCmd.PersistentFlags().Int16VarP(&dynomitePort, "dynomite-port", "p", 22222, "dynomite admin port")
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
