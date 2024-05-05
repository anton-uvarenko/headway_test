FROM golang:1.22.2-alpine3.19 as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o ./app ./main.go

FROM alpine
COPY --from=builder /app/app /bin/nasa
CMD tail -f /dev/null
