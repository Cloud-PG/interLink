package common

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodUID struct {
	UID string `json:"podUID"`
}

type PodStatus struct {
	PodStatus string `json:"podStatus"`
}

type StatusResponse struct {
	PodUID
	PodStatus
	ReturnVal string `json:"returnVal"`
}

type StatusRequest struct {
	PodUIDs []PodUID `json:"podUIDs"`
}

type CreateRequest struct {
	Container v1.Container
	Metadata  metav1.ObjectMeta
	Pod       v1.Pod
}

type InterLinkConfig struct {
	Interlinkurl  string
	Sidecarurl    string
	Interlinkport string
	Sidecarport   string
	Service       string
}

var InterLinkConfigInst InterLinkConfig
