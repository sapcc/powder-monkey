package cmd

import (
	"fmt"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/powder-monkey/dynomite"
	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Get and Set state of dynomite",
}

var stateGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get state of dynomite",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := dynomite.GetState(dynomiteHost, dynomitePort)
		if err != nil {
			logg.Fatal(err.Error())
		}
		fmt.Printf("Dynomite [%s] State: %v\n", dynomiteHost, state)
	},
}

var stateSetCmd = &cobra.Command{
	Use:       "set [normal|standby|writes_only|resuming]",
	Short:     "Set state of dynomite",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"normal", "standby", "writes_only", "resuming"},
	Run: func(cmd *cobra.Command, args []string) {
		var state dynomite.State
		state = dynomite.State(args[0])
		result, err := dynomite.SetState(dynomiteHost, dynomitePort, state)
		if err != nil {
			logg.Fatal(err.Error())
		}
		fmt.Printf("Dynomite [%s] State set to %v: %s\n", dynomiteHost, state, result)
	},
}

func init() {
	stateCmd.AddCommand(stateGetCmd)
	stateCmd.AddCommand(stateSetCmd)
	rootCmd.AddCommand(stateCmd)
}
