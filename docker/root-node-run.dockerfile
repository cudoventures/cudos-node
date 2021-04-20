FROM golang:alpine

RUN apk add --no-cache make bash

WORKDIR /usr/cudos

COPY ./project-node ./

RUN make

CMD ["/bin/bash", "-c", "cudos-noded start"] 