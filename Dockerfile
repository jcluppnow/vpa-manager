FROM golang:1.22 AS builder
WORKDIR /go/src/github.com/jcluppnow/vpa-manager

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /vpa-manager

FROM scratch
COPY --from=builder /vpa-manager /
USER 10002
ENTRYPOINT ["/vpa-manager"]
