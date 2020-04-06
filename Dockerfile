FROM scratch
WORKDIR /lhproxy
ENV HOME /lhproxy
ADD ./build/out/linux-amd64/lhproxy/lhproxy .
CMD [ "./lhproxy", "help" ]
