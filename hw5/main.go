package main

import (
	"flag"
	"hw5/server"
	"hw5/utils"
	"log"
	"os"
	"os/signal"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	flagRootDir := flag.String("rootdir", "./www", "root dir of the server")
	flagServAddr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	lg := utils.NewLogger()

	db, err := gorm.Open("mysql", "admin:admin@tcp(172.20.0.2:3306)/gorm?parseTime=true")
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}
	defer db.Close()

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
