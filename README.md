# Overview
This project aims to enable a communication between a Kubernetes VitualKubelet and a container manager, like for example Docker.
The project is based on KNoC, for reference, check https://github.com/CARV-ICS-FORTH/knoc
Everything is thought to be modular and it's divided in different layers. These layers are summarized in the following drawing:
https://excalidraw.com/#room=27d4efeab01b19377b97,i-2RMSKaE6JMTT0QkeEfEQ

# Components
- Kubernetes API server: 
That's your local K8S instance running on your server. K8S talks to the next layer through its own API

- Virtual kubelet (Knoc):
Knoc is a Virtual Kubeled provider able to allow Pods' registration to the K8S cluster. We have implemented 3 more functions able to communicate with the InterLink layer; these functions are called createRequest, deleteRequest and statusRequest, which calls through a REST API to the InterLink layer. CreateRequest uses a POST, deleteRequest uses a DELETE, statusRequest uses a GET.

- InterLink:
This is the layer managing the communication with the plug-ins. We began implementing a Mock module, to return dummy answers, and then moved towards a Docker plugin, using a library to emulate a shell to call the Docker CLI commands to implement containers creation, deletion and status querying. We chose to not use Docker API to extend modularity and porting to other managers, since we can think to use a job workload queue like Slurm.

- Sidecars
Basically, that's the name we refer to each plug-in talking with the InterLink layer. Each Sidecar is inependent and separately talks with the InterLink layer.

# Build and Usage
Requirements: 
- Golang >= 1.18.9 (might work with older version, but didn't test)
- A working Kubernetes instance
- An already set up KNoC environment

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

A quick start-up command for the VK executable is given by the following example:
```
./bin/vk -- --nodename *A-NAME* --provider knoc --provider-config ./scripts/cfg.json --startup-timeout 10s --klog.v "2" --kubeconfig *PATH-TO-KUBECONFIG.YAML* --klog.logtostderr --log-level debug
```

# Debug
To debug, we found out Delve debugger is pretty handful. To run a debug session, install delve debugger by running;
```
go install github.com/go-delve/delve/cmd/dlv@latest
```

Then, assuming $COMMAND is the normal string you would use to run your executable, run the following:
```
dlv debug $COMMAND
```

For example, based on the previous example:
```
dlv debug . -- --nodename *A-NAME* --provider knoc --provider-config ./scripts/cfg.json --startup-timeout 10s --klog.v "2" --kubeconfig *PATH-TO-KUBECONFIG.YAML* --klog.logtostderr --log-level debug
```

The only difference is you have to pass the path to your main and not the path to your executable