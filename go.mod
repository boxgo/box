module github.com/boxgo/box

go 1.14

replace (
	github.com/coreos/bbolt v1.3.5 => go.etcd.io/coreos/bbolt v1.3.5
	github.com/coreos/etcd v3.3.22+incompatible => go.etcd.io/etcd v3.3.22+incompatible
)
