FROM golang
MAINTAINER Brian Morton "brian@mmm.hm"

RUN apt-get -y update && apt-get -y install pkg-config cmake
ADD https://github.com/libgit2/libgit2/archive/v0.22.2.tar.gz /tmp/
RUN cd /tmp && tar xfvz v0.22.2.tar.gz
RUN mkdir /tmp/libgit2-0.22.2/build
RUN cd /tmp/libgit2-0.22.2/build && cmake .. && cmake --build .
RUN cp -R /tmp/libgit2-0.22.2/build/* /usr/lib/
ENV PKG_CONFIG_PATH /tmp/libgit2-0.22.2/build
ENV C_INCLUDE_PATH /tmp/libgit2-0.22.2/include

ADD . /go/src/github.com/bmorton/builder
RUN cd /go/src/github.com/bmorton/builder && go get -v -d
RUN go install github.com/bmorton/builder
ENTRYPOINT ["/go/bin/builder"]
EXPOSE 3000
