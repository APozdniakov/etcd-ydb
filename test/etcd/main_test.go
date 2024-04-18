package etcd_test

import (
	"os"
	"testing"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

const endpoint = "etcd:2379"

var client *etcd.Client

func Init() {
	var err error
	client, err = etcd.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}
