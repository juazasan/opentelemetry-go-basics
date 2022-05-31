package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/juazasan/opentelemetry-go-basics/pkg/runtime"
)

func main() {
	rt := runtime.NewRuntime()
	rt.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop
	rt.Shutdown()
}
