#!/usr/bin/env bash

function start_etcd {
    ssh master '/home/user/services/etcd/start.sh' 2> /dev/null
    sleep 1
}

function stop_etcd {
    ssh master '/home/user/services/etcd/stop.sh' 2> /dev/null
}

function start_ydb {
    ssh master '/home/user/services/ydb/ydbd/brrr.sh' 2> /dev/null
}

function stop_ydb {
    ssh master '/home/user/services/ydb/ydbd/stop.sh' 2> /dev/null
}


## $1: operation
## $2: total
## $3: read-ratio
## $4: txn-ops
function get_args {
    ETCDCTL_FLAGS="--key-size=11_700 --val-size=11_700 --key-space-size=20_000_000_000 --total=$2"
    if [[ "$1" == "put" ]]; then
        echo "$ETCDCTL_FLAGS"
    elif [[ "$1" == "range" ]]; then
        echo "$ETCDCTL_FLAGS"
    elif [[ "$1" == "txn" ]]; then
        echo "$ETCDCTL_FLAGS --read-ratio=$3 --txn-ops=$4"
    fi
}

## $1: target
function endpoint {
    if [[ "$1" == "etcd" ]]; then
        echo "158.160.22.130:2379"
    elif [[ "$1" == "ydb" ]]; then
        echo "158.160.22.130:2136"
    fi
}

## $1: target
## $2: operation
## $3: endpoint
## $4: read-ratio
## $5: txn-ops
function run {
    COUNTER=1
    start_$1
    for TOTAL in "40_000" "120_000" "80_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000"; do
        echo $(get_args "$2" "$TOTAL" "$4" "$5")
        time go run . --clients=1000 --conns=100 --endpoint="$3" "$2" $(get_args "$2" "$TOTAL" "$4" "$5") > "result/$1/$2/$4/$5/$COUNTER.json"
        if [[ "$2" == "put" ]]; then
            time go run . --clients=1000 --conns=100 --endpoint="$3" "range" $(get_args "range" "$TOTAL" "$4" "$5") > "result/$1/range/$4/$5/$COUNTER.json"
        fi
        if [[ "$1" == "etcd" ]]; then
            etcdctl --endpoints=$3 endpoint status -w table
        fi
        COUNTER=$(($COUNTER +1))
    done
    stop_$1
}

function main() {
    for TARGET in ydb; do
        for OPERATION in put; do
            if [[ $OPERATION == "put" ]]; then
                mkdir -p result/$TARGET/$OPERATION
                run $TARGET $OPERATION $(endpoint $TARGET)
            elif [[ $OPERATION == "txn" ]]; then
                for READ_RATIO in 0 0.25 0.5 0.75; do
                    for TXN_OPS in 1 8 64; do
                        mkdir -p result/$TARGET/$OPERATION/$READ_RATIO/$TXN_OPS
                        run $TARGET $OPERATION $(endpoint $TARGET) $READ_RATIO $TXN_OPS
                    done
                done
            fi
        done
    done
}

main
