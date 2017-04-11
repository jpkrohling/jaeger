FROM alpine:3.5
MAINTAINER Juraci Paixão Kröhling <juraci.docker@kroehling.de>

ADD cmd/standalone/standalone-linux /jaeger
ENTRYPOINT ["/jaeger", "-span-storage.type", "memory", "-query.port", "3001"]
