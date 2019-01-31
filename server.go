package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func initGenerator() (*snowflakeEnvConfig, *snowflakeGenerator) {
	config, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	machineID, err := initMachineID(config)
	if err != nil {
		panic(err)
	}

	return config, newSnowflakeGenerator(config, machineID)
}

func initRoutes(routeGroup *gin.RouterGroup, generator *snowflakeGenerator) {
	// Health Check
	routeGroup.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, healthCheckResponse{Status: "OK"})
	})

	// Single Snowflake ID
	routeGroup.GET("id", func(c *gin.Context) {
		id := generator.NextID()
		c.JSON(200, &snowflakeResponse{
			ID:       id,
			IDString: strconv.FormatUint(uint64(id), 10),
		})
	})

	// Multiple Snowflake IDs
	routeGroup.GET("/ids/:count", func(c *gin.Context) {
		str := c.Params.ByName("count")
		count, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			c.String(400, "text/plain", "Count must be a number")
			return
		}

		var i int64
		ids := make([]*snowflakeResponse, count)
		for i = 0; i < count; i++ {
			id := generator.NextID()
			ids[i] = &snowflakeResponse{
				ID:       id,
				IDString: strconv.FormatUint(uint64(id), 10),
			}
		}

		c.JSON(200, ids)
	})
}

func main() {
	// Set up the generator and router
	router := gin.Default()
	config, generator := initGenerator()
	initRoutes(router.Group(config.BasePath), generator)

	// Print out the configuration for debugging purposes
	fmt.Println()
	fmt.Println("Config:")
	fmt.Println("==========")
	fmt.Println("Machine ID:", generator.machineID)
	fmt.Println("Base Path:", config.BasePath)
	fmt.Println("Epoch:", config.Epoch)
	fmt.Println("Port:", config.Port)
	fmt.Println("Redis Machine ID Prefix:", config.RedisMachineIDPrefix)
	fmt.Println("Redis URI:", config.RedisURI)
	fmt.Println("Use Milliseconds:", config.UseMilliseconds)
	fmt.Println()

	// Listen for requests
	fmt.Printf("Server listening on %v\n", config.Port)
	router.Run(":" + strconv.Itoa(config.Port))
}
