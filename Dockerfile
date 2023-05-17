# docker build -t etcd-ydb:dev -f Dockerfile .

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    clang-12 \
    cmake \
    gdb \
    git \
    lld-12 \
    lldb-12 \
    make \
    ninja-build \
    python3 \
    python3-pip \
    valgrind

RUN pip3 install conan==1.59.0

ARG UID=1000
RUN useradd -m -u ${UID} -s /bin/bash builder
