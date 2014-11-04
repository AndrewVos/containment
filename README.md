# containment

Docker container management for large amounts of servers

This is very much a work in progress. I haven't even decided whether it's a good idea.

## Creating an ec2 instance

Best is to follow the instructions at https://docs.docker.com/installation/amazon/

## Installation

go get github.com/AndrewVos/containment

## Example configuration

```yaml
clusters:
  -
    name: containment-test
    hosts:
      -
        address: 54.170.230.47
        user: ubuntu
      -
        address: 54.73.192.243
        user: ubuntu
      -
        address: 54.246.65.24
        user: ubuntu

containers:
  -
    image: luisbebop/docker-sinatra-hello-world
    clusters:
      - containment-test
    ports:
      - 80:5000
```

## Usage

```
$ containment update luisbebop/docker-sinatra-hello-world
[ubuntu@54.246.65.24] Pulling repository luisbebop/docker-sinatra-hello-world
[ubuntu@54.246.65.24] Status: Image is up to date for luisbebop/docker-sinatra-hello-world:latest
[ubuntu@54.73.192.243] Pulling repository luisbebop/docker-sinatra-hello-world
[ubuntu@54.73.192.243] Status: Image is up to date for luisbebop/docker-sinatra-hello-world:latest
[ubuntu@54.170.230.47] Pulling repository luisbebop/docker-sinatra-hello-world
[ubuntu@54.170.230.47] Status: Image is up to date for luisbebop/docker-sinatra-hello-world:latest

$ containment start luisbebop/docker-sinatra-hello-world
[ubuntu@54.170.230.47] b3a0598b7e310c52397b3f8548732633ce76ab42f31cc4df9563219b297ed256
[ubuntu@54.73.192.243] f5ef5baa460f8e5d29ba1689aa4bbdd927c96316faae9a69c12b398926078e84
[ubuntu@54.246.65.24] ac2371e344cda0a5b66b7d587281001bd3ca1917004967d3aa11b63cfb8002e0

$ containment status luisbebop/docker-sinatra-hello-world
[ubuntu@54.246.65.24] luisbebop/docker-sinatra-hello-world running
[ubuntu@54.73.192.243] luisbebop/docker-sinatra-hello-world running
[ubuntu@54.170.230.47] luisbebop/docker-sinatra-hello-world running

$ containment stop luisbebop/docker-sinatra-hello-world
[ubuntu@54.170.230.47] luisbebop-docker-sinatra-hello-world
[ubuntu@54.170.230.47] luisbebop-docker-sinatra-hello-world
[ubuntu@54.73.192.243] luisbebop-docker-sinatra-hello-world
[ubuntu@54.73.192.243] luisbebop-docker-sinatra-hello-world
[ubuntu@54.246.65.24] luisbebop-docker-sinatra-hello-world
[ubuntu@54.246.65.24] luisbebop-docker-sinatra-hello-world
```
