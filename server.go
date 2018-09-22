package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
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

	generator := makeSnowflakeGenerator(machineID, epoch)
	router := gin.Default()
	routeGroup := router.Group(getBasePath())

	// Health Check
	routeGroup.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, healthCheckResponse{Status: "OK"})
	})

	// Specs
	routeGroup.GET("/spec.yaml", func(c *gin.Context) {
		c.Header("Content-Type", "application/x-yaml")
		c.File("specs/spec.yaml")
	})

	routeGroup.GET("/spec.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.File("specs/spec.json")
	})

	// Single Snowflake ID
	routeGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, makeSnowflakeResponse(generator.NextID()))
	})

	// Multiple Snowflake IDs
	routeGroup.GET("/:count", func(c *gin.Context) {
		str := c.Params.ByName("count")
		count, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			c.String(400, "text/plain", "Count must be a number")
			return
		}

		var i int64
		ids := make([]*snowflakeResponse, count)
		for i = 0; i < count; i++ {
			ids[i] = makeSnowflakeResponse(generator.NextID())
		}

		c.JSON(200, ids)
	})

	// Listen
	fmt.Printf("Server listening on %s\n", port)
	router.Run(port)
}
