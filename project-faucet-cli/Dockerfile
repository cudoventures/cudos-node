FROM golang:alpine AS builder
WORKDIR /src/app/
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o=/usr/local/bin/faucet ./cmd/faucet

FROM alpine
COPY --from=builder /usr/local/bin/faucet /usr/local/bin/faucet
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
CMD [ "faucet" ]