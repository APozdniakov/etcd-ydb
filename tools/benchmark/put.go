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

var putCmd = &cobra.Command{
	Use:  "put",
	RunE: putFunc,
}

var (
	putTotal        uint64
	putRateLimit    uint64
	putKeySize      uint64
	putValSize      uint64
	putKeySpaceSize uint64
)

func init() {
	RootCmd.AddCommand(putCmd)
	putCmd.Flags().Uint64Var(&putTotal, "total", 10000, "Total number of requests")
	putCmd.Flags().Uint64Var(&putRateLimit, "rate-limit", math.MaxUint64, "Maximum requests per second")
	putCmd.Flags().Uint64Var(&putKeySize, "key-size", 8, "Key size of request")
	putCmd.Flags().Uint64Var(&putValSize, "val-size", 8, "Value size of request")
	putCmd.Flags().Uint64Var(&putKeySpaceSize, "key-space-size", 1, "Maximum possible keys")
}

func putFunc(_ *cobra.Command, _ []string) error {
	clients, err := newClients()
	if err != nil {
		return err
	}
	limit := rate.NewLimiter(rate.Limit(putRateLimit), 1)

	bar := pb.New64(int64(putTotal))
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
		key, value := []byte(strings.Repeat("-", int(putKeySize))), strings.Repeat("-", int(putValSize))
		for range putTotal {
			j := 0
			for n := rand.Uint64() % putKeySpaceSize; n > 0; n /= 10 {
				key[j] = byte('0' + n%10)
				j++
			}
			slices.Reverse(key[:j])
			op := &etcd.PutRequest{Key: string(key), Value: value}
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
