package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	types "github.com/cloud-pg/interlink/pkg/common"
	"github.com/cloud-pg/interlink/pkg/interlink"
	"gopkg.in/yaml.v3"
)

var Url string

func main() {

	yfile, err := ioutil.ReadFile("/etc/interlink/InterLinkConfig.yaml")
	yaml.Unmarshal(yfile, &types.InterLinkConfigInst)

	if os.Getenv("INTERLINKURL") != "" {
		types.InterLinkConfigInst.Interlinkurl = os.Getenv("INTERLINKURL")
	}

	if os.Getenv("SIDECARURL") != "" {
		types.InterLinkConfigInst.Sidecarurl = os.Getenv("SIDECARURL")
	}

	if os.Getenv("INTERLINKPORT") != "" {
		types.InterLinkConfigInst.Interlinkport = os.Getenv("INTERLINKPORT")
	}

	if os.Getenv("SIDECARSERVICE") != "" {
		if os.Getenv("SIDECARSERVICE") != "docker" && os.Getenv("SIDECARSERVICE") != "slurm" {
			fmt.Println("export SIDECARSERVICE as docker or slurm")
			return
		}
		types.InterLinkConfigInst.Sidecarservice = os.Getenv("SIDECARSERVICE")
	} else if types.InterLinkConfigInst.Sidecarservice != "docker" && types.InterLinkConfigInst.Sidecarservice != "slurm" {
		fmt.Println("Set \"docker\" or \"slurm\" in config file or export SIDECARSERVICE as ENV")
		return
	}

	if os.Getenv("SIDECARPORT") != "" && os.Getenv("SIDECARSERVICE") == "" {
		types.InterLinkConfigInst.Sidecarport = os.Getenv("SIDECARPORT")
		types.InterLinkConfigInst.Sidecarservice = "Custom Service"
	} else {
		switch types.InterLinkConfigInst.Sidecarservice {
		case "docker":
			types.InterLinkConfigInst.Sidecarport = "4000"

		case "slurm":
			types.InterLinkConfigInst.Sidecarport = "4001"

		default:
			fmt.Println("Define in InterLinkConfig.yaml one service between docker and slurm")
			return
		}
	}

	mutex := http.NewServeMux()
	mutex.HandleFunc("/status", interlink.StatusHandler)
	mutex.HandleFunc("/create", interlink.CreateHandler)
	mutex.HandleFunc("/delete", interlink.DeleteHandler)

	fmt.Println(types.InterLinkConfigInst)

	err = http.ListenAndServe(":"+types.InterLinkConfigInst.Interlinkport, mutex)
	if err != nil {
		log.Fatal(err)
	}
}
