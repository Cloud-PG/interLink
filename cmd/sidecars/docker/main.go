package main

import (
	"log"
	"net/http"

	mock "github.com/cloud-pg/interlink/pkg/sidecars/docker"
)

func main() {

	mutex := http.NewServeMux()
	mutex.HandleFunc("/status", mock.StatusHandler)
	mutex.HandleFunc("/create", mock.CreateHandler)
	mutex.HandleFunc("/delete", mock.DeleteHandler)

	err := http.ListenAndServe(":4000", mutex)
	if err != nil {
		log.Fatal(err)
	}
}
