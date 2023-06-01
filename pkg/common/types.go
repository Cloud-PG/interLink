package common

import (
	v1 "k8s.io/api/core/v1"
)

const (
	RUNNING = 0
	STOP    = 1
	UNKNOWN = 2
	SBATCH  = "/usr/bin/sbatch"
	SCANCEL = "/usr/bin/sbatch"
)

type PodName struct {
	Name string `json:"podname"`
}

type PodStatus struct {
	PodStatus uint `json:"podStatus"`
}

type StatusResponse struct {
	PodName   []PodName   `json:"podname"`
	PodStatus []PodStatus `json:"podstatus"`
	ReturnVal string      `json:"returnVal"`
}

type Request struct {
	Pods map[string]*v1.Pod `json:"pods"`
}

type InterLinkConfig struct {
	Interlinkurl   string `yaml:"InterlinkURL"`
	Sidecarurl     string `yaml:"SidecarURL"`
	Interlinkport  string `yaml:"InterlinkPort"`
	Sidecarport    string
	Sidecarservice string `yaml:"SidecarService"`
	Commandprefix  string `yaml:"CommandPrefix"`
	Tsocks         bool   `yaml:"Tsocks"`
	Tsockspath     string `yaml:"TsocksPath"`
	Tsocksconfig   string `yaml:"TsocksConfig"`
	Tsockslogin    string `yaml:"TsocksLoginNode"`
	set            bool
}
