package main

import (
	_ "routers"
	"os"
	"errors"
	"utils"
	"mirbase"
	"github.com/stainberg/koalart"
	"os/signal"
	"syscall"
	"log"
	"net/http"
)

func StaticServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	staticHandler := http.FileServer(http.Dir("/home/rbo/"))
	staticHandler.ServeHTTP(w, r)
	return
}

func main() {
	if len(os.Args) == 1 {
		utils.LoadConfig("")
	} else if len(os.Args) == 2 {
		utils.LoadConfig(os.Args[1])
	} else {
		panic(errors.New("params error"))
	}
	mirbase.InitClient()
	go func() {
		koala.Run(utils.Conf.HttpConf.RestAPIPort)
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	log.Println("Shutdown signal received, exiting...")
}
