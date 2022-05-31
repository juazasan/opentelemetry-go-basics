package main

import (
	"context"

	"github.com/juazasan/opentelemetry-go-basics/pkg/client"
)

func main() {
	client := client.NewClient("http://localhost:9000/hello")
	client.Run(context.Background())
	client.Shutdown()
}
