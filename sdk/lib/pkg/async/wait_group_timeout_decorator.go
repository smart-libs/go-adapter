package task

//
//import (
//	"sync"
//	"time"
//)
//
//type TimeoutWaitGroupDecorator struct {
//	sync.WaitGroup
//}
//
//// WaitOrTimeout waits for the group Wait() completion and returns false or timeout and returns true
//func (w *TimeoutWaitGroupDecorator) WaitOrTimeout(timeout time.Duration) bool {
//	c := make(chan struct{})
//	go func() {
//		defer close(c)
//		w.WaitGroup.Wait()
//	}()
//	select {
//	case <-c:
//		return false // completed normally
//	case <-time.After(timeout):
//		return true // timed out
//	}
//}
