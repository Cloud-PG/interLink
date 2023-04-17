all: interlink vk sidecars

interlink:
	go build -o bin/interlink cmd/interlink/main.go

vk:
	go build -o bin/vk

sidecars:
	go build -o bin/docker-sd cmd/sidecars/docker/main.go
	go build -o bin/slurm-sd cmd/sidecars/slurm/main.go

clean:
	rm -rf ./bin