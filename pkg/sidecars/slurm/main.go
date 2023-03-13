package slurm

import (
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker get status

}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker create container

}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	//call to docker delete container

}
