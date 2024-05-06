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

var rangeCmd = &cobra.Command{
	Use:  "range",
	RunE: rangeFunc,
}

var (
	rangeTotal        uint64
	rangeKeySize      uint64
	rangeValSize      uint64
	rangeKeySpaceSize uint64
)

func init() {
	RootCmd.AddCommand(rangeCmd)
	rangeCmd.Flags().Uint64Var(&rangeTotal, "total", 10000, "Total number of range requests")
	rangeCmd.Flags().Uint64Var(&rangeKeySize, "key-size", 8, "Key size of range request")
	rangeCmd.Flags().Uint64Var(&rangeValSize, "val-size", 8, "Value size of range request")
	rangeCmd.Flags().Uint64Var(&rangeKeySpaceSize, "key-space-size", 1, "Maximum possible keys")
}

func rangeFunc(_ *cobra.Command, _ []string) error {
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

	bar := pb.New64(int64(rangeTotal))
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
		key := []byte(strings.Repeat("-", int(rangeKeySize)))
		for range rangeTotal {
			j := 0
			for n := rand.Uint64() % rangeKeySpaceSize; n > 0; n /= 10 {
				key[j] = byte('0' + n%10)
				j++
			}
			slices.Reverse(key[:j])
			op := &etcd.RangeRequest{Key: string(key)}
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
