package cmd

import (
	"time"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/powder-monkey/dynomite"

	"github.com/spf13/cobra"
)

var (
	container string
	every     time.Duration
)

var backupCmd = &cobra.Command{
	Use:   "backup [prefix]",
	Short: "Trigger dynomite backend backup and upload to dump file with the specified prefix",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dyno := dynomite.NewDynomiteRedis(dynomiteHost, dynomitePort, backendPort, backendPassword)

		if every != 0 {
			err := dyno.BackupEvery(every, container, args[0])
			if err != nil {
				logg.Fatal(err.Error())
			}
		} else {
			err := dyno.Backup(container, args[0])
			if err != nil {
				logg.Fatal(err.Error())
			}
		}

		logg.Info("Backup succesful")
	},
}

func init() {
	backupCmd.PersistentFlags().StringVar(&container, "container", "db_backup", "Container to store the backup")
	backupCmd.PersistentFlags().DurationVar(&every, "every", 0, "Run the backup periodicly")

	backendCmd.AddCommand(backupCmd)
}
