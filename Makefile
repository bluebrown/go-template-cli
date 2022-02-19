ver=1.0.0

build:
	podman build -t docker.io/bluebrown/tpl -t docker.io/bluebrown/tpl:$(ver) .

publish:
	podman push docker.io/bluebrown/tpl
	podman push docker.io/bluebrown/tpl:$(ver)

local:
	podman run --rm --volume $(CURDIR):/src --workdir /src golang go build -o tpl .
	sudo mv -f tpl /usr/local/bin/tpl


.PHONY: build publish local
