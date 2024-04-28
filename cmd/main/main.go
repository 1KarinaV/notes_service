package main

import (
	"go.uber.org/zap"
	"log"
	"notes_service/internal/app"
)

func main() {
	logger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatal()
	}

	zap.ReplaceGlobals(logger)

	a, err := app.New()
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	err = a.Init()
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	err = a.Run()
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}
