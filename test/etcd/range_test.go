package etcd_test

import (
	"fmt"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

func fillRangeRequest(revision *int64, request *etcd.RangeRequest) {
	request.Revision += *revision
}

func fillRangeResponse(revision *int64, response *etcd.RangeResponse) {
	response.Revision += *revision
	for _, kv := range response.Kvs {
		fillKeyValue(revision, kv)
	}
}

func runRangeTestCase(t *testing.T, client *etcd.Client, request etcd.RangeRequest, expected *etcd.RangeResponse, expectedErr error) {
	t.Helper()
	fillRangeRequest(revision, &request)
	fmt.Printf(" request = %#v\n", request)
	actual, err := etcd.Range(client, request)

	if expectedErr != nil {
		assert.ErrorIs(t, err, expectedErr)
		return
	} else if !assert.NoError(t, err) {
		return
	}

	if expected == nil {
		assert.Nil(t, actual)
		return
	} else if !assert.NotNil(t, actual) {
		return
	}

	initRevision(actual.Revision)
	fillRangeResponse(revision, expected)
	fmt.Printf("  actual = %#v\n", actual)
	fmt.Printf("expected = %#v\n", expected)
	assert.Equal(t, expected.Revision, actual.Revision)
	assert.Equal(t, expected.Count, actual.Count)
	assert.Equal(t, expected.More, actual.More)
	sort.Slice(expected.Kvs, func(i, j int) bool { return expected.Kvs[i].Key < expected.Kvs[j].Key })
	sort.Slice(actual.Kvs, func(i, j int) bool { return actual.Kvs[i].Key < actual.Kvs[j].Key })
	assert.True(t, slices.EqualFunc(expected.Kvs, actual.Kvs, func(e, a *etcd.KeyValue) bool { compareKeyValue(t, e, a); return true }))
}

func TestRange(t *testing.T) {
	client, err := etcd.NewClient(endpoint)
	require.NoError(t, err)

	for _, tc := range []struct {
		name      string
		testcases []TestCase
	}{
		{
			name: "SetUp",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "a", Value: "a"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "range_key1", Value: "range_value1"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "range_key2", Value: "range_value2"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "range_key3", Value: "range_value3"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "z", Value: "z"},
					response: &etcd.PutResponse{},
				},
			},
		},
		{
			name: "Basic NilKey",
			testcases: []TestCase{
				{
					request: etcd.RangeRequest{},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic EmptyKey",
			testcases: []TestCase{
				{
					request: etcd.RangeRequest{Key: ""},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: "range_key1"},
					response: &etcd.RangeResponse{Count: 1, Kvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}}},
				},
			},
		},
		{
			name: "Basic Prefix",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: "range_", RangeEnd: getPrefix("range_")},
					response: &etcd.RangeResponse{Count: 3, Kvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}, {Key: "range_key2", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "range_value2"}, {Key: "range_key3", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "range_value3"}}},
				},
			},
		},
		{
			name: "Basic FromKey",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: "range_", RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 4, Kvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}, {Key: "range_key2", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "range_value2"}, {Key: "range_key3", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "range_value3"}, {Key: "z", Version: 1, Value: "z"}}},
				},
			},
		},
		{
			name: "Basic All",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: emptyKey, RangeEnd: emptyKey},
					response: &etcd.RangeResponse{Count: 5, Kvs: []*etcd.KeyValue{{Key: "a", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "a"}, {Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}, {Key: "range_key2", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "range_value2"}, {Key: "range_key3", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "range_value3"}, {Key: "z", Version: 1, Value: "z"}}},
				},
			},
		},
		{
			name: "Limit Total-2",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: "range_", RangeEnd: getPrefix("range_"), Limit: 1},
					response: &etcd.RangeResponse{Count: 3, More: true, Kvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}}},
				},
			},
		},
		{
			name: "Limit Total-1",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: "range_", RangeEnd: getPrefix("range_"), Limit: 2},
					response: &etcd.RangeResponse{Count: 3, More: true, Kvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}, {Key: "range_key2", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "range_value2"}}},
				},
			},
		},
		{
			name: "Limit Total",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: "range_", RangeEnd: getPrefix("range_"), Limit: 3},
					response: &etcd.RangeResponse{Count: 3, More: false, Kvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}, {Key: "range_key2", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "range_value2"}, {Key: "range_key3", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "range_value3"}}},
				},
			},
		},
		{
			name: "Limit Total+1",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: "range_", RangeEnd: getPrefix("range_"), Limit: 4},
					response: &etcd.RangeResponse{Count: 3, More: false, Kvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}, {Key: "range_key2", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "range_value2"}, {Key: "range_key3", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "range_value3"}}},
				},
			},
		},
		{
			name: "Limit Total+2",
			testcases: []TestCase{
				{
					request:  etcd.RangeRequest{Key: "range_", RangeEnd: getPrefix("range_"), Limit: 5},
					response: &etcd.RangeResponse{Count: 3, More: false, Kvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value1"}, {Key: "range_key2", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "range_value2"}, {Key: "range_key3", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "range_value3"}}},
				},
			},
		},
		{
			name: "TearDown",
			testcases: []TestCase{
				{
					request:  etcd.DeleteRequest{Key: "z", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{{Key: "z", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "z"}}},
				},
				{
					request:  etcd.DeleteRequest{Key: "range_", RangeEnd: getPrefix("range_"), PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 3, PrevKvs: []*etcd.KeyValue{{Key: "range_key1", ModRevision: -5, CreateRevision: -5, Version: 1, Value: "range_value1"}, {Key: "range_key2", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "range_value2"}, {Key: "range_key3", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "range_value3"}}},
				},
				{
					request:  etcd.DeleteRequest{Key: "a", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"}}},
				},
			},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
