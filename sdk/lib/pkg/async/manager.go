package task

//
//import (
//	"context"
//	"fmt"
//	"gitlab.com/route/b2b-core/shared/go-logger/pkg/logging"
//	"os"
//	"sync"
//	"sync/atomic"
//	"time"
//)
//
//type (
//	// Func specifies the signature a task should provide to be managed.
//	// The task should use the stopRequested() method as the main indicator the task should stop or not. It also shall pass the
//	// given context to all calls that demand one.
//	// The manager will change the stopRequested() to true as the first alternative to stop the task. Next the manager
//	// triggers a timer that once finished will cancel the task context as a second alternative to stop it.
//	Func func(ctx context.Context, stopRequested func() bool) error
//
//	Manager interface {
//		// StartTask creates a GO routine to execute the given taskFunc argument. The context provided in the call
//		// will be the parent of the new one created for the task which means if the parent is cancelled the task context
//		// will be also cancelled.
//		StartTask(ctx context.Context, taskFunc Func) error
//		// StopTasks stops all running tasks and returns the error returned by each taskFunc started. The second output
//		// argument returns any error occurred in the StopTasks method.
//		StopTasks(ctx context.Context) (tasksErrors []error, stopError error)
//	}
//
//	ManagerOptions struct {
//		// ID to be used to identify the Manager in the log
//		ID string
//		// time to wait for graceful shutdown before abandoning go routines
//		StopTimeout time.Duration
//		// number of intervals to check all tasks have stopped (interval time = StopTimeout / StopCheckingIntervals)
//		StopCheckingIntervals int
//	}
//
//	defaultManager struct {
//		options ManagerOptions
//		// managementStarted is used to identify whether the main loop is running or not
//		managementStarted atomic.Bool
//		stopFlag          bool
//		dataLocker        sync.Mutex
//		TimeoutWaitGroupDecorator
//		// shutdownChannel is used by StopAllTasks() to stop tasks and management
//		shutdownChannel chan os.Signal
//		managedTasks    HandlerList
//	}
//)
//
//const (
//	stopByContextCancellation = 1
//	stopByShutdownMsg         = 2
//	stopByStopTasks           = 3
//)
//
//func NewManager(options ManagerOptions) Manager {
//	return &defaultManager{options: options}
//}
//
//func (m *defaultManager) StopTasks(ctx context.Context) (tasksErrors []error, stopError error) {
//	logger := logging.FromContext(ctx)
//	panicActionDecorator(ctx, m.options.ID, func(a any) {
//		logger.InfoF("%s: panic: %v", m.options.ID, a)
//	}, func() {
//		tasksErrors, stopError = m.doStopTasks(ctx)
//	})
//	return
//}
//
//func (m *defaultManager) doStopTasks(ctx context.Context) (tasksErrors []error, stopError error) {
//	if m.managementStarted.Load() == false {
//		return nil, fmt.Errorf("%s: not started", m.options.ID)
//	}
//
//	const defaultIntervals = 3
//	const defaultStopTimeout = 5 * time.Second
//
//	logger := logging.FromContext(ctx)
//	calcIntervalWaitTime := func(stopTimeout time.Duration, intervals int) time.Duration {
//		return time.Duration((stopTimeout.Milliseconds() / int64(intervals)) * int64(time.Millisecond))
//	}
//	waitForTaskCompletion := func() {
//		stopTimeout := coalesce(m.options.StopTimeout, defaultStopTimeout)
//		intervals := coalesce(m.options.StopCheckingIntervals, defaultIntervals)
//		intervalWaitTime := calcIntervalWaitTime(stopTimeout, intervals)
//		timedOut := true
//		runningCounter := m.managedTasks.countRunning()
//		for ; runningCounter > 0 && timedOut && intervals > 0; runningCounter = m.managedTasks.countRunning() {
//			intervals--
//			logger.InfoF("%s: running instances=[%d], waiting %s for next check %d",
//				m.options.ID, runningCounter, intervalWaitTime, intervals)
//			timedOut = m.WaitOrTimeout(intervalWaitTime) // false when all tasks call Done()
//		}
//
//		logger.InfoF("%s: Cancelling tasks ********************", m.options.ID)
//		m.managedTasks.cancelAll()
//		finalMsg := "%s: stop succeeded, remaining running instances=[%d] of total=[%d]"
//		if intervals == 0 && timedOut {
//			finalMsg = "%s: stop timed out, remaining running instances=[%d] of total=[%d]"
//		}
//
//		logger.InfoF(finalMsg, m.options.ID, runningCounter, len(m.managedTasks))
//	}
//
//	logger.InfoF("%s: m.stop() ********************", m.options.ID)
//	m.stop(stopByStopTasks)
//	logger.InfoF("%s: doWithLock ********************", m.options.ID)
//	m.doWithLock(func() {
//		logger.InfoF("%s: Waiting tasks ********************", m.options.ID)
//		waitForTaskCompletion()
//		logger.InfoF("%s: Getting errors ********************", m.options.ID)
//		tasksErrors = m.managedTasks.getErrors()
//	})
//	logger.InfoF("%s: Leaving ********************", m.options.ID)
//	return
//}
//
//func (m *defaultManager) doWithLock(action func()) {
//	m.dataLocker.Lock()
//	defer m.dataLocker.Unlock()
//	action()
//}
//
//func (m *defaultManager) start(ctx context.Context) {
//	// if already started return
//	start := m.managementStarted.CompareAndSwap(false, true)
//	if !start {
//		return
//	}
//	m.doWithLock(func() {
//		m.stopFlag = false
//		m.shutdownChannel = make(chan os.Signal, 1) // cap = 1 to avoid deadlock
//	})
//	doNothingOnPanic := func(any) {}
//	go panicActionDecorator(ctx, m.options.ID, doNothingOnPanic, func() {
//		doNothingOnShutdown := func() {}
//		stopMessageActionDecorator(ctx, fmt.Sprintf("%s:MainLoop", m.options.ID), doNothingOnShutdown, func() {
//			// Wait for signals to stop
//			for {
//				select {
//				case <-ctx.Done():
//					m.stop(stopByContextCancellation)
//					return
//				case <-m.shutdownChannel:
//					m.stop(stopByShutdownMsg)
//					return
//				}
//			}
//		})
//	})
//}
//
//func (m *defaultManager) createCtxLogger(ctx context.Context, id int) (context.Context, logging.Logger) {
//	ctxVar := map[string]any{
//		m.options.ID: id,
//	}
//	newCtx, newLogger := logging.NewContextLoggerDecorator(ctx, ctxVar, logging.FromContext(ctx), context.WithValue)
//	// this is performed as WA for a log issue that uses 2 different context key to retrieve logger
//	return context.WithValue(newCtx, logging.LoggerKey, newLogger), newLogger
//}
//
//func (m *defaultManager) StartTask(ctx context.Context, taskFunc Func) error {
//	if taskFunc == nil {
//		return fmt.Errorf("%s: invalid nil Func", m.options.ID)
//	}
//	if m.managementStarted.Load() == false {
//		m.start(ctx)
//	}
//
//	createTaskID := func(id int) string { return fmt.Sprintf("%s[%d]", m.options.ID, id) }
//
//	nextID := len(m.managedTasks) + 1
//	handler := Handler{
//		id:        createTaskID(nextID),
//		waitGroup: &m.WaitGroup,
//	}
//	taskCtx, _ := m.createCtxLogger(ctx, nextID)
//	taskCtx, handler.taskCancelFunc = context.WithCancel(taskCtx)
//
//	m.doWithLock(func() {
//		m.managedTasks = append(m.managedTasks, &handler)
//	})
//
//	handler.start(taskCtx, taskFunc)
//	logging.FromContext(ctx).InfoF("%s: Started instance=[%d]", m.options.ID, nextID)
//	return nil
//}
//
//func coalesce[T comparable](v1, v2 T) T {
//	var zero T
//	if v1 == zero {
//		return v2
//	}
//	return v1
//}
//
//func (m *defaultManager) stop(reason int) {
//	m.doWithLock(func() {
//		if !m.stopFlag {
//			m.stopFlag = true
//			if reason != stopByContextCancellation {
//				m.shutdownChannel <- os.Kill
//			}
//			m.managedTasks.stopAll()
//		}
//	})
//}
//
//var _ Manager = &defaultManager{}
