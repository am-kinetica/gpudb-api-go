package gpudb

import (
	"context"
	"strconv"

	"github.com/hamba/avro"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/multierr"
)

// NewDefaultInsertRecordsOptions
//
//	@return *InsertRecordsOptions
func NewDefaultInsertRecordsOptions() *InsertRecordsOptions {
	return &InsertRecordsOptions{
		UpdateOnExistingPk:     false,
		IgnoreExistingPk:       false,
		ReturnRecordIDs:        false,
		TruncateStrings:        false,
		ReturnIndividualErrors: false,
		AllowPartialBatch:      false,
	}
}

// InsertRecordsRaw //
//
//	@receiver gpudb
//	@param ctx
//	@param table
//	@param data
//	@return *InsertRecordsResponse
//	@return error
func (gpudb *Gpudb) InsertRecordsRaw(
	ctx context.Context,
	table string, data []interface{}) (*InsertRecordsResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.InsertRecordsRaw()")
	defer childSpan.End()

	return gpudb.InsertRecordsRawWithOpts(childCtx, table, data, NewDefaultInsertRecordsOptions())
}

// InsertRecordsRawWithOpts //
//
//	@receiver gpudb
//	@param ctx
//	@param table
//	@param data
//	@param options
//	@return *InsertRecordsResponse
//	@return error
func (gpudb *Gpudb) InsertRecordsRawWithOpts(
	ctx context.Context,
	table string, data []interface{}, options *InsertRecordsOptions) (*InsertRecordsResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
		errors    []error
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.InsertRecordsRawWithOpts()")
	defer childSpan.End()
	showTableResult, err := gpudb.ShowTableRawWithOpts(context.TODO(), table, &ShowTableOptions{
		ForceSynchronous:   true,
		GetSizes:           false,
		ShowChildren:       false, // needs to be false for tables
		NoErrorIfNotExists: false,
		GetColumnInfo:      true,
	})
	if err != nil {
		return nil, err
	}

	typeSchema := showTableResult.TypeSchemas[0]
	recordSchema, recordErr := avro.Parse(typeSchema)
	if recordErr != nil {
		return nil, err
	}

	// Convert the records into a byte array according to the Avro schema
	buffer := make([][]byte, len(data))
	for i := 0; i < len(data); i++ {
		buf, err := avro.Marshal(recordSchema, data[i])
		if err != nil {
			errors = append(errors, err)
		}
		buffer[i] = buf
	}

	gpudb.mutex.Lock()
	mapOptions := gpudb.buildInsertRecordsOptionsMap(childCtx, options)
	gpudb.mutex.Unlock()

	response := InsertRecordsResponse{}
	request := InsertRecordsRequest{
		TableName:    table,
		List:         buffer,
		ListString:   [][]string{},
		ListEncoding: "binary",
		Options:      *mapOptions,
	}
	err = gpudb.submitRawRequest(
		childCtx, "/insert/records",
		&Schemas.insertRecordsRequest, &Schemas.insertRecordsResponse,
		&request, &response)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return &response, nil
}

func (gpudb *Gpudb) buildInsertRecordsOptionsMap(ctx context.Context, options *InsertRecordsOptions) *map[string]string {
	var (
		// childCtx context.Context
		childSpan trace.Span
	)

	_, childSpan = gpudb.tracer.Start(ctx, "gpudb.buildInsertRecordsOptionsMap()")
	defer childSpan.End()

	mapOptions := make(map[string]string)

	mapOptions["update_on_existing_pk"] = strconv.FormatBool(options.UpdateOnExistingPk)
	mapOptions["ignore_existing_pk"] = strconv.FormatBool(options.IgnoreExistingPk)
	mapOptions["return_record_ids"] = strconv.FormatBool(options.ReturnRecordIDs)
	mapOptions["truncate_strings"] = strconv.FormatBool(options.TruncateStrings)
	mapOptions["return_individual_errors"] = strconv.FormatBool(options.ReturnIndividualErrors)
	mapOptions["allow_partial_batch"] = strconv.FormatBool(options.AllowPartialBatch)

	return &mapOptions
}

// InsertRecordsMap //
//
//	@receiver gpudb
//	@param ctx
//	@param table
//	@param data
//	@return *InsertRecordsResponse
//	@return error
func (gpudb *Gpudb) InsertRecordsMap(
	ctx context.Context,
	table string, data []interface{}) (*InsertRecordsResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.GetRecordsMap()")
	defer childSpan.End()

	return gpudb.InsertRecordsMapWithOpts(childCtx, table, data, NewDefaultInsertRecordsOptions())
}

// InsertRecordsMapWithOpts
//
//	@receiver gpudb
//	@param ctx
//	@param table
//	@param data
//	@param options
//	@return *InsertRecordsResponse
//	@return error
func (gpudb *Gpudb) InsertRecordsMapWithOpts(
	ctx context.Context,
	table string, data []interface{}, options *InsertRecordsOptions) (*InsertRecordsResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.InsertRecordsMapWithOpts()")
	defer childSpan.End()

	raw, err := gpudb.InsertRecordsRawWithOpts(childCtx, table, data, options)
	if err != nil {
		return nil, err
	}
	// start := time.Now()
	return raw, nil
	// TODO
}

// InsertRecordsStruct
//
//	@receiver gpudb
//	@param ctx
//	@param table
//	@param data
//	@return *InsertRecordsResponse
//	@return error
func (gpudb *Gpudb) InsertRecordsStruct(
	ctx context.Context,
	table string, data []interface{}) (*InsertRecordsResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.InsertRecordsStruct()")
	defer childSpan.End()

	return gpudb.InsertRecordsStructWithOpts(childCtx, table, data, NewDefaultInsertRecordsOptions())
}

// InsertRecordsStructWithOpts
//
//	@receiver gpudb
//	@param ctx
//	@param table
//	@param data
//	@param options
//	@return *InsertRecordsResponse
//	@return error
func (gpudb *Gpudb) InsertRecordsStructWithOpts(
	ctx context.Context,
	table string, data []interface{}, options *InsertRecordsOptions) (*InsertRecordsResponse, error) {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = gpudb.tracer.Start(ctx, "gpudb.InsertRecordsStructWithOpts()")
	defer childSpan.End()

	response, err := gpudb.InsertRecordsMapWithOpts(childCtx, table, data, options)
	if err != nil {
		return nil, err
	}

	return response, nil
}
