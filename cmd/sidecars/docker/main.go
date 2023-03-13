package main

import (
	"log"
	"net/http"

	docker "github.com/cloud-pg/interlink/pkg/sidecars/docker"
)

func main() {

	mutex := http.NewServeMux()
	mutex.HandleFunc("/status", docker.StatusHandler)
	mutex.HandleFunc("/create", docker.CreateHandler)
	mutex.HandleFunc("/delete", docker.DeleteHandler)

	err := http.ListenAndServe(":4000", mutex)
	if err != nil {
		log.Fatal(err)
	}
}
