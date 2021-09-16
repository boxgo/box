package pubsub_test

import (
	"sync/atomic"
	"testing"

	"github.com/boxgo/box/pkg/pubsub"
)

func Benchmark_PubSub(b *testing.B) {
	var cnt1, cnt2, cnt3, total int64

	eventbus := pubsub.New(
		pubsub.WithBufferSize(1000),
	)
	eventbus.Subscribe("topic1", func(data pubsub.Message) {
		atomic.AddInt64(&cnt1, 1)
	})
	eventbus.Subscribe("topic1", func(data pubsub.Message) {
		atomic.AddInt64(&cnt2, 1)
	})
	eventbus.Subscribe("topic2", func(data pubsub.Message) {
		atomic.AddInt64(&cnt3, 1)
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt64(&total, 1)
			eventbus.Publish("topic1", "1234")
		}
	})

	eventbus.Wait()

	if total != cnt1 {
		b.Fatalf("Publish: %d, topic1 handler1: %d", total, cnt1)
	}
	if total != cnt2 {
		b.Fatalf("Publish: %d, topic1 handler2: %d", total, cnt2)
	}
	if cnt3 != 0 {
		b.Fatalf("Publish: %d, topic2 handler1: %d", total, cnt3)
	}
}
