# LHProxy

**Last Hope Proxy** does a ``TCP`` **encrypted** tunnel over ``HTTP`` (not ``HTTPS`` or ``CONNECT``).

For **Linux**, **Windows** and **Mac**

[![LHProxy](https://travis-ci.org/murer/lhproxy.svg)](https://circleci.com/gh/murer/lhproxy)

### Download

Download from <a href="https://github.com/murer/lhproxy/releases">Github Releases</a>.

### Docker

```shell
docker run -it murer/lhproxy:latest lhproxy help
```

### Basics

Start the server somewhere

```shell
LHPROXY_SECRET=myweaksecret lhproxy server 0.0.0.0:8080
```

Start your tunnel from the client

```shell
LHPROXY_SECRET=myweaksecret lhproxy client pipe http http://yourserver:8080 localhost:22
```
