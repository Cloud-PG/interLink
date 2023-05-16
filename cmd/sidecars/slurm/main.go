package main

import (
	"log"
	"net/http"

	commonIL "github.com/cloud-pg/interlink/pkg/common"
	slurm "github.com/cloud-pg/interlink/pkg/sidecars/slurm"
)

func main() {

	commonIL.NewInterLinkConfig()

	mutex := http.NewServeMux()
	mutex.HandleFunc("/status", slurm.StatusHandler)
	mutex.HandleFunc("/submit", slurm.SubmitHandler)
	mutex.HandleFunc("/stop", slurm.StopHandler)

	err := http.ListenAndServe(":4001", mutex)
	if err != nil {
		log.Fatal(err)
	}
}
