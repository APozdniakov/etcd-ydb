# etcd-ydb

## Requirements

- docker 24.0.7
- docker compose 2.21.0

## Build

*TODO*: introduce [dev containers](https://containers.dev/)

```bash
mkdir -p cmake_build
docker build -t etcd-ydb/dev -f .devcontainer/Dockerfile .

docker run \
  --rm \
  -it \
  -u 1000:1000 \
  -v "$(pwd)":/tmp/etcd-ydb \
  etcd-ydb/dev \
  cmake \
  -DCMAKE_BUILD_TYPE=Release \
  -G Ninja \
  -DCMAKE_TOOLCHAIN_FILE=/tmp/etcd-ydb/clang.toolchain \
  -S /tmp/etcd-ydb \
  -B /tmp/etcd-ydb/cmake_build

docker run \
  --rm \
  -it \
  -u 1000:1000 \
  -v "$(pwd)":/tmp/etcd-ydb \
  etcd-ydb/dev \
  cmake \
  --build /tmp/etcd-ydb/cmake_build
```

## Run

```bash
docker compose up --build
```
