FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /subscriber

ENV OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=host.docker.internal:7021

CMD ["/subscriber"]
