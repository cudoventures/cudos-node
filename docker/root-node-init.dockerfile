FROM golang:alpine

# install jq to parse json within bash scripts
# RUN curl -o /usr/local/bin/jq http://stedolan.github.io/jq/download/linux64/jq && \
#   chmod +x /usr/local/bin/jq
RUN apk add --no-cache jq make bash

WORKDIR /usr/cudos

COPY ./project-node ./

RUN make

RUN chmod +x ./init-root.sh

CMD ["/bin/bash", "-c", "./init-root.sh&& cudos-noded start"] 
# CMD ["sleep","infinity"]