FROM --platform=$BUILDPLATFORM golang:1.22 AS builder
WORKDIR /go/src/github.com/jcluppnow/vpa-manager

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /vpa-manager

FROM scratch
COPY --from=builder /vpa-manager /
USER 10002
ENTRYPOINT ["/vpa-manager"]
