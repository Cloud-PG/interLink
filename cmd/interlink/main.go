package main

import (
	"log"
	"net/http"

	"github.com/cloud-pg/interlink/pkg/interlink"

	v1 "k8s.io/api/core/v1"
)

type statusRequest struct {
	Pods []jsonRequest `json:"pods"`
}

type jsonRequest struct {
	v1.Pod
}

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
