FROM golang as builder
WORKDIR /workspace
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o tpl .

FROM gcr.io/distroless/static:nonroot as final
WORKDIR /
COPY --from=builder /workspace/tpl .
USER 65532:65532
ENTRYPOINT ["/tpl"]
