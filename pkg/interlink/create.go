package interlink

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker create container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:4000/create", reader)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	returnValue, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(returnValue))

	w.Write(returnValue)
}
