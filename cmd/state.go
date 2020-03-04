package cmd

import (
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
		dyno := dynomite.NewDynomite(dynomiteHost, dynomitePort)
		state, err := dyno.GetState()
		if err != nil {
			logg.Fatal(err.Error())
		}
		logg.Info("Dynomite [%s] State: %v", dyno.Host, state)
	},
}

var stateSetCmd = &cobra.Command{
	Use:       "set [normal|standby|writes_only|resuming]",
	Short:     "Set state of dynomite",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"normal", "standby", "writes_only", "resuming"},
	Run: func(cmd *cobra.Command, args []string) {
		state, err := dynomite.StrToState(args[0])
		if err != nil {
			logg.Fatal(err.Error())
		}
		dyno := dynomite.NewDynomite(dynomiteHost, dynomitePort)
		result, err := dyno.SetState(state)
		if err != nil {
			logg.Fatal(err.Error())
		}
		logg.Info("Dynomite [%s] State set to %v: %s", dyno.Host, state, result)
	},
}

func init() {
	stateCmd.AddCommand(stateGetCmd)
	stateCmd.AddCommand(stateSetCmd)
	rootCmd.AddCommand(stateCmd)
}
