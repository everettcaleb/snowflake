package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	bits10 uint64 = (1 << 10) - 1
	bits12 uint64 = (1 << 12) - 1
	bits41 uint64 = (1 << 41) - 1
	bit12  uint64 = (1 << 12)
)

var now = func() time.Time {
	return time.Now()
}

func timestamp() uint64 {
	return uint64(now().Unix())
}

func msTimestamp() uint64 {
	return uint64(now().UnixNano() / 1e6)
}

type snowflakeID uint64

type snowflakeIDParts struct {
	id        snowflakeID
	counter   uint64
	machineID uint64
	timestamp uint64
}

func newSnowflakeGenerator(config *snowflakeEnvConfig, machineID uint64) *snowflakeGenerator {
	epoch := config.Epoch
	tsfunc := timestamp

	if config.UseMilliseconds {
		epoch = epoch * 1000
		tsfunc = msTimestamp
	}

	return &snowflakeGenerator{
		counter:       0,
		epoch:         epoch,
		lastTimestamp: 0,
		machineID:     machineID,
		timestampFunc: tsfunc,
	}
}

type snowflakeGenerator struct {
	mutex         sync.Mutex
	counter       uint64
	epoch         uint64
	lastTimestamp uint64
	machineID     uint64
	timestampFunc func() uint64
}

func splitSnowflakeID(id snowflakeID) *snowflakeIDParts {
	rid := uint64(id)
	sid := &snowflakeIDParts{id: id}

	sid.counter = rid & bits12
	rid = rid &^ bits12 >> 12

	sid.machineID = rid & bits10
	rid = rid &^ bits10 >> 10

	sid.timestamp = rid

	return sid
}

// NextID generates the next ID in the sequence
func (s *snowflakeGenerator) NextID() snowflakeID {
	// Time since 1970-01-01 in milliseconds
	timestamp := s.timestampFunc()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if timestamp < s.lastTimestamp {
		timestamp = s.lastTimestamp
	}

	if timestamp == s.lastTimestamp {
		s.counter = (s.counter + 1) % bit12
		if s.counter == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = s.timestampFunc()
			}
		}
	} else {
		s.counter = 0
	}

	s.lastTimestamp = timestamp

	td := timestamp - s.epoch
	if td >= bits41 {
		panic(fmt.Errorf("Timestamp epoch delta exceeded 41 bits"))
	}
	return snowflakeID(((td & bits41) << 22) + ((s.machineID & bits10) << 12) + s.counter)
}
