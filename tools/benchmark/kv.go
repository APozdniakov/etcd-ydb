package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
	"github.com/ydb-platform/etcd-ydb/pkg/report"
)

var kvCmd = &cobra.Command{
	Use:  "kv",
	RunE: kvFunc,
}

var (
	kvTotal     uint64
	kvReadRatio float64

	kvRangeLimit int64

	kvKeySize      uint64
	kvValSize      uint64
	kvKeySpaceSize uint64
)

func init() {
	RootCmd.AddCommand(kvCmd)
	kvCmd.Flags().Uint64Var(&kvTotal, "total", 10000, "Total number of kv requests")
	kvCmd.Flags().Float64Var(&kvReadRatio, "read-ratio", 0.5, "Read/all ops ratio")
	kvCmd.Flags().Int64Var(&kvRangeLimit, "limit", 1000, "Read operation range result limit")
	kvCmd.Flags().Uint64Var(&kvKeySize, "key-size", 8, "Key size of kv request")
	kvCmd.Flags().Uint64Var(&kvValSize, "val-size", 8, "Value size of kv request")
	kvCmd.Flags().Uint64Var(&kvKeySpaceSize, "key-space-size", 1, "Maximum possible keys")
}

func kvFunc(_ *cobra.Command, _ []string) error {
	client, err := etcd.NewClient(endpoint)
	if err != nil {
		return err
	}

	bar := pb.New64(int64(kvTotal))
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
		key, value := []byte(strings.Repeat("-", int(kvKeySize))), strings.Repeat("-", int(kvValSize))
		for range kvTotal {
			var op etcd.Request
			if rand.Float64() < kvReadRatio {
				op = &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Limit: kvRangeLimit}
			} else {
				i := 0
				for n := rand.Uint64()%kvKeySpaceSize; n > 0; n /= 10 {
					key[i] = byte('0' + n%10)
					i++
				}
				slices.Reverse(key[:i])
				op = &etcd.PutRequest{Key: string(key), Value: value}
			}
			ops <- op
		}
		close(ops)
	}()

	rc := rep.Run()
	wg.Wait()
	close(rep.Results())
	bar.Finish()
	fmt.Printf("%#v\n", <-rc)
	return nil
}
