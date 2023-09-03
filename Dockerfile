# docker build -t etcd-ydb:dev -f Dockerfile .

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    clang-14 \
    cmake \
    gdb \
    git \
    lld-14 \
    lldb-14 \
    make \
    ninja-build \
    python3 \
    python3-pip \
    valgrind

RUN pip3 install conan==1.59.0

ARG UID=1000
RUN useradd -m -u ${UID} -s /bin/bash builder
