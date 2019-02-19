package middleware

import "context"

type codecProcessorBridge struct {
	ctx     context.Context
	message interface{}
}
