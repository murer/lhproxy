#!/bin/bash -xe

export LHPROXY_SECRET="$(cat /root/secret.txt)"
export HTTP_PROXY=http://lhproxy_it_squid:3128/
ssh -o 'ProxyCommand lhproxy client pipe http http://lhproxy_it_squid:8080/ %h:%p' localhost whoami
