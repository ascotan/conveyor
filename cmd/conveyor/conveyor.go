package main

import (
	"github.com/ascotan/conveyor/internal/log"
	"net/http"
)

func main() {
	log.Logger.Info("Starting conveyor on http://localhost:8081/")
	err := http.ListenAndServe(":8081", http.FileServer(http.Dir(".")))
	log.Logger.Error(err.Error())
}
