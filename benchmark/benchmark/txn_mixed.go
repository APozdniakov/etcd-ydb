package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"

	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/pkg/v3/report"
)

// mixeTxnCmd represents the mixedTxn command
var mixedTxnCmd = &cobra.Command{
	Use:   "txn-mixed key [end-range]",
	Short: "Benchmark a mixed load of txn-put & txn-range.",

	Run: mixedTxnFunc,
}

var (
	mixedTxnRate        int
	mixedTxnTotal       int
	mixedTxnEndKey      string
	mixedTxnRangeLimit  int64
	mixedTxnConsistency string
	mixedTxnOpsPerTxn   int
	mixedTxnReadRatio   float64
)

func init() {
	RootCmd.AddCommand(mixedTxnCmd)
	mixedTxnCmd.Flags().IntVar(&keySize, "key-size", 8, "Key size of mixed txn")
	mixedTxnCmd.Flags().IntVar(&valSize, "val-size", 8, "Value size of mixed txn")
	mixedTxnCmd.Flags().IntVar(&keySpaceSize, "key-space-size", 1, "Maximum possible keys")
	mixedTxnCmd.Flags().IntVar(&mixedTxnRate, "rate", 0, "Maximum txns per second (0 is no limit)")
	mixedTxnCmd.Flags().IntVar(&mixedTxnTotal, "total", 10000, "Total number of txn requests")
	mixedTxnCmd.Flags().StringVar(&mixedTxnEndKey, "end-key", "",
	"Read operation range end key. By default, we do full range query with the default limit of 1000.")
	mixedTxnCmd.Flags().Int64Var(&mixedTxnRangeLimit, "limit", 1000, "Read operation range result limit")
	mixedTxnCmd.Flags().StringVar(&mixedTxnConsistency, "consistency", "l", "Linearizable(l) or Serializable(s)")
	mixedTxnCmd.Flags().IntVar(&mixedTxnOpsPerTxn, "txn-ops", 1, "Number of puts per txn")
	mixedTxnCmd.Flags().Float64Var(&mixedTxnReadRatio, "read-ratio", 0.5, "Read/all ops ratio")
}

func mixedTxnFunc(cmd *cobra.Command, _ []string) {
	if keySpaceSize <= 0 {
		fmt.Fprintf(os.Stderr, "expected positive --key-space-size, got (%v)", keySpaceSize)
		os.Exit(1)
	}

	requests := make(chan []v3.Op, totalClients)
	if mixedTxnRate == 0 {
		mixedTxnRate = math.MaxInt32
	}
	limit := rate.NewLimiter(rate.Limit(mixedTxnRate), 1)
	clients := mustCreateClients(totalClients, totalConns)
	k, v := make([]byte, keySize), string(mustRandBytes(valSize))

	bar = pb.New(mixedTxnTotal)
	bar.Start()

	r := newReport()
	for i := range clients {
		wg.Add(1)
		go func(c *v3.Client) {
			defer wg.Done()
			for ops := range requests {
				limit.Wait(context.Background())

				st := time.Now()
				_, err := c.Txn(context.TODO()).Then(ops...).Commit()
				r.Results() <- report.Result{Err: err, Start: st, End: time.Now()}
				bar.Increment()
			}
		}(clients[i])
	}

	go func() {
		for i := 0; i < mixedTxnTotal; i++ {
			ops := make([]v3.Op, mixedTxnOpsPerTxn)
			for j := 0; j < len(ops); j++ {
				if j < int(mixedTxnReadRatio * float64(len(ops))) {
					opts := []v3.OpOption{v3.WithRange(mixedTxnEndKey), v3.WithPrefix(), v3.WithLimit(mixedTxnRangeLimit)}
					if mixedTxnConsistency == "s" {
						opts = append(opts, v3.WithSerializable())
					}
					ops[j] = v3.OpGet("", opts...)
				} else {
					binary.PutVarint(k, int64(rand.Intn(keySpaceSize)))
					ops[j] = v3.OpPut(string(k), v)
				}
			}
			rand.Shuffle(len(ops), func(i, j int) { ops[i], ops[j] = ops[j], ops[i] })
			requests <- ops
		}
		close(requests)
	}()

	rc := r.Run()
	wg.Wait()
	close(r.Results())
	bar.Finish()
	fmt.Println(<-rc)
}
