package gpudb

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func NewDefaultShowSystemStatusOptions() *ShowSystemStatusOptions {
	return &ShowSystemStatusOptions{}
}

func (gpudb *Gpudb) buildShowSystemStatusOptionsMap(
	ctx context.Context, options *ShowSystemStatusOptions) *map[string]string {
	var (
		// childCtx context.Context
		childSpan trace.Span
	)

	_, childSpan = gpudb.tracer.Start(ctx, "gpudb.buildShowSystemStatusOptionsMap()")
	defer childSpan.End()

	mapOptions := make(map[string]string)

	return &mapOptions
}

func (gpudb *Gpudb) ShowSystemStatusRaw(ctx context.Context) (*ShowSystemStatusResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.ShowSystemStatusRaw()")
	defer childSpan.End()

	return gpudb.ShowSystemStatusRawWithOpts(childCtx, NewDefaultShowSystemStatusOptions())
}

func (gpudb *Gpudb) ShowSystemStatusRawWithOpts(
	ctx context.Context, options *ShowSystemStatusOptions) (*ShowSystemStatusResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.ShowSystemStatusRawWithOpts()")
	defer childSpan.End()

	mapOptions := gpudb.buildShowSystemStatusOptionsMap(childCtx, options)
	response := ShowSystemStatusResponse{}
	request := ShowSystemStatusRequest{Options: *mapOptions}
	err := gpudb.submitRawRequest(
		childCtx, "/show/system/status",
		&Schemas.showSystemStatusRequest, &Schemas.showSystemStatusResponse,
		&request, &response)

	return &response, err
}
