ver=1.0.0
bindir=bin
container_cli=docker

binaries: bindir dynmaic static

dynmaic:
	go build -a -o tpl-$(ver)-amd64  .
	tar -czf tpl-amd64.tar.gz tpl-$(ver)-amd64
	mv tpl-amd64.tar.gz bin/
	rm tpl-$(ver)-amd64

static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o tpl-$(ver)-amd64-static .
	tar -czf tpl-amd64-static.tar.gz tpl-$(ver)-amd64-static
	mv tpl-amd64-static.tar.gz bin/
	rm tpl-$(ver)-amd64-static

bindir:
	rm -rf $(bindir)
	mkdir -p $(bindir)

image:
	$(container_cli) build -t bluebrown/tpl -t bluebrown/tpl:$(ver) .

push:
	$(container_cli) push bluebrown/tpl
	$(container_cli) push bluebrown/tpl:$(ver)


.PHONY: binaries dynamic static bindir image push
