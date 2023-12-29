# benchmark

## common

- `--key-size int`         Key size of put request (default 8)
- `--key-space-size int`   Maximum possible keys (default 1)
- `--rate int`             Maximum requests per second (0 is no limit)
- `--total int`            Total number of requests (default 10000)
- `--val-size int`         Value size of request (default 8)

```bash
export ETCDCTL_FLAGS="--total=1_000_000 --key-size=1_000 --key-space-size=20_000_000_000 --val-size=1_000"
```

Флаг `--total` подбирался так, чтобы после первого запуска БД заполнялась наполовину,
после второго запуска -- примерно до объема в 8 Гб,
а после третьего -- выше этого рекомендованного предела.
При постобработке заметна деградация производительности между 2 и 3 замерами

## range

- `--consistency string`   Linearizable(l) or Serializable(s) (default "l")
- `--count-only`           Only returns the count of keys
- `--limit int`            Maximum number of results to return from range request (0 is no limit)
- `--rate int`             Maximum range requests per second (0 is no limit)
- `--total int`            Total number of range requests (default 10000)

```bash
go run ./tools/benchmark --clients=1000 --conns=100 range foo --consistency=s --total=1_000_000
```

## put

- `--compact-index-delta int`     Delta between current revision and compact revision (e.g. current revision 10000, compact at 9000) (default 1000)
- `--compact-interval duration`   Interval to compact database (do not duplicate this with etcd's 'auto-compaction-retention' flag) (e.g. --compact-interval=5m compacts every 5-minute)
- `--sequential-keys`             Use sequential keys

```bash
go run ./tools/benchmark --clients=1000 --conns=100 put $ETCDCTL_FLAGS --compact-index-delta=1000
```

## txn-put

- `--txn-ops int`          Number of puts per txn (default 1)

```bash
go run ./tools/benchmark --clients=1000 --conns=100 txn-put $ETCDCTL_FLAGS --txn-ops=4
```

## txn-mixed

- `--consistency string`   Linearizable(l) or Serializable(s) (default "l")
- `--end-key string`       Read operation range end key. By default, we do full range query with the default limit of 1000.
- `--limit int`            Read operation range result limit (default 1000)
- `--rw-ratio float`       Read/write ops ratio (default 1)

```bash
go run ./tools/benchmark --clients=1000 --conns=100 txn-mixed $ETCDCTL_FLAGS --consistency=s --rw-ratio=0.2
```
