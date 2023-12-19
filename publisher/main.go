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
		message := "publishing from publisher"
		if err := redisClient.Publish(ctx, "publisher-chan", message).Err(); err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)

	})

	subscriber := redisClient.Subscribe(ctx, "subscriber-channel")
	defer subscriber.Close()
	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(msg.Payload)
		}
	}()

	e.Logger.Fatal(e.Start(":3000"))

}
