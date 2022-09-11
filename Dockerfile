FROM golang:1.19 AS builder

WORKDIR /src
COPY go.mod go.sum installer.sh tools.go ./
RUN go mod download && ./installer.sh
COPY . .
RUN buf generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /app -a -ldflags '-w -s' ./collector

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /benchmark-japan-vps-collector
ENTRYPOINT ["/benchmark-japan-vps-collector"]