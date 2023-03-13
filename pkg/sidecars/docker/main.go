package docker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	exec "github.com/alexellis/go-execute/pkg/v1"
	types "github.com/cloud-pg/interlink/pkg/common"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker get status

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req types.StatusRequest
	json.Unmarshal(bodyBytes, &req)

	cmd := []string{"ps -aqf \"name="}
	for _, pod := range req.PodUIDs {
		cmd[0] += " " + pod.UID
	}

	cmd[0] += "\""

	shell := exec.ExecTask{
		Command: "docker",
		Args:    cmd,
		Shell:   true,
	}

	execReturn, err := shell.Execute()
	execReturn.Stdout = strings.ReplaceAll(execReturn.Stdout, "\n", "")

	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(execReturn.Stdout))
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker create container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req types.CreateRequest
	json.Unmarshal(bodyBytes, &req)

	cmd := []string{"run", "-d", "--name", req.Container.Name}
	for _, args := range req.Container.Args {
		cmd = append(cmd, args)
	}

	cmd = append(cmd, req.Container.Image)

	for _, command := range req.Container.Command {
		cmd = append(cmd, command)
	}

	shell := exec.ExecTask{
		Command: "docker",
		Args:    cmd,
		Shell:   true,
	}

	execReturn, err := shell.Execute()
	if err != nil {
		log.Fatal(err)
	}

	shell = exec.ExecTask{
		Command: "docker",
		Args:    []string{"ps", "-aqf", "name=^" + req.Container.Name + "$"},
		Shell:   true,
	}

	execReturn, err = shell.Execute()
	execReturn.Stdout = strings.ReplaceAll(execReturn.Stdout, "\n", "")

	w.Write([]byte(execReturn.Stdout))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker delete container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req types.PodUID
	json.Unmarshal(bodyBytes, &req)

	cmd := []string{"stop", req.UID}
	shell := exec.ExecTask{
		Command: "docker",
		Args:    cmd,
		Shell:   true,
	}
	execReturn, err := shell.Execute()

	cmd = []string{"rm", execReturn.Stdout}
	shell = exec.ExecTask{
		Command: "docker",
		Args:    cmd,
		Shell:   true,
	}
	execReturn, err = shell.Execute()
	execReturn.Stdout = strings.ReplaceAll(execReturn.Stdout, "\n", "")

	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(execReturn.Stdout))
}
