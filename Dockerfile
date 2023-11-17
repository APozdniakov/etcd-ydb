# syntax=docker/dockerfile:1

FROM ghcr.io/userver-framework/ubuntu-userver-build-base:v2 AS builder

RUN apt-get update && apt-get install -y \
    clang-14 \
    cmake \
    git \
    lld-14 \
    llvm-14 \
    make \
    m4 \
    ninja-build \
    pkg-config \
    python3 \
    python3-pip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && pip3 install conan==1.62.0

WORKDIR /conan
COPY ./conan ./
RUN conan config install .

# INFO [pavelbezpravel]:
# We copy the existing third_party/userver instead of cloning the repository,
# because in this case we can manually checkout to commits that are safe and reliable.
# I hope userver will be added to the conan-center-index soon.

WORKDIR /userver-framework
COPY third_party/userver ./userver
RUN conan create \
    /userver-framework/userver \
    --build=missing \
    -pr:b=default

WORKDIR /etcd-ydb
COPY ./ ./

RUN cmake \
    -DCMAKE_BUILD_TYPE=Release \
    -G Ninja \
    -DCMAKE_TOOLCHAIN_FILE=clang.toolchain \
    -S . \
    -B cmake_build \
    && cmake \
    --build cmake_build \
    -j "$(nproc)"

FROM ubuntu:22.04 AS runner
WORKDIR /etcd-ydb
COPY --from=builder /etcd-ydb/cmake_build/etcd-ydb ./
COPY --from=builder /etcd-ydb/configs ./configs/
EXPOSE 2379 2380
ENTRYPOINT ["./etcd-ydb"]
CMD ["-c", "/etcd-ydb/configs/static_config.yaml"]
