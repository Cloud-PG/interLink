package main

import (
	"log"
	"net/http"

	"github.com/cloud-pg/interlink/pkg/interlink"
)

func main() {

	mutex := http.NewServeMux()
	mutex.HandleFunc("/status", interlink.StatusHandler)
	mutex.HandleFunc("/create", interlink.CreateHandler)
	mutex.HandleFunc("/delete", interlink.DeleteHandler)

	err := http.ListenAndServe(":3000", mutex)
	if err != nil {
		log.Fatal(err)
	}
}
