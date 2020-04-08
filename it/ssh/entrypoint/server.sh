#!/bin/sh -xe

export LHPROXY_SECRET="$(cat /root/secret.txt)"

/usr/sbin/sshd -D -e 1> /var/log/sshd.out 2>&1 &
squid -N -d1 1> /var/log/squid.out 2>&1 &
lhproxy server 0.0.0.0:8080 1> /var/log/lhproxy.log 2>&1 &
sleep 1

curl --proxy http://localhost:3128/ -f http://127.0.0.1:8080/version.txt

find /var/log -type f | xargs tail -f
