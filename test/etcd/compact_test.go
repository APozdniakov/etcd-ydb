package etcd_test

import (
	"testing"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

func fillCompactRequest(revision *int64, request *etcd.CompactRequest) {
	request.Revision += *revision
}

func fillCompactResponse(revision *int64, response *etcd.CompactResponse) {
	response.Revision += *revision
}

func TestCompact(t *testing.T) {
	for _, tc := range []struct {
		name      string
		testcases []TestCase
	}{
		{
			name: "SetUp",
			testcases: []TestCase{
				{
					request:  &etcd.PutRequest{Key: "compact_key", Value: "compact_value1", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "compact_key", Value: "compact_value2", PrevKv: true},
					response: &etcd.PutResponse{PrevKv: &etcd.KeyValue{Key: "compact_key", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "compact_value1"}},
				},
			},
		},
		{
			name: "Basic",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "compact_key", Revision: -1},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "compact_key", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "compact_value1"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "compact_key"},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "compact_key", CreateRevision: -1, Version: 2, Value: "compact_value2"},
						},
					},
				},
				{
					request:  &etcd.CompactRequest{},
					response: &etcd.CompactResponse{},
				},
				{
					request: &etcd.RangeRequest{Key: "compact_key", Revision: -1},
					err:     rpctypes.ErrGRPCCompacted,
				},
				{
					request: &etcd.RangeRequest{Key: "compact_key"},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "compact_key", CreateRevision: -1, Version: 2, Value: "compact_value2"},
						},
					},
				},
			},
		},
		{
			name: "FutureRevision",
			testcases: []TestCase{
				{
					request: &etcd.CompactRequest{Revision: 1},
					err:     rpctypes.ErrGRPCFutureRev,
				},
			},
		},
		{
			name: "TearDown",
			testcases: []TestCase{
				{
					request: &etcd.DeleteRequest{Key: "compact_", RangeEnd: etcd.GetPrefix("compact_"), PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 1,
						PrevKvs: []*etcd.KeyValue{
							{Key: "compact_key", ModRevision: -1, CreateRevision: -2, Version: 2, Value: "compact_value2"},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
