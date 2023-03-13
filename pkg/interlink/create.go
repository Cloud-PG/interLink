package interlink

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	types "github.com/cloud-pg/interlink/pkg/common"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodPost, types.InterLinkConfigInst.Sidecarurl+":"+types.InterLinkConfigInst.Sidecarport+"/create", reader)
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
