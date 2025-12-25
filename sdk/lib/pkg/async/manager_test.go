package task

//
//import (
//	"context"
//	"github.com/stretchr/testify/assert"
//	"gitlab.com/route/b2b-core/shared/go-logger/pkg/logging"
//	"gitlab.com/route/b2b-core/shared/go-logger/pkg/logging/console"
//	"testing"
//	"time"
//)
//
//func Test_Manager_UsingStopTasks(t *testing.T) {
//	options := ManagerOptions{
//		ID:                    "XX",
//		StopTimeout:           time.Second,
//		StopCheckingIntervals: 2,
//	}
//	manager := NewManager(options)
//	//var logEntries []logging.LogEntry
//	//logger := logging.NewLoggerMock(logging.Config{QueueSize: -1}, func(entry logging.LogEntry) {
//	//	logEntries = append(logEntries, entry)
//	//})
//	logger := console.New(console.Config{
//		Config: logging.Config{
//			QueueSize: -1,
//		},
//		BufferSize: -1,
//	})
//	ctx := logging.NewContext(context.Background(), logger)
//	max := 3
//	for i := 0; i < max; i++ {
//		err := manager.StartTask(ctx, func(ctx context.Context, stopRequested func() bool) error {
//			for !stopRequested() {
//				logging.FromContext(ctx).InfoF("Waiting 500ms")
//				time.Sleep(200 * time.Millisecond)
//			}
//			return nil
//		})
//
//		assert.NoError(t, err)
//	}
//
//	tErrors, err := manager.StopTasks(ctx)
//	assert.NoError(t, err)
//	if assert.Equal(t, max, len(tErrors)) {
//		assert.NoError(t, tErrors[0])
//	}
//}
//
//func Test_Manager_UsingContextCancel(t *testing.T) {
//	options := ManagerOptions{
//		ID:                    "XX",
//		StopTimeout:           time.Second,
//		StopCheckingIntervals: 2,
//	}
//	manager := NewManager(options)
//
//	//var logEntries []logging.LogEntry
//	//logger := logging.NewLoggerMock(logging.Config{QueueSize: -1}, func(entry logging.LogEntry) {
//	//	logEntries = append(logEntries, entry)
//	//})
//	logger := console.New(console.Config{
//		Config: logging.Config{
//			QueueSize: -1,
//		},
//		BufferSize: -1,
//	})
//	ctx, cancel := context.WithCancel(context.Background())
//	ctx = logging.NewContext(ctx, logger)
//	max := 30
//	for i := 0; i < max; i++ {
//		err := manager.StartTask(ctx, func(ctx context.Context, stopRequested func() bool) error {
//			for !stopRequested() {
//				logging.FromContext(ctx).InfoF("Waiting 500ms")
//				time.Sleep(200 * time.Millisecond)
//			}
//			return nil
//		})
//
//		assert.NoError(t, err)
//	}
//
//	// Cancel first
//	cancel()
//	tErrors, err := manager.StopTasks(ctx)
//	assert.NoError(t, err)
//	if assert.Equal(t, max, len(tErrors)) {
//		assert.NoError(t, tErrors[0])
//	}
//}
