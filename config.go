package main

import (
	"time"

	"github.com/everettcaleb/envconfig"
)

type snowflakeEnvConfig struct {
	BasePath             string `env:"APP_BASE_PATH"`
	Epoch                uint64 `env:"SNOWFLAKE_EPOCH"`
	Port                 int    `env:"PORT"`
	RedisURI             string `env:"REDIS_URI" required:"true"`
	RedisMachineIDPrefix string `env:"REDIS_MACHINE_ID_PREFIX"`
	UseMilliseconds      bool   `env:"SNOWFLAKE_USE_MILLISECONDS"`
}

func defaultConfig(epoch time.Time) *snowflakeEnvConfig {
	return &snowflakeEnvConfig{
		BasePath:             "/",
		Epoch:                uint64(epoch.Unix()),
		Port:                 8080,
		RedisMachineIDPrefix: "snowflake:machine:",
		RedisURI:             "redis://localhost:6379",
		UseMilliseconds:      false,
	}
}

func loadEnvConfig() (*snowflakeEnvConfig, error) {
	// Parse the default epoch
	epoch, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
	if err != nil {
		return nil, err
	}

	// Set the defaults
	config := defaultConfig(epoch)

	// Load from the environment
	err = envconfig.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
