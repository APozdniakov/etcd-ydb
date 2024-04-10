package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type CompactRequest struct {
	Revision int64
	Physical bool
}

func (CompactRequest) Request() {}

func serializeCompactRequest(request CompactRequest) *etcdserverpb.CompactionRequest {
	return &etcdserverpb.CompactionRequest{
		Revision: request.Revision,
		Physical: request.Physical,
	}
}

type CompactResponse struct {
	Revision int64
}

func (CompactResponse) Response() {}

func deserializeCompactResponse(response *etcdserverpb.CompactionResponse) *CompactResponse {
	if response == nil {
		return nil
	}
	return &CompactResponse{
		Revision: response.Header.Revision,
	}
}

func Compact(client *Client, request CompactRequest) (*CompactResponse, error) {
	response, err := client.Compact(serializeCompactRequest(request))
	if err != nil {
		return nil, err
	}
	return deserializeCompactResponse(response), nil
}
