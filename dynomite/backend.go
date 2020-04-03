package dynomite

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sapcc/go-bits/logg"

	"github.com/gomodule/redigo/redis"
)

// Redis represent the dynomite backend as Redis instance
type Redis struct {
	Host     string
	Port     int16
	connPool *redis.Pool
}

// NewRedis creates a new Redis struct with initialized ConnectionPool
func NewRedis(host string, port int16, password string) *Redis {
	dialops := []redis.DialOption{
		redis.DialConnectTimeout(3 * time.Second),
		redis.DialReadTimeout(3 * time.Second),
		redis.DialWriteTimeout(3 * time.Second),
	}

	if password != "" {
		dialops = append(dialops, redis.DialPassword(password))
	}

	connection := fmt.Sprintf("%s:%s", host, strconv.FormatInt(int64(port), 10))
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", connection, dialops...)
		},
	}

	return &Redis{
		Host:     host,
		Port:     port,
		connPool: pool,
	}
}

// ReplicationOffsets represents the current state of replication of master and slave
type ReplicationOffsets struct {
	Master int64
	Slave  int64
}

// Ping checks liveness od Redis
func (r Redis) Ping() (bool, error) {
	conn := r.connPool.Get()
	defer conn.Close()

	pong, err := redis.String(conn.Do("PING"))
	if err != nil {
		return false, err
	}

	return (pong == "PONG"), nil
}

// WaitFor waits for a succesful Ping to the backend during the specified timeout
func (r Redis) WaitFor(timeout time.Duration) error {
	timer := time.After(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ping, _ := r.Ping(); ping {
				return nil
			}
			logg.Info("Backend %s did not yet succesfully ping", r.Host)
		case <-timer:
			return fmt.Errorf("Pinging backend %s timed out", r.Host)
		}
	}
}

// Role returns the current role of Redis master/slave
func (r Redis) Role() (string, error) {
	conn := r.connPool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("INFO", "replication"))
	if err != nil {
		return "", err
	}

	var role, master, connectedSlaves string
	for _, line := range strings.Split(result, "\r\n") {
		if strings.HasPrefix(line, "role:") {
			role = strings.TrimPrefix(line, "role:")
		} else if strings.HasPrefix(line, "master_host:") {
			master = strings.TrimPrefix(line, "master_host:")
		} else if strings.HasPrefix(line, "connected_slaves:") {
			connectedSlaves = strings.TrimPrefix(line, "connected_slaves:")
		}
	}

	if role == "slave" {
		return fmt.Sprintf("%s (Master is %s)", role, master), nil
	}
	return fmt.Sprintf("%s with %s connected slaves", role, connectedSlaves), nil
}

// Replicate activates replication from the given master
func (r Redis) Replicate(master Redis) (bool, error) {
	masterPing, err := master.Ping()
	if err != nil {
		return false, err
	}
	if !masterPing {
		return false, fmt.Errorf("Master System %s is not ready", master.Host)
	}

	conn := r.connPool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("REPLICAOF", master.Host, master.Port))
	if err != nil {
		return false, err
	}

	if strings.HasPrefix(result, "OK") {
		return true, nil
	}

	return false, fmt.Errorf("Replication could not be setup: %s", result)
}

// StopReplication activates replication from the given master
func (r Redis) StopReplication() error {
	conn := r.connPool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("REPLICAOF", "NO", "ONE"))
	if err != nil {
		return err
	}

	if strings.HasPrefix(result, "OK") {
		return nil
	}

	return fmt.Errorf("Replication could not be stopped: %s", result)
}

// ReplicationOffsets determines the ReplicationOffset difference between master and slave
func (r Redis) ReplicationOffsets(slaveHost string) (*ReplicationOffsets, error) {
	conn := r.connPool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("INFO", "replication"))
	if err != nil {
		return nil, err
	}
	logg.Debug("%v", result)

	/*
		# Replication
		role:master
		connected_slaves:1
		slave0:ip=127.0.0.1,port=22122,state=online,offset=1288,lag=1
		master_replid:af226365937302a504735a6a9a881758421680af
		master_replid2:0000000000000000000000000000000000000000
		master_repl_offset:1288
		second_repl_offset:-1
		repl_backlog_active:1
		repl_backlog_size:1048576
		repl_backlog_first_byte_offset:29
		repl_backlog_histlen:1260
	*/

	replOffsets := &ReplicationOffsets{}

	for _, line := range strings.Split(result, "\r\n") {
		if strings.HasPrefix(line, "slave") {
			// slave0:ip=127.0.0.1,port=22122,state=online,offset=1288,lag=1

			// Get the values
			values := strings.SplitN(line, ":", 2)[1]
			// Get key values
			kv := strings.Split(values, ",")

			ip := strings.TrimPrefix(kv[0], "ip=")
			if ip == slaveHost {
				offset := strings.TrimPrefix(kv[3], "offset=")
				replOffsets.Slave, err = strconv.ParseInt(offset, 10, 64)
				if err != nil {
					return nil, err
				}
				logg.Info("Slave Offset %d", replOffsets.Slave)
			}
		} else if strings.HasPrefix(line, "master_repl_offset") {
			master := strings.TrimPrefix(line, "master_repl_offset:")
			replOffsets.Master, err = strconv.ParseInt(master, 10, 64)
			if err != nil {
				return nil, err
			}
			logg.Info("Master Offset %d", replOffsets.Master)
		}
		if replOffsets.Master != 0 && replOffsets.Slave != 0 {
			// Parsed everything
			break
		}
	}

	return replOffsets, nil
}

// Warmup does a simple replication from a master backend without dealing with dynomite states
func (r Redis) Warmup(master Redis, accecptedDiff int64, timeout time.Duration, slaveHost string) (bool, error) {
	err := r.WaitFor(2 * time.Minute)
	if err != nil {
		return false, fmt.Errorf("Warmup failed: %s", err.Error())
	}

	// Backend to replicate from master
	_, err = r.Replicate(master)
	if err != nil {
		return false, fmt.Errorf("Warmup failed - Initiate Replication: %s", err.Error())
	}
	logg.Info("Replication setup from %s", master.Host)

	timer := time.After(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			replOffsets, err := master.ReplicationOffsets(slaveHost)
			if err != nil {
				logg.Error("Warmup failed - Get Replication Offset: %s", err.Error())
			}

			if replOffsets.Master > 0 && replOffsets.Slave > 0 {
				diff := replOffsets.Master - replOffsets.Slave
				logg.Info("Current replication offset diff: %d", diff)

				// Accecpted Diff reached
				if diff < accecptedDiff {
					// Stop Sync
					err = r.StopReplication()
					if err != nil {
						return false, fmt.Errorf("Warmup failed - Stopping Replication: %s", err.Error())
					}
					logg.Info("Replication stopped")
					return true, nil
				}
			} else {
				logg.Info("Replication not yet progessed. Waiting")
			}
		case <-timer:
			return false, fmt.Errorf("Warmup timed out")
		}
	}
}

// DBSize retrieves the number of Keys in the Backend
func (r Redis) DBSize() (int64, error) {
	conn := r.connPool.Get()
	defer conn.Close()

	var size int64

	size, err := redis.Int64(conn.Do("DBSIZE"))
	if err != nil {
		return size, err
	}
	logg.Debug("DBSIZE: %d", size)

	return size, nil
}

// LastSave return the Time of the last save
func (r Redis) LastSave() (time.Time, error) {
	conn := r.connPool.Get()
	defer conn.Close()

	var lastSave time.Time

	result, err := redis.Int64(conn.Do("LASTSAVE"))
	if err != nil {
		return lastSave, err
	}

	lastSave = time.Unix(result, 0)
	logg.Debug("LASTSAVE: %v", lastSave)

	return lastSave, nil
}

// BGSave triggers Background Save and watches the status
func (r Redis) BGSave(timeout time.Duration) (time.Time, error) {
	currentTime := time.Now()

	conn := r.connPool.Get()
	defer conn.Close()

	var currentSave time.Time

	// Getting Last Save Timestamp
	lastSave, err := r.LastSave()
	if err != nil {
		return currentSave, err
	}

	// Trigger Background SAVE
	result, err := redis.String(conn.Do("BGSAVE"))
	if err != nil {
		return currentSave, err
	}
	logg.Info("BGSAVE %s", result)

	timer := time.After(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Getting Current Save Timestamp
			currentSave, err := r.LastSave()
			if err != nil {
				logg.Error(err.Error())
				break
			}

			// LastSave was updated --> BGSave finished
			if currentSave.After(lastSave) {
				logg.Info("BGSAVE finished in %s", time.Now().Sub(currentTime))
				return currentSave, nil
			}
		case <-timer:
			return currentSave, fmt.Errorf("BGSAVE timed out")
		}
	}
}
