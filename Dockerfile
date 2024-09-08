FROM golang:1.23-alpine3.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o apibot ./cmd

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/apibot /app/apibot
COPY --from=builder /app/cred /app/cred
COPY --from=builder /app/config /app/config
ENTRYPOINT ["/app/apibot"]