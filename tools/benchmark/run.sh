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
## $3: txn-ops
function get_args {
    ETCDCTL_FLAGS="--rate-limit=10_000 --key-size=11_700 --val-size=11_700 --total=$2"
    if [[ "$1" == "put" ]]; then
        echo "$ETCDCTL_FLAGS"
    elif [[ "$1" == "range" ]]; then
        echo "$ETCDCTL_FLAGS"
    elif [[ "$1" == "mixed" ]]; then
        echo "$ETCDCTL_FLAGS --read-ratio=0.75"
    elif [[ "$1" == "txn-put" ]]; then
        echo "$ETCDCTL_FLAGS --txn-ops=$3"
    elif [[ "$1" == "txn-range" ]]; then
        echo "$ETCDCTL_FLAGS --txn-ops=$3"
    elif [[ "$1" == "txn-mixed" ]]; then
        echo "$ETCDCTL_FLAGS --txn-ops=$3 --read-ratio=0.75"
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
## $3: total
## $4: txn-ops
## $5: counter
function run {
    echo "go run . --clients=100 --conns=10 --endpoint=$(endpoint $1) $2 $(get_args $2 $3 $4) > result/$1/$2/$4/$5.json"
    time  go run . --clients=100 --conns=10 --endpoint="$(endpoint $1)" "$2" $(get_args "$2" "$3" "$4") > "result/$1/$2/$4/$5.json"
}

## $1: target
## $2: operation
## $3: txn-ops
function fill {
    echo "START $1"
    start_$1

    mkdir -p result/$1/$2/$3/
    if [[ "$2" == "put" ]]; then
        mkdir -p result/$1/range/
    elif [[ "$2" == "txn-put" ]]; then
        mkdir -p result/$1/txn-range/$3/
    fi

    COUNTER=1
    for TOTAL in "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000" "40_000"; do
        run "$1" "$2" "$TOTAL" "$3" "$COUNTER"
        if [[ "$2" == "put" ]]; then
            run "$1" "range" "$TOTAL" "$3" "$COUNTER"
        elif [[ "$2" == "txn-put" ]]; then
            run "$1" "txn-range" "$TOTAL" "$3" "$COUNTER"
        fi
        if [[ "$1" == "etcd" ]]; then
            echo "etcdctl --endpoints=$(endpoint $1) endpoint status -w table"
            etcdctl --endpoints="$(endpoint $1)" endpoint status -w table
        fi
        COUNTER=$(($COUNTER +1))
    done

    echo "STOP $1"
    stop_$1
}

function main() {
    for TARGET in etcd ydb; do
        for OPERATION in put mixed txn-put txn-mixed; do
            if [[ $OPERATION == "put" ]]; then
                fill $TARGET $OPERATION
            elif [[ $OPERATION == "mixed" ]]; then
                fill $TARGET $OPERATION
            elif [[ $OPERATION == "txn-put" ]]; then
                for TXN_OPS in 1 8 64; do
                    fill $TARGET $OPERATION $TXN_OPS
                done
            elif [[ $OPERATION == "txn-mixed" ]]; then
                for TXN_OPS in 1 8 64; do
                    fill $TARGET $OPERATION $TXN_OPS
                done
            fi
        done
    done
}

main
