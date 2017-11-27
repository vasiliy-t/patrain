FROM patrain/base:latest as builder
WORKDIR /go/src/github.com/vasiliy-t/patrain
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o producer/producer github.com/vasiliy-t/patrain/producer

FROM patrain/base:latest
ENV CONSUL_HTTP_ADDR=http://172.19.0.12:8500

WORKDIR /app
COPY --from=builder /go/src/github.com/vasiliy-t/patrain/producer/producer ./
CMD ["./producer"]
