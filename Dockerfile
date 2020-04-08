FROM scratch AS lhproxy_scratch
WORKDIR /usr/local/bin
ENV PATH /usr/local/bin
COPY ./build/out/linux-amd64/lhproxy/lhproxy .
WORKDIR /
CMD [ "lhproxy", "help" ]

FROM alpine:3.11 AS lhproxy_alpine
RUN mkdir /lhproxy && chmod 777 /lhproxy
WORKDIR /lhproxy
COPY --from=lhproxy_scratch /usr/local/bin/lhproxy /usr/local/bin
CMD [ "lhproxy", "help" ]
