FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /subscriber

ENV OTEL_TRACES_ENDPOINT=host.docker.internal:7021

CMD ["/subscriber"]
