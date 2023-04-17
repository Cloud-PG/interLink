package virtualkubelet

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	common "github.com/CARV-ICS-FORTH/knoc/common"
	commonIL "github.com/cloud-pg/interlink/pkg/common"

	"github.com/containerd/containerd/log"
	v1 "k8s.io/api/core/v1"
)

func createRequest(jsonBody []byte) {
	request := commonIL.CreateRequest{}
	json.Unmarshal(jsonBody, &request)
	var req *http.Request
	var err error

	reader := bytes.NewReader(jsonBody)
	req, err = http.NewRequest(http.MethodPost, commonIL.InterLinkConfigInst.Interlinkurl+":"+commonIL.InterLinkConfigInst.Interlinkport+"/create", reader)

	if err != nil {
		log.L.Error(err)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.L.Error(err)
	}
}

func deleteRequest(jsonBody []byte) []byte {
	var returnValue, _ = json.Marshal(commonIL.PodStatus{PodStatus: "UNDEFINED"})

	reader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodDelete, commonIL.InterLinkConfigInst.Interlinkurl+":"+commonIL.InterLinkConfigInst.Interlinkport+"/delete", reader)
	if err != nil {
		log.L.Error(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.L.Error(err)
	}

	returnValue, _ = ioutil.ReadAll(resp.Body)
	var response commonIL.PodStatus
	json.Unmarshal(returnValue, &response)

	return returnValue
}

func statusRequest(jsonBody []byte) []byte {
	var request commonIL.StatusRequest
	var returnValue []byte
	var response []commonIL.StatusResponse

	returnValue, _ = json.Marshal(response)
	json.Unmarshal(jsonBody, &request)

	reader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodGet, commonIL.InterLinkConfigInst.Interlinkurl+":"+commonIL.InterLinkConfigInst.Interlinkport+"/status", reader)
	if err != nil {
		log.L.Error(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.L.Error(err)
	}

	returnValue, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(returnValue, &response)
	return returnValue
}

func RemoteExecution(p *VirtualKubeletProvider, ctx context.Context, mode int8, imageLocation string, pod *v1.Pod, container v1.Container) error {
	var err error
	var jsonVar commonIL.CreateRequest

	switch mode {
	case common.CREATE:
		//v1.Pod used only for secrets and volumes management; TO BE IMPLEMENTED
		if err != nil {
			return err
		}

		jsonVar = commonIL.CreateRequest{Container: container, Pod: *pod}
		jsonBytes, _ := json.Marshal(jsonVar)
		createRequest(jsonBytes)
		break

	case common.DELETE:
		//request := types.PodUID{UID: string(container.Name)}
		jsonBytes, _ := json.Marshal(container)
		returnVal := deleteRequest(jsonBytes)
		log.G(ctx).Infof(string(returnVal))

		if err != nil {
			return err
		}
		break
	}
	return nil
}

func checkPodsStatus(p *VirtualKubeletProvider, ctx context.Context) {
	if len(p.pods) == 0 {
		return
	}
	var jsonBytes []byte
	var returnVal []byte
	var podsList commonIL.StatusRequest

	for _, pod := range p.pods {
		for _, container := range pod.Spec.Containers {
			podsList.PodUIDs = append(podsList.PodUIDs, commonIL.PodUID{UID: container.Name})
		}

	}

	jsonBytes, err := json.Marshal(podsList)

	if err != nil {
		log.L.Error(err)
	}

	returnVal = statusRequest(jsonBytes)
	log.G(ctx).Infof(string(returnVal))
}
