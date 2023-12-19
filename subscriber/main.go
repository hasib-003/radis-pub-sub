package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
)

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	e := echo.New()

	e.POST("/publish", func(c echo.Context) error {
		message := "publishing from subscriber"
		err := redisClient.Publish(ctx, "subscriber-channel", message)
		if err != nil {
			return err.Err()
		}
		return c.NoContent(http.StatusOK)
	})

	subscriber := redisClient.Subscribe(ctx, "publisher-chan")
	defer subscriber.Close()

	go func() {

		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println("Error receiving message:", err)
				continue
			}
			fmt.Println(msg.Payload)
		}
	}()

	e.Logger.Fatal(e.Start(":3001"))

}
