package middleware

import (
	"context"

	"github.com/lovoo/goka"
	"github.com/go-kit/kit/log"

	"github.com/mrwinstead/laminar"
)

func BuildHandleError(l log.Logger) laminar.ProcessingErrorHandler {
	return func(ctx context.Context, gCtx goka.Context, err error) uint {
		_ := l.Log(
			"Message", "error encountered during processing",
			"Error", err.Error(), "Key", gCtx.Key(), "Offset", gCtx.Offset(),
			"Topic", gCtx.Topic())
		return laminar.ErrorHandlingContinue
	}
}
