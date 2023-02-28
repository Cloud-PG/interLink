package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	types "github.com/cloud-pg/interlink/pkg/common"
)

func main() {

	mutex := http.NewServeMux()
	mutex.HandleFunc("/status", statusHandler)
	mutex.HandleFunc("/create", createHandler)
	mutex.HandleFunc("/delete", deleteHandler)

	err := http.ListenAndServe(":4000", mutex)
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

	var req types.StatusRequest
	json.Unmarshal(bodyBytes, &req)

	var resp []types.StatusResponse

	for _, pod := range req.PodUIDs {
		temp := types.StatusResponse{}
		temp.UID = pod.UID
		temp.PodStatus.PodStatus = "RUNNING"
		temp.ReturnVal = "200: OK"
		resp = append(resp, temp)
	}

	fmt.Println(resp)
	respBytes, _ := json.Marshal(resp)
	w.Write(respBytes)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker create container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req types.CreateRequest
	json.Unmarshal(bodyBytes, &req)
	fmt.Println("CREATE " + req.UID + " pod")
	w.Write([]byte("201: OK"))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker delete container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req types.PodUID
	json.Unmarshal(bodyBytes, &req)
	fmt.Println("DELETE " + req.UID + " pod")
	var resp types.PodStatus
	resp.PodStatus = "TERMINATING"

	respJson, _ := json.Marshal(resp)
	w.Write(respJson)
}
