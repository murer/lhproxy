FROM alpine:3.11 AS base

RUN apk add --update openssh squid curl

ENV HOME /root
WORKDIR /root

RUN ssh-keygen -A && \
  mkdir .ssh && \
  ssh-keygen -f .ssh/id_rsa -N '' && \
  cp .ssh/id_rsa.pub .ssh/authorized_keys && \
  cat /etc/ssh/*.pub | sed 's/^/localhost /g' > .ssh/known_hosts && \
  echo "Port 22" > /etc/ssh/sshd_config && \
  echo "ListenAddress 127.0.0.1" >> /etc/ssh/sshd_config && \
  echo "AuthorizedKeysFile ~/.ssh/authorized_keys" >> /etc/ssh/sshd_config && \
  echo "PermitRootLogin yes" >> /etc/ssh/sshd_config && \
  sed -i 's/^root:!:/root:*:/g' /etc/shadow && \
  echo "visible_hostname lavaburst" > /etc/squid/squid.conf && \
  echo "http_access allow all" >> /etc/squid/squid.conf && \
  echo "http_port 3128" >> /etc/squid/squid.conf && \
  echo "coredump_dir /var/cache/squid" >> /etc/squid/squid.conf && \
  echo "	hosts_file /etc/hosts" >> /etc/squid/squid.conf && \
  dd if=/dev/random count=1 bs=12 | base64 > /root/secret.txt

COPY entrypoint /root/entrypoint
RUN chmod +x /root/entrypoint/*.sh

FROM base AS final
COPY --from=lhproxy/lhproxy:dev /lhproxy/lhproxy /usr/local/bin
#COPY build/out/linux-amd64/lhproxy/lhproxy /usr/local/bin
