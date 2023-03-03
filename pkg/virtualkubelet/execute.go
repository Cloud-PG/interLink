package virtualkubelet

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	common "github.com/CARV-ICS-FORTH/knoc/common"
	types "github.com/cloud-pg/interlink/pkg/common"

	"github.com/containerd/containerd/log"
	v1 "k8s.io/api/core/v1"
)

func createRequest(jsonBody []byte) {
	request := types.CreateRequest{}
	json.Unmarshal(jsonBody, &request)

	reader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/create", reader)
	if err != nil {
		log.L.Error(err)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.L.Error(err)
	}
}

func deleteRequest(jsonBody []byte) []byte {
	request := types.CreateRequest{}
	var returnValue, _ = json.Marshal(types.PodStatus{PodStatus: "UNDEFINED"})
	json.Unmarshal(jsonBody, &request)

	reader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:3000/delete", reader)
	if err != nil {
		log.L.Error(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.L.Error(err)
	}

	returnValue, _ = ioutil.ReadAll(resp.Body)
	var response types.PodStatus
	json.Unmarshal(returnValue, &response)

	return returnValue
}

func statusRequest(jsonBody []byte) []byte {
	var request types.StatusRequest
	var returnValue []byte
	var response []types.StatusResponse

	returnValue, _ = json.Marshal(response)
	json.Unmarshal(jsonBody, &request)

	reader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/status", reader)
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
	var jsonVar types.CreateRequest

	switch mode {
	case common.CREATE:
		//v1.Pod used only for secrets and volumes management; TO BE IMPLEMENTED
		if err != nil {
			return err
		}

		jsonVar = types.CreateRequest{container}
		jsonBytes, _ := json.Marshal(jsonVar)
		createRequest(jsonBytes)
		break

	case common.DELETE:
		request := types.PodUID{UID: string(container.Name)}
		jsonBytes, _ := json.Marshal(request)
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
	var podsList types.StatusRequest

	for _, pod := range p.pods {
		for _, container := range pod.Spec.Containers {
			podsList.PodUIDs = append(podsList.PodUIDs, types.PodUID{UID: container.Name})
		}

	}

	jsonBytes, err := json.Marshal(podsList)

	if err != nil {
		log.L.Error(err)
	}

	returnVal = statusRequest(jsonBytes)
	log.G(ctx).Infof(string(returnVal))
}
