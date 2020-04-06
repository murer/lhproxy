#!/bin/bash -xe

cd "$(dirname "$0")"

docker rm -f lhproxy_it_squid || true
docker rm -f lhproxy_it_pipe || true

cd ..
docker build -t lhproxy/it:dev -f it/Dockerfile .
cd -
docker run -d --rm --label lhproxy_dev --name lhproxy_it_squid \
  -p 3128:3128 -h lhproxy_it_squid lhproxy/it:dev /root/server.sh

docker run -it --rm --label lhproxy_dev --name lhproxy_it_pipe \
  --network host -e "HTTP_PROXY=http://localhost:3128/" \
  -e "LHPROXY_SECRET=C3hbthuSzAJjknn8" lhproxy/it:dev \
  ssh -o 'ProxyCommand lhproxy client pipe http http://lhproxy_it_squid:8080/ %h:%p' localhost whoami
echo "SUCCESS"
