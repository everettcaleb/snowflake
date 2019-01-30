package main

import (
	"fmt"
	"math/rand"
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
	fmt.Println(1)

	// Randomly pick a machine ID and then try to set Redis
	// to see if it's available (and reserve it)
	var machineID uint64
	var key string
	for {
		fmt.Println(2)
		// Randomly pick one and build the key
		machineID = rand.Uint64() % 1024
		key = config.RedisMachineIDPrefix + string(machineID)

		// Attempt to set it (SETNX fails if it's already set)
		var setNxResult int
		err = conn.Do(radix.Cmd(&setNxResult, "SETNX", key, "1"))
		if err != nil {
			return 0, err
		}

		// A zero value means that ID is already taken
		if setNxResult == 0 {
			continue
		}

		// Let's mark it to expire in case this service dies otherwise
		// there's a risk to reserve every machine ID by rapidly
		// recycling services
		err = conn.Do(radix.Cmd(nil, "EXPIRE", key, string(machineIDLifeSeconds)))
		if err != nil {
			return 0, err
		}

		// Break out of the loop, since we have an ID now
		break
	}

	// Kick off a goroutine loop that renews the machine ID
	go func() {
		for {
			renewMachineID(conn, key)
		}
	}()

	return machineID, nil
}

func renewMachineID(conn radix.Conn, key string) {
	// Sleep for half the time so we don't flood Redis
	time.Sleep(machineIDLifeSeconds / 2)

	// Renew the reservation
	conn.Do(radix.Cmd(nil, "EXPIRE", key, string(machineIDLifeSeconds)))
}
