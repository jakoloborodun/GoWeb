package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hw7/conf"
	"hw7/server"
	"hw7/utils"
	"os"
	"os/signal"
	"time"
)

func main() {
	stopChan := make(chan os.Signal)

	cfgPath, err := conf.ParseConfigFlag()
	if err != nil {
		panic(fmt.Sprintf("can't find config file: %s", err))
	}
	cfg, err := conf.NewConfig(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("can't read config file: %s", err))
	}

	lg, err := utils.NewLogger(cfg)
	if err != nil {
		panic(fmt.Sprintf("can't configure logger: %s", err))
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	dbUri := fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}
	lg.Info("Connect successful to db")
	defer client.Disconnect(ctx)

	db := client.Database(cfg.Database.Name)

	srv := server.New(lg, cfg.Server.RootDir, db, cfg)
	srv.ParseTemplates()

	// Start the server in a new goroutine
	go func() {
		err := srv.Start(cfg.Server.Host + ":" + cfg.Server.Port)
		if err != nil {
			lg.WithError(err).Fatal("can't start the server")
		}
	}()

	signal.Notify(stopChan, os.Interrupt, os.Kill)
	interrupt := <-stopChan
	lg.Printf("Server is shutting down due to %+v\n", interrupt)
}
