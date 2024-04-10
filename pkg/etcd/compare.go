package etcd

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type Compare struct {
	Key string
	RangeEnd string
	// Result etcdserverpb.Compare_CompareResult
	// Target etcdserverpb.Compare_CompareTarget
	// TargetUnion etcdserverpb.isCompare_TargetUnion
}

func serializeCompare(compare Compare) *etcdserverpb.Compare {
	return &etcdserverpb.Compare{
		Key: []byte(compare.Key),
		RangeEnd: []byte(compare.RangeEnd),
	}
}

type CompareBuilder struct {
	cmp *etcdserverpb.Compare
}

func NewCompareBuilder() *CompareBuilder {
	return &CompareBuilder{
		cmp: &etcdserverpb.Compare{},
	}
}

func (builder *CompareBuilder) WithKey(key []byte) *CompareBuilder {
	builder.cmp.Key = key
	return builder
}

func (builder *CompareBuilder) WithRangeEnd(rangeEnd []byte) *CompareBuilder {
	builder.cmp.RangeEnd = rangeEnd
	return builder
}

func (builder *CompareBuilder) Equals() *CompareBuilder {
	builder.cmp.Result = etcdserverpb.Compare_EQUAL
	return builder
}

func (builder *CompareBuilder) Greater() *CompareBuilder {
	builder.cmp.Result = etcdserverpb.Compare_GREATER
	return builder
}

func (builder *CompareBuilder) Less() *CompareBuilder {
	builder.cmp.Result = etcdserverpb.Compare_LESS
	return builder
}

func (builder *CompareBuilder) NotEquals() *CompareBuilder {
	builder.cmp.Result = etcdserverpb.Compare_NOT_EQUAL
	return builder
}

func (builder *CompareBuilder) Version(version int64) *CompareBuilder {
	builder.cmp.Target = etcdserverpb.Compare_VERSION
	builder.cmp.TargetUnion = &etcdserverpb.Compare_Version{Version: version}
	return builder
}

func (builder *CompareBuilder) CreateRev(rev int64) *CompareBuilder {
	builder.cmp.Target = etcdserverpb.Compare_CREATE
	builder.cmp.TargetUnion = &etcdserverpb.Compare_CreateRevision{CreateRevision: rev}
	return builder
}

func (builder *CompareBuilder) ModRev(rev int64) *CompareBuilder {
	builder.cmp.Target = etcdserverpb.Compare_MOD
	builder.cmp.TargetUnion = &etcdserverpb.Compare_ModRevision{ModRevision: rev}
	return builder
}

func (builder *CompareBuilder) Value(value []byte) *CompareBuilder {
	builder.cmp.Target = etcdserverpb.Compare_VALUE
	builder.cmp.TargetUnion = &etcdserverpb.Compare_Value{Value: value}
	return builder
}

func (builder CompareBuilder) Build() *etcdserverpb.Compare {
	return builder.cmp
}
