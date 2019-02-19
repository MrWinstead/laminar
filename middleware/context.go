package middleware

import (
	"context"
	"time"

	"github.com/lovoo/goka"

	"github.com/mrwinstead/laminar"
)

const (
	ContextKeyRawMessage = iota
)

var (
	_ laminar.MiddlewareContext = &Context{}
)

type Context struct {
	ctx        context.Context
	failError  error
	rawMessage []byte
	value      interface{}
	wrapped    goka.Context
}

func NewMiddlewareContext(gCtx goka.Context, rawMessage []byte) *Context {
	created := &Context{
		ctx:        gCtx.Context(),
		wrapped:    gCtx,
		rawMessage: rawMessage,
	}
	return created
}

func (ec *Context) Topic() goka.Stream {
	return ec.wrapped.Topic()
}

func (ec *Context) Key() string {
	return ec.wrapped.Key()
}

func (ec *Context) Partition() int32 {
	return ec.wrapped.Partition()
}

func (ec *Context) Offset() int64 {
	return ec.wrapped.Offset()
}

func (ec *Context) Value() interface{} {
	return ec.wrapped.Value()
}

func (ec *Context) SetValue(value interface{}) {
	ec.value = value
}

func (ec *Context) Delete() {
	ec.wrapped.Delete()
}

func (ec *Context) Timestamp() time.Time {
	return ec.wrapped.Timestamp()
}

func (ec *Context) Join(topic goka.Table) interface{} {
	return ec.wrapped.Join(topic)
}

func (ec *Context) Lookup(topic goka.Table, key string) interface{} {
	return ec.wrapped.Lookup(topic, key)
}

func (ec *Context) Emit(topic goka.Stream, key string, value interface{}) {
	ec.wrapped.Emit(topic, key, value)
}

func (ec *Context) Loopback(key string, value interface{}) {
	ec.wrapped.Loopback(key, value)
}

func (ec *Context) Fail(err error) {
	ec.failError = err
}

func (ec *Context) Context() context.Context {
	return ec.ctx
}

func (ec *Context) RawMessage() []byte {
	return ec.ctx.Value(ContextKeyRawMessage).([]byte)
}
