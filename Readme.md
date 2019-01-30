# Simple reverse proxy

[![Build Status](https://travis-ci.com/Broke116/go-reverse-proxy.svg?branch=master)](https://travis-ci.com/Broke116/go-reverse-proxy)

Learning project to understand how reverse proxy works. Later then it will be used for building a load balancer.
[Inspired by this article](https://hackernoon.com/writing-a-reverse-proxy-in-just-one-line-with-go-c1edfa78c84b)

## How to build
```
# Build the project
$ go build -o reverseproxy.exe

# For development purposes opening a multiple http servers using http-server npm package
$ http-server -p 1331
$ http-server -p 1332
$ http-server -p 1333
```