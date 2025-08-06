package retriever

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-zookeeper/zk"
)

const (
	zkPath    = "/id-generator"
	rangeSize = 100000000000
)

type Config struct {
	ZkAddress string `envconfig:"ZK_ADDRESS" default:"localhost"`
	ZkPort    int    `envconfig:"ZK_PORT" default:"2181"`
}

type Retriever interface {
	GetBase() (int, error)
}

type Zookeeper struct {
	client       *zk.Conn
	endRange     int
	currentRange int
}

var _ Retriever = (*Zookeeper)(nil)

func New(cfg Config) (*Zookeeper, error) {
	conn, _, err := zk.Connect([]string{cfg.ZkAddress + ":" + strconv.Itoa(cfg.ZkPort)}, time.Second*25)
	if err != nil {
		return nil, fmt.Errorf("connect to zookeeper: %w", err)
	}

	exists, _, err := conn.Exists(zkPath)
	if err != nil {
		return nil, fmt.Errorf("check if node exists in zookeeper: %w", err)
	}

	if !exists {
		_, err = conn.Create(zkPath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return nil, fmt.Errorf("create new path: %w", err)
		}
		slog.Debug("Created base path", "path", zkPath)
	}

	return &Zookeeper{client: conn}, nil
}

func (z *Zookeeper) allocateIDRange() (int, int, error) {
	for {
		data, stat, err := z.client.Get(zkPath)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to read zk path: %w", err)
		}

		var currentID int

		if len(data) != 0 {
			currentID, err = strconv.Atoi(string(data))
			if err != nil {
				return 0, 0, fmt.Errorf("invalid data in znode: %w", err)
			}
		}

		newID := currentID + rangeSize

		_, err = z.client.Set(zkPath, []byte(strconv.Itoa(newID)), stat.Version)
		if err == zk.ErrBadVersion {
			continue
		} else if err != nil {
			return 0, 0, fmt.Errorf("failed to update znode: %w", err)
		}

		return newID, newID + rangeSize, nil
	}
}

func (z *Zookeeper) GetBase() (int, error) {
	if z.currentRange >= z.endRange {
		start, end, err := z.allocateIDRange()
		if err != nil {
			return 0, fmt.Errorf("allocating new range: %w", err)
		}
		z.endRange = end
		z.currentRange = start
		return start, nil
	}

	z.currentRange++

	return z.currentRange, nil
}
