package middleware_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrwinstead/laminar/processing_context"
)

func TestNewEnhancedContext(t *testing.T) {
	created := processing_context.NewEnhancedContext(nil)
	assert.NotNil(t, created)
}
