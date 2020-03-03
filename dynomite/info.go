package dynomite

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// InfoReport contains the information as reported by the /info endpoint of
// a Dynomite instance.
type InfoReport struct {
	Uptime       uint64 `json:"uptime"`
	Rack         string `json:"rack"`
	DC           string `json:"dc"`
	DynomiteRack struct {
		ClientConnections     uint64 `json:"client_connections"`
		ClientReadRequests    uint64 `json:"client_read_requests"`
		ClientWriteRequests   uint64 `json:"client_write_requests"`
		ClientDroppedRequests uint64 `json:"client_dropped_requests"`
	} `json:"dynomite-rack"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// We need this because Dynomite reports the information about a rack with the
// key: "dynomite-$rack".
func (r *InfoReport) UnmarshalJSON(b []byte) error {
	// decode the original data into a temporary map
	var a interface{}
	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}
	m, ok := a.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected a value of type: map[string]interface{}, got: %T", a)
	}

	v, ok := m["dc"]
	if !ok {
		return errors.New("no key found with name: 'dc'")
	}

	dc, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected a value of type: string, got: %T", v)
	}

	key := "dynomite-" + dc

	// replace the map key with a generic one
	v, ok = m[key]
	if !ok {
		return fmt.Errorf("no key found with name: %q", key)
	}

	m["dynomite-rack"] = v
	delete(m, key)

	// encode back to JSON
	b, err = json.Marshal(m)
	if err != nil {
		return err
	}

	// Unmarshal into a InfoReport
	type t InfoReport // temporary type wrapper
	var d t
	err = json.Unmarshal(b, &d)
	if err != nil {
		return err
	}

	*r = InfoReport(d)

	return nil
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
