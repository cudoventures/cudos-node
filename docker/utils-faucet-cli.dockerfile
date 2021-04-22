FROM golang:alpine

RUN apk add --no-cache make

WORKDIR /usr/faucet-cli

COPY ./project-faucet-cli ./

RUN make

RUN chmod +x run-docker.sh

COPY --from=cudos-root-node /go/bin/cudos-noded /go/bin/cudos-noded

CMD ["sh", "run-docker.sh"]