package etcd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// "go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

func fillTxnRequest(revision *int64, request *etcd.TxnRequest) {}

func fillTxnResponse(revision *int64, response *etcd.TxnResponse) {}

func runTxnTestCase(t *testing.T, client *etcd.Client, request etcd.TxnRequest, expected *etcd.TxnResponse, expectedErr error) {
	t.Helper()
	fillTxnRequest(revision, &request)
	fmt.Printf(" request = %#v\n", request)
	actual, err := etcd.Txn(client, request)

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
	fillTxnResponse(revision, expected)
	fmt.Printf("  actual = %#v\n", actual)
	fmt.Printf("expected = %#v\n", expected)
	assert.Equal(t, expected.Revision, actual.Revision)
}

func TestTxn(t *testing.T) {
	client, err := etcd.NewClient(endpoint)
	require.NoError(t, err)

	for _, tc := range []struct {
		name      string
		testcases []TestCase
	}{
		{
			name:      "Basic",
			testcases: []TestCase{},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
