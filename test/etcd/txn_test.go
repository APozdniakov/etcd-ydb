package etcd_test

import (
	"math"
	"testing"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

func fillTxnRequest(revision *int64, request *etcd.TxnRequest) {
	for _, compare := range request.Compare {
		if compare.ModRevision != nil {
			if *compare.ModRevision == math.MinInt64 {
				*compare.ModRevision = 0
			} else {
				*compare.ModRevision += *revision
			}
		}
		if compare.CreateRevision != nil {
			*compare.CreateRevision += *revision
		}
	}
	for _, success := range request.Success {
		switch success := success.(type) {
		case *etcd.DeleteRequest:
			fillDeleteRequest(revision, success)
		case *etcd.PutRequest:
			fillPutRequest(revision, success)
		case *etcd.RangeRequest:
			fillRangeRequest(revision, success)
		case *etcd.TxnRequest:
			fillTxnRequest(revision, success)
		default:
			panic("unknown request type")
		}
	}
	for _, failure := range request.Failure {
		switch failure := failure.(type) {
		case *etcd.DeleteRequest:
			fillDeleteRequest(revision, failure)
		case *etcd.PutRequest:
			fillPutRequest(revision, failure)
		case *etcd.RangeRequest:
			fillRangeRequest(revision, failure)
		case *etcd.TxnRequest:
			fillTxnRequest(revision, failure)
		default:
			panic("unknown request type")
		}
	}
}

func fillTxnResponse(revision *int64, response *etcd.TxnResponse) {
	response.Revision += *revision
	for _, response := range response.Responses {
		switch response := response.(type) {
		case *etcd.DeleteResponse:
			fillDeleteResponse(revision, response)
		case *etcd.PutResponse:
			fillPutResponse(revision, response)
		case *etcd.RangeResponse:
			fillRangeResponse(revision, response)
		case *etcd.TxnResponse:
			fillTxnResponse(revision, response)
			response.Revision = 0
		default:
			panic("unknown response type")
		}
	}
}

func TestTxn(t *testing.T) {
	for _, tc := range []struct {
		name      string
		testcases []TestCase
	}{
		{
			name: "SetUp",
			testcases: []TestCase{
				{
					request: &etcd.TxnRequest{
						Compare: []etcd.Compare{},
						Success: []etcd.Request{
							&etcd.PutRequest{Key: "txn_key1", Value: "txn_value1", PrevKv: true},
						},
						Failure: []etcd.Request{},
					},
					response: &etcd.TxnResponse{
						Succeeded: true,
						Responses: []etcd.Response{
							&etcd.PutResponse{},
						},
					},
				},
			},
		},
		{
			name: "Basic",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "txn_", RangeEnd: etcd.GetPrefix("txn_")},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "txn_key1", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "txn_value1"},
						},
					},
				},
				{
					request: &etcd.TxnRequest{
						Compare: []etcd.Compare{},
						Success: []etcd.Request{
							&etcd.RangeRequest{Key: "txn_", RangeEnd: etcd.GetPrefix("txn_")},
							&etcd.DeleteRequest{Key: "z", PrevKv: true},
						},
						Failure: []etcd.Request{
							&etcd.PutRequest{Key: "txn_key2", Value: "txn_value2"},
						},
					},
					response: &etcd.TxnResponse{
						Succeeded: true,
						Responses: []etcd.Response{
							&etcd.RangeResponse{
								Count: 1,
								Kvs: []*etcd.KeyValue{
									{Key: "txn_key1", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "txn_value1"},
								},
							},
							&etcd.DeleteResponse{Deleted: 0, PrevKvs: []*etcd.KeyValue{}},
						},
					},
				},
				{
					request: &etcd.RangeRequest{Key: "txn_", RangeEnd: etcd.GetPrefix("txn_")},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "txn_key1", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "txn_value1"},
						},
					},
				},
			},
		},
		{
			name: "Check Absence",
			testcases: []TestCase{
				{
					request: &etcd.TxnRequest{
						Compare: []etcd.Compare{etcd.Compare{Key: "txn_key1"}.Equal().SetModRevision(math.MinInt64)},
						Success: []etcd.Request{},
						Failure: []etcd.Request{},
					},
					response: &etcd.TxnResponse{
						Succeeded: false,
						Responses: []etcd.Response{},
					},
				},
				{
					request: &etcd.TxnRequest{
						Compare: []etcd.Compare{etcd.Compare{Key: "txn_key2"}.Equal().SetModRevision(math.MinInt64)},
						Success: []etcd.Request{},
						Failure: []etcd.Request{},
					},
					response: &etcd.TxnResponse{
						Succeeded: true,
						Responses: []etcd.Response{},
					},
				},
			},
		},
		{
			name: "Empty Compare",
			testcases: []TestCase{
				{
					request:  &etcd.TxnRequest{Compare: []etcd.Compare{}, Success: []etcd.Request{}, Failure: []etcd.Request{}},
					response: &etcd.TxnResponse{Succeeded: true, Responses: []etcd.Response{}},
				},
			},
		},
		{
			name: "Basic Compare",
			testcases: []TestCase{
				{
					request: &etcd.TxnRequest{
						Compare: []etcd.Compare{etcd.Compare{Key: "txn_key1"}.Greater().SetValue("txn_value")},
						Success: []etcd.Request{},
						Failure: []etcd.Request{},
					},
					response: &etcd.TxnResponse{Succeeded: true, Responses: []etcd.Response{}},
				},
			},
		},
		{
			name: "Mixed Compare",
			testcases: []TestCase{
				{
					request: &etcd.TxnRequest{
						Compare: []etcd.Compare{
							etcd.Compare{Key: "txn_key1"}.Less().SetModRevision(1),
							etcd.Compare{Key: "txn_key1"}.NotEqual().SetCreateRevision(1),
							etcd.Compare{Key: "txn_key1"}.Greater().SetVersion(0),
							etcd.Compare{Key: "txn_key1"}.Equal().SetValue("txn_value1"),
						},
						Success: []etcd.Request{},
						Failure: []etcd.Request{},
					},
					response: &etcd.TxnResponse{Succeeded: true, Responses: []etcd.Response{}},
				},
			},
		},
		{
			name: "DuplicateKey",
			testcases: []TestCase{
				{
					request: &etcd.RangeRequest{Key: "txn_", RangeEnd: etcd.GetPrefix("txn_")},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "txn_key1", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "txn_value1"},
						},
					},
				},
				{
					request: &etcd.TxnRequest{
						Success: []etcd.Request{
							&etcd.PutRequest{Key: "txn_key2", Value: "txn_value0", PrevKv: true},
							&etcd.PutRequest{Key: "txn_key2", Value: "txn_value1", PrevKv: true},
						},
					},
					err: rpctypes.ErrGRPCDuplicateKey,
				},
				{
					request: &etcd.RangeRequest{Key: "txn_", RangeEnd: etcd.GetPrefix("txn_")},
					response: &etcd.RangeResponse{
						Count: 1,
						Kvs: []*etcd.KeyValue{
							{Key: "txn_key1", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "txn_value1"},
						},
					},
				},
			},
		},
		{
			name: "Recursive",
			testcases: []TestCase{
				{
					request: &etcd.TxnRequest{
						Compare: []etcd.Compare{},
						Success: []etcd.Request{
							&etcd.TxnRequest{
								Compare: []etcd.Compare{},
								Success: []etcd.Request{
									&etcd.RangeRequest{Key: "txn_", RangeEnd: etcd.GetPrefix("txn_")},
								},
								Failure: []etcd.Request{},
							},
						},
						Failure: []etcd.Request{},
					},
					response: &etcd.TxnResponse{
						Succeeded: true,
						Responses: []etcd.Response{
							&etcd.TxnResponse{
								Succeeded: true,
								Responses: []etcd.Response{
									&etcd.RangeResponse{
										Count: 1,
										Kvs: []*etcd.KeyValue{
											{Key: "txn_key1", ModRevision: 0, CreateRevision: 0, Version: 1, Value: "txn_value1"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TearDown",
			testcases: []TestCase{
				{
					request: &etcd.TxnRequest{
						Success: []etcd.Request{
							&etcd.DeleteRequest{Key: "txn_", RangeEnd: etcd.GetPrefix("txn_"), PrevKv: true},
						},
					},
					response: &etcd.TxnResponse{
						Succeeded: true,
						Responses: []etcd.Response{
							&etcd.DeleteResponse{
								Deleted: 1,
								PrevKvs: []*etcd.KeyValue{
									{Key: "txn_key1", ModRevision: -1, CreateRevision: -1, Version: 1, Value: "txn_value1"},
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, runTest(client, tc.testcases))
	}
}
