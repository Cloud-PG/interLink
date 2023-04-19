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
	request := commonIL.Request{}
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
	var returnValue, _ = json.Marshal(commonIL.PodStatus{PodStatus: commonIL.UNKNOWN})

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

func statusRequest(podsList commonIL.Request) []byte {
	var returnValue []byte
	var response []commonIL.StatusResponse

	bodyBytes, err := json.Marshal(podsList)
	reader := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodGet, commonIL.InterLinkConfigInst.Interlinkurl+":"+commonIL.InterLinkConfigInst.Interlinkport+"/status", reader)
	if err != nil {
		log.L.Error(err)
	}

	log.L.Println(string(bodyBytes))

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
	var jsonVar commonIL.Request

	switch mode {
	case common.CREATE:
		//v1.Pod used only for secrets and volumes management; TO BE IMPLEMENTED
		if err != nil {
			return err
		}

		jsonVar = commonIL.Request{Pods: map[string]*v1.Pod{
			pod.Name: pod,
		}}

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
	var returnVal []byte
	var PodsList commonIL.Request
	PodsList.Pods = p.pods

	returnVal = statusRequest(PodsList)
	log.G(ctx).Infof(string(returnVal))
}
