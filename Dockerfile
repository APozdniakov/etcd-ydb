# syntax=docker/dockerfile:1

FROM ghcr.io/userver-framework/ubuntu-userver-build-base:v2

RUN apt-get update && apt-get install -y \
    clang-14 \
    cmake \
    git \
    lld-14 \
    lldb-14 \
    llvm-14 \
    make \
    ninja-build \
    pkg-config \
    python3 \
    python3-pip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && pip3 install conan==1.61.0

WORKDIR /etcd-ydb
COPY ./ ./

RUN mkdir -p cmake_build \
    && cmake \
    -DCMAKE_BUILD_TYPE=Release \
    -G Ninja \
    -DCMAKE_TOOLCHAIN_FILE=clang.toolchain \
    -S . \
    -B cmake_build \
    && cmake \
    --build cmake_build \
    -j $(nproc)

WORKDIR cmake_build
ENTRYPOINT ["./etcd-ydb"]
