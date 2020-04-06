FROM alpine:3.11
RUN mkdir /lhproxy && chmod 777 /lhproxy
ENV HOME /lhproxy
ADD ./build/out/linux-amd64/lhproxy/lhproxy /usr/local/bin
CMD [ "lhproxy", "help" ]
