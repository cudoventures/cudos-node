FROM golang:1.15

WORKDIR /usr/cudos

COPY ./ ./

RUN make

RUN chmod +x ./init-root.sh

CMD ["/bin/sh", "-c", "./init-network.sh&& cudos-noded start"] 