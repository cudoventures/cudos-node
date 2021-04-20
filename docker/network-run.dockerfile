FROM golang:alpine

RUN apk add --no-cache make bash

WORKDIR /usr/cudos

COPY ./project-node ./

RUN make

CMD ["/bin/sh", "-c", "cudos-noded start"] 