FROM golang:1.12.6-stretch

LABEL maintainer="contact@zekro.de"

### ARGS

### PREPARATION

RUN apt-get update -y &&\
    apt-get install -y \
        git \
        dos2unix

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash - &&\
        apt-get install -y nodejs &&\
        npm install -g @angular/cli

ENV PATH="${GOPATH}/bin:${PATH}"

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${GOPATH}/src/github.com/vplan2/vplan2019

ADD . .

RUN mkdir -p /etc/vplan/certs &&\
    mkdir -p /etc/vplan/config

RUN dos2unix ./scripts/*.sh &&\
    chmod +x ./scripts/*.sh

### COMPILE BACK END

RUN dep ensure -v

ENV LDFLAGS="github.com/vplan2/vplan2019/internal/ldflags"

RUN go build -v -o /var/vplan/server \
        -ldflags "\
            -X ${LDFLAGS}.AppVersion=$(git describe --tags) \
            -X ${LDFLAGS}.AppCommit=$(git rev-parse HEAD) \
            -X ${LDFLAGS}.GoVersion=$(go version | sed -e 's/ /_/g') \
            -X ${LDFLAGS}.Release=TRUE" \
        ./cmd/server

### COMPILE FRONT END

RUN cd web &&\
    npm install &&\
    ng build --prod=true

### EXPOSE

EXPOSE 8080

CMD ./scripts/docker-starter.sh \
        "/var/vplan/server -c /etc/vplan/config/config.yml -web /var/vplan/web" \
        /etc/vplan/config/config.yml \
        ./config/docker.config.yml
