FROM alpine

RUN apk update

RUN mkdir -p /run/docker/plugins /mnt/state

#COPY bin/openebs /bin/openebs

CMD ["bin/sh"]
