package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type DeleteRequest struct {
	Key      string
	RangeEnd string
	PrevKv   bool
}

func (DeleteRequest) Request() {}

func serializeDeleteRequest(request *DeleteRequest) *etcdserverpb.DeleteRangeRequest {
	if request == nil {
		return nil
	}
	return &etcdserverpb.DeleteRangeRequest{
		Key:      []byte(request.Key),
		RangeEnd: []byte(request.RangeEnd),
		PrevKv:   request.PrevKv,
	}
}

type DeleteResponse struct {
	Revision int64
	Deleted  int64
	PrevKvs  []*KeyValue
}

func (DeleteResponse) Response() {}

func (response DeleteResponse) GetRevision() int64 {
	return response.Revision
}

func (response DeleteResponse) IsWrite() bool {
	return response.Deleted != 0
}

func deserializeDeleteResponse(response *etcdserverpb.DeleteRangeResponse) *DeleteResponse {
	if response == nil {
		return nil
	}
	result := &DeleteResponse{
		Revision: response.Header.Revision,
		Deleted:  response.Deleted,
		PrevKvs:  make([]*KeyValue, 0, len(response.PrevKvs)),
	}
	for _, prev_kv := range response.PrevKvs {
		result.PrevKvs = append(result.PrevKvs, deserializeKeyValue(prev_kv))
	}
	return result
}

func Delete(client *Client, request *DeleteRequest) (*DeleteResponse, error) {
	response, err := client.Delete(serializeDeleteRequest(request))
	if err != nil {
		return nil, err
	}
	return deserializeDeleteResponse(response), nil
}
