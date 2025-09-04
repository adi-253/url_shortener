package main

import (
	"github.com/adi-253/url_shortener/api"
	"net/http"
	"log/slog"
	"os"

)

func main(){
	logger:= slog.New(slog.NewTextHandler(os.Stderr, nil))
	server := api.InitServer()
	
	http.ListenAndServe(":8080", server.Router)
	logger.Info("Server Stopped")

}