package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const bits41 uint64 = 2199023255551
const bits10 uint64 = 1023
const bits12Mod uint64 = 4096

var counter uint64
var lastTimestamp uint64

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

func msTimestamp() uint64 {
	return uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
}

func generateID(machineID uint64, epoch uint64) uint64 {
	// Time since 1970-01-01 in milliseconds
	timestamp := msTimestamp()
	if timestamp < lastTimestamp {
		timestamp = lastTimestamp
	}

	if timestamp == lastTimestamp {
		counter = (counter + 1) % bits12Mod
		if counter == 0 {
			for timestamp <= lastTimestamp {
				timestamp = msTimestamp()
			}
		}
	} else {
		counter = 0
	}

	lastTimestamp = timestamp
	return (((timestamp - epoch) & bits41) << 22) + machineID + counter
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

type snowflakeResponse struct {
	ID       uint64 `json:"id"`
	IDString string `json:"idStr"`
}

func main() {
	router := mux.NewRouter().PathPrefix(getBasePath()).Subrouter()

	machineID, err := getMachineID()
	if err != nil {
		panic(err)
	}

	epoch, err := getEpoch()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Epoch(ms): %d\n", epoch)

	port, err := getListenerPort()
	if err != nil {
		panic(err)
	}

	// Health Check
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(healthCheckResponse{Status: "OK"})
	}).Methods("GET")

	// Specs
	router.HandleFunc("/spec.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "specs/spec.yaml")
	}).Methods("GET")

	router.HandleFunc("/spec.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "specs/spec.json")
	}).Methods("GET")

	// Single Snowflake ID
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		id := generateID(machineID, epoch)
		json.NewEncoder(w).Encode(snowflakeResponse{ID: id, IDString: fmt.Sprintf("%d", id)})
	}).Methods("GET")

	// Multiple Snowflake IDs
	router.HandleFunc("/{count}", func(w http.ResponseWriter, r *http.Request) {
		str := mux.Vars(r)["count"]
		count, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			http.Error(w, "Count must be a number", http.StatusBadRequest)
			return
		}

		var i int64
		ids := make([]snowflakeResponse, count)
		for i = 0; i < count; i++ {
			id := generateID(machineID, epoch)
			ids[i] = snowflakeResponse{ID: id, IDString: strconv.FormatUint(id, 10)}
		}

		json.NewEncoder(w).Encode(ids)
	}).Methods("GET")

	// Listen
	fmt.Printf("Server listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
