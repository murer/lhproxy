# LHProxy

**LastHopeProxy** is a encrypted port forward through HTTP

It is nice when you need to bypass proxy that decrypt your HTTPS.

Some companies apply this technique (HTTPS decryption) on their employees.

So, we are going to create a encrypted tunnel over plain HTTP.

For **Linux**, **Windows** and **Mac**

[![Build Status](https://travis-ci.org/murer/lhproxy.svg?branch=master)](https://travis-ci.org/murer/lhproxy)


```
 +-----------------------+   The real connection            
 |        LHProxy        |   from YOU to the SERVER         
 |    Last Hope Proxy    ---------------------------+       
 |                       |   Maybe SSH, IRC, HTTP   |       
 +-----|-----------------+   or whatever            |       
       |                                            |       
       |                          +-----------------|-----+
       |                          |      Unreachable      |
       |                          |         SERVER        |
       |                          |                       |
       |                          +-----------|-----------+
       | The forwarded HTTP                   |             
       | request/response with                |             
       | a bunch of unreadable                              
       | binary data inside                Impossible       
       |                                     Path           
       |                                                    
       |                                      |             
       |     +-----------------------+        |             
       |     | A Very Boring Firwall |        |             
       +------       or Proxy        ---------+             
             |                       |                      
             +-----------|-----------+                      
  POST http://lhproxy/   |                                  
  Encrypted Request Body |                                  
                         |                                  
  200 OK HTTP/1.1        |                                  
  Encrypted Reponse Body |                                  
                         |                                  
             +-----------|-----------+                      
             |          YOU          |                      
             |       the Client      |                      
             |                       |                      
             +-----------------------+                      
```

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
LHPROXY_SECRET=myweaksecret lhproxy client pipe http http://yourserver:8080 google:80
```

Send it

```http
GET / HTTP/1.1
Host: google.com
```

You will get something like:

```http
HTTP/1.1 301 Moved Permanently
Location: http://www.google.com/
Content-Type: text/html; charset=UTF-8

<HTML>...</HTML>
```
