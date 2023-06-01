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
- Docker
- Sbatch, Scancel and Squeue (Slurm environment) for the Slurm sidecar

## :fast_forward: Quick Start
- Fastest way to start using interlink, is by deploying a VK in Kubernetes using the prebuilt image:
    ```bash
    kubectl create ns vk
    kubectl kustomize ./kustomizations
    kubectl apply -n vk -k ./kustomizations
    ```

- Build InterLink and Sidecars binaries by simply using make:
    ```bash
    make all
    ```
    Output files will be created within the bin folder.

- Now you have your VK running and you have built needed binaries, specify in the configuration file named InterLinkConfig.yaml, located under ./config, which service (Slurm/Docker for the moment) you want to use. You only have to set the SidecarService to either "docker" or "slurm". Check the [InterLink Config File](#wrench-interlink-config-file) section for a detailed explanation of each value in the file.
- Run your InterLink and Sidecar executables. You are now running:
    - A Virtual Kubelet
    - The InterLink service
    - A Sidecar
- Submit a YAML to your K8S cluster to test it. You could try:
    ```bash
    kubectl apply -f examples/interlink_mock/payloads/busyecho_k8s.yaml -n vk
    ```
Note: I will soon update the quick start section to only use docker images / k8s deployments

## :hammer: Building from sources
It is possible you need to perform some adjustments or any modification to the source code and you want to rebuild it. You can both binaries, Docker images and even customize your own Kubernetes deployment. 
### Binaries
Building standalone binaries is way simpler and all you need is a simple
 ```bash
make all
```
You will find all VK, InterLink and Sidecars binaries in the bin folder. Replace all with vk/interlink/sidecars to only build the respective component.

### Docker images
Building Docker Images is still simple, but requires 'a little' more effort.
- First of all, login into your Docker Hub account
    ```bash
    docker login
    ```
- Then you can build and push your new images to your Docker Hub. Remember to specify the correct Dockerfile, according to your needs; here's an example with the Virtual Kubelet image:
    ```bash
    docker build -t *your docker hub username*/vk:latest -f Dockerfile.vk .
    docker push *your docker hub username*/vk:latest
    ```
- After pushing the image, edit the deployment.yaml file, located inside the kustomization sub-folder, to reflect the new image name. Check the [Kustomizing your Virtual Kubelet](#wrench-kustomizing-your-Virtual-Kubelet) section for more informations on how to customize your VK deployment.

### :wrench: Kustomizing your Virtual Kubelet
Since ideally the Virtual Kubelet runs into a Docker Container orchestred by a Kubernetes cluster, it is possible to customize your deployment by editing configuration files within the kustomizations directory:
- kustomization.yaml: here you can specify resource files and generate configMaps
- deployment.yaml: that's the main file you want to edit. Nested into spec -> template -> spec -> containers you can find these fields:
    - name: the container name
    - image: Here you can specify which image to use, if you need another one. 
    - args: These are the arguments passed to the VK binary running inside the container.
    - env: Environment Variables used by kubelet and by the VK itself. Check the ENVS list for a detailed explanation on how to set them.
- knoc-cfg.json: it's the config file for the VK itself. Here you can specify how many resources to allocate for the VK. Note that the name specified here for the VK must match the name given in the others config files.
- InterLinkConfig.yaml: configuration file for the inbound/outbound communication (and not only) to/from the InterLink module. For a detailed explanation of all fields, check the [InterLink Config File](#wrench-interlink-config-file) section.
If you perform any change to the listed files, you will have to
```bash
kubectl apply -n vk -k ./kustomizations
```

### :wrench: InterLink Config file
something something