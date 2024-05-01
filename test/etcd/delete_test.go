package etcd_test

import (
	"testing"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

func fillDeleteRequest(revision *int64, request *etcd.DeleteRequest) {}

func fillDeleteResponse(revision *int64, response *etcd.DeleteResponse) {
	response.Revision = *revision
	for _, prev_kv := range response.PrevKvs {
		prev_kv.ModRevision += *revision
		prev_kv.CreateRevision += *revision
	}
}

func TestDelete(t *testing.T) {
	for _, tc := range []struct {
		name      string
		testcases []TestCase
	}{
		{
			name: "Basic EmptyKey",
			testcases: []TestCase{
				{
					request: &etcd.DeleteRequest{},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic Single",
			testcases: []TestCase{
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key1", Value: "delete_value1", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key2", Value: "delete_value2", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 2,
						Kvs: []*etcd.KeyValue{
							{Key: "delete_key1", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "delete_value1"},
							{Key: "delete_key2", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "delete_value2"},
						},
					},
				},
				{
					request:  &etcd.DeleteRequest{Key: "delete_key1"},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "delete_key2", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "delete_value2"},
						},
					},
				},
				{
					request:  &etcd.DeleteRequest{Key: "delete_key1"},
					response: &etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "delete_key2", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "delete_value2"},
						},
					},
				},
				{
					request:  &etcd.DeleteRequest{Key: "delete_key2"},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
			},
		},
		{
			name: "Basic Prefix",
			testcases: []TestCase{
				{
					request:  &etcd.PutRequest{Key: "delete_key1", Value: "delete_value1"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key2", Value: "delete_value2"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key3", Value: "delete_value3"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "z", Value: "z"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "a", Value: "a"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.DeleteRequest{Key: "delete_", RangeEnd: getPrefix("delete_")},
					response: &etcd.DeleteResponse{Deleted: 3, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "z"},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "a"},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{}},
				},
			},
		},
		{
			name: "Basic FromKey",
			testcases: []TestCase{
				{
					request:  &etcd.PutRequest{Key: "delete_key1", Value: "delete_value1"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key2", Value: "delete_value2"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key3", Value: "delete_value3"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "z", Value: "z"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "a", Value: "a"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.DeleteRequest{Key: "delete_", RangeEnd: emptyKey},
					response: &etcd.DeleteResponse{Deleted: 4, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "z"},
					response: &etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "a"},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{}},
				},
			},
		},
		{
			name: "Basic All",
			testcases: []TestCase{
				{
					request:  &etcd.PutRequest{Key: "delete_key1", Value: "delete_value1"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key2", Value: "delete_value2"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key3", Value: "delete_value3"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "z", Value: "z"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "a", Value: "a"},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.DeleteRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.DeleteResponse{Deleted: 5, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "z"},
					response: &etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "a"},
					response: &etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
				},
			},
		},
		{
			name: "PrevKv Single",
			testcases: []TestCase{
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key", Value: "delete_value", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "delete_key", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "delete_value"},
						},
					},
				},
				{
					request: &etcd.DeleteRequest{Key: "delete_key", PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 1,
						PrevKvs: []*etcd.KeyValue{
							{Key: "delete_key", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "delete_value"},
						},
					},
				},
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "delete_key", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
			},
		},
		{
			name: "PrevKv Prefix",
			testcases: []TestCase{
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key1", Value: "delete_value1", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key2", Value: "delete_value2", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key3", Value: "delete_value3", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "z", Value: "z", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "a", Value: "a", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 5,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "a"},
							{Key: "delete_key1", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "delete_value1"},
							{Key: "delete_key2", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "delete_value2"},
							{Key: "delete_key3", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "delete_value3"},
							{Key: "z", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "z"},
						},
					},
				},
				{
					request: &etcd.DeleteRequest{Key: "delete_", RangeEnd: getPrefix("delete_"), PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 3,
						PrevKvs: []*etcd.KeyValue{
							{Key: "delete_key1", ModRevision: -5, CreateRevision: -5, Version: 1, Value: "delete_value1"},
							{Key: "delete_key2", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "delete_value2"},
							{Key: "delete_key3", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "delete_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 2,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "a"},
							{Key: "z", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "z"},
						},
					},
				},
				{
					request: &etcd.DeleteRequest{Key: "z", PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 1,
						PrevKvs: []*etcd.KeyValue{
							{Key: "z", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "z"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "a"},
						},
					},
				},
				{
					request: &etcd.DeleteRequest{Key: "a", PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 1,
						PrevKvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "a"},
						},
					},
				},
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
			},
		},
		{
			name: "PrevKv FromKey",
			testcases: []TestCase{
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key1", Value: "delete_value1", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key2", Value: "delete_value2", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key3", Value: "delete_value3", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "z", Value: "z", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "a", Value: "a", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 5,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "a"},
							{Key: "delete_key1", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "delete_value1"},
							{Key: "delete_key2", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "delete_value2"},
							{Key: "delete_key3", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "delete_value3"},
							{Key: "z", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "z"},
						},
					},
				},
				{
					request: &etcd.DeleteRequest{Key: "delete_", RangeEnd: emptyKey, PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 4,
						PrevKvs: []*etcd.KeyValue{
							{Key: "delete_key1", ModRevision: -5, CreateRevision: -5, Version: 1, Value: "delete_value1"},
							{Key: "delete_key2", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "delete_value2"},
							{Key: "delete_key3", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "delete_value3"},
							{Key: "z", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "z"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "a"},
						},
					},
				},
				{
					request:  &etcd.DeleteRequest{Key: "z", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "a"},
						},
					},
				},
				{
					request: &etcd.DeleteRequest{Key: "a", PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 1,
						PrevKvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "a"},
						},
					},
				},
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
			},
		},
		{
			name: "PrevKv All",
			testcases: []TestCase{
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key1", Value: "delete_value1", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key2", Value: "delete_value2", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "delete_key3", Value: "delete_value3", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "z", Value: "z", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "a", Value: "a", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request: &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{
						Count: 5,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "a"},
							{Key: "delete_key1", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "delete_value1"},
							{Key: "delete_key2", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "delete_value2"},
							{Key: "delete_key3", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "delete_value3"},
							{Key: "z", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "z"},
						},
					},
				},
				{
					request: &etcd.DeleteRequest{Key: emptyKey, RangeEnd: emptyKey, PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 5,
						PrevKvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "a"},
							{Key: "delete_key1", ModRevision: -5, CreateRevision: -5, Version: 1, Value: "delete_value1"},
							{Key: "delete_key2", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "delete_value2"},
							{Key: "delete_key3", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "delete_value3"},
							{Key: "z", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "z"},
						},
					},
				},
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "z", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.DeleteRequest{Key: "a", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
			},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
