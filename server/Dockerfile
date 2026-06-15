FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/whiteboard-server ./cmd/whiteboard-server

FROM alpine:3.20

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /out/whiteboard-server /app/whiteboard-server
COPY configs /app/configs

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/whiteboard-server", "--config", "/app/configs/config.yaml"]
