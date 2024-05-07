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

var txnPutCmd = &cobra.Command{
	Use:  "txn-put",
	RunE: txnPutFunc,
}

var (
	txnPutTotal        uint64
	txnPutKeySize      uint64
	txnPutValSize      uint64
	txnPutKeySpaceSize uint64
	txnPutOpsPerTxn    uint64
)

func init() {
	RootCmd.AddCommand(txnPutCmd)
	txnPutCmd.Flags().Uint64Var(&txnPutTotal, "total", 10000, "Total number of txn requests")
	txnPutCmd.Flags().Uint64Var(&txnPutKeySize, "key-size", 8, "Key size of txn")
	txnPutCmd.Flags().Uint64Var(&txnPutValSize, "val-size", 8, "Value size of txn")
	txnPutCmd.Flags().Uint64Var(&txnPutKeySpaceSize, "key-space-size", 1, "Maximum possible keys")
	txnPutCmd.Flags().Uint64Var(&txnPutOpsPerTxn, "txn-ops", 1, "Number of puts per txn")
}

func txnPutFunc(_ *cobra.Command, _ []string) error {
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

	txnPutTotal /= txnPutOpsPerTxn
	bar := pb.New64(int64(txnPutTotal))
	bar.Start()

	ops := make(chan etcd.Request, totalClients)
	rep := report.NewReport(totalClients)
	var wg sync.WaitGroup
	for i := range clients {
		wg.Add(1)
		go func(client *etcd.Client) {
			defer wg.Done()
			for op := range ops {
				start := time.Now()
				_, err := etcd.Do(client, op)
				rep.Results() <- report.Result{TotalTime: time.Since(start), Err: err}
				bar.Increment()
			}
		}(clients[i])
	}

	go func() {
		key, value := []byte(strings.Repeat("-", int(txnPutKeySize))), strings.Repeat("-", int(txnPutValSize))
		for range txnPutTotal {
			var compare []etcd.Compare
			if txnPutOpsPerTxn == 1 {
				compare = []etcd.Compare{etcd.Compare{Key: string(key)}.Equal().SetModRevision(0)}
			}

			success := make([]etcd.Request, txnPutOpsPerTxn)
			for i := range success {
				j := 0
				for n := rand.Uint64() % txnPutKeySpaceSize; n > 0; n /= 10 {
					key[j] = byte('0' + n%10)
					j++
				}
				slices.Reverse(key[:j])
				success[i] = &etcd.PutRequest{Key: string(key), Value: value}
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
