package interlink

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	types "github.com/cloud-pg/interlink/pkg/common"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req *http.Request
	reader := bytes.NewReader(bodyBytes)

	switch types.InterLinkConfigInst.Sidecarservice {
	case "docker":
		req, err = http.NewRequest(http.MethodPost, types.InterLinkConfigInst.Sidecarurl+":"+types.InterLinkConfigInst.Sidecarport+"/delete", reader)

	case "slurm":
		req, err = http.NewRequest(http.MethodPost, types.InterLinkConfigInst.Sidecarurl+":"+types.InterLinkConfigInst.Sidecarport+"/stop", reader)

	default:
		break
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	returnValue, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Deleted container " + string(returnValue))

	w.Write(returnValue)
}
