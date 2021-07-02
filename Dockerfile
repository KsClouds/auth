FROM alpine

COPY ./auth /tmp/auth

WORKDIR /tmp/

RUN chmod +x auth