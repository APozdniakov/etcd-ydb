package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type RangeRequest struct {
	Key               string
	RangeEnd          string
	Limit             int64
	Revision          int64
	SortTarget        etcdserverpb.RangeRequest_SortTarget
	SortOrder         etcdserverpb.RangeRequest_SortOrder
	KeysOnly          bool
	CountOnly         bool
	MinModRevision    int64
	MaxModRevision    int64
	MinCreateRevision int64
	MaxCreateRevision int64
}

func (request RangeRequest) OrderByKey() *RangeRequest {
	request.SortTarget = etcdserverpb.RangeRequest_KEY
	return &request
}

func (request RangeRequest) OrderByModRevision() *RangeRequest {
	request.SortTarget = etcdserverpb.RangeRequest_MOD
	return &request
}

func (request RangeRequest) OrderByCreateRevision() *RangeRequest {
	request.SortTarget = etcdserverpb.RangeRequest_CREATE
	return &request
}

func (request RangeRequest) OrderByVersion() *RangeRequest {
	request.SortTarget = etcdserverpb.RangeRequest_VERSION
	return &request
}

func (request RangeRequest) OrderByValue() *RangeRequest {
	request.SortTarget = etcdserverpb.RangeRequest_VALUE
	return &request
}

func (request *RangeRequest) Ascending() *RangeRequest {
	request.SortOrder = etcdserverpb.RangeRequest_ASCEND
	return request
}

func (request *RangeRequest) Descending() *RangeRequest {
	request.SortOrder = etcdserverpb.RangeRequest_DESCEND
	return request
}

func (RangeRequest) Request() {}

func serializeRangeRequest(request *RangeRequest) *etcdserverpb.RangeRequest {
	if request == nil {
		return nil
	}
	return &etcdserverpb.RangeRequest{
		Key:               []byte(request.Key),
		RangeEnd:          []byte(request.RangeEnd),
		Limit:             request.Limit,
		Revision:          request.Revision,
		SortOrder:         request.SortOrder,
		SortTarget:        request.SortTarget,
		KeysOnly:          request.KeysOnly,
		CountOnly:         request.CountOnly,
		MinModRevision:    request.MinModRevision,
		MaxModRevision:    request.MaxModRevision,
		MinCreateRevision: request.MinCreateRevision,
		MaxCreateRevision: request.MaxCreateRevision,
	}
}

type RangeResponse struct {
	Revision int64
	Count    int64
	More     bool
	Kvs      []*KeyValue
}

func (RangeResponse) Response() {}

func (response RangeResponse) GetRevision() int64 {
	return response.Revision
}

func (RangeResponse) IsWrite() bool {
	return false
}

func deserializeRangeResponse(response *etcdserverpb.RangeResponse) *RangeResponse {
	if response == nil {
		return nil
	}
	result := &RangeResponse{
		Revision: response.Header.Revision,
		More:     response.More,
		Count:    response.Count,
		Kvs:      make([]*KeyValue, 0, len(response.Kvs)),
	}
	for _, kv := range response.Kvs {
		result.Kvs = append(result.Kvs, deserializeKeyValue(kv))
	}
	return result
}

func Range(client *Client, request *RangeRequest) (*RangeResponse, error) {
	response, err := client.Range(serializeRangeRequest(request))
	if err != nil {
		return nil, err
	}
	return deserializeRangeResponse(response), nil
}
