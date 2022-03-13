FROM golang
WORKDIR /workspace
RUN echo 'nobody:*:65534:65534:nobody:/_nonexistent:/bin/false' >passwd
COPY . ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix static -a -o tpl .

FROM scratch
COPY --from=0 /workspace/passwd /etc/passwd
COPY --from=0 /workspace/tpl /tpl
ENTRYPOINT ["/tpl"]
USER nobody
