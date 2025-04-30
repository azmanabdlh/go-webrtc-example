FROM golang:1.24.2-alpine AS builder

WORKDIR /myapp
COPY . .

RUN make build

FROM alpine:latest
WORKDIR /myapp

COPY --from=builder /myapp/build .

EXPOSE 8000
CMD ["./build"]