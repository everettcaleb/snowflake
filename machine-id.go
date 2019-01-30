package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v3"
)

const machineIDLifeSeconds = 120

func initMachineID(config *snowflakeEnvConfig) (uint64, error) {
	// Connect to Redis
	conn, err := radix.Dial("tcp", config.RedisURI)
	if err != nil {
		return 0, err
	}

	// Randomly pick a machine ID and then try to set Redis
	// to see if it's available (and reserve it)
	var machineID uint64
	var key string
	var setNxResult int

	for setNxResult == 0 {
		// Randomly pick one and build the key
		machineID = rand.Uint64() % 1024
		key = config.RedisMachineIDPrefix + strconv.FormatUint(machineID, 10)

		// Attempt to set it (SETNX fails if it's already set)
		// This will set "setNxResult" to 0 which causes the loop to try again
		err = conn.Do(radix.Cmd(&setNxResult, "SETNX", key, "1"))
		if err != nil {
			return 0, err
		}
	}

	// Kick off a goroutine loop that renews the machine ID
	go func() {
		// Make sure we close the connection if the goroutine dies or something weird
		defer conn.Close()

		// Renew the expiration time on every iteration
		for {
			// Let's mark it to expire in case this service dies otherwise
			// there's a risk to reserve every machine ID by rapidly recycling services
			conn.Do(radix.Cmd(nil, "EXPIRE", key, strconv.Itoa(machineIDLifeSeconds)))

			// Sleep for half the time so we don't flood Redis
			time.Sleep(machineIDLifeSeconds / 2 * time.Second)
		}
	}()

	return machineID, nil
}
