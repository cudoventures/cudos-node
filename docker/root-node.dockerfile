FROM golang:alpine as cudos-root-node-builder

RUN apk add --no-cache jq make bash

WORKDIR /usr/cudos-builder

COPY ./project-node ./

RUN make

RUN chmod +x ./init-root.sh

RUN sed -i 's/\r$//' ./init-root.sh

RUN /bin/bash ./init-root.sh

FROM golang:alpine

WORKDIR /usr/cudos

RUN apk add --no-cache bash

COPY --from=cudos-root-node-builder /go/bin/cudos-noded /go/bin/cudos-noded

COPY --from=cudos-root-node-builder /usr/cudos-builder/cudos-data /usr/cudos/cudos-data

CMD ["/bin/bash", "-c", "cudos-noded start"] 