FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./devops ./cmd/api


FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/devops .
CMD ["./devops"]

EXPOSE 3000
