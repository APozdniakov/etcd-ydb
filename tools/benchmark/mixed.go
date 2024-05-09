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

var mixedCmd = &cobra.Command{
	Use:  "mixed",
	RunE: mixedFunc,
}

var (
	mixedTotal     uint64
	mixedRateLimit uint64
	mixedKeySize   uint64
	mixedValSize   uint64
	mixedReadRatio float64
)

func init() {
	RootCmd.AddCommand(mixedCmd)
	mixedCmd.Flags().Uint64Var(&mixedTotal, "total", 10000, "Total number of requests")
	mixedCmd.Flags().Uint64Var(&mixedRateLimit, "rate-limit", math.MaxUint64, "Maximum requests per second")
	mixedCmd.Flags().Uint64Var(&mixedKeySize, "key-size", 8, "Key size of request")
	mixedCmd.Flags().Uint64Var(&mixedValSize, "val-size", 8, "Value size of request")
	mixedCmd.Flags().Float64Var(&mixedReadRatio, "read-ratio", 0.5, "Read/all ops ratio")
}

func mixedFunc(_ *cobra.Command, _ []string) error {
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
	limit := rate.NewLimiter(rate.Limit(mixedRateLimit), 1)

	bar := pb.New64(int64(mixedTotal))
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
		key, value := []byte(strings.Repeat("-", int(mixedKeySize))), strings.Repeat("-", int(mixedValSize))
		for range mixedTotal {
			j := 0
			for n := rand.Uint64(); n > 0; n /= 10 {
				key[j] = byte('0' + n%10)
				j++
			}
			slices.Reverse(key[:j])

			var op etcd.Request
			if rand.Float64() < mixedReadRatio {
				op = &etcd.RangeRequest{Key: string(key)}
			} else {
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
	stats := <-rc
	fmt.Fprintf(os.Stderr, "%#v\n", stats)
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
