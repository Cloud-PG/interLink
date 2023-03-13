package main

import (
	"log"
	"net/http"

	slurm "github.com/cloud-pg/interlink/pkg/sidecars/docker"
)

func main() {

	mutex := http.NewServeMux()
	mutex.HandleFunc("/status", slurm.StatusHandler)
	mutex.HandleFunc("/create", slurm.CreateHandler)
	mutex.HandleFunc("/delete", slurm.DeleteHandler)

	err := http.ListenAndServe(":4001", mutex)
	if err != nil {
		log.Fatal(err)
	}
}
