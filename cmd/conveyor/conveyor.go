package main

import (
	"github.com/ascotan/conveyor/internal/log"
	"net/http"
)

func main() {
	log.Logger.Info("Starting conveyor")
	http.ListenAndServe(":8081", http.FileServer(http.Dir(".")))
}
