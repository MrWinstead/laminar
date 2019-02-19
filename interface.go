package laminar

import (
	"context"

	"github.com/lovoo/goka"
)

const (
	// ErrorHandlingFailProcessor instructs laminar to fully fail the processor
	// when the error is encountered
	ErrorHandlingFailProcessor = iota

	// ErrorHandlingContinue ignores the error from goka's perspective
	ErrorHandlingContinue

	// ErrorHandlingDeleteValue deletes the value which caused the error using
	// goka.Context.Delete()
	ErrorHandlingDeleteValue
)

// CodecBeforeDecodeFunc is called before the Decode function of a codec to
// enrich the context.Context or modify the raw message
type CodecBeforeDecodeFunc func(ctx context.Context) ([]byte, context.Context,
	error)

// CodecAfterEncodeFunc is called after the Encode function of a codec to
// modify the raw message
type CodecAfterEncodeFunc func(ctx context.Context, data []byte) ([]byte, error)

// ProcessingBeforeFunc is called before a ProcessingCallback and can enrich the
// process's context
type ProcessingBeforeFunc func(ctx MiddlewareContext, msg interface{}) (
	context.Context, error)

// ProcessingAfterFunc may be used to modify the value returned from a
// ProcessingCallback.
type ProcessingAfterFunc func(ctx MiddlewareContext, fromLog,
	toLog interface{}) (interface{}, error)

// ProcessingErrorHandler is notified when a processing error is encountered
// and may return ErrorHandling* options to instruct laminar of how to handle
// the error
type ProcessingErrorHandler func(ctx MiddlewareContext, err error) uint

// ProcessingFinalizer will always be called after all other functions after
// processing a single message
type ProcessingFinalizer func(ctx MiddlewareContext, err error, msg interface{})

// MiddlewareContext extends a goka.Context with elements useful for crafting
// middlewares
type MiddlewareContext interface {
	goka.Context
	RawMessage() []byte
}

// ProcessingHandler assembles Processing* functions together to be presented to
// goka
type ProcessingHandler interface {
	ProcessCallback() goka.ProcessCallback
}

// MiddlewareEnablingCodec expounds a Goka codec with additional processing
// contextual information
type MiddlewareEnablingCodec interface {
	DecodeEx(ctx context.Context, rawMessage []byte) (interface{}, error)
	EncodeEx(ctx context.Context, msg interface{}) ([]byte, error)
}
