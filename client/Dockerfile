FROM node:8.3

RUN useradd --user-group --create-home --shell /bin/false app &&\
  npm install --global yarn@latest

ENV HOME=/home/app

COPY package.json $HOME/playground/
RUN chown -R app:app $HOME/*

USER app
WORKDIR $HOME/playground
RUN yarn install
