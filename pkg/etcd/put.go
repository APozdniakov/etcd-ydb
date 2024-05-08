package etcd

import (
	"context"

	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type PutRequest struct {
	Key         string
	Value       string
	PrevKv      bool
	IgnoreValue bool
}

func (PutRequest) Request() {}

func serializePutRequest(request *PutRequest) *etcdserverpb.PutRequest {
	if request == nil {
		return nil
	}
	return &etcdserverpb.PutRequest{
		Key:         []byte(request.Key),
		Value:       []byte(request.Value),
		PrevKv:      request.PrevKv,
		IgnoreValue: request.IgnoreValue,
	}
}

type PutResponse struct {
	Revision int64
	PrevKv   *KeyValue
}

func (PutResponse) Response() {}

func (response PutResponse) GetRevision() int64 {
	return response.Revision
}

func (PutResponse) IsWrite() bool {
	return true
}

func deserializePutResponse(response *etcdserverpb.PutResponse) *PutResponse {
	if response == nil {
		return nil
	}
	return &PutResponse{
		Revision: response.Header.Revision,
		PrevKv:   deserializeKeyValue(response.PrevKv),
	}
}

func Put(ctx context.Context, client *Client, request *PutRequest) (*PutResponse, error) {
	response, err := client.Put(ctx, serializePutRequest(request))
	if err != nil {
		return nil, err
	}
	return deserializePutResponse(response), nil
}
