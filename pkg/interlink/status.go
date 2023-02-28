package interlink

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker get status
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodGet, "http://localhost:4000/status", reader)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("1st Layer: forwarding GET to 2nd Layer")
	returnValue, _ := ioutil.ReadAll(resp.Body)

	w.Write(returnValue)
}
