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
	Standby    State = "standby"
	WritesOnly State = "writes_only"
	Resuming   State = "resuming"
)

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

// GetState retrieves dynomite internal state
func (dyno Dynomite) GetState() (State, error) {
	url := fmt.Sprintf("http://%s:%d/state/get_state", dyno.Host, dyno.Port)
	resp, err := netClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	st := string(body)
	st = strings.TrimSuffix(st, "\n")
	if strings.HasPrefix(st, "State: ") {
		st = strings.ReplaceAll(st, "State: ", "")
	}

	return StrToState(st)
}

// SetState sets Dynomite's internal state and returns the response status,
// if successful.
func (dyno Dynomite) SetState(state State) (string, error) {
	url := fmt.Sprintf("http://%s:%d/state/%s", dyno.Host, dyno.Port, state)
	// Thats seems to be strange, but setting a state is actually a GET
	// https://github.com/Netflix/dynomite/wiki/REST#state
	resp, err := netClient.Get(url)
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

// StrToState validates and converts a string to a State.
func StrToState(str string) (State, error) {
	given := State(strings.ToLower(str))
	switch given {
	case Normal, Standby, WritesOnly, Resuming:
		return given, nil
	default:
		return "", fmt.Errorf("Invalid state: %s", str)
	}
}
