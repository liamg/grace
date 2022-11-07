FROM ubuntu:22.04

RUN apt-get update -y && apt-get install -y git build-essential golang ca-certificates rsync
RUN git clone --depth 1 https://github.com/torvalds/linux.git /linux
COPY . /src
WORKDIR /src
RUN make test