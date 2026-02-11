package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithPartitionLeaseDuration(t *testing.T) {
	opts := NewQueueOptions(
		WithPartitionLeaseDuration(2 * time.Second),
	)

	require.Equal(t, 2*time.Second, opts.partitionLeaseDuration)
	require.Equal(t, 2*time.Second, opts.shadowPartitionLeaseDuration)
	require.Equal(t, 2*time.Second, opts.backlogNormalizeLeaseDuration)
}

func TestWithMinWorkersFree(t *testing.T) {
	opts := NewQueueOptions(
		WithMinWorkersFree(17),
	)

	require.Equal(t, int64(17), opts.minWorkersFree)

	opts = NewQueueOptions(
		WithMinWorkersFree(-2),
	)

	require.Equal(t, int64(0), opts.minWorkersFree)
}
