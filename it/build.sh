#!/bin/bash -xe

cd "$(dirname "$0")"

cleanup() {
  docker rm -f lhproxy_it_squid || true
  docker rm -f lhproxy_it_pipe || true
}
trap cleanup EXIT

cleanup

cd ..
docker build -t lhproxy/it:dev -f it/Dockerfile .
cd -
docker run -d --rm --label lhproxy_dev --name lhproxy_it_squid \
  -p 3128:3128 -h lhproxy_it_squid lhproxy/it:dev /root/entrypoint/server.sh

docker run -it --rm --label lhproxy_dev --name lhproxy_it_pipe \
  --network host lhproxy/it:dev /root/entrypoint/pipe.sh

echo "SUCCESS"
