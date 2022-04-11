FROM golang:1.18
WORKDIR /workspace
RUN echo 'nobody:*:65534:65534:nobody:/_nonexistent:/bin/false' >passwd
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ARG ldflags
RUN go build \
    -installsuffix static \
    -ldflags="${ldflags}" \
    -a ./cmd/tpl

FROM scratch
COPY --from=0 /workspace/passwd /etc/passwd
COPY --from=0 /workspace/tpl /tpl
ENTRYPOINT ["/tpl"]
USER nobody
