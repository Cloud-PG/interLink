package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	mutex.HandleFunc("/status", statusHandler)
	mutex.HandleFunc("/create", createHandler)
	mutex.HandleFunc("/delete", deleteHandler)

	err := http.ListenAndServe(":3000", mutex)
	if err != nil {
		log.Fatal(err)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker get status
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req statusRequest
	json.Unmarshal(bodyBytes, &req)

	for _, pod := range req.Pods {
		fmt.Println("GET " + pod.Name + " status")
	}

	w.Write([]byte("200: OK"))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker create container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req jsonRequest
	json.Unmarshal(bodyBytes, &req)
	fmt.Println("CREATE " + req.Name + " pod")
	w.Write([]byte("201: OK"))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker delete container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req jsonRequest
	json.Unmarshal(bodyBytes, &req)
	fmt.Println("DELETE " + req.Name + " pod")
	w.Write([]byte("202: OK"))
}
