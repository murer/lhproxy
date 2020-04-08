#!/bin/sh -xe

export LHPROXY_SECRET="$(cat /root/secret.txt)"
ssh -o 'ProxyCommand lhproxy client pipe http -p http://lhproxy_it_squid:3128/ http://127.0.0.1:8080/ %h:%p' localhost whoami

export HTTP_PROXY=http://lhproxy_it_squid:3128/
ssh -o 'ProxyCommand lhproxy client pipe http http://lhproxy_it_squid:8080/ %h:%p' localhost whoami
