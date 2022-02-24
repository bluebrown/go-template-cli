ver=1.0.0
bindir=bin

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

build:
	podman build -t bluebrown/tpl -t bluebrown/tpl:$(ver) .

publish:
	podman push bluebrown/tpl
	podman push bluebrown/tpl:$(ver)


.PHONY: binaries dynamic static build publish
