package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type PutRequest struct {
	Key string
	Value string
	PrevKv bool
	IgnoreValue bool
}

func (PutRequest) Request() {}

func serializePutRequest(request PutRequest) *etcdserverpb.PutRequest {
	return &etcdserverpb.PutRequest{
		Key: []byte(request.Key),
		Value: []byte(request.Value),
		PrevKv: request.PrevKv,
		IgnoreValue: request.IgnoreValue,
	}
}

type PutResponse struct {
	Revision int64
	PrevKv *KeyValue
}

func (PutResponse) Response() {}

func deserializePutResponse(response *etcdserverpb.PutResponse) *PutResponse {
	if response == nil {
		return nil
	}
	return &PutResponse{
		Revision: response.Header.Revision,
		PrevKv: deserializeKeyValue(response.PrevKv),
	}
}

func Put(client *Client, request PutRequest) (*PutResponse, error) {
	response, err := client.Put(serializePutRequest(request))
	if err != nil {
		return nil, err
	}
	return deserializePutResponse(response), nil
}
