package dynomite

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// State representing the working mode of dynomite
type State string

// Possible states of dynomite
const (
	Normal     State = "normal"
	Standby          = "standby"
	WritesOnly       = "writes_only"
	Resuming         = "resuming"
)

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

// GetState retrieves dynomite internal state
func (dyno Dynomite) GetState() (State, error) {
	var state State

	url := fmt.Sprintf("http://%s:%d/state/get_state", dyno.Host, dyno.Port)
	resp, err := netClient.Get(url)
	if err != nil {
		return state, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	st := string(body)
	st = strings.TrimSuffix(st, "\n")

	if strings.HasPrefix(st, "State: ") {
		st = strings.ToLower((strings.ReplaceAll(st, "State: ", "")))
	}
	state = State(st)
	err = validState(state)
	return state, err
}

// SetState sets dynomites internal state
func (dyno Dynomite) SetState(state State) (string, error) {
	var result string

	err := validState(state)
	if err != nil {
		return result, err
	}
	url := fmt.Sprintf("http://%s:%d/state/%s", dyno.Host, dyno.Port, state)
	// Thats seems to be strange, but setting a state is actually a GET
	// https://github.com/Netflix/dynomite/wiki/REST#state
	resp, err := netClient.Get(url)
	if err != nil {
		return result, err
	}

	return resp.Status, nil
}

func validState(state State) error {
	switch state {
	case Normal, Standby, WritesOnly, Resuming:
		return nil
	}
	return fmt.Errorf("Invalid State '%s'", state)
}
