package main

import (
	"context"

	"github.com/MohammadAsDev/geo_tracker/src/config"
	"github.com/MohammadAsDev/geo_tracker/src/usecases"
	"github.com/MohammadAsDev/geo_tracker/src/usecases/tracker"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	tracker := usecases.NewTrackingSystem(context.Background(), tracker.DEFAULT, config)
	tracker.StartTrackingSystem()

}
