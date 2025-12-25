package task

//
//import (
//	"context"
//	appevent "gitlab.com/route/b2b-core/shared/go-app/pkg/app/event"
//	apppanic "gitlab.com/route/b2b-core/shared/go-app/pkg/app/panic"
//)
//
//func panicActionDecorator(ctx context.Context, actionID string, panicAction func(any), action func()) {
//	// if there is an event manager, then it will be reported if this action panics.
//	eventManager := appevent.FromContext(ctx)
//	defer apppanic.Invoke(func(arg any, stack []byte) {
//		if eventManager != nil {
//			eventManager.Notify(appevent.Panic{Arg: arg, Stack: stack, Where: actionID})
//		}
//		panicAction(arg)
//	})()
//
//	action()
//}
