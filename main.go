package main

import (
	// "encoding/json"
	// "fmt"
	"log/slog"
	"os"
	// "net/http"

	// "github.com/gorilla/mux"
	
)

func main(){
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger.Info("hi")

}