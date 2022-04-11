ver=0.1.0
bindir=bin
container_cli=docker
cmd_dir=./cmd/tpl/

ldflags=-X 'main.version=$(ver)' -X 'main.commit=$(shell git rev-parse HEAD)'

binaries: bindir dynmaic static

dynmaic:
	go build -ldflags="$(ldflags)" -a -o tpl-$(ver)-amd64 $(cmd_dir)
	tar -czf tpl-amd64.tar.gz tpl-$(ver)-amd64
	mv tpl-amd64.tar.gz bin/
	rm tpl-$(ver)-amd64

static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(ldflags)" -a -o tpl-$(ver)-amd64-static $(cmd_dir)
	tar -czf tpl-amd64-static.tar.gz tpl-$(ver)-amd64-static
	mv tpl-amd64-static.tar.gz bin/
	rm tpl-$(ver)-amd64-static

bindir:
	rm -rf $(bindir)
	mkdir -p $(bindir)

image:
	$(container_cli) build --file assets/container/Dockerfile --build-arg ldflags="$(ldflags)" -t bluebrown/tpl -t bluebrown/tpl:$(ver) .

push:
	$(container_cli) push bluebrown/tpl
	$(container_cli) push bluebrown/tpl:$(ver)

install: binaries
	tar -xzf bin/tpl-amd64.tar.gz
	sudo mv tpl-$(ver)-amd64 /usr/local/bin/tpl


.PHONY: binaries dynamic static bindir image push install
