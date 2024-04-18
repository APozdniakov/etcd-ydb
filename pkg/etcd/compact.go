package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type CompactRequest struct {
	Revision int64
}

func (CompactRequest) Request() {}

func serializeCompactRequest(request *CompactRequest) *etcdserverpb.CompactionRequest {
	if request == nil {
		return nil
	}
	return &etcdserverpb.CompactionRequest{
		Revision: request.Revision,
	}
}

type CompactResponse struct {
	Revision int64
}

func (CompactResponse) Response() {}

func (response CompactResponse) GetRevision() int64 {
	return response.Revision
}

func (CompactResponse) IsWrite() bool {
	return false
}

func deserializeCompactResponse(response *etcdserverpb.CompactionResponse) *CompactResponse {
	if response == nil {
		return nil
	}
	return &CompactResponse{
		Revision: response.Header.Revision,
	}
}

func Compact(client *Client, request *CompactRequest) (*CompactResponse, error) {
	response, err := client.Compact(serializeCompactRequest(request))
	if err != nil {
		return nil, err
	}
	return deserializeCompactResponse(response), nil
}
