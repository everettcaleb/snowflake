package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getMachineID() (uint64, error) {
	str := os.Getenv("MACHINE_ID")
	if len(str) < 1 {
		return 0, fmt.Errorf("MACHINE_ID environment variable must be set to a number 0-1023 or HOST to extract it (for Kubernetes StatefulSet, etc)")
	}

	var machineID uint64
	if str == "HOST" {
		host, err := os.Hostname()
		if err != nil {
			return 0, fmt.Errorf("MACHINE_ID environment variable was set to HOST but could not get hostname")
		}

		// Extract machine ID from the last piece, ex: "snowflake-0" where 0 is the machine ID
		parts := strings.Split(host, "-")
		machineID, err = strconv.ParseUint(parts[len(parts)-1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("Hostname must end in a number, ex: snowflake-0 if MACHINE_ID is HOST")
		}
	} else {
		var err error
		machineID, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("MACHINE_ID environment variable must be set to a number 0-1023 or HOST to extract it (for Kubernetes StatefulSet, etc)")
		}
	}

	if machineID >= 1024 {
		return 0, fmt.Errorf("MACHINE_ID was not in the range 0-1023")
	}
	return machineID << 12, nil
}

func getEpoch() (uint64, error) {
	str := os.Getenv("SNOWFLAKE_EPOCH")
	if len(str) < 1 {
		return 0, fmt.Errorf("SNOWFLAKE_EPOCH needs to be set to a Unix timestamp")
	}

	epoch, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("SNOWFLAKE_EPOCH is not a number")
	}
	return epoch * 1000, nil
}

func getListenerPort() (string, error) {
	str := os.Getenv("PORT")
	if len(str) < 1 {
		return ":8080", nil
	}

	port, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return "", fmt.Errorf("PORT environment variable is not a number")
	}

	return fmt.Sprintf(":%d", port), nil
}

func getBasePath() string {
	str := os.Getenv("APP_BASE_PATH")
	if len(str) < 1 {
		return "/v1/snowflake"
	}
	return str
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

type snowflakeResponse struct {
	ID       uint64 `json:"id"`
	IDString string `json:"idStr"`
}

func makeSnowflakeResponse(id uint64) *snowflakeResponse {
	return &snowflakeResponse{ID: id, IDString: fmt.Sprintf("%d", id)}
}
