FROM node:buster as builder

RUN apt update

RUN apt install build-essential libudev-dev zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev wget python3 -y

RUN curl https://install.meteor.com/ | sh

WORKDIR /usr/explorer

COPY ./project-explorer ./source

RUN cd ./source && npm i

RUN cd ./source && meteor build ../output/ --architecture os.linux.x86_64 --server-only --allow-superuser

RUN cd ./output && tar -zxvf ./source.tar.gz

RUN cd ./output/bundle/programs/server&& npm i

FROM node:buster

WORKDIR /usr/explorer

COPY --from=builder /usr/explorer/output/bundle ./

COPY --from=builder /usr/explorer/source/run-docker.sh ./

RUN chmod +x ./run-docker.sh

# CMD ["sh", "-c", "meteor --settings default_settings.json --allow-superuser"] 
CMD ["sh", "./run-docker.sh"]