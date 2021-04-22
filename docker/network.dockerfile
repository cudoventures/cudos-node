FROM golang:alpine as cudos-network-builder

RUN apk add --no-cache make bash

WORKDIR /usr/cudos-builder

COPY ./project-node ./

RUN make

ARG MONIKER

COPY --from=cudos-root-node /usr/cudos/cudos-data /usr/cudos-root-data

RUN chmod +x ./init-network.sh

RUN sed -i 's/\r$//' ./init-network.sh

RUN /bin/bash ./init-network.sh

FROM golang:alpine

WORKDIR /usr/cudos

RUN apk add --no-cache bash

COPY --from=cudos-network-builder /go/bin/cudos-noded /go/bin/cudos-noded

COPY --from=cudos-network-builder /usr/cudos-builder/cudos-data /usr/cudos/cudos-data

CMD ["/bin/sh", "-c", "cudos-noded start"] 