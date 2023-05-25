# etcd-ydb

## Build Requirements

- cmake 3.22+
- clang-12
- lld-12
- python3.8
- pip3
- ninja 1.10+
- conan 1.59.0
- docker (optional)

## How to Build

### local

```bash
mkdir -p cmake_build
cmake -DCMAKE_BUILD_TYPE=Debug \
  -G Ninja \
  -DCMAKE_TOOLCHAIN_FILE=clang.toolchain \
  -S . \
  -B cmake_build

cmake --build cmake_build
```

### dockerized

```bash
mkdir -p cmake_build
docker build -t etcd-ydb:dev -f Dockerfile .

docker run --rm -it \
  -u 1000:1000 \
  -v $(pwd):/tmp/etcd-ydb \
  -v $(conan config home):/home/builder/.conan \
  etcd-ydb:dev \
  cmake -DCMAKE_BUILD_TYPE=Debug \
  -G Ninja \
  -DCMAKE_TOOLCHAIN_FILE=/tmp/etcd-ydb/clang.toolchain \
  -S /tmp/etcd-ydb \
  -B /tmp/etcd-ydb/cmake_build

docker run --rm -it \
  -u 1000:1000 \
  -v $(pwd):/tmp/etcd-ydb \
  -v $(conan config home):/home/builder/.conan \
  etcd-ydb:dev \
  cmake --build /tmp/etcd-ydb/cmake_build
```
