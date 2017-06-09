FROM alpine

RUN apk update

RUN mkdir -p /run/docker/plugins /mnt/state

COPY bin/openebsdriver /bin/openebsdriver

CMD ["bin/sh"]
