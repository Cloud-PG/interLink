package common

import v1 "k8s.io/api/core/v1"

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
	v1.Container
}
