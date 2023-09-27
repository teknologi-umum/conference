package main

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5"
)

type PGXTracer struct{}

func (t PGXTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	span := sentry.StartSpan(ctx, "pgx.query", sentry.WithTransactionName("PGX TraceQuery"))
	if span == nil {
		return ctx
	}

	span.SetContext("data", sentry.Context{
		"sql":  data.SQL,
		"args": data.Args,
	})

	return span.Context()
}

func (t PGXTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	span := sentry.SpanFromContext(ctx)
	if span == nil {
		return
	}

	span.SetData("command_tag", data.CommandTag.String())

	if data.Err != nil {
		span.Status = sentry.SpanStatusInternalError
		span.SetData("error", data.Err.Error())
	}

	span.Finish()
}
