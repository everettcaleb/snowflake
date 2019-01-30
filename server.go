package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	machineID, err := initMachineID(config)
	if err != nil {
		panic(err)
	}

	generator := newSnowflakeGenerator(config, machineID)
	router := gin.Default()
	routeGroup := router.Group(config.BasePath)

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
	routeGroup.GET("id", func(c *gin.Context) {
		id := generator.NextID()
		c.JSON(200, &snowflakeResponse{
			ID:       id,
			IDString: string(id),
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
				IDString: string(id),
			}
		}

		c.JSON(200, ids)
	})

	// Listen
	fmt.Printf("Server listening on %d\n", config.Port)
	router.Run(":" + string(config.Port))
}
