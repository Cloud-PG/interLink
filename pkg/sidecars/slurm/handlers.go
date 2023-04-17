package slurm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	exec "github.com/alexellis/go-execute/pkg/v1"
	types "github.com/cloud-pg/interlink/pkg/common"
	v1 "k8s.io/api/core/v1"
)

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	//call to slurm create container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	var req types.CreateRequest
	json.Unmarshal(bodyBytes, &req)
	if err != nil {
		log.Print(err)
		return
	}

	container := req.Container
	metadata := req.Pod.ObjectMeta

	log.Print("create_container")
	commstr1 := []string{"singularity", "exec"}

	envs := prepare_envs(container)

	image := ""
	//mounts := prepare_mounts(container)
	if strings.HasPrefix(container.Image, "/") {
		if image_uri, ok := metadata.Annotations["slurm-job.knoc.io/image-root"]; ok {
			log.Print(image_uri)
			image = image_uri + container.Image
		} else {
			log.Print("image-uri annotation not specified for path in remote filesystem")
		}
	} else {
		image = "docker://" + container.Image
	}
	image = container.Image

	singularity_command := append(commstr1, envs...)
	//singularity_command = append(singularity_command, mounts...)
	singularity_command = append(singularity_command, image)
	singularity_command = append(singularity_command, container.Command...)
	singularity_command = append(singularity_command, container.Args...)

	log.Println("Generating Slurm script")
	path := produce_slurm_script(container, metadata, singularity_command)
	log.Println("Submitting Slurm job")
	out := slurm_batch_submit(path)
	handle_jid(container, out)
	log.Print(out)
	log.Println(path)

	w.Write([]byte(path))
}

func StopHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	var req v1.Container
	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		log.Print(err)
		return
	}

	delete_container(req)

	log.Print("delete slurm job")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	cmd := []string{"--me"}
	shell := exec.ExecTask{
		Command: "squeue",
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
