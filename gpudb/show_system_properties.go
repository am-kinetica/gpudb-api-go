package gpudb

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func NewDefaultShowSystemPropertiesOptions() *ShowSystemPropertiesOptions {
	return &ShowSystemPropertiesOptions{
		Properties: "",
	}
}

func (gpudb *Gpudb) buildShowSystemPropertiesOptionsMap(ctx context.Context, options *ShowSystemPropertiesOptions) *map[string]string {
	var (
		// childCtx context.Context
		childSpan trace.Span
	)

	_, childSpan = gpudb.tracer.Start(ctx, "gpudb.buildShowSystemPropertiesOptionsMap()")
	defer childSpan.End()

	mapOptions := make(map[string]string)
	if options.Properties != "" {
		mapOptions["properties"] = options.Properties
	}

	return &mapOptions
}

func (gpudb *Gpudb) ShowSystemPropertiesRaw(ctx context.Context) (*ShowSystemPropertiesResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.ShowSystemPropertiesRaw()")
	defer childSpan.End()

	return gpudb.ShowSystemPropertiesRawWithOpts(childCtx, NewDefaultShowSystemPropertiesOptions())
}

func (gpudb *Gpudb) ShowSystemPropertiesRawWithOpts(
	ctx context.Context, options *ShowSystemPropertiesOptions) (*ShowSystemPropertiesResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.ShowSystemPropertiesRawWithOpts()")
	defer childSpan.End()

	mapOptions := gpudb.buildShowSystemPropertiesOptionsMap(childCtx, options)
	response := ShowSystemPropertiesResponse{}
	request := ShowSystemPropertiesRequest{Options: *mapOptions}
	err := gpudb.submitRawRequest(
		childCtx, "/show/system/properties",
		&Schemas.showSystemPropertiesRequest, &Schemas.showSystemPropertiesResponse,
		&request, &response)

	return &response, err
}
