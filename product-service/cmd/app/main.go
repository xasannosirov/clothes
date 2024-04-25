package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"product-service/internal/app"
	"product-service/internal/pkg/config"

	"go.uber.org/zap"
)

func main() {
	config := config.New()

	app, err := app.NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	// consumer init
	// go func() {
	// 	// fmt.Println("worked app.BrokerConsumer.Run()")
	// 	app.BrokerConsumer.Run()
	// }()

	//running
	go func() {
		if err := app.Run(); err != nil {
			app.Logger.Error("app run", zap.Error(err))
		}
	}()

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	app.Logger.Info("Product service stops !")

	// app stops
	app.Stop()
}
