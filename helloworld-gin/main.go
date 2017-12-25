package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request
		c.Set("logger", "mylogger")
		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func main() {
	// r := gin.Default()
	r := gin.New()
	r.Use(Logger())
	r.GET("/ping", func(c *gin.Context) {
		fmt.Println(c.Get("logger"))
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
