package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var InterLinkConfigInst InterLinkConfig

func NewInterLinkConfig() {
	if InterLinkConfigInst.set == false {
		var path string
		if os.Getenv("INTERLINKCONFIGPATH") != "" {
			path = os.Getenv("INTERLINKCONFIGPATH")
		} else {
			path = "$HOME/.config/InterLinkConfig.yaml"
		}

		if _, err := os.Stat(path); err != nil {
			log.Println("File " + path + " doesn't exist. You can set a custom path by exporting INTERLINKCONFIGPATH. Exiting...")
			os.Exit(-1)
		}

		yfile, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("Error opening config file, exiting...")
			os.Exit(1)
		}
		yaml.Unmarshal(yfile, &InterLinkConfigInst)

		if os.Getenv("INTERLINKURL") != "" {
			InterLinkConfigInst.Interlinkurl = os.Getenv("INTERLINKURL")
		}

		if os.Getenv("SIDECARURL") != "" {
			InterLinkConfigInst.Sidecarurl = os.Getenv("SIDECARURL")
		}

		if os.Getenv("INTERLINKPORT") != "" {
			InterLinkConfigInst.Interlinkport = os.Getenv("INTERLINKPORT")
		}

		if os.Getenv("SIDECARSERVICE") != "" {
			if os.Getenv("SIDECARSERVICE") != "docker" && os.Getenv("SIDECARSERVICE") != "slurm" {
				fmt.Println("export SIDECARSERVICE as docker or slurm")
				return
			}
			InterLinkConfigInst.Sidecarservice = os.Getenv("SIDECARSERVICE")
		} else if InterLinkConfigInst.Sidecarservice != "docker" && InterLinkConfigInst.Sidecarservice != "slurm" {
			fmt.Println("Set \"docker\" or \"slurm\" in config file or export SIDECARSERVICE as ENV")
			return
		}

		if os.Getenv("SIDECARPORT") != "" && os.Getenv("SIDECARSERVICE") == "" {
			InterLinkConfigInst.Sidecarport = os.Getenv("SIDECARPORT")
			InterLinkConfigInst.Sidecarservice = "Custom Service"
		} else {
			switch InterLinkConfigInst.Sidecarservice {
			case "docker":
				InterLinkConfigInst.Sidecarport = "4000"

			case "slurm":
				InterLinkConfigInst.Sidecarport = "4001"

			default:
				fmt.Println("Define in InterLinkConfig.yaml one service between docker and slurm")
				return
			}
		}
		InterLinkConfigInst.set = true
	}
}
