package context

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	keyTraceId = "trace_id"
)

type Context struct {
	context.Context
	TraceId uuid.UUID
}

func NewContext() *Context {
	traceId := uuid.New()
	ctx := context.WithValue(context.Background(), keyTraceId, traceId)

	return &Context{
		Context: ctx,
		TraceId: traceId,
	}
}

func FromGinC(ginC *gin.Context) *Context {
	traceId, ok := ginC.Get(keyTraceId)
	if !ok {
		traceId = uuid.New()
		ginC.Set(keyTraceId, traceId.(string))
	}
	parsedTraceId, err := uuid.Parse(traceId.(string))
	if err != nil {
		print(err)
	}
	return &Context{
		Context: ginC,
		TraceId: parsedTraceId,
	}
}
