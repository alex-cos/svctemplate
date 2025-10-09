
FROM alpine:3.22.2 AS base
RUN apk update && \
  apk add --no-cache bash git wget curl zip unzip gzip 

FROM base AS dev
# install updates and tools
RUN apk update && \
  apk add --no-cache make rsync gcc g++ musl-dev binutils pkgconfig \
  autoconf automake build-base binutils cmake libgcc libtool linux-headers
# download and install go verion 1.24.7
ENV GO_VERSION=1.24.7
RUN curl -sSLO https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz && \
  rm -rf /usr/local/go && \
  tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
  rm go$GO_VERSION.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

FROM dev AS build
WORKDIR /work
COPY . .
RUN make

FROM base
COPY --from=build /work/webserv /
COPY --from=build /work/config.yaml /
	
ENTRYPOINT ["/webserv"]
