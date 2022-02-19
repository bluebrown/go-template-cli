ver=1.0.0

build:
	podman build -t bluebrown/tpl -t bluebrown/tpl:$(ver) .

publish:
	podman push bluebrown/tpl
	podman push bluebrown/tpl:$(ver)

local:
	podman run --rm --volume $(CURDIR):/src --workdir /src golang go build -o tpl .
	sudo mv -f tpl /usr/local/bin/tpl


.PHONY: build publish local
