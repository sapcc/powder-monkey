package dynomite

import (
	"fmt"
	"time"

	"github.com/sapcc/go-bits/logg"
)

// Warmup starts the warmup process
// 1. Set Dynomite in Standby mode
// 2. Set Dynomite Backend as replica of master
// 3. Wait for accecptable replication offset difference between master and replica or timeout
// 4. Set Dynomite State to write_only
// 5. Stop replication
// 6. Set Dynomite State to resuming
// 7. Set Dynomite State to normal
func (dyno Dynomite) Warmup(master Redis, accecptedDiff int64, timeout time.Duration) (bool, error) {
	// Check backend alive
	ping, err := dyno.Backend.Ping()
	if err != nil {
		return false, fmt.Errorf("Warmup failed: %s", err.Error())
	}
	if !ping {
		return false, fmt.Errorf("Warmup failed: Backend did not ping")
	}

	// Set State standby
	_, err = dyno.SetState(Standby)
	if err != nil {
		return false, fmt.Errorf("Warmup failed: %s", err.Error())
	}
	logg.Info("Setting state %s", Standby)

	// Backend to replicate from master
	_, err = dyno.Backend.Replicate(master)
	if err != nil {
		return false, fmt.Errorf("Warmup failed: %s", err.Error())
	}
	logg.Info("Replication setup from %s", master.Host)

	timer := time.After(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			diff, err := dyno.Backend.ReplicationOffset(master)
			if err != nil {
				logg.Error("Warmup failed: %s", err.Error())
			}
			logg.Info("Current replication offset diff: %d", diff)

			// Accecpted Diff reached
			if diff < accecptedDiff {
				// Set State writes_only
				_, err = dyno.SetState(WritesOnly)
				if err != nil {
					return false, fmt.Errorf("Warmup failed: %s", err.Error())
				}
				logg.Info("Setting state %s", WritesOnly)

				// Stop Sync
				err = dyno.Backend.StopReplication()
				if err != nil {
					return false, fmt.Errorf("Warmup failed: %s", err.Error())
				}
				logg.Info("Replication stopped")

				// Set State resuming
				_, err = dyno.SetState(Resuming)
				if err != nil {
					return false, fmt.Errorf("Warmup failed: %s", err.Error())
				}
				logg.Info("Setting state %s", Resuming)

				// sleep 15s for the flushing to catch up
				time.Sleep(15 * time.Second)

				// Set State Normal
				_, err = dyno.SetState(Normal)
				if err != nil {
					return false, fmt.Errorf("Warmup failed: %s", err.Error())
				}
				logg.Info("Setting state %s", Normal)
				return true, nil
			}
		case <-timer:
			return false, fmt.Errorf("Warmup timedout")
		}
	}
}
