package cmd

import (
	"fmt"
	"time"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/powder-monkey/dynomite"

	"github.com/spf13/cobra"
)

var (
	backendPort       int16
	backendPassword   string
	masterBackend     string
	masterBackendPort int16
	acceptedDiff      int64
	timeoutMinutes    int
	replicaAnnounceIP string
)

var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Interact with dynomite backend",
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping dynomite backend",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := rootCmd.PersistentFlags().GetString("dynomite-host")
		redis := dynomite.NewRedis(host, backendPort, backendPassword)

		state, err := redis.Ping()
		if err != nil {
			logg.Fatal(err.Error())
		}

		fmt.Printf("Dynomite backend [%s] ping: %v\n", host, state)
	},
}

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Get Role of dynomite backend",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := rootCmd.PersistentFlags().GetString("dynomite-host")
		redis := dynomite.NewRedis(host, backendPort, backendPassword)

		role, err := redis.Role()
		if err != nil {
			logg.Fatal(err.Error())
		}

		fmt.Printf("Dynomite backend [%s]: role %s\n", host, role)
	},
}

var warmupCmd = &cobra.Command{
	Use:   "warmup [master]",
	Short: "Warmup dynomite backend",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := rootCmd.PersistentFlags().GetString("dynomite-host")
		port, _ := rootCmd.PersistentFlags().GetInt16("dynomite-port")
		dyno := dynomite.NewDynomiteRedis(host, port, backendPort, backendPassword)
		master := dynomite.NewRedis(args[0], masterBackendPort, backendPassword)
		slaveHost := host
		if replicaAnnounceIP != "" {
			slaveHost = replicaAnnounceIP
		}

		result, err := dyno.Warmup(*master, acceptedDiff, time.Duration(timeoutMinutes)*time.Minute, slaveHost)
		if err != nil {
			logg.Fatal(err.Error())
		}

		fmt.Printf("Dynomite backend [%s] warmup from [%s] done: %v\n", host, master.Host, result)
	},
}

func init() {
	backendCmd.PersistentFlags().Int16Var(&backendPort, "backend-port", 22122, "dynomite backend port")
	backendCmd.PersistentFlags().StringVar(&backendPassword, "backend-password", "", "dynomite backend password")
	backendCmd.AddCommand(pingCmd)
	backendCmd.AddCommand(roleCmd)

	warmupCmd.PersistentFlags().IntVar(&timeoutMinutes, "timeout-minutes", 5, "Time in minutes until the Warmup times out")
	warmupCmd.PersistentFlags().Int64Var(&acceptedDiff, "accepted-diff", 100, "Accepted difference for replication offset between master and replica")
	warmupCmd.PersistentFlags().Int16Var(&masterBackendPort, "master-backend-port", 22122, "master backend port")
	warmupCmd.PersistentFlags().StringVar(&replicaAnnounceIP, "replica-announce-ip", "", "external IP announced to the master")
	backendCmd.AddCommand(warmupCmd)

	rootCmd.AddCommand(backendCmd)
}
