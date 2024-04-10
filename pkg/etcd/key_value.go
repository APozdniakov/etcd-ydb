package etcd

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type KeyValue struct {
	Key string
	ModRevision int64
	CreateRevision int64
	Version int64
	Value string
}

func deserializeKeyValue(kv *mvccpb.KeyValue) *KeyValue {
	if kv == nil {
		return nil
	}
	return &KeyValue{
		Key: string(kv.Key),
		ModRevision: kv.ModRevision,
		CreateRevision: kv.CreateRevision,
		Version: kv.Version,
		Value: string(kv.Value),
	}
}
