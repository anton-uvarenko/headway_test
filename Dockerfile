FROM golang:1.22.2-alpine3.19 as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o ./app ./cmd/app/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/app .
CMD ./app
