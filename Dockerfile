FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /x402fuel .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /x402fuel /usr/local/bin/x402fuel
EXPOSE 8420
ENTRYPOINT ["x402fuel"]
CMD ["serve"]