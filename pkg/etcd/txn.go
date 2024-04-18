package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type TxnRequest struct {
	Compare []Compare
	Success []Request
	Failure []Request
}

func (TxnRequest) Request() {}

func serializeRequestOp(request Request) *etcdserverpb.RequestOp {
	switch r := request.(type) {
	case *DeleteRequest:
		return &etcdserverpb.RequestOp{Request: &etcdserverpb.RequestOp_RequestDeleteRange{RequestDeleteRange: serializeDeleteRequest(r)}}
	case *PutRequest:
		return &etcdserverpb.RequestOp{Request: &etcdserverpb.RequestOp_RequestPut{RequestPut: serializePutRequest(r)}}
	case *RangeRequest:
		return &etcdserverpb.RequestOp{Request: &etcdserverpb.RequestOp_RequestRange{RequestRange: serializeRangeRequest(r)}}
	case *TxnRequest:
		return &etcdserverpb.RequestOp{Request: &etcdserverpb.RequestOp_RequestTxn{RequestTxn: serializeTxnRequest(r)}}
	default:
		panic("unknown request type")
	}
}

func serializeTxnRequest(request *TxnRequest) *etcdserverpb.TxnRequest {
	if request == nil {
		return nil
	}
	result := &etcdserverpb.TxnRequest{
		Compare: make([]*etcdserverpb.Compare, 0, len(request.Compare)),
		Success: make([]*etcdserverpb.RequestOp, 0, len(request.Success)),
		Failure: make([]*etcdserverpb.RequestOp, 0, len(request.Failure)),
	}
	for _, compare := range request.Compare {
		result.Compare = append(result.Compare, serializeCompare(compare))
	}
	for _, success := range request.Success {
		result.Success = append(result.Success, serializeRequestOp(success))
	}
	for _, failure := range request.Failure {
		result.Failure = append(result.Failure, serializeRequestOp(failure))
	}
	return result
}

type TxnResponse struct {
	Revision  int64
	Succeeded bool
	Responses []Response
}

func (TxnResponse) Response() {}

func (response TxnResponse) GetRevision() int64 {
	return response.Revision
}

func (response TxnResponse) IsWrite() bool {
	result := false
	for _, response := range response.Responses {
		result = result || response.IsWrite()
	}
	return result
}

func deserializeRequestOp(response *etcdserverpb.ResponseOp) Response {
	if deleteResponse := response.GetResponseDeleteRange(); deleteResponse != nil {
		return deserializeDeleteResponse(deleteResponse)
	} else if putResponse := response.GetResponsePut(); putResponse != nil {
		return deserializePutResponse(putResponse)
	} else if rangeResponse := response.GetResponseRange(); rangeResponse != nil {
		return deserializeRangeResponse(rangeResponse)
	} else if txnResponse := response.GetResponseTxn(); txnResponse != nil {
		return deserializeTxnResponse(txnResponse)
	} else {
		panic("unknown response type")
	}
}

func deserializeTxnResponse(response *etcdserverpb.TxnResponse) *TxnResponse {
	if response == nil {
		return nil
	}
	result := &TxnResponse{
		Revision:  response.Header.Revision,
		Succeeded: response.Succeeded,
		Responses: make([]Response, 0, len(response.Responses)),
	}
	for _, response := range response.Responses {
		result.Responses = append(result.Responses, deserializeRequestOp(response))
	}
	return result
}

func Txn(client *Client, request *TxnRequest) (*TxnResponse, error) {
	response, err := client.Txn(serializeTxnRequest(request))
	if err != nil {
		return nil, err
	}
	return deserializeTxnResponse(response), nil
}
