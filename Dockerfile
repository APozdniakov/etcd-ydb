# syntax=docker/dockerfile:1

FROM kel.osll.ru:8083/docker/etcd-ydb-build-base:0.1.0 AS builder

WORKDIR /etcd-ydb
COPY ./ ./

RUN cmake \
    --preset=release \
    && cmake \
    --build \
    --preset=release

FROM ubuntu:22.04 AS runner
WORKDIR /etcd-ydb
COPY --from=builder /etcd-ydb/cmake-build-release/etcd-ydb ./
COPY --from=builder /etcd-ydb/configs ./configs/
EXPOSE 22379 22380
ENTRYPOINT ["./etcd-ydb"]
CMD ["-c", "/etcd-ydb/configs/static_config.yaml"]
