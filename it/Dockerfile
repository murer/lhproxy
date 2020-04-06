FROM alpine:3.11 AS base

RUN apk add --update openssh squid bash vim

ENV HOME /root
WORKDIR /root

RUN ssh-keygen -A
RUN mkdir .ssh
RUN ssh-keygen -f .ssh/id_rsa -N ''
RUN cp .ssh/id_rsa.pub .ssh/authorized_keys
RUN cat /etc/ssh/*.pub | sed 's/^/localhost /g' > .ssh/known_hosts
RUN echo "Port 22" > /etc/ssh/sshd_config
RUN echo "ListenAddress 127.0.0.1" >> /etc/ssh/sshd_config
RUN echo "AuthorizedKeysFile ~/.ssh/authorized_keys" >> /etc/ssh/sshd_config
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
RUN sed -i 's/^root:!:/root:*:/g' /etc/shadow

RUN echo "visible_hostname lavaburst" > /etc/squid/squid.conf
RUN echo "http_access allow all" >> /etc/squid/squid.conf
RUN echo "http_port 3128" >> /etc/squid/squid.conf
RUN echo "coredump_dir /var/cache/squid" >> /etc/squid/squid.conf

RUN dd if=/dev/random count=1 bs=12 | base64 > /root/secret.txt

COPY it/entrypoint /root/entrypoint
RUN chmod +x /root/entrypoint/*.sh

FROM base AS final
COPY build/out/linux-amd64/lhproxy/lhproxy /usr/local/bin
