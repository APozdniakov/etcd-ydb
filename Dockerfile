# syntax=docker/dockerfile:1

FROM kel.osll.ru:8083/docker/etcd-ydb-build-base:0.1.0 AS builder

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
EXPOSE 22379 22380
ENTRYPOINT ["./etcd-ydb"]
CMD ["-c", "/etcd-ydb/configs/static_config.yaml"]
