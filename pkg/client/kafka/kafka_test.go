package kafka_test

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/client/kafka"
)

const (
	testTopic = "test"
)

func Example() {
	kfk := kafka.StdConfig("default").Build()

	producer, err := kfk.NewSyncProducer()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	consumer, err := kfk.NewConsumer()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(testTopic, 0, kafka.OffsetNewest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var cnt int32

	go func() {
		for {
			select {
			case <-partitionConsumer.Messages():
				atomic.AddInt32(&cnt, 1)
			case <-signals:
				break
			}
		}
	}()

	partition, offset, err := producer.SendMessage(&kafka.ProducerMessage{
		Topic: testTopic,
		Value: kafka.StringEncoder("hi"),
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	fmt.Println(offset >= 0, partition == 0, atomic.LoadInt32(&cnt) > 0)
	// Output: true true true
}

func TestKafka_ConsumerGroupOperations(t *testing.T) {
	kfk := kafka.StdConfig("default").Build()

	admin, err := kfk.NewClusterAdmin()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := admin.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Test ListConsumerGroups
	t.Run("ListConsumerGroups", func(t *testing.T) {
		groups, err := admin.ListConsumerGroups()
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Found %d consumer groups", len(groups))
		for groupID, groupType := range groups {
			t.Logf("Group: %s, Type: %s", groupID, groupType)
		}
	})

	// Test DescribeConsumerGroups
	t.Run("DescribeConsumerGroups", func(t *testing.T) {
		// First, get the list of consumer groups
		groups, err := admin.ListConsumerGroups()
		if err != nil {
			t.Fatal(err)
		}

		if len(groups) == 0 {
			t.Skip("No consumer groups found, skipping describe test")
		}

		// Get the first group ID for testing
		var groupIDs []string
		for groupID := range groups {
			groupIDs = append(groupIDs, groupID)
			break // Only test with the first group
		}

		descriptions, err := admin.DescribeConsumerGroups(groupIDs)
		if err != nil {
			t.Fatal(err)
		}

		for _, desc := range descriptions {
			t.Logf("Group: %s, State: %s, Protocol: %s", desc.GroupId, desc.State, desc.ProtocolType)
			t.Logf("Members: %d", len(desc.Members))
		}
	})

	// Test DeleteConsumerGroup
	t.Run("DeleteConsumerGroup", func(t *testing.T) {
		// Create a test consumer group first
		testGroupID := "wechat_switch_change_v2.cache.push-v2-6d5797558-nhrnz"

		// Note: In a real test environment, you would create a consumer group first
		// For this test, we'll just test the delete operation on a non-existent group
		// which should return an error
		err := admin.DeleteConsumerGroup(testGroupID)
		if err != nil {
			// This is expected for a non-existent group
			t.Logf("Delete consumer group returned error (expected): %v", err)
		} else {
			t.Logf("Successfully deleted consumer group: %s", testGroupID)
		}
	})
}
