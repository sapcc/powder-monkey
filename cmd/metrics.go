package cmd

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/go-bits/httpee"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/powder-monkey/dynomite"

	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Exposes metrics at :9682",
	Run: func(cmd *cobra.Command, args []string) {
		dyno := dynomite.NewDynomiteRedis(dynomiteHost, dynomitePort, backendPort, backendPassword)

		prometheus.MustRegister(dynomite.NewCollector(dyno))

		// this port has been allocated for Powder Monkey
		// See: https://github.com/prometheus/prometheus/wiki/Default-port-allocations
		listenAddr := ":9682"
		http.HandleFunc("/", metriclandingPageHandler)
		http.Handle("/metrics", promhttp.Handler())
		logg.Info("listening on " + listenAddr)
		err := httpee.ListenAndServeContext(httpee.ContextWithSIGINT(context.Background()), listenAddr, nil)
		if err != nil {
			logg.Fatal(err.Error())
		}
	},
}

func init() {
	metricsCmd.PersistentFlags().Int16Var(&backendPort, "backend-port", 22122, "dynomite backend port")
	metricsCmd.PersistentFlags().StringVar(&backendPassword, "backend-password", "", "dynomite backend password")

	rootCmd.AddCommand(metricsCmd)
}

func metriclandingPageHandler(w http.ResponseWriter, r *http.Request) {
	pageBytes := []byte(`<html>
<head><title>Powder Monkey</title></head>
<body>
<h1>Powder Monkey</h1>
<p><a href="/metrics">Metrics</a></p>
<p><a href="https://github.com/sapcc/powder-monkey">Source Code</a></p>
</body>
</html>`)

	_, err := w.Write(pageBytes)
	if err != nil {
		logg.Error(err.Error())
	}
}
