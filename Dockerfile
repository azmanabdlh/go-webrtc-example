FROM golang:1.24.2-alpine AS builder

WORKDIR /myapp
COPY . .

RUN go build -o build *.go

FROM alpine:latest
WORKDIR /myapp

COPY --from=builder /myapp/build .
COPY --from=builder /myapp/public public

EXPOSE 8000
CMD ["./build"]