package etcd

import (
	"context"
	"os"

	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
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
	kv       etcdserverpb.KVClient
}

func NewClient(endpoint string) (*Client, error) {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))

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
		kv:       etcdserverpb.NewKVClient(conn),
	}, nil
}

func (client *Client) Range(ctx context.Context, request *etcdserverpb.RangeRequest) (*etcdserverpb.RangeResponse, error) {
	return client.kv.Range(ctx, request, client.callOpts...)
}

func (client *Client) Put(ctx context.Context, request *etcdserverpb.PutRequest) (*etcdserverpb.PutResponse, error) {
	return client.kv.Put(ctx, request, client.callOpts...)
}

func (client *Client) Delete(ctx context.Context, request *etcdserverpb.DeleteRangeRequest) (*etcdserverpb.DeleteRangeResponse, error) {
	return client.kv.DeleteRange(ctx, request, client.callOpts...)
}

func (client *Client) Txn(ctx context.Context, request *etcdserverpb.TxnRequest) (*etcdserverpb.TxnResponse, error) {
	return client.kv.Txn(ctx, request, client.callOpts...)
}

func (client *Client) Compact(ctx context.Context, request *etcdserverpb.CompactionRequest) (*etcdserverpb.CompactionResponse, error) {
	return client.kv.Compact(ctx, request, client.callOpts...)
}

func Do(ctx context.Context, client *Client, request Request) (Response, error) {
	switch r := request.(type) {
	case *CompactRequest:
		return Compact(ctx, client, r)
	case *DeleteRequest:
		return Delete(ctx, client, r)
	case *PutRequest:
		return Put(ctx, client, r)
	case *RangeRequest:
		return Range(ctx, client, r)
	case *TxnRequest:
		return Txn(ctx, client, r)
	default:
		panic("unknown request type")
	}
}
