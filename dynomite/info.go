package dynomite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// InfoReport contains the information as reported by the /info endpoint of
// a Dynomite instance.
type InfoReport struct {
	Uptime uint64 `json:"uptime"`
	Rack   string `json:"rack"`
	DC     string `json:"dc"`
	Pool   struct {
		ClientConnections     uint64 `json:"client_connections"`
		ClientReadRequests    uint64 `json:"client_read_requests"`
		ClientWriteRequests   uint64 `json:"client_write_requests"`
		ClientDroppedRequests uint64 `json:"client_dropped_requests"`
	} `json:"dyn_o_mite"`
}

// Info returns the InfoReport for a Dynomite instance.
func (dyno Dynomite) Info() (*InfoReport, error) {
	url := fmt.Sprintf("http://%s:%d/info", dyno.Host, dyno.Port)
	resp, err := netClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var report InfoReport
	err = json.Unmarshal(body, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}
