package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
	"github.com/ydb-platform/etcd-ydb/pkg/report"
)

var txnCmd = &cobra.Command{
	Use:  "txn",
	RunE: txnFunc,
}

var (
	txnTotal     uint64
	txnReadRatio float64

	txnRangeLimit int64

	txnKeySize      uint64
	txnValSize      uint64
	txnKeySpaceSize uint64

	txnOpsPerTxn uint64
)

func init() {
	RootCmd.AddCommand(txnCmd)
	txnCmd.Flags().Uint64Var(&txnTotal, "total", 10000, "Total number of txn requests")
	txnCmd.Flags().Float64Var(&txnReadRatio, "read-ratio", 0.5, "Read/all ops ratio")
	txnCmd.Flags().Int64Var(&txnRangeLimit, "limit", 1000, "Read operation range result limit")
	txnCmd.Flags().Uint64Var(&txnKeySize, "key-size", 8, "Key size of txn")
	txnCmd.Flags().Uint64Var(&txnValSize, "val-size", 8, "Value size of txn")
	txnCmd.Flags().Uint64Var(&txnKeySpaceSize, "key-space-size", 1, "Maximum possible keys")
	txnCmd.Flags().Uint64Var(&txnOpsPerTxn, "txn-ops", 1, "Number of puts per txn")
}

func txnFunc(_ *cobra.Command, _ []string) error {
	client, err := etcd.NewClient(endpoint)
	if err != nil {
		return err
	}

	bar := pb.New64(int64(txnTotal))
	bar.Start()

	ops := make(chan etcd.Request, totalClients)
	rep := report.NewReport(totalClients)
	var wg sync.WaitGroup
	for range totalClients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for op := range ops {
				start := time.Now()
				_, err := etcd.Do(client, op)
				rep.Results() <- report.Result{TotalTime: time.Since(start), Err: err}
				bar.Increment()
			}
		}()
	}

	go func() {
		key, value := []byte(strings.Repeat("-", int(txnKeySize))), strings.Repeat("-", int(txnValSize))
		for range txnTotal {
			success := make([]etcd.Request, txnOpsPerTxn)
			for i := range success {
				if i < int(txnReadRatio*float64(len(success))) {
					success[i] = &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Limit: txnRangeLimit}
				} else {
					i := 0
					for n := rand.Uint64() % kvKeySpaceSize; n > 0; n /= 10 {
						key[i] = byte('0' + n%10)
						i++
					}
					slices.Reverse(key[:i])
					success[i] = &etcd.PutRequest{Key: string(key), Value: value}
				}
			}
			rand.Shuffle(len(success), func(i, j int) { success[i], success[j] = success[j], success[i] })
			op := &etcd.TxnRequest{Success: success}
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
