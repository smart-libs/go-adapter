package task

//
//type (
//	HandlerList []*Handler
//)
//
//func (h HandlerList) forEach(action func(h *Handler)) {
//	for _, handler := range h {
//		action(handler)
//	}
//}
//
//func (h HandlerList) cancelAll() { h.forEach(func(h *Handler) { h.cancel() }) }
//func (h HandlerList) stopAll()   { h.forEach(func(h *Handler) { h.stop() }) }
//
//func (h HandlerList) countRunning() (counter int) {
//	h.forEach(func(h *Handler) {
//		if h.isRunning() {
//			counter++
//		}
//	})
//	return
//}
//
//func (h HandlerList) getErrors() (result []error) {
//	h.forEach(func(h *Handler) { result = append(result, h.getError()) })
//	return
//}
