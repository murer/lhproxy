FROM scratch
WORKDIR /lhproxy
ENV PATH /lhproxy
COPY ./build/out/linux-amd64/lhproxy/lhproxy .
CMD [ "lhproxy", "help" ]
