FROM golang:1.15

WORKDIR /usr/cudos

COPY ./ ./

RUN make

CMD ["/bin/sh", "-c", "./init-network.sh&& cudos-noded start"] 