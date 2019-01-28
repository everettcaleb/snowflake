package main

import (
	"fmt"
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

func msTimestamp() uint64 {
	return uint64(now().UnixNano() / 1e6)
}

type snowflakeGenerator struct {
	counter       uint64
	lastTimestamp uint64
	machineID     uint64
	epoch         uint64
}

type snowflakeID struct {
	id        uint64
	timestamp uint64
	machineID uint64
	counter   uint64
}

// makeSnowflakeGenerator creates a snowflake ID generator instance that tracks generator state
func makeSnowflakeGenerator(machineID uint64, epoch uint64) *snowflakeGenerator {
	return &snowflakeGenerator{
		machineID: machineID,
		epoch:     epoch,
	}
}

func splitSnowflakeID(id uint64) *snowflakeID {
	sid := &snowflakeID{id: id}

	sid.counter = id & bits12
	id = id &^ bits12 >> 12

	sid.machineID = id & bits10
	id = id &^ bits10 >> 10

	sid.timestamp = id

	return sid
}

// NextID generates the next ID in the sequence
func (s *snowflakeGenerator) NextID() uint64 {
	// Time since 1970-01-01 in milliseconds
	timestamp := msTimestamp()
	if timestamp < s.lastTimestamp {
		timestamp = s.lastTimestamp
	}

	if timestamp == s.lastTimestamp {
		s.counter = (s.counter + 1) % bit12
		if s.counter == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = msTimestamp()
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
	return ((td & bits41) << 22) + ((s.machineID & bits10) << 12) + s.counter
}
