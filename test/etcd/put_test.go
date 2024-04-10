package etcd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

func fillPutRequest(revision *int64, request *etcd.PutRequest) {}

func fillPutResponse(revision *int64, response *etcd.PutResponse) {
	response.Revision = *revision
	if response.PrevKv != nil {
		fillKeyValue(revision, response.PrevKv)
	}
}

func runPutTestCase(t *testing.T, client *etcd.Client, request etcd.PutRequest, expected *etcd.PutResponse, expectedErr error) {
	t.Helper()
	fillPutRequest(revision, &request)
	fmt.Printf(" request = %#v\n", request)
	actual, err := etcd.Put(client, request)

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

	initRevision(actual.Revision - 1)
	(*revision)++
	fillPutResponse(revision, expected)
	fmt.Printf("  actual = %#v\n", actual)
	fmt.Printf("expected = %#v\n", expected)
	assert.Equal(t, expected.Revision, actual.Revision)
	compareKeyValue(t, expected.PrevKv, actual.PrevKv)
}

func TestPut(t *testing.T) {
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
					request: etcd.PutRequest{},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic EmptyKey",
			testcases: []TestCase{
				{
					request: etcd.PutRequest{Key: ""},
					err:     rpctypes.ErrGRPCEmptyKey,
				},
			},
		},
		{
			name: "Basic",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "put_key", Value: "put_value1"},
					response: &etcd.PutResponse{},
				},
			},
		},
		{
			name: "PrevKv",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "put_key", Value: "put_value2", PrevKv: true},
					response: &etcd.PutResponse{PrevKv: &etcd.KeyValue{Key: "put_key", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "put_value1"}},
				},
			},
		},
		{
			name: "IgnoreValue",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "put_key", IgnoreValue: true},
					response: &etcd.PutResponse{},
				},
			},
		},
		{
			name: "IgnoreValue PrevKv",
			testcases: []TestCase{
				{
					request:  etcd.PutRequest{Key: "put_key", PrevKv: true, IgnoreValue: true},
					response: &etcd.PutResponse{PrevKv: &etcd.KeyValue{Key: "put_key", ModRevision: -1, CreateRevision: -3, Version: 3, Value: "put_value2"}},
				},
			},
		},
		{
			name: "IgnoreValue UnknownKey",
			testcases: []TestCase{
				{
					request: etcd.PutRequest{Key: "put_unknown", IgnoreValue: true},
					err:     rpctypes.ErrGRPCKeyNotFound,
				},
			},
		},
		{
			name: "IgnoreValue WithValue",
			testcases: []TestCase{
				{
					request: etcd.PutRequest{Key: "put_key", Value: "put_value", IgnoreValue: true},
					err:     rpctypes.ErrGRPCValueProvided,
				},
			},
		},
		{
			name: "TearDown",
			testcases: []TestCase{
				{
					request:  etcd.DeleteRequest{Key: "put_", RangeEnd: getPrefix("put_")},
					response: &etcd.DeleteResponse{Deleted: 1},
				},
			},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
