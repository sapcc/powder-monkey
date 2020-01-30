package cmd

import (
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/powder-monkey/dynomite"

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore [prefix]",
	Short: "Triggers dynomite backend restore of the last dump file with the specified prefix",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := dynomite.Restore(container, args[0])
		if err != nil {
			logg.Fatal(err.Error())
		}
	},
}

var listBackups = &cobra.Command{
	Use:   "listbackups [prefix]",
	Short: "Lists the available dump files with the specified prefix",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := dynomite.ListBackups(container, args[0], 25)
		if err != nil {
			logg.Error(err.Error())
		}
	},
}

func init() {
	restoreCmd.PersistentFlags().StringVar(&container, "container", "db_backup", "Container to store the backup")

	backendCmd.AddCommand(restoreCmd)
	backendCmd.AddCommand(listBackups)
}
