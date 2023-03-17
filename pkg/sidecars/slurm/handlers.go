package slurm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	types "github.com/cloud-pg/interlink/pkg/common"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	//call to slurm get status

}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	//call to slurm create container

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var req types.CreateRequest
	json.Unmarshal(bodyBytes, &req)

	container := req.Container
	metadata := req.Metadata

	log.Print("create_container")
	commstr1 := []string{"singularity", "exec"}

	envs := prepare_envs(container)

	image := ""
	mounts := prepare_mounts(container)
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
	log.Print("questo c.VolumeMounts")
	log.Print(container.VolumeMounts)
	log.Print("fatto")

	singularity_command := append(commstr1, envs...)
	singularity_command = append(singularity_command, mounts...)
	singularity_command = append(singularity_command, image)
	singularity_command = append(singularity_command, container.Command...)
	singularity_command = append(singularity_command, container.Args...)

	path := produce_slurm_script(c, singularity_command)
	out := slurm_batch_submit(path, c)
	handle_jid(c, out)
	log.Debugln(singularity_command)
	log.Infoln(out)
}

func StopHandler(w http.ResponseWriter, r *http.Request) {
	//call to slurm delete container

}
