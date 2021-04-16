FROM node:buster

# RUN apk add --no-cache curl bash

RUN curl https://install.meteor.com/ | sh

WORKDIR /usr/explorer

COPY ./third_party/big-dipper ./

RUN meteor npm install --save

CMD ["sh", "-c", "meteor --settings default_settings.json --allow-superuser"] 
# CMD ["sleep", "infinity"]