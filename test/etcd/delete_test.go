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

func fillDeleteRequest(revision *int64, request *etcd.DeleteRequest) {}

func fillDeleteResponse(revision *int64, response *etcd.DeleteResponse) {
	response.Revision = *revision
	for _, prev_kv := range response.PrevKvs {
		fillKeyValue(revision, prev_kv)
	}
}

func runDeleteTestCase(t *testing.T, client *etcd.Client, revision *int64, request etcd.DeleteRequest, expected *etcd.DeleteResponse, expectedErr error) {
	t.Helper()
	fillDeleteRequest(revision, &request)
	fmt.Printf(" request = %#v\n", request)
	actual, err := etcd.Delete(client, request)

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

	if expected.Deleted == 0 {
		initRevision(actual.Revision)
	} else {
		initRevision(actual.Revision - 1)
		(*revision)++
	}
	fillDeleteResponse(revision, expected)
	fmt.Printf("  actual = %#v\n", actual)
	fmt.Printf("expected = %#v\n", expected)
	assert.Equal(t, expected.Revision, actual.Revision)
	assert.Equal(t, expected.Deleted, actual.Deleted)
	sort.Slice(expected.PrevKvs, func(i, j int) bool { return expected.PrevKvs[i].Key < expected.PrevKvs[j].Key })
	sort.Slice(actual.PrevKvs, func(i, j int) bool { return actual.PrevKvs[i].Key < actual.PrevKvs[j].Key })
	assert.True(t, slices.EqualFunc(expected.PrevKvs, actual.PrevKvs, func(e, a *etcd.KeyValue) bool { compareKeyValue(t, e, a); return true }))
}

func TestDelete(t *testing.T) {
	client, err := etcd.NewClient(endpoint)
	require.NoError(t, err)

	for _, tc := range []struct {
		name      string
		testcases []TestCase
	}{
		{
			name: "Basic NilKey",
			testcases: []TestCase{
				{
					request: etcd.DeleteRequest{},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic EmptyKey",
			testcases: []TestCase{
				{
					request: etcd.DeleteRequest{Key: ""},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic Single",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "delete_key", Value: "delete_value", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.DeleteRequest{Key: "delete_key"},
					response: &etcd.DeleteResponse{Deleted: 1},
				},
			},
		},
		{
			name: "Basic Prefix",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "a", Value: "a"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key1", Value: "delete_value1"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key2", Value: "delete_value2"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key3", Value: "delete_value3"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "z", Value: "z"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.DeleteRequest{Key: "delete_", RangeEnd: getPrefix("delete_")},
					response: &etcd.DeleteResponse{Deleted: 3},
				},
				{
					request:  etcd.DeleteRequest{Key: "z"},
					response: &etcd.DeleteResponse{Deleted: 1},
				},
				{
					request:  etcd.DeleteRequest{Key: "a"},
					response: &etcd.DeleteResponse{Deleted: 1},
				},
			},
		},
		{
			name: "Basic FromKey",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "a", Value: "a"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key1", Value: "delete_value1"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key2", Value: "delete_value2"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key3", Value: "delete_value3"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "z", Value: "z"},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.DeleteRequest{Key: "delete_", RangeEnd: emptyKey},
					response: &etcd.DeleteResponse{Deleted: 4},
				},
				{
					request:  etcd.DeleteRequest{Key: "z"},
					response: &etcd.DeleteResponse{Deleted: 0},
				},
				{
					request:  etcd.DeleteRequest{Key: "a"},
					response: &etcd.DeleteResponse{Deleted: 1},
				},
			},
		},
		{
			name: "PrevKv Single",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "delete_key", Value: "delete_value", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.DeleteRequest{Key: "delete_key", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{{Key: "delete_key", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "delete_value"}}},
				},
			},
		},
		{
			name: "PrevKv Prefix",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "a", Value: "a", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key1", Value: "delete_value1", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key2", Value: "delete_value2", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key3", Value: "delete_value3", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "z", Value: "z", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.DeleteRequest{Key: "delete_", RangeEnd: getPrefix("delete_"), PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 3, PrevKvs: []*etcd.KeyValue{{Key: "delete_key1", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "delete_value1"}, {Key: "delete_key2", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "delete_value2"}, {Key: "delete_key3", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "delete_value3"}}},
				},
				{
					request:  etcd.DeleteRequest{Key: "z", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{{Key: "z", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "z"}}},
				},
				{
					request:  etcd.DeleteRequest{Key: "a", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{{Key: "a", ModRevision: -7, CreateRevision: -7, Version: 1, Value: "a"}}},
				},
			},
		},
		{
			name: "PrevKv FromKey",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "a", Value: "a", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key1", Value: "delete_value1", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key2", Value: "delete_value2", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "delete_key3", Value: "delete_value3", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.PutRequest{Key: "z", Value: "z", PrevKv: true},
					response: &etcd.PutResponse{},
				},
				{
					request:  etcd.DeleteRequest{Key: "delete_", RangeEnd: emptyKey, PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 4, PrevKvs: []*etcd.KeyValue{{Key: "delete_key1", ModRevision: -4, CreateRevision: -4, Version: 1, Value: "delete_value1"}, {Key: "delete_key2", ModRevision: -3, CreateRevision: -3, Version: 1, Value: "delete_value2"}, {Key: "delete_key3", ModRevision: -2, CreateRevision: -2, Version: 1, Value: "delete_value3"}, {Key: "z", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "z"}}},
				},
				{
					request:  etcd.DeleteRequest{Key: "z", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 0},
				},
				{
					request:  etcd.DeleteRequest{Key: "a", PrevKv: true},
					response: &etcd.DeleteResponse{Deleted: 1, PrevKvs: []*etcd.KeyValue{{Key: "a", ModRevision: -6, CreateRevision: -6, Version: 1, Value: "a"}}},
				},
			},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
