FROM golang:1.15

WORKDIR /usr/cudos

COPY ./ ./

RUN make

CMD ["/bin/sh", "-c", "cudos-noded start"] 