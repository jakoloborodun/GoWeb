package main

import (
	"context"
	"flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hw7/server"
	"hw7/utils"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	flagRootDir := flag.String("rootdir", "./www", "root dir of the server")
	flagServAddr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	lg := utils.NewLogger()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}
	defer client.Disconnect(ctx)

	db := client.Database("gomongo")

	srv := server.New(lg, *flagRootDir, db)
	srv.ParseTemplates()

	go func() {
		err := srv.Start(*flagServAddr)
		if err != nil {
			lg.WithError(err).Fatal("can't start the server")
		}
	}()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan
	log.Print("Shutting down...")
}
