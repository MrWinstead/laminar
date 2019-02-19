package middleware

import (
	"github.com/go-kit/kit/log"
	"github.com/lovoo/goka"

	"github.com/mrwinstead/laminar"
)

type ProcessingHandlerOption func(mh *ProcessingHandler)

func ProcessingHandlerBefore(f laminar.ProcessingBeforeFunc) ProcessingHandlerOption {
	return func(mh *ProcessingHandler) {
		mh.beforeFuncs = append(mh.beforeFuncs, f)
	}
}

func ProcessingHandlerAfter(f laminar.ProcessingAfterFunc) ProcessingHandlerOption {
	return func(mh *ProcessingHandler) {
		mh.afterFuncs = append(mh.afterFuncs, f)
	}
}

func ProcessingHandlerError(f laminar.ProcessingErrorHandler) ProcessingHandlerOption {
	return func(mh *ProcessingHandler) {
		mh.errorHandler = f
	}
}

func ProcessingHandlerFinalizer(f laminar.ProcessingFinalizer) ProcessingHandlerOption {
	return func(mh *ProcessingHandler) {
		mh.finalizers = append(mh.finalizers, f)
	}
}

type ProcessingHandler struct {
	afterFuncs   []laminar.ProcessingAfterFunc
	beforeFuncs  []laminar.ProcessingBeforeFunc
	errorHandler laminar.ProcessingErrorHandler
	finalizers   []laminar.ProcessingFinalizer
	logger       log.Logger
	wrapped      goka.ProcessCallback
}

func NewProcessingHandler(wrapped goka.ProcessCallback, l log.Logger, opts ...ProcessingHandlerOption,
) laminar.ProcessingHandler {
	created := &ProcessingHandler{
		afterFuncs:  []laminar.ProcessingAfterFunc{},
		beforeFuncs: []laminar.ProcessingBeforeFunc{},
		finalizers:  []laminar.ProcessingFinalizer{},
		logger:      l,
		wrapped:     wrapped,
	}

	for _, opt := range opts {
		opt(created)
	}

	if nil == created.errorHandler {
		created.errorHandler = BuildHandleError(l)
	}
	return created
}

func (ph *ProcessingHandler) ProcessCallback() goka.ProcessCallback {
	return func(gCtx goka.Context, msgIface interface{}) {
		bridge := msgIface.(*codecProcessorBridge)
		ctx := bridge.ctx
		middlewareCtx := NewMiddlewareContext(gCtx,
			ctx.Value(ContextKeyRawMessage).([]byte))
		var encounteredErr error

		defer func() {
			recovered := recover()
			if recoveredAsErr, ok := recovered.(error); ok {
				encounteredErr = recoveredAsErr
			}

			if nil != encounteredErr {
				continueChoice := ph.errorHandler(middlewareCtx, encounteredErr)
				switch continueChoice {
				case laminar.ErrorHandlingContinue:
				case laminar.ErrorHandlingDeleteValue:
					gCtx.Delete()
				case laminar.ErrorHandlingFailProcessor:
					gCtx.Fail(middlewareCtx.failError)
				}
			}

			for _, finalizer := range ph.finalizers {
				finalizer(middlewareCtx, encounteredErr, bridge.message)
			}
		}()

		for _, before := range ph.beforeFuncs {
			middlewareCtx.ctx, encounteredErr = before(middlewareCtx,
				bridge.message)
			if nil != encounteredErr {
				return
			}
		}

		ph.wrapped(middlewareCtx, bridge.message)
		if nil != middlewareCtx.failError {
			return
		}

		for _, after := range ph.afterFuncs {
			middlewareCtx.value, encounteredErr = after(middlewareCtx,
				bridge.message, middlewareCtx.value)
			if nil != encounteredErr {
				return
			}
		}

		if nil != middlewareCtx.value {
			gCtx.SetValue(middlewareCtx.value)
		}
	}
}
