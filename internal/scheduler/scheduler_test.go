package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/test"

	"github.com/stretchr/testify/require"
)

func TestScheduler(t *testing.T) {
	t.Parallel()
	t.Run("Start", func(t *testing.T) {
		now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		setFixedTime(now)

		er := &mockEntryReader{
			Entries: []*entry{
				{
					Job:    &mockJob{},
					Next:   now,
					Logger: test.NewLogger(),
				},
				{
					Job:    &mockJob{},
					Next:   now.Add(time.Minute),
					Logger: test.NewLogger(),
				},
			},
		}

		schedulerInstance := newScheduler(newSchedulerArgs{
			EntryReader: er,
			LogDir:      testHomeDir,
			Logger:      test.NewLogger(),
		})

		go func() {
			_ = schedulerInstance.Start(context.Background())
		}()

		time.Sleep(time.Second + time.Millisecond*100)
		schedulerInstance.Stop()

		require.Equal(t, int32(1), er.Entries[0].Job.(*mockJob).RunCount.Load())
		require.Equal(t, int32(0), er.Entries[1].Job.(*mockJob).RunCount.Load())
	})
	t.Run("Restart", func(t *testing.T) {
		now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		setFixedTime(now)

		entryReader := &mockEntryReader{
			Entries: []*entry{
				{
					EntryType: entryTypeRestart,
					Job:       &mockJob{},
					Next:      now,
					Logger:    test.NewLogger(),
				},
			},
		}

		schedulerInstance := newScheduler(newSchedulerArgs{
			EntryReader: entryReader,
			LogDir:      testHomeDir,
			Logger:      test.NewLogger(),
		})

		go func() {
			_ = schedulerInstance.Start(context.Background())
		}()
		defer schedulerInstance.Stop()

		time.Sleep(time.Second + time.Millisecond*100)
		require.Equal(t, int32(1), entryReader.Entries[0].Job.(*mockJob).RestartCount.Load())
	})
	t.Run("NextTick", func(t *testing.T) {
		now := time.Date(2020, 1, 1, 1, 0, 50, 0, time.UTC)
		setFixedTime(now)
		schedulerInstance := newScheduler(newSchedulerArgs{
			EntryReader: &mockEntryReader{},
			LogDir:      testHomeDir,
			Logger:      test.NewLogger(),
		})
		next := schedulerInstance.nextTick(now)
		require.Equal(t, time.Date(2020, 1, 1, 1, 1, 0, 0, time.UTC), next)
	})
	t.Run("FixedTime", func(t *testing.T) {
		fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		setFixedTime(fixedTime)
		require.Equal(t, fixedTime, now())

		// Reset
		setFixedTime(time.Time{})
		require.NotEqual(t, fixedTime, now())
	})
}
