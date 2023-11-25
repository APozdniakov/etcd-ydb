# etcd-ydb

## Prerequisites

### git clone

Clone with `--recurse-submodules` option or run `git submodule update --init`

### docker

Append this to `/etc/docker/daemon.json`:
```
{
    "insecure-registries": [
        "kel.osll.ru:8083"
    ]
}
```

Run:
```bash
sudo systemctl restart docker.service
```

### devcontainers/cli

```bash
npm install -g @devcontainers/cli
```

## Requirements

- docker 24.0.7
- docker compose 2.21.0
- devcontainers/cli (optional)

## Build

### Docker

You can build a project manually with `docker`. You have to build an image and run a containers with `cmake` commands:
```bash
docker pull kel.osll.ru:8083/docker/etcd-ydb-build-base:0.1.0
mkdir -p cmake_build

docker run \
  --rm \
  -it \
  -u 1000:1000 \
  -v "$(pwd)":/workspaces/etcd-ydb \
  kel.osll.ru:8083/docker/etcd-ydb-build-base:0.1.0 \
  cmake \
  -DCMAKE_BUILD_TYPE=Release \
  -G Ninja \
  -DCMAKE_TOOLCHAIN_FILE=/workspaces/etcd-ydb/clang.toolchain \
  -S /workspaces/etcd-ydb \
  -B /workspaces/etcd-ydb/cmake_build

docker run \
  --rm \
  -it \
  -u 1000:1000 \
  -v "$(pwd)":/workspaces/etcd-ydb \
  kel.osll.ru:8083/docker/etcd-ydb-build-base:0.1.0 \
  cmake \
  --build /workspaces/etcd-ydb/cmake_build
```

### Dev Container Cli

Also, you can build a project with [devcontainers](https://containers.dev/) in an easy and convenient way.  
Your IDE (e.g. Clion) or code editor (e.g. VS Code) can run and attach to devcontainer.  

You can use devcontainers/cli to set up environment and build the project manually via bash:
```bash
devcontainer up --workspace-folder .

devcontainer exec --workspace-folder . \
  cmake \
  -DCMAKE_BUILD_TYPE=Release \
  -G Ninja \
  -DCMAKE_TOOLCHAIN_FILE=/workspaces/etcd-ydb/clang.toolchain \
  -S /workspaces/etcd-ydb \
  -B /workspaces/etcd-ydb/cmake_build

devcontainer exec --workspace-folder . \
  cmake \
  --build /workspaces/etcd-ydb/cmake_build
```

## Run

You can run existing binary:
```bash
cmake_buid/etcd-ydb -c configs/static_config.yaml
```

Also you can run dockerized application via docker compose:
```bash
docker compose up --build
```
