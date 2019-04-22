# Simple reverse proxy

This is a learning project to understand how things are handled on the background when sending a request to server.

[![Build Status](https://travis-ci.com/Broke116/go-reverse-proxy.svg?branch=master)](https://travis-ci.com/Broke116/go-reverse-proxy)

Learning project to understand how reverse proxy works. Later then it will be used for building a load balancer.
[Inspired by this article](https://hackernoon.com/writing-a-reverse-proxy-in-just-one-line-with-go-c1edfa78c84b)

Mini load balancer application in Go was demonstrated by Rob Pike. [Slide](https://talks.golang.org/2012/waza.slide#1). So that I inspired the implementation from his slides. 
After covering this awesome slide, I wanted to create a reverse proxy/load balancer server to learn and understand the concepts of
concurrency, worker pools and networking.
On top of that I mixed load balancer with reverse proxy and worker pool implementation.

## How to build
```
# Build the project
$ go build -o reverseproxy.exe

# For development purposes opening a multiple http servers using http-server npm package
$ http-server -p 1331
$ http-server -p 1332
```