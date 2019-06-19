# Run vPlan 2 using Docker

## Images on Docker Hub

[**Here**](https://cloud.docker.com/u/zekro/repository/docker/zekro/vplan2019) you can find pre-built Docker images.

<a href="https://cloud.docker.com/u/zekro/repository/docker/zekro/vplan2019"><img alt="Docker Automated build" src="https://img.shields.io/docker/automated/zekro/vplan2019.svg?color=cyan&logo=docker&logoColor=cyan&style=for-the-badge"></a>

Simply pull them using the following command:
```
# docker pull zekro/vplan2019:latest
```

## Running the Image

To run the image, you need to bind the container port `8080` to a prefered port on your host system. Also, you need to mount the both container pathes `/etc/vplan/config` and `/etc/vplan/certs`:

*I) Using the `docker run` command:*
```
# docker run \
    -p 443:8080 \
    -v /home/vplan/config:/etc/vplan/config \
    -v /home/vplan/certs:/etc/vplan/certs \
    -d \
    zekro/vplan2019:latest
```

*II) Using docker-compose:*
```yml
version: '3'

services:
  # ...
  vplan:
    image: 'zekro/vplan2019:latest'
    ports:
      - '443:8080'
    volumes:
      - '/home/vplan/config:/etc/vplan/config'
      - '/home/vplan/certs:/etc/vplan/certs'
```

## Self-build the Docker Image

First, you need to clone the repository locally:
```
$ git clone https://github.com/zekroTJA/vplan2019 --branch master --depth 5
```

Then, build the image using `docker build`:
```
$ cd vplan2019
# docker build . -t vplan2019
```
