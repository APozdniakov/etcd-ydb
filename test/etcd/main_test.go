package etcd_test

import (
	"flag"
	"os"
	"testing"
)

var (
	endpoint string
)

func Init() {
	flag.StringVar(&endpoint, "endpoint", "etcd:2379", "gRPC endpoint") // TODO: fix
	flag.Parse()
}

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}
