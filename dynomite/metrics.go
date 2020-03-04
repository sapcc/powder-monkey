// Copyright 2020 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dynomite

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/go-bits/logg"
)

// >>>>>>
// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");

type typedDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

func (d *typedDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}

// <<<<<<

func (d *typedDesc) describe(ch chan<- *prometheus.Desc) {
	ch <- d.desc
}

// Collector implements the prometheus.Collector interface.
type Collector struct {
	dyno *Dynomite

	state                 typedDesc
	dbSize                typedDesc
	uptime                typedDesc
	clientConnections     typedDesc
	clientReadRequests    typedDesc
	clientWriteRequests   typedDesc
	clientDroppedRequests typedDesc
}

// NewCollector creates a new Collector.
func NewCollector(dyno *Dynomite) *Collector {
	return &Collector{
		dyno: dyno,
		state: typedDesc{
			desc: prometheus.NewDesc(
				"dynomite_state",
				"State as reported by Dynomite.",
				[]string{"state", "rack", "dc", "ip_address"}, nil),
			valueType: prometheus.GaugeValue,
		},
		dbSize: typedDesc{
			desc: prometheus.NewDesc(
				"dynomite_db_size",
				"Key database size as reported by the Redis backend.",
				[]string{"rack", "dc", "ip_address"}, nil),
			valueType: prometheus.GaugeValue,
		},
		uptime: typedDesc{
			desc: prometheus.NewDesc(
				"dynomite_uptime",
				"Uptime as reported by Dynomite info.",
				[]string{"rack", "dc", "ip_address"}, nil),
			valueType: prometheus.GaugeValue,
		},
		clientConnections: typedDesc{
			desc: prometheus.NewDesc(
				"dynomite_client_connections",
				"Client connections as reported by Dynomite info.",
				[]string{"rack", "dc", "ip_address"}, nil),
			valueType: prometheus.GaugeValue,
		},
		clientReadRequests: typedDesc{
			desc: prometheus.NewDesc(
				"dynomite_client_read_requests",
				"Client read requests as reported by Dynomite info.",
				[]string{"rack", "dc", "ip_address"}, nil),
			valueType: prometheus.GaugeValue,
		},
		clientWriteRequests: typedDesc{
			desc: prometheus.NewDesc(
				"dynomite_client_write_requests",
				"Client write requests as reported by Dynomite info.",
				[]string{"rack", "dc", "ip_address"}, nil),
			valueType: prometheus.GaugeValue,
		},
		clientDroppedRequests: typedDesc{
			desc: prometheus.NewDesc(
				"dynomite_client_dropped_requests",
				"Client dropped requests as reported by Dynomite info.",
				[]string{"rack", "dc", "ip_address"}, nil),
			valueType: prometheus.GaugeValue,
		},
	}
}

// Describe implements the prometheus.Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.state.describe(ch)
	c.dbSize.describe(ch)
	c.uptime.describe(ch)
	c.clientConnections.describe(ch)
	c.clientReadRequests.describe(ch)
	c.clientWriteRequests.describe(ch)
	c.clientDroppedRequests.describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ip := os.Getenv("DYNO_INSTANCE")
	if ip == "" {
		logg.Error("could not get ip address from env variable: DYNO_INSTANCE")
	}

	var rack, dc string
	ir, err := c.dyno.Info()
	if err != nil {
		logg.Error(err.Error())
	} else {
		rack = ir.Rack
		dc = ir.DC

		ch <- c.uptime.mustNewConstMetric(float64(ir.Uptime), rack, dc, ip)
		ch <- c.clientConnections.mustNewConstMetric(float64(ir.Pool.ClientConnections), rack, dc, ip)
		ch <- c.clientReadRequests.mustNewConstMetric(float64(ir.Pool.ClientReadRequests), rack, dc, ip)
		ch <- c.clientWriteRequests.mustNewConstMetric(float64(ir.Pool.ClientWriteRequests), rack, dc, ip)
		ch <- c.clientDroppedRequests.mustNewConstMetric(float64(ir.Pool.ClientDroppedRequests), rack, dc, ip)
	}

	stateVal := 1 // until proven otherwise
	state, err := c.dyno.GetState()
	if err != nil {
		stateVal = 0
		logg.Error(err.Error())
	}

	if state != Normal {
		stateVal = 0
	}
	ch <- c.state.mustNewConstMetric(float64(stateVal), string(state), rack, dc, ip)

	size, err := c.dyno.Backend.DBSize()
	if err != nil {
		logg.Error(err.Error())
	} else {
		ch <- c.dbSize.mustNewConstMetric(float64(size), rack, dc, ip)
	}
}
