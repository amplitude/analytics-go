package main

import (
	"github.com/amplitude/Amplitude-Go/amplitude"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config := amplitude.NewConfig("your-app-id")
	client := amplitude.NewClient(config)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/analytics", func(c *gin.Context) {

		// Create a BaseEvent instance
		event := amplitude.Event{
			EventOptions: amplitude.EventOptions{DeviceID: "gin-device-id", UserID: "gin-user-id"},
			EventType:    "gin-event-type",
		}

		// Track an event
		client.Track(event)
		defer client.Shutdown()

		c.String(http.StatusOK, "Amplitude analytics")
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
