package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chungweeeei/ShortenURL/data"
	"github.com/chungweeeei/ShortenURL/helpers"
)

const (
	serverPort = "80"
)

func main() {

	db := initDB()
	rdb := initRedis()

	// init log instance
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	// id generator
	generator, _ := helpers.NewGenerator(1)

	app := Config{
		DB:            db,
		RDB:           rdb,
		Generator:     generator,
		Model:         data.New(db),
		InfoLog:       infoLog,
		Errorlog:      errorLog,
		ErrorChan:     make(chan error),
		ErrorDoneChan: make(chan bool),
	}

	go app.listenForErrors()

	go app.listenForShutdown()

	app.serve()
}

func (app *Config) serve() {

	// start running http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting Shorten URL service...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) listenForShutdown() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) listenForErrors() {

	for {
		select {
		case err := <-app.ErrorChan:
			app.Errorlog.Println(err)
		case <-app.ErrorDoneChan:
			// close the goroutine
			return
		}
	}

}

func (app *Config) shutdown() {

	// perform any cleanup tasks
	app.InfoLog.Println("Would run cleanup tasks...")

	// notify "listenForErrors" channel to close
	app.ErrorDoneChan <- true

	// shutdown
	app.InfoLog.Println("closing channels and shutting down application...")
	close(app.ErrorChan)
	close(app.ErrorDoneChan)
}
