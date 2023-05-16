# Overview
This project aims to enable a communication between a Kubernetes VitualKubelet and a container manager, like for example Docker.
The project is based on KNoC, for reference, check https://github.com/CARV-ICS-FORTH/knoc
Everything is thought to be modular and it's divided in different layers. These layers are summarized in the following drawing:

![drawing](imgs/interLink%20schematic.svg "InterLink schematic")

# Components
- Virtual kubelet:
We have implemented 3 more functions able to communicate with the InterLink layer; these functions are called createRequest, deleteRequest and statusRequest, which calls through a REST API to the InterLink layer. Request uses a POST, deleteRequest uses a DELETE, statusRequest uses a GET.

- InterLink:
This is the layer managing the communication with the plug-ins. We began implementing a Mock module, to return dummy answers, and then moved towards a Docker plugin, using a library to emulate a shell to call the Docker CLI commands to implement containers creation, deletion and status querying. We chose to not use Docker API to extend modularity and porting to other managers, since we can think to use a job workload queue like Slurm.

- Sidecars
Basically, that's the name we refer to each plug-in talking with the InterLink layer. Each Sidecar is inependent and separately talks with the InterLink layer.

# End2end example

Please refer to this [repository](https://github.com/Cloud-PG/interLink/blob/main/README.md)

# Build and Usage
Requirements: 
- Golang >= 1.18.9 (might work with older version, but didn't test)
- A working Kubernetes instance
- An already set up KNoC environment
- Docker for the Docker Sidecar
- Sbatch, Scancel and Squeue (Slurm environment) for the Slurm sidecar

Build the components by running:
```
make all
```
Output files will be created within the bin folder.

Remember to correctly set-up Environment Variables (or the InterLinkConfig.yaml file. ENVS have priority over config file) according to the service you want to use!

```
List of Environment Variables:
$INTERLINKURL -> the URL to contact the InterLink executable. No need to specify a port here
$INTERLINKPORT -> the InterLink listening port. Default is 3000
$INTERLINKCONFIGPATH -> your config file path
$SIDECARURL -> the URL to allow InterLink to communicate with the Sidecar module (docker, slurm, etc). No need to specify port here
$SIDECARPORT -> the Sidecar listening port. Docker default is 4000, Slurm default is 4001
$SIDECARSERVICE -> can be "docker" or "slurm" only (for the moment). If SIDECARPORT is not set, will set Sidecar Port in the code to default settings.
$TSOCKS -> true or false, to use tsocks library allowing proxy networking. Working on Slurm sidecar at the moment. Remember to properly configure your TSOCKS instance 
$TSOCKSPATH -> path to your tsocks library
```

ENVS and config naming matches, so you will just find the config names to be the lowercases of the ENVS naming

Give exec permissions and run all of them, then test by submitting a YAML to your K8S cluster. For example, you can run
```
kubectl apply -f examples/busyecho_k8s.yaml
```

A quick start-up command for the VK executable is given by the following example:
```
./bin/vk -- --nodename *A-NAME* --provider knoc --provider-config ./scripts/cfg.json --startup-timeout 10s --klog.v "2" --kubeconfig *PATH-TO-KUBECONFIG.YAML* --klog.logtostderr --log-level debug
```

For the other executables, you can just normally run them.

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
dlv debug . -- --nodename *A-NAME* --provider knoc --provider-config ./config/cfg.json --startup-timeout 10s --klog.v "2" --kubeconfig *PATH-TO-KUBECONFIG.YAML* --klog.logtostderr --log-level debug
```

The only difference is you have to pass the path to your main and not the path to your executable
