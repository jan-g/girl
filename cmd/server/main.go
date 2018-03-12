package main

import (
	"context"
	"os"
	"time"

	"github.com/jan-g/girl/model"
	"github.com/jan-g/girl/server"
	"github.com/sirupsen/logrus"
)

func main() {
	me := os.Args[1]
	peer := os.Args[2]

	if level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		logrus.SetLevel(level)
	}

	epoch := time.Now().Unix()
	limiter, control := model.NewLimiter(me, epoch)

	control.AddLimit("foo", 10, 1, 1)
	control.AddLimit("bar", 10, 10, 1)

	spi := limiter.(model.LimiterSPI)
	srv, err := server.NewServer("tcp", me, spi)
	if err != nil {
		panic(err)
	}
	srv.AddPeer(peer)
	srv.Start(context.Background())
	srv.Serve()
}
