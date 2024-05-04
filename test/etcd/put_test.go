package etcd_test

import (
	"testing"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

func fillPutRequest(revision *int64, request *etcd.PutRequest) {}

func fillPutResponse(revision *int64, response *etcd.PutResponse) {
	response.Revision = *revision
	if response.PrevKv != nil {
		response.PrevKv.ModRevision += *revision
		response.PrevKv.CreateRevision += *revision
	}
}

func TestPut(t *testing.T) {
	for _, tc := range []struct {
		name      string
		testcases []TestCase
	}{
		{
			name: "SetUp",
			testcases: []TestCase{
				{
					request:  &etcd.PutRequest{Key: "put_key", Value: "put_value0", PrevKv: true},
					response: &etcd.PutResponse{},
				},
			},
		},
		{
			name: "Basic EmptyKey",
			testcases: []TestCase{
				{
					request: &etcd.PutRequest{},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_")},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "put_value0"},
						},
					},
				},
				{
					request:  &etcd.PutRequest{Key: "put_key", Value: "put_value1"},
					response: &etcd.PutResponse{},
				},
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_"), Revision: -1},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "put_value0"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_"), Revision: 0},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: 0, CreateRevision: -1, Version: 2, Value: "put_value1"},
						},
					},
				},
			},
		},
		{
			name: "PrevKv",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_")},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: 0, CreateRevision: -1, Version: 2, Value: "put_value1"},
						},
					},
				},
				{
					request:  &etcd.PutRequest{Key: "put_key", Value: "put_value2", PrevKv: true},
					response: &etcd.PutResponse{PrevKv: &etcd.KeyValue{Key: "put_key", ModRevision: -1, CreateRevision: -2, Version: 2, Value: "put_value1"}},
				},
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_"), Revision: -1},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: -1, CreateRevision: -2, Version: 2, Value: "put_value1"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_"), Revision: 0},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: 0, CreateRevision: -2, Version: 3, Value: "put_value2"},
						},
					},
				},
			},
		},
		{
			name: "IgnoreValue",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_")},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: 0, CreateRevision: -2, Version: 3, Value: "put_value2"},
						},
					},
				},
				{
					request:  &etcd.PutRequest{Key: "put_key", IgnoreValue: true},
					response: &etcd.PutResponse{},
				},
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_"), Revision: -1},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: -1, CreateRevision: -3, Version: 3, Value: "put_value2"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_"), Revision: 0},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: 0, CreateRevision: -3, Version: 4, Value: "put_value2"},
						},
					},
				},
			},
		},
		{
			name: "IgnoreValue PrevKv",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_")},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: 0, CreateRevision: -3, Version: 4, Value: "put_value2"},
						},
					},
				},
				{
					request:  &etcd.PutRequest{Key: "put_key", IgnoreValue: true, PrevKv: true},
					response: &etcd.PutResponse{PrevKv: &etcd.KeyValue{Key: "put_key", ModRevision: -1, CreateRevision: -4, Version: 4, Value: "put_value2"}},
				},
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_"), Revision: -1},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: -1, CreateRevision: -4, Version: 4, Value: "put_value2"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_"), Revision: 0},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "put_key", ModRevision: 0, CreateRevision: -4, Version: 5, Value: "put_value2"},
						},
					},
				},
			},
		},
		{
			name: "IgnoreValue UnknownKey",
			testcases: []TestCase{
				{
					request: &etcd.PutRequest{Key: "put_unknown", IgnoreValue: true},
					err:     rpctypes.ErrGRPCKeyNotFound,
				},
			},
		},
		{
			name: "IgnoreValue WithValue",
			testcases: []TestCase{
				{
					request: &etcd.PutRequest{Key: "put_key", Value: "put_value", IgnoreValue: true},
					err:     rpctypes.ErrGRPCValueProvided,
				},
			},
		},
		{
			name: "TearDown",
			testcases: []TestCase{
				{
					request:  &etcd.DeleteRequest{Key: "put_", RangeEnd: etcd.GetPrefix("put_")},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{}},
				},
			},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
