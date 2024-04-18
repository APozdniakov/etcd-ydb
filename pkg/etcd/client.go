package etcd

import (
	"context"

	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Request interface {
	Request()
}

type Response interface {
	Response()
	GetRevision() int64
	IsWrite() bool
}

var defaultCallOpts = []grpc.CallOption{
	grpc.WaitForReady(true),
	grpc.MaxCallSendMsgSize(2 * 1024 * 1024),
	grpc.MaxCallRecvMsgSize(4 * 1024 * 1024),
}

type Client struct {
	endpoint string
	callOpts []grpc.CallOption
	ctx      context.Context
	kv       etcdserverpb.KVClient
}

func NewClient(endpoint string) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		endpoint: endpoint,
		callOpts: defaultCallOpts,
		ctx:      context.Background(),
		kv:       etcdserverpb.NewKVClient(conn),
	}, nil
}

func (client *Client) Range(request *etcdserverpb.RangeRequest) (*etcdserverpb.RangeResponse, error) {
	return client.kv.Range(client.ctx, request, client.callOpts...)
}

func (client *Client) Put(request *etcdserverpb.PutRequest) (*etcdserverpb.PutResponse, error) {
	return client.kv.Put(client.ctx, request, client.callOpts...)
}

func (client *Client) Delete(request *etcdserverpb.DeleteRangeRequest) (*etcdserverpb.DeleteRangeResponse, error) {
	return client.kv.DeleteRange(client.ctx, request, client.callOpts...)
}

func (client *Client) Txn(request *etcdserverpb.TxnRequest) (*etcdserverpb.TxnResponse, error) {
	return client.kv.Txn(client.ctx, request, client.callOpts...)
}

func (client *Client) Compact(request *etcdserverpb.CompactionRequest) (*etcdserverpb.CompactionResponse, error) {
	return client.kv.Compact(client.ctx, request, client.callOpts...)
}

func Do(client *Client, request Request) (Response, error) {
	switch r := request.(type) {
	case *CompactRequest:
		return Compact(client, r)
	case *DeleteRequest:
		return Delete(client, r)
	case *PutRequest:
		return Put(client, r)
	case *RangeRequest:
		return Range(client, r)
	case *TxnRequest:
		return Txn(client, r)
	default:
		panic("unknown request type")
	}
}
