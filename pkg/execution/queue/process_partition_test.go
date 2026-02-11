package queue

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestShouldProcessPartitionInParallel(t *testing.T) {
	fnID := uuid.New()
	accountID := uuid.New()

	t.Run("function override enables parallel", func(t *testing.T) {
		opts := NewQueueOptions(
			WithDisableFifoForFunctions(map[string]struct{}{
				fnID.String(): {},
			}),
		)
		processor := &queueProcessor{QueueOptions: opts}
		partition := &QueuePartition{
			ID:         fnID.String(),
			FunctionID: &fnID,
			AccountID:  accountID,
		}

		require.True(t, processor.shouldProcessPartitionInParallel(partition))
	})

	t.Run("account override enables parallel", func(t *testing.T) {
		opts := NewQueueOptions(
			WithDisableFifoForAccounts(map[string]struct{}{
				accountID.String(): {},
			}),
		)
		processor := &queueProcessor{QueueOptions: opts}
		partition := &QueuePartition{
			ID:         fnID.String(),
			FunctionID: &fnID,
			AccountID:  accountID,
		}

		require.True(t, processor.shouldProcessPartitionInParallel(partition))
	})

	t.Run("system queue mapping enables parallel", func(t *testing.T) {
		queueName := "system-queue"
		opts := NewQueueOptions(
			WithKindToQueueMapping(map[string]string{
				queueName: "some-kind",
			}),
		)
		processor := &queueProcessor{QueueOptions: opts}
		partition := &QueuePartition{
			ID:        queueName,
			QueueName: &queueName,
			AccountID: accountID,
		}

		require.True(t, processor.shouldProcessPartitionInParallel(partition))
	})

	t.Run("default partition remains fifo", func(t *testing.T) {
		opts := NewQueueOptions()
		processor := &queueProcessor{QueueOptions: opts}
		partition := &QueuePartition{
			ID:         fnID.String(),
			FunctionID: &fnID,
			AccountID:  accountID,
		}

		require.False(t, processor.shouldProcessPartitionInParallel(partition))
	})
}

func TestQueueOptionOverrides(t *testing.T) {
	t.Run("partition lease duration overrides all lease durations", func(t *testing.T) {
		opts := NewQueueOptions(WithPartitionLeaseDuration(2500 * time.Millisecond))
		require.Equal(t, 2500*time.Millisecond, opts.partitionLeaseDuration)
		require.Equal(t, 2500*time.Millisecond, opts.shadowPartitionLeaseDuration)
		require.Equal(t, 2500*time.Millisecond, opts.backlogNormalizeLeaseDuration)
	})

	t.Run("min workers free can be overridden", func(t *testing.T) {
		opts := NewQueueOptions(WithMinWorkersFree(2))
		require.Equal(t, int64(2), opts.minWorkersFree)
	})
}
