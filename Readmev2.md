# InterLink
## :information_source: Overview

This project aims to enable a communication between a Kubernetes VitualKubelet and a container manager, like for example Docker.
The project is based on KNoC, for reference, check https://github.com/CARV-ICS-FORTH/knoc
Everything is thought to be modular and it's divided in different layers. These layers are summarized in the following drawing:

![drawing](imgs/InterLink.svg)

## :information_source: Components

- Virtual kubelet:
We have implemented 3 more functions able to communicate with the InterLink layer; these functions are called createRequest, deleteRequest and statusRequest, which calls through a REST API to the InterLink layer. Request uses a POST, deleteRequest uses a DELETE, statusRequest uses a GET.

- InterLink:
This is the layer managing the communication with the plug-ins. We began implementing a Mock module, to return dummy answers, and then moved towards a Docker plugin, using a library to emulate a shell to call the Docker CLI commands to implement containers creation, deletion and status querying. We chose to not use Docker API to extend modularity and porting to other managers, since we can think to use a job workload queue like Slurm.

- Sidecars: 
Basically, that's the name we refer to each plug-in talking with the InterLink layer. Each Sidecar is inependent and separately talks with the InterLink layer.

## :grey_exclamation: Requirements
- Golang >= 1.18.9 (might work with older version, but didn't test)
- A working Kubernetes instance
- An already set up KNoC environment
- Docker for the Docker Sidecar
- Sbatch, Scancel and Squeue (Slurm environment) for the Slurm sidecar

## :fast_forward: Quick Start
Fastest way to start using interlink, is by deploying a VK in Kubernetes using the prebuilt image:
```bash
kubectl create ns vk
kubectl kustomize ./kustomizations
kubectl apply -n vk -k ./kustomizations
```

Build InterLink and Sidecars binaries by simply using make:
```bash
make all
```

Now you have your VK running, specify in the configuration file named InterLinkConfig.yaml, located under ./config, which service (Slurm/Docker for the moment) you want to use 