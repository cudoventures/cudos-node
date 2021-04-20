FROM golang:1.15

# install jq to parse json within bash scripts
RUN curl -o /usr/local/bin/jq http://stedolan.github.io/jq/download/linux64/jq && \
  chmod +x /usr/local/bin/jq

WORKDIR /usr/cudos

COPY ./ ./

RUN make

RUN chmod +x ./init-root.sh

CMD ["/bin/sh", "-c", "./init-root.sh&& cudos-noded start"] 
# CMD ["sleep","infinity"]