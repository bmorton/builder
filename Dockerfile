FROM golang
MAINTAINER Brian Morton "brian@mmm.hm"

RUN apt-get -y update && apt-get -y install pkg-config cmake libssl-dev libssh-dev
ADD https://github.com/libgit2/libgit2/archive/v0.22.2.tar.gz /tmp/
RUN cd /tmp && tar xfvz v0.22.2.tar.gz
RUN mkdir /tmp/libgit2-0.22.2/build
RUN cd /tmp/libgit2-0.22.2/build && cmake .. && cmake --build .
RUN cp -R /tmp/libgit2-0.22.2/build/* /usr/lib/
ENV PKG_CONFIG_PATH /tmp/libgit2-0.22.2/build
ENV C_INCLUDE_PATH /tmp/libgit2-0.22.2/include

COPY . /go/src/github.com/bmorton/builder
RUN cd /go/src/github.com/bmorton/builder && go get -v -d
RUN go build --ldflags '-extldflags "-static"'
RUN go install github.com/bmorton/builder

RUN mkdir -p /dist/db
WORKDIR /dist
RUN cp /go/bin/builder .
RUN cp -R /go/src/github.com/bmorton/builder/static ./static
RUN cp /go/src/github.com/bmorton/builder/Dockerfile.run Dockerfile

CMD tar -cf - .
