package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type Compare struct {
	Key            string
	Result         etcdserverpb.Compare_CompareResult
	ModRevision    *int64
	CreateRevision *int64
	Version        *int64
	Value          *string
}

func (compare Compare) Equal() Compare {
	compare.Result = etcdserverpb.Compare_EQUAL
	return compare
}

func (compare Compare) Greater() Compare {
	compare.Result = etcdserverpb.Compare_GREATER
	return compare
}

func (compare Compare) Less() Compare {
	compare.Result = etcdserverpb.Compare_LESS
	return compare
}

func (compare Compare) NotEqual() Compare {
	compare.Result = etcdserverpb.Compare_NOT_EQUAL
	return compare
}

func (compare Compare) SetModRevision(modRevision int64) Compare {
	compare.ModRevision = &modRevision
	return compare
}

func (compare Compare) SetCreateRevision(createRevision int64) Compare {
	compare.CreateRevision = &createRevision
	return compare
}

func (compare Compare) SetVersion(version int64) Compare {
	compare.Version = &version
	return compare
}

func (compare Compare) SetValue(value string) Compare {
	compare.Value = &value
	return compare
}

func serializeCompare(compare Compare) *etcdserverpb.Compare {
	result := &etcdserverpb.Compare{
		Key:      []byte(compare.Key),
		Result:   compare.Result,
	}
	switch {
	case compare.ModRevision != nil:
		result.Target = etcdserverpb.Compare_MOD
		result.TargetUnion = &etcdserverpb.Compare_ModRevision{ModRevision: *compare.ModRevision}
	case compare.CreateRevision != nil:
		result.Target = etcdserverpb.Compare_CREATE
		result.TargetUnion = &etcdserverpb.Compare_CreateRevision{CreateRevision: *compare.CreateRevision}
	case compare.Version != nil:
		result.Target = etcdserverpb.Compare_VERSION
		result.TargetUnion = &etcdserverpb.Compare_Version{Version: *compare.Version}
	case compare.Value != nil:
		result.Target = etcdserverpb.Compare_VALUE
		result.TargetUnion = &etcdserverpb.Compare_Value{Value: []byte(*compare.Value)}
	default:
		panic("expected one of compare target")
	}
	return result
}
