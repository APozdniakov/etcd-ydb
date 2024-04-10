package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type SortOrder int32

const (
	NONE SortOrder = iota
	ASCEND
	DESCEND
)

func serializeSortOrder(sort_order SortOrder) (etcdserverpb.RangeRequest_SortOrder) {
	switch sort_order {
	case NONE: return etcdserverpb.RangeRequest_NONE
	case ASCEND: return etcdserverpb.RangeRequest_ASCEND
	case DESCEND: return etcdserverpb.RangeRequest_DESCEND
	default: panic("unknown etcdserverpb.RangeRequest_SortOrder")
	}
}

type SortTarget int32

const (
	KEY SortTarget = iota
	MOD
	CREATE
	VERSION
	VALUE
)

func serializeSortTarget(sort_target SortTarget) (etcdserverpb.RangeRequest_SortTarget) {
	switch sort_target {
	case KEY: return etcdserverpb.RangeRequest_KEY
	case MOD: return etcdserverpb.RangeRequest_MOD
	case CREATE: return etcdserverpb.RangeRequest_CREATE
	case VERSION: return etcdserverpb.RangeRequest_VERSION
	case VALUE: return etcdserverpb.RangeRequest_VALUE
	default: panic("unknown etcdserverpb.RangeRequest_SortTarget")
	}
}

type RangeRequest struct {
	Key string
	RangeEnd string
	Limit int64
	Revision int64
	SortOrder SortOrder
	SortTarget SortTarget
	KeysOnly bool
	CountOnly bool
	MinModRevision int64
	MaxModRevision int64
	MinCreateRevision int64
	MaxCreateRevision int64
}

func (RangeRequest) Request() {}

func serializeRangeRequest(request RangeRequest) *etcdserverpb.RangeRequest {
	return &etcdserverpb.RangeRequest{
		Key: []byte(request.Key),
		RangeEnd: []byte(request.RangeEnd),
		Limit: request.Limit,
		Revision: request.Revision,
		SortOrder: serializeSortOrder(request.SortOrder),
		SortTarget: serializeSortTarget(request.SortTarget),
		KeysOnly: request.KeysOnly,
		CountOnly: request.CountOnly,
		MinModRevision: request.MinModRevision,
		MaxModRevision: request.MaxModRevision,
		MinCreateRevision: request.MinCreateRevision,
		MaxCreateRevision: request.MaxCreateRevision,
	}
}

type RangeResponse struct {
	Revision int64
	Count int64
	More bool
	Kvs []*KeyValue
}

func (RangeResponse) Response() {}

func deserializeRangeResponse(response *etcdserverpb.RangeResponse) *RangeResponse {
	if response == nil {
		return nil
	}
	result := &RangeResponse{
		Revision: response.Header.Revision,
		More: response.More,
		Count: response.Count,
		Kvs: make([]*KeyValue, 0, len(response.Kvs)),
	}
	for _, kv := range response.Kvs {
		result.Kvs = append(result.Kvs, deserializeKeyValue(kv))
	}
	return result
}

func Range(client *Client, request RangeRequest) (*RangeResponse, error) {
	response, err := client.Range(serializeRangeRequest(request))
	if err != nil {
		return nil, err
	}
	return deserializeRangeResponse(response), nil
}
