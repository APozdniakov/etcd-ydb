package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
	"github.com/ydb-platform/etcd-ydb/pkg/report"
)

var txnMixedCmd = &cobra.Command{
	Use:  "txn-mixed",
	RunE: txnMixedFunc,
}

var (
	txnMixedTotal        uint64
	txnMixedRateLimit    uint64
	txnMixedKeySize      uint64
	txnMixedValSize      uint64
	txnMixedKeySpaceSize uint64
	txnMixedOpsPerTxn    uint64
	txnMixedReadRatio    float64
)

func init() {
	RootCmd.AddCommand(txnMixedCmd)
	txnMixedCmd.Flags().Uint64Var(&txnMixedTotal, "total", 10000, "Total number of requests")
	txnMixedCmd.Flags().Uint64Var(&txnMixedRateLimit, "rate-limit", math.MaxUint64, "Maximum requests per second")
	txnMixedCmd.Flags().Uint64Var(&txnMixedKeySize, "key-size", 8, "Key size of request")
	txnMixedCmd.Flags().Uint64Var(&txnMixedValSize, "val-size", 8, "Value size of request")
	txnMixedCmd.Flags().Uint64Var(&txnMixedKeySpaceSize, "key-space-size", 1, "Maximum possible keys")
	txnMixedCmd.Flags().Uint64Var(&txnMixedOpsPerTxn, "txn-ops", 1, "Number of ops per txn")
	txnMixedCmd.Flags().Float64Var(&txnMixedReadRatio, "read-ratio", 0.5, "Read/all ops ratio")
}

func txnMixedFunc(_ *cobra.Command, _ []string) error {
	clients, err := newClients()
	if err != nil {
		return err
	}
	limit := rate.NewLimiter(rate.Limit(txnMixedRateLimit), 1)

	txnMixedTotal = uint64(float64(txnMixedTotal) / (1 - txnMixedReadRatio))
	txnMixedTotal /= txnMixedOpsPerTxn
	bar := pb.New64(int64(txnMixedTotal))
	bar.Start()

	ops := make(chan etcd.Request, totalClients)
	rep := report.NewReport(totalClients)
	var wg sync.WaitGroup
	for i := range clients {
		wg.Add(1)
		go func(client *etcd.Client) {
			defer wg.Done()
			for op := range ops {
				limit.Wait(context.Background())

				start := time.Now()
				_, err := etcd.Do(context.Background(), client, op)
				rep.Results() <- report.Result{TotalTime: time.Since(start), Err: err}
				bar.Increment()
			}
		}(clients[i])
	}

	go func() {
		key, value := []byte(strings.Repeat("-", int(txnMixedKeySize))), strings.Repeat("-", int(txnMixedKeySize))
		for range txnMixedTotal {
			success := make([]etcd.Request, txnMixedOpsPerTxn)
			for i := range success {
				j := 0
				for n := rand.Uint64() % txnMixedKeySpaceSize; n > 0; n /= 10 {
					key[j] = byte('0' + n%10)
					j++
				}
				slices.Reverse(key[:j])
				if rand.Float64() < txnMixedReadRatio {
					success[i] = &etcd.RangeRequest{Key: string(key)}
				} else {
					success[i] = &etcd.PutRequest{Key: string(key), Value: value}
				}
			}
			rand.Shuffle(len(success), func(i, j int) { success[i], success[j] = success[j], success[i] })

			var compare []etcd.Compare
			if txnMixedOpsPerTxn == 1 {
				compare = []etcd.Compare{etcd.Compare{Key: string(key)}.Equal().SetModRevision(0)}
			}

			op := &etcd.TxnRequest{Compare: compare, Success: success}
			ops <- op
		}
		close(ops)
	}()

	rc := rep.Run()
	wg.Wait()
	close(rep.Results())
	bar.Finish()
	stats := <-rc
	fmt.Fprintf(os.Stderr, "%#v\n", stats)
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
