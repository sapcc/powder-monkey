package dynomite

import (
	"bufio"
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

// Role returns the current role of Redis master/slave
func (r Redis) Role() (string, error) {
	conn := r.connPool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("INFO", "replication"))
	if err != nil {
		return "", err
	}

	var role, master, connectedSlaves string
	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "role:") {
			role = line[5:]
		} else if strings.HasPrefix(line, "master_host:") {
			master = line[12:]
		} else if strings.HasPrefix(line, "connected_slaves:") {
			connectedSlaves = line[17:]
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

// ReplicationOffset determines the ReplicationOffset difference between master and slave
func (r Redis) ReplicationOffset(master Redis) (int64, error) {
	conn := master.connPool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("INFO", "replication"))
	if err != nil {
		return 0, err
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

	scanner := bufio.NewScanner(strings.NewReader(result))
	var masterOffset, slaveOffset int64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "slave") {
			// slave0:ip=127.0.0.1,port=22122,state=online,offset=1288,lag=1

			// Get the values
			values := strings.SplitN(line, ":", 2)[1]
			// Get key values
			kv := strings.Split(values, ",")

			ip := strings.SplitN(kv[0], "=", 2)[1]
			if ip == r.Host {
				offset := strings.SplitN(kv[3], "=", 2)[1]
				slaveOffset, err = strconv.ParseInt(offset, 10, 64)
				if err != nil {
					return 0, err
				}
				logg.Info("Slave Offset %d", slaveOffset)
			}
		} else if strings.HasPrefix(scanner.Text(), "master_repl_offset") {
			master := strings.Split(scanner.Text(), ":")
			masterOffset, err = strconv.ParseInt(master[1], 10, 64)
			if err != nil {
				return 0, err
			}
			logg.Info("Master Offset %d", masterOffset)
		}
		if slaveOffset != 0 && masterOffset != 0 {
			break
		}
	}
	return (masterOffset - slaveOffset), nil
}
