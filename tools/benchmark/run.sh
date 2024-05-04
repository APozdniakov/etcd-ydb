function getArgs {
    ETCDCTL_FLAGS="--key-size=975 --val-size=975 --key-space-size=20_000_000_000 --total=$2"
    if [[ "$1" == "range" ]]; then
        echo "--total=$2"
    elif [[ "$1" == "put" ]]; then
        echo "$ETCDCTL_FLAGS --compact-index-delta=1000"
    elif [[ "$1" == "txn-mixed" ]]; then
        echo "$ETCDCTL_FLAGS --txn-ops=4 --read-ratio=0.25"
    elif [[ "$1" == "txn-put" ]]; then
        echo "$ETCDCTL_FLAGS --txn-ops=4"
    fi
}

echo $(getArgs "$1" "1_000_000")
time go run . --clients=1000 "$1" $(getArgs "$1" "1_000_000") > "../$1/1.txt"
time go run . --clients=1000 "$1" $(getArgs "$1" "1_000_000") > "../$1/2.txt"
echo $(getArgs "$1" "100_000")
time go run . --clients=1000 "$1" $(getArgs "$1" "100_000") > "../$1/3.txt"
