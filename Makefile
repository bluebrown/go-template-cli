ver=0.3.0
bin_dir=bin
cmd_dir=./cmd/tpl/
container_cli=docker

ldflags=-X 'main.version=$(ver)' -X 'main.commit=$(shell git rev-parse HEAD)'

binaries: bindir dynmaic static

dynmaic:
	go build -ldflags="$(ldflags)" -a -o $(bin_dir)/tpl-linux-amd64 $(cmd_dir)

static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(ldflags)" -a -o $(bin_dir)/tpl-linux-amd64-static $(cmd_dir)

bindir:
	rm -rf $(bin_dir)
	mkdir -p $(bin_dir)

image:
	$(container_cli) build --file assets/container/Dockerfile --build-arg ldflags="$(ldflags)" -t bluebrown/tpl -t bluebrown/tpl:$(ver) .

push:
	$(container_cli) push bluebrown/tpl
	$(container_cli) push bluebrown/tpl:$(ver)

install: binaries
	sudo mv $(bin_dir)/tpl-linux-amd64 /usr/local/$(bin_dir)/tpl


.PHONY: binaries dynamic static bindir image push install

.PHONY: vet
vet:
	go vet -race $(cmd_dir)

.PHONY: test
test:
	go test -cover $(cmd_dir)
