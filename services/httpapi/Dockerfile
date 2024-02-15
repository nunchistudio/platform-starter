FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /httpapi
EXPOSE 8080

ENV OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=host.docker.internal:7021

CMD ["/httpapi"]
