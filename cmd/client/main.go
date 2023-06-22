package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"pow/internal/client"
	"pow/internal/pkg/config"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	log.Info("start client")

	configInst, err := config.Load("config/config.yaml")
	if err != nil {
		log.Errorf("error loading config: %v", err)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)

	address := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)

	err = client.Run(ctx, address, log)
	if err != nil {
		log.Errorf("client error: %v", err)
	}
}
