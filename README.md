# overview
This project aims to enable a communication between a Kubernetes VitualKubelet and a container manager, like for example Docker.
The project is based on KNOC, for reference, check 
Everything is thought to be modular and it's divided in different layers. These layers are summarized in the following drawing:
https://excalidraw.com/#room=27d4efeab01b19377b97,i-2RMSKaE6JMTT0QkeEfEQ

# components
- Kubernetes API server: 
That's your local K8S instance running on your server. K8S talks to the next layer through its own API

- Virtual kubelet (Knoc):
Knoc is a Virtual Kubeled provider able to allow Pods' registration to the K8S cluster. We have implemented 3 more functions able to communicate with the InterLink layer; these functions are called createRequest, deleteRequest and statusRequest, which calls through a REST API to the InterLink layer. CreateRequest uses a POST, deleteRequest uses a DELETE, statusRequest uses a GET.

- InterLink:
This is the layer managing the communication with the plug-ins. We began implementing a Mock module, to return dummy answers, and then moved towards a Docker plugin, using a library to emulate a shell to call the Docker CLI commands to implement containers creation, deletion and status querying. We chose to not use Docker API to extend modularity and porting to other managers, since we can think to use a job workload queue like Slurm.

- Sidecars
Basically, that's the name we refer to each plug-in talking with the InterLink layer. Each Sidecar is inependent and separately talks with the InterLink layer.

# usage
Requirements: 
- Golang >= 1.18.9 (might work with older version, but didn't test)
- A working Kubernetes instance

build the 3 components by running:
```
go build -o bin/vk
go build -o bin/interlink cmd/interlink/main.go
go build -o bin/docker-sd cmd/sidecars/docker/main.go
```
Three output files called vk, interlink and docker-sd will be created within the bin folder
Give exec permissions and run all of them, then test by submitting a YAML to your K8S cluster. For example, you can run
```
kubectl apply -f examples/busyecho_k8s.yaml
```

# interlink

export KUBERNETES_SERVICE_HOST=127.0.0.1
export KUBERNETES_SERVICE_PORT=2345 
