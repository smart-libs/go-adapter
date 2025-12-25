package task

//
//import (
//	"context"
//	"fmt"
//	"sync"
//)
//
//type (
//	Handler struct {
//		dataLocker     sync.Mutex
//		id             string
//		state          State
//		stopFlag       bool
//		terminationErr error
//		taskCancelFunc func()
//		waitGroup      *sync.WaitGroup
//	}
//)
//
//func (h *Handler) onPanic(panicArg any)                     { h.setStopped(fmt.Errorf("%s: panic: %v", h.id, panicArg)) }
//func (h *Handler) getError() (err error)                    { h.doWithLock(func() { err = h.terminationErr }); return }
//func (h *Handler) stop()                                    { h.doWithLock(func() { h.stopFlag = true }) }
//func (h *Handler) cancel()                                  { h.doWithLock(func() { h.taskCancelFunc() }) }
//func (h *Handler) start(ctx context.Context, taskFunc Func) { go h.doStart(ctx, taskFunc) }
//func (h *Handler) isStopRequested() bool                    { return h.isWithLock(func() bool { return h.stopFlag }) }
//func (h *Handler) isStopped() bool                          { return h.isWithLock(func() bool { return h.state == stopped }) }
//
//func (h *Handler) setRunning() {
//	h.doWithLock(func() {
//		if h.state != running {
//			h.waitGroup.Add(1)
//			h.state = running
//		}
//	})
//}
//
//func (h *Handler) setStopped(err error) {
//	h.doWithLock(func() {
//		h.state = stopped
//		h.terminationErr = err
//		h.waitGroup.Done()
//	})
//}
//
//func (h *Handler) isWithLock(action func() bool) bool {
//	h.dataLocker.Lock()
//	defer h.dataLocker.Unlock()
//	return action()
//}
//
//func (h *Handler) doWithLock(action func()) {
//	h.dataLocker.Lock()
//	defer h.dataLocker.Unlock()
//	action()
//}
//
//func (h *Handler) isRunning() bool {
//	return h.isWithLock(func() bool { return h.state == running || h.state == stopping })
//}
//
//func (h *Handler) doStart(taskCtx context.Context, taskFunc Func) {
//	panicActionDecorator(taskCtx, h.id, h.onPanic,
//		func() {
//			h.setRunning()
//			h.setStopped(taskFunc(taskCtx, h.isStopRequested))
//		},
//	)
//}
