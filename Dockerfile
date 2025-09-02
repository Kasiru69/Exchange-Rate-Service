# Use official Go image with Go 1.24 alpine
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 1001 appgroup && adduser -u 1001 -G appgroup -s /sbin/nologin -D appuser

WORKDIR /home/appuser

COPY --from=builder /app/main .

#COPY --from=builder /app/.env.example .

RUN chown -R appuser:appgroup /home/appuser

USER appuser

EXPOSE 8080

CMD ["./main"]
