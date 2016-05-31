FROM phusion/baseimage

COPY ./bin/artemisd /artemisd
COPY ./artemis.docker.conf /artemis.docker.conf

CMD ["./artemisd", "--config=./artemis.docker.conf"]

EXPOSE 80