package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"

	"pow/internal/config"
	"pow/internal/server"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	configInst, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	ctx := context.WithValue(context.Background(), "config", configInst)

	rand.Seed(time.Now().UnixNano())

	serverAddress := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)
	err = server.Run(ctx, serverAddress, log)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
