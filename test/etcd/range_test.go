package etcd_test

import (
	"testing"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

func fillRangeRequest(revision *int64, request *etcd.RangeRequest) {
	if request.Revision != 0 {
		request.Revision += *revision
	}
	if request.MinModRevision != 0 {
		request.MinModRevision += *revision
	}
	if request.MaxModRevision != 0 {
		request.MaxModRevision += *revision
	}
	if request.MinCreateRevision != 0 {
		request.MinCreateRevision += *revision
	}
	if request.MaxCreateRevision != 0 {
		request.MaxCreateRevision += *revision
	}
}

func fillRangeResponse(revision *int64, response *etcd.RangeResponse) {
	response.Revision += *revision
	for _, kv := range response.Kvs {
		kv.ModRevision += *revision
		kv.CreateRevision += *revision
	}
}

func TestRange(t *testing.T) {
	for _, tc := range []struct {
		name      string
		testcases []TestCase
	}{
		{
			name: "SetUp",
			testcases: []TestCase{
				{
					request:  &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.PutRequest{Key: "a", Value: "a", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "range_key3", Value: "range_value3", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "range_key2", Value: "range_value2", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "range_key1", Value: "range_value1", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  &etcd.PutRequest{Key: "range_key1", IgnoreValue: true, PrevKv: true},
					response: &etcd.PutResponse{PrevKv: &etcd.KeyValue{Key: "range_key1", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "range_value1"}},
				},
				{
					request:  &etcd.PutRequest{Key: "range_key2", IgnoreValue: true, PrevKv: true},
					response: &etcd.PutResponse{PrevKv: &etcd.KeyValue{Key: "range_key2", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value2"}},
				},
				{
					request:  &etcd.PutRequest{Key: "range_key3", IgnoreValue: true, PrevKv: true},
					response: &etcd.PutResponse{PrevKv: &etcd.KeyValue{Key: "range_key3", ModRevision: -5, CreateRevision: -5, Version: 1, Value: "range_value3"}},
				},
				{
					request:  &etcd.PutRequest{Key: "z", Value: "z", PrevKv: true},
					response: &etcd.PutResponse{},
				},
			},
		},
		{
			name: "Basic EmptyKey",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "range_key1"},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "range_key1"},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
						},
					},
				},
			},
		},
		{
			name: "Basic Range",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "range_key2", RangeEnd: "range_key3"},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
						},
					},
				},
			},
		},
		{
			name: "Basic Prefix",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")},
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
			},
		},
		{
			name: "Basic FromKey",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.EmptyKey},
					response: &etcd.RangeResponse{
						Count: 4,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
							{Key: "z", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "z"},
						},
					},
				},
			},
		},
		{
			name: "Basic All",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey},
					response: &etcd.RangeResponse{
						Count: 5,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
							{Key: "z", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "z"},
						},
					},
				},
			},
		},
		{
			name: "Limit",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), Limit: 1},
					response: &etcd.RangeResponse{
						Count: 3,
						More:  true,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), Limit: 2},
					response: &etcd.RangeResponse{
						Count: 3,
						More:  true,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), Limit: 3},
					response: &etcd.RangeResponse{
						Count: 3,
						More:  false,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), Limit: 4},
					response: &etcd.RangeResponse{
						Count: 3,
						More:  false,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), Limit: 5},
					response: &etcd.RangeResponse{
						Count: 3,
						More:  false,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
			},
		},
		{
			name: "Revision",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: -7},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: -6},
					response: &etcd.RangeResponse{
						Count: 2,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
							{Key: "range_key3", ModRevision: -6, CreateRevision: -6, Version: 1, Value: "range_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: -5},
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
							{Key: "range_key2", ModRevision: -5, CreateRevision: -5, Version: 1, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -6, CreateRevision: -6, Version: 1, Value: "range_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: -4},
					response: &etcd.RangeResponse{
						Count: 4,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
							{Key: "range_key1", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -5, CreateRevision: -5, Version: 1, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -6, CreateRevision: -6, Version: 1, Value: "range_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: -3},
					response: &etcd.RangeResponse{
						Count: 4,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -5, CreateRevision: -5, Version: 1, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -6, CreateRevision: -6, Version: 1, Value: "range_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: -2},
					response: &etcd.RangeResponse{
						Count: 4,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -6, CreateRevision: -6, Version: 1, Value: "range_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: -1},
					response: &etcd.RangeResponse{
						Count: 4,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: 0},
					response: &etcd.RangeResponse{
						Count: 5,
						Kvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
							{Key: "z", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "z"},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey, Revision: 1},
					err:     rpctypes.ErrGRPCFutureRev,
				},
			},
		},
		{
			name: "Sort",
			testcases: []TestCase{
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByKey().Ascending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByKey().Descending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByModRevision().Ascending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByModRevision().Descending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByCreateRevision().Ascending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByCreateRevision().Descending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByVersion().Ascending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByVersion().Descending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByValue().Ascending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
						},
					},
				},
				{
					request: etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_")}.OrderByValue().Descending(),
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2, Value: "range_value3"},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2, Value: "range_value2"},
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2, Value: "range_value1"},
						},
					},
				},
			},
		},
		{
			name: "KeysOnly",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), KeysOnly: true},
					response: &etcd.RangeResponse{
						Count: 3,
						Kvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -3, CreateRevision: -4, Version: 2},
							{Key: "range_key2", ModRevision: -2, CreateRevision: -5, Version: 2},
							{Key: "range_key3", ModRevision: -1, CreateRevision: -6, Version: 2},
						},
					},
				},
			},
		},
		{
			name: "CountOnly",
			testcases: []TestCase{
				{
					request:  &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), CountOnly: true},
					response: &etcd.RangeResponse{Count: 3, Kvs: []*etcd.KeyValue{}},
				},
				{
					request:  &etcd.RangeRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), KeysOnly: true, CountOnly: true},
					response: &etcd.RangeResponse{Count: 3, Kvs: []*etcd.KeyValue{}},
				},
			},
		},
		{
			name: "TearDown",
			testcases: []TestCase{
				{
					request: &etcd.DeleteRequest{Key: "a", PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 1,
						PrevKvs: []*etcd.KeyValue{
							{Key: "a", ModRevision: -8, CreateRevision: -8, Version: 1, Value: "a"},
						},
					},
				},
				{
					request: &etcd.DeleteRequest{Key: "range_", RangeEnd: etcd.GetPrefix("range_"), PrevKv: true},
					response: &etcd.DeleteResponse{
						Deleted: 3,
						PrevKvs: []*etcd.KeyValue{
							{Key: "range_key1", ModRevision: -5, CreateRevision: -6, Version: 2, Value: "range_value1"},
							{Key: "range_key2", ModRevision: -4, CreateRevision: -7, Version: 2, Value: "range_value2"},
							{Key: "range_key3", ModRevision: -3, CreateRevision: -8, Version: 2, Value: "range_value3"},
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
					request:  &etcd.RangeRequest{Key: etcd.EmptyKey, RangeEnd: etcd.EmptyKey},
					response: &etcd.RangeResponse{Count: 0, Kvs: []*etcd.KeyValue{}},
				},
			},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
