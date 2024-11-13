FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY .. .
RUN go build -o app ./src/cmd/calculator/calculator.go

FROM alpine

WORKDIR /app
COPY --from=builder /build/app /app/app
COPY ./src/config ./config

CMD ["./app", "--config=./config/config.yml"]