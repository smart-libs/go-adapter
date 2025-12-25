package task

//
//import (
//	"context"
//	"gitlab.com/route/b2b-core/shared/go-logger/pkg/logging"
//)
//
//func stopMessageActionDecorator(ctx context.Context, actionID string, shutdownAction func(), mainAction func()) {
//	logger := logging.FromContext(ctx)
//	defer func() {
//		logger.InfoF("%s: stop request was received!", actionID)
//		shutdownAction()
//		logger.InfoF("%s: stop completed!", actionID)
//	}()
//
//	mainAction()
//}
