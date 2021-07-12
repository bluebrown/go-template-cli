FROM golang as builder
WORKDIR /src
COPY . .
RUN go build .

FROM alpine
COPY --from=builder /src/jpipe /run/app/jpipe
ENTRYPOINT [ "/run/app/jpipe" ]
ARG UID=8080
ARG USER="jpipe"
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home /run/app/jpipe \
    --no-create-home \
    --uid "$UID" \
    "$USER"
USER $USER
