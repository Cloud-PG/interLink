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
			path = "/etc/interlink/InterLinkConfig.yaml"
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
				os.Exit(-1)
			}
			InterLinkConfigInst.Sidecarservice = os.Getenv("SIDECARSERVICE")
		} else if InterLinkConfigInst.Sidecarservice != "docker" && InterLinkConfigInst.Sidecarservice != "slurm" {
			fmt.Println("Set \"docker\" or \"slurm\" in config file or export SIDECARSERVICE as ENV")
			os.Exit(-1)
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
				os.Exit(-1)
			}
		}

		if os.Getenv("SBATCHPATH") != "" {
			InterLinkConfigInst.Sbatchpath = os.Getenv("SBATCHPATH")
		}

		if os.Getenv("SCANCELPATH") != "" {
			InterLinkConfigInst.Scancelpath = os.Getenv("SCANCELPATH")
		}

		if os.Getenv("TSOCKS") != "" {
			if os.Getenv("TSOCKS") != "true" && os.Getenv("TSOCKS") != "false" {
				fmt.Println("export TSOCKS as true or false")
				os.Exit(-1)
			}
			if os.Getenv("TSOCKS") == "true" {
				InterLinkConfigInst.Tsocks = true
			} else {
				InterLinkConfigInst.Tsocks = false
			}
		}

		if os.Getenv("TSOCKSPATH") != "" {
			path := os.Getenv("TSOCKSPATH")
			if _, err := os.Stat(path); err != nil {
				log.Println("File " + path + " doesn't exist. You can set a custom path by exporting TSOCKSPATH. Exiting...")
				os.Exit(-1)
			}

			InterLinkConfigInst.Tsockspath = path
		}

		InterLinkConfigInst.set = true
	}
}
