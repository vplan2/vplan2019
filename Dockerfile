FROM golang:1.12.6-stretch

LABEL maintainer="contact@zekro.de"

### ARGS

ARG ZOLAVERSION="0.7.0"

### PREPARATION

RUN apt-get update -y &&\
    apt-get install -y \
        git \
        dos2unix

ENV PATH="${GOPATH}/bin:${PATH}"

RUN go get -u github.com/golang/dep/cmd/dep

RUN wget https://github.com/getzola/zola/releases/download/v${ZOLAVERSION}/zola-v${ZOLAVERSION}-x86_64-unknown-linux-gnu.tar.gz \
        -O ./zola.tar.gz &&\
    tar -xzvf zola.tar.gz &&\
    chmod +x ./zola &&\
    mv ./zola /usr/bin/zola

WORKDIR ${GOPATH}/src/github.com/zekroTJA/vplan2019

ADD . .

RUN mkdir -p /etc/vplan/certs &&\
    mkdir -p /etc/vplan/config

RUN dos2unix ./scripts/*.sh &&\
    chmod +x ./scripts/*.sh

### COMPILE BACK END

RUN dep ensure -v

ENV LDFLAGS="github.com/zekroTJA/vplan2019/internal/ldflags"

RUN go build -v -o /var/vplan/server \
        -ldflags "\
            -X ${LDFLAGS}.AppVersion=$(git describe --tags) \
            -X ${LDFLAGS}.AppCommit=$(git rev-parse HEAD) \
            -X ${LDFLAGS}.GoVersion=$(go version | sed -e 's/ /_/g') \
            -X ${LDFLAGS}.Release=TRUE" \
        ./cmd/server

### COMPILE FRONT END

RUN cp ./config/frontend.release.toml \
        ./web/config.toml

RUN cd web &&\
        zola build

RUN mv ./web/public /var/vplan/web

### EXPOSE

EXPOSE 8080

CMD ./scripts/docker-starter.sh \
        "/var/vplan/server -c /etc/vplan/config/config.yml -web /etc/vplan/web" \
        /etc/vplan/config/config.yml \
        ./config/docker.config.yml
