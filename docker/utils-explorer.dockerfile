FROM node:buster as builder

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

# CMD ["sh", "-c", "meteor --settings default_settings.json --allow-superuser"] 
CMD ["node", "main.js"]