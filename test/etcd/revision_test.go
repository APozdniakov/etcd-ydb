package etcd_test

var (
	revision *int64
)

func initRevision(actual int64) {
	if revision == nil {
		revision = &actual
	}
}
