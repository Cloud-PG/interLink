package common

import (
	v1 "k8s.io/api/core/v1"
)

const (
	RUNNING = 0
	STOP    = 1
	UNKNOWN = 2
	SBATCH  = "/opt/slurm/current/bin/sbatch"
	SCANCEL = "/opt/slurm/current/bin/sbatch"
)

type PodUID struct {
	UID string `json:"podUID"`
}

type PodStatus struct {
	PodStatus uint `json:"podStatus"`
}

type StatusResponse struct {
	PodUID    []PodUID    `json:"poduid"`
	PodStatus []PodStatus `json:"podstatus"`
	ReturnVal string      `json:"returnVal"`
}

type Request struct {
	Pods map[string]*v1.Pod `json:"pods"`
}

type InterLinkConfig struct {
	Interlinkurl   string
	Sidecarurl     string
	Interlinkport  string
	Sidecarport    string
	Sidecarservice string
	set            bool
}
