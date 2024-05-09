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

var txnRangeCmd = &cobra.Command{
	Use:  "txn-range",
	RunE: txnRangeFunc,
}

var (
	txnRangeTotal     uint64
	txnRangeRateLimit uint64
	txnRangeKeySize   uint64
	txnRangeValSize   uint64
	txnRangeOpsPerTxn uint64
)

func init() {
	RootCmd.AddCommand(txnRangeCmd)
	txnRangeCmd.Flags().Uint64Var(&txnRangeTotal, "total", 10000, "Total number of requests")
	txnRangeCmd.Flags().Uint64Var(&txnRangeRateLimit, "rate-limit", math.MaxUint64, "Maximum requests per second")
	txnRangeCmd.Flags().Uint64Var(&txnRangeKeySize, "key-size", 8, "Key size of request")
	txnRangeCmd.Flags().Uint64Var(&txnRangeValSize, "val-size", 8, "Value size of request")
	txnRangeCmd.Flags().Uint64Var(&txnRangeOpsPerTxn, "txn-ops", 1, "Number of ops per txn")
}

func txnRangeFunc(_ *cobra.Command, _ []string) error {
	conns := make([]*etcd.Client, totalConns)
	for i := range conns {
		conn, err := etcd.NewClient(endpoint)
		if err != nil {
			return err
		}
		conns[i] = conn
	}
	clients := make([]*etcd.Client, totalClients)
	for i := range clients {
		clients[i] = conns[i%len(conns)]
	}
	limit := rate.NewLimiter(rate.Limit(txnRangeRateLimit), 1)

	txnRangeTotal /= txnRangeOpsPerTxn
	bar := pb.New64(int64(txnRangeTotal))
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
		key := []byte(strings.Repeat("-", int(txnRangeKeySize)))
		for range txnRangeTotal {
			success := make([]etcd.Request, txnRangeOpsPerTxn)
			for i := range success {
				j := 0
				for n := rand.Uint64(); n > 0; n /= 10 {
					key[j] = byte('0' + n%10)
					j++
				}
				slices.Reverse(key[:j])
				success[i] = &etcd.RangeRequest{Key: string(key)}
			}

			var compare []etcd.Compare
			if txnRangeOpsPerTxn == 1 {
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
