package etcd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

var revision *int64

type TestCase struct {
	request  etcd.Request
	response etcd.Response
	err      error
}

func runTest(client *etcd.Client, tcs []TestCase) func(*testing.T) {
	return func(t *testing.T) {
		for _, tc := range tcs {
			runTestCase(t, client, tc)
		}
	}
}

func runTestCase(t *testing.T, client *etcd.Client, tc TestCase) {
	t.Helper()
	fillRequest(revision, tc.request)
	fmt.Printf(" request = %#v\n", tc.request)
	actual, err := etcd.Do(context.Background(), client, tc.request)

	if tc.err != nil {
		assert.ErrorIs(t, err, tc.err)
		return
	} else if !assert.NoError(t, err) {
		return
	}

	if tc.response == nil {
		assert.Nil(t, actual)
		return
	} else if !assert.NotNil(t, actual) {
		return
	}

	if revision != nil && actual.IsWrite() {
		(*revision)++
	}
	if revision == nil {
		currentRevision := actual.GetRevision()
		revision = &currentRevision
	}
	fillResponse(revision, tc.response)
	fmt.Printf("  actual = %#v\n", actual)
	fmt.Printf("expected = %#v\n", tc.response)
	assert.Equal(t, tc.response, actual)
}

func fillRequest(revisoin *int64, request etcd.Request) {
	switch request := request.(type) {
	case *etcd.CompactRequest:
		fillCompactRequest(revisoin, request)
	case *etcd.DeleteRequest:
		fillDeleteRequest(revisoin, request)
	case *etcd.PutRequest:
		fillPutRequest(revisoin, request)
	case *etcd.RangeRequest:
		fillRangeRequest(revisoin, request)
	case *etcd.TxnRequest:
		fillTxnRequest(revision, request)
	default:
		panic("unknown request type")
	}
}

func fillResponse(revision *int64, response etcd.Response) {
	switch response := response.(type) {
	case *etcd.CompactResponse:
		fillCompactResponse(revision, response)
	case *etcd.DeleteResponse:
		fillDeleteResponse(revision, response)
	case *etcd.PutResponse:
		fillPutResponse(revision, response)
	case *etcd.RangeResponse:
		fillRangeResponse(revision, response)
	case *etcd.TxnResponse:
		fillTxnResponse(revision, response)
	default:
		panic("unknown response type")
	}
}
