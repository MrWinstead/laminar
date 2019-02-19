package middleware

import (
	"context"

	"github.com/lovoo/goka"

	"github.com/mrwinstead/laminar"
)

var (
	_ goka.Codec = &Codec{}
)

type CodecOption func(i *Codec)

func WithCodecBeforeDecodeFunc(before laminar.CodecBeforeDecodeFunc) CodecOption {
	return func(i *Codec) {
		i.beforeCodecFuncs = append(i.beforeCodecFuncs, before)
	}
}

func WithCodecAfterEncodeFunc(after laminar.CodecAfterEncodeFunc) CodecOption {
	return func(i *Codec) {
		i.afterCodecFuncs = append(i.afterCodecFuncs, after)
	}
}

type Codec struct {
	rootContext  context.Context
	wrappedCodec laminar.MiddlewareEnablingCodec

	beforeCodecFuncs []laminar.CodecBeforeDecodeFunc
	afterCodecFuncs  []laminar.CodecAfterEncodeFunc
}

func NewCodec(rootContext context.Context, userCodec laminar.MiddlewareEnablingCodec,
	opts ...CodecOption) goka.Codec {
	created := &Codec{
		rootContext:      rootContext,
		wrappedCodec:     userCodec,
		beforeCodecFuncs: []laminar.CodecBeforeDecodeFunc{},
		afterCodecFuncs:  []laminar.CodecAfterEncodeFunc{},
	}

	for _, opt := range opts {
		opt(created)
	}

	return created
}

func (i *Codec) Decode(data []byte) (interface{}, error) {
	var beforeFuncErr error
	processedData := data
	ctx := context.WithValue(i.rootContext, ContextKeyRawMessage, data)
	for _, beforeFunc := range i.beforeCodecFuncs {
		processedData, ctx, beforeFuncErr = beforeFunc(ctx)
		if nil != beforeFuncErr {
			return nil, beforeFuncErr
		}
	}

	decoded, decodeErr := i.wrappedCodec.DecodeEx(ctx, processedData)
	if nil != decodeErr {
		return nil, decodeErr
	}

	bridge := &codecProcessorBridge{
		ctx:     ctx,
		message: decoded,
	}

	return bridge, nil
}

func (i *Codec) Encode(valueIface interface{}) ([]byte, error) {
	value := valueIface.(*codecProcessorBridge)

	serialized, encodeErr := i.wrappedCodec.EncodeEx(value.ctx, value.message)
	if nil != encodeErr {
		return nil, encodeErr
	}

	var afterFuncErr error
	for _, afterFunc := range i.afterCodecFuncs {
		serialized, afterFuncErr = afterFunc(value.ctx, serialized)
		if nil != afterFuncErr {
			return nil, afterFuncErr
		}
	}

	return serialized, nil
}
