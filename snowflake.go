package main

import "time"

const bits41 uint64 = 2199023255551
const bits10 uint64 = 1023
const bits12Mod uint64 = 4096

func msTimestamp() uint64 {
	return uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
}

type snowflakeGenerator struct {
	counter       uint64
	lastTimestamp uint64
	machineID     uint64
	epoch         uint64
}

// MakeSnowflakeGenerator creates a snowflake ID generator instance that tracks generator state
func makeSnowflakeGenerator(machineID uint64, epoch uint64) *snowflakeGenerator {
	return &snowflakeGenerator{
		counter:       0,
		lastTimestamp: 0,
		machineID:     machineID,
		epoch:         epoch,
	}
}

// NextID generates the next ID in the sequence
func (s *snowflakeGenerator) NextID() uint64 {
	// Time since 1970-01-01 in milliseconds
	timestamp := msTimestamp()
	if timestamp < s.lastTimestamp {
		timestamp = s.lastTimestamp
	}

	if timestamp == s.lastTimestamp {
		s.counter = (s.counter + 1) % bits12Mod
		if s.counter == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = msTimestamp()
			}
		}
	} else {
		s.counter = 0
	}

	s.lastTimestamp = timestamp
	return (((timestamp - s.epoch) & bits41) << 22) + s.machineID + s.counter
}
