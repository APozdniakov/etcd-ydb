# syntax=docker/dockerfile:1.2

FROM ubuntu:22.04 AS builder

RUN apt-get update && DEBIAN_FRONTEND=nointeractive apt-get install -y --no-install-recommends \
    ccache \
    clangd-14 \
    clang-14 \
    cmake \
    curl \
    gcc \
    gdb \
    git \
    g++ \
    iwyu \
    lld-14 \
    lldb-14 \
    llvm-14 \
    make \
    m4 \
    ninja-build \
    pkg-config \
    python3 \
    python3-dev \
    python3-jinja2 \
    python3-pip \
    python3-protobuf \
    python3-virtualenv \
    python3-voluptuous \
    tzdata \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && pip3 install --no-cache-dir conan==1.62.0

WORKDIR /etcdctl
RUN curl -L https://github.com/etcd-io/etcd/releases/download/v3.5.10/etcd-v3.5.10-linux-amd64.tar.gz -o /etcdctl/etcd-v3.5.10-linux-amd64.tar.gz \
    && tar xzvf /etcdctl/etcd-v3.5.10-linux-amd64.tar.gz --strip-components=1 --no-same-owner \
    && rm -f /etcdctl/etcd-v3.5.10-linux-amd64.tar.gz \
    && install -o root -g root -m 0755 etcdctl /usr/local/bin/etcdctl

RUN groupadd builder \
    && useradd -m -s /bin/bash -g builder builder
USER builder:builder

WORKDIR /conan
RUN --mount=type=bind,source=./conan,target=./ \
    conan config install .

# INFO [pavelbezpravel]:
# We copy the existing third_party/userver instead of cloning the repository,
# because in this case we can manually checkout to commits that are safe and reliable.
# I hope userver will be added to the conan-center-index soon.

WORKDIR /userver-framework
RUN --mount=type=bind,source=./third_party/userver,target=./userver \
    conan create \
    /userver-framework/userver \
    --build=missing \
    -pr=etcd-ydb/default \
    -pr:b=etcd-ydb/default \
    -pr:h=etcd-ydb/default \
    && conan remove "*" --builds --src --force

WORKDIR /home/builder
EXPOSE 22379 22380
