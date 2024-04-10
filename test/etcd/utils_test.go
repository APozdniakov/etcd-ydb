package etcd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

var emptyKey = string([]byte{0})

func getPrefix(key string) string {
	range_end := strings.Clone(key)
	for i := len(range_end) - 1; i >= 0; i-- {
		if range_end[i] < 0xff {
			return string(append([]byte(range_end[:i]), range_end[i]+1))
		}
	}
	return emptyKey
}

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
	switch request := tc.request.(type) {
	case etcd.CompactRequest:
		switch response := tc.response.(type) {
		case nil:
			runCompactTestCase(t, client, request, nil, tc.err)
		case *etcd.CompactResponse:
			runCompactTestCase(t, client, request, response, tc.err)
		default:
			assert.Fail(t, "unknown response type")
		}
	case etcd.DeleteRequest:
		switch response := tc.response.(type) {
		case nil:
			runDeleteTestCase(t, client, revision, request, nil, tc.err)
		case *etcd.DeleteResponse:
			runDeleteTestCase(t, client, revision, request, response, tc.err)
		default:
			assert.Fail(t, "unknown response type")
		}
	case etcd.PutRequest:
		switch response := tc.response.(type) {
		case nil:
			runPutTestCase(t, client, request, nil, tc.err)
		case *etcd.PutResponse:
			runPutTestCase(t, client, request, response, tc.err)
		default:
			assert.Fail(t, "unknown response type")
		}
	case etcd.RangeRequest:
		switch response := tc.response.(type) {
		case nil:
			runRangeTestCase(t, client, request, nil, tc.err)
		case *etcd.RangeResponse:
			runRangeTestCase(t, client, request, response, tc.err)
		default:
			assert.Fail(t, "unknown response type")
		}
	case etcd.TxnRequest:
		switch response := tc.response.(type) {
		case nil:
			runTxnTestCase(t, client, request, nil, tc.err)
		case *etcd.TxnResponse:
			runTxnTestCase(t, client, request, response, tc.err)
		default:
			assert.Fail(t, "unknown response type")
		}
	default:
		assert.Fail(t, "unknown request type")
	}
}

func fillKeyValue(revision *int64, kv *etcd.KeyValue) {
	kv.ModRevision += *revision
	kv.CreateRevision += *revision
}

func compareKeyValue(t *testing.T, expected, actual *etcd.KeyValue) {
	if expected == nil {
		assert.Nil(t, actual)
		return
	}
	if !assert.NotNil(t, actual) {
		return
	}
	fmt.Printf("  actual = %#v\n", actual)
	fmt.Printf("expected = %#v\n", expected)
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.ModRevision, actual.ModRevision)
	assert.Equal(t, expected.CreateRevision, actual.CreateRevision)
	assert.Equal(t, expected.Version, actual.Version)
	assert.Equal(t, expected.Value, actual.Value)
}
