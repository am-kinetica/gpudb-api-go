package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/am-kinetica/gpudb-api-go/example"
	"github.com/am-kinetica/gpudb-api-go/gpudb"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func main() {
	endpoint := "https://172.17.0.2:8082/gpudb-0" //os.Args[1]
	username := "admin"                           //os.Args[2]
	password := "Kinetica1."                      //os.Args[3]
	Logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.TODO()
	options := gpudb.GpudbOptions{Username: username, Password: password, ByPassSslCertCheck: true}
	// fmt.Println("Options", options)
	gpudbInst := gpudb.NewWithOptions(ctx, endpoint, &options)

	var anyValue []any
	intSlice := []int{1, 2, 3, 4}
	anyValue = append(anyValue, intSlice)

	fmt.Println(anyValue...)

	// runShowExplainVerboseAnalyseSqlStatement(gpudbInst)
	// runShowExplainVerboseSqlStatement(gpudbInst)
	// runShowStoredProcedureDDL(gpudbInst)
	// runShowTableDDL(gpudbInst)
	// runShowTableResourcesBySchema1(gpudbInst)
	// runAggregateGroupBy1(gpudbInst)
	// runAggregateGroupBy2(gpudbInst)
	// runAggregateGroupBy3(gpudbInst)
	// runAggregateStatistics1(gpudbInst)
	// runAggregateStatistics2(gpudbInst)
	// runShowResourceGroups1(gpudbInst)
	// runShowSecurity1(gpudbInst)
	// runShowResourceStatistics1(gpudbInst)
	// runShowResourceStatistics2(gpudbInst)
	// runShowResourceStatistics3(gpudbInst)
	// runShowSqlProc1(gpudbInst)
	// runShowSqlProc2(gpudbInst)
	// runShowSystemStatus1(gpudbInst)
	// runShowSystemTiming1(gpudbInst)
	// runShowSystemProperties1(gpudbInst)
	// runExecuteSql1(gpudbInst)
	// runExecuteSql2(gpudbInst)
	// runExecuteSql3(gpudbInst)
	// runExecuteSql4(gpudbInst)
	// runExecuteSql5(gpudbInst)
	// runShowSchema1(gpudbInst)
	// runShowTable1(gpudbInst)
	// runGetRecords1(gpudbInst)
	// runGetRecords2(gpudbInst)
	// runGetRecords3(gpudbInst)
	// runGetRecords4(gpudbInst)
	runInsertRecords(Logger, gpudbInst)
	// runCreateResourceGroup(gpudbInst, "lucid_test")
	// runDeleteResourceGroup(gpudbInst, "lucid_test")
	// runCreateSchema(gpudbInst, "lucid_test")
	// runDropSchema(gpudbInst, "lucid_test")

}

func runShowExplainVerboseAnalyseSqlStatement(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	// string field
	result, err := gpudbInst.ShowExplainVerboseAnalyseSqlStatement(context.TODO(), `select * from 
(select vendor_id, payment_type, sum(fare_amount) as sum_fare1 from demo.nyctaxi group by vendor_id, payment_type) a
inner join
(select vendor_id, rate_code_id, sum(trip_distance) as sum_trip1 from demo.nyctaxi group by vendor_id, rate_code_id) b
on a.vendor_id = b.vendor_id
`)
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	//fmt.Println(*result)
	for _, plan := range *result.Plans {
		fmt.Println(plan)
		requestMap, _ := plan.JsonRequestMap()
		fmt.Println("JSON", requestMap)
	}

	for _, plan := range *result.Plans {
		fmt.Println(plan)
		depPlans, _ := plan.FindDependentPlans()
		fmt.Println("DependentPlans", depPlans)
	}
	fmt.Println("ShowExplainVerboseAnalyseSqlStatement", duration.Milliseconds(), " ms")
}

func runShowExplainVerboseSqlStatement(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	// string field
	result, err := gpudbInst.ShowExplainVerboseSqlStatement(context.TODO(), "select count(*) from demo.nyctaxi")
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(*result)
	fmt.Println("ShowExplainVerboseSqlStatement", duration.Milliseconds(), " ms")
}

func runShowStoredProcedureDDL(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	// string field
	result, err := gpudbInst.ShowStoredProcedureDDL(context.TODO(), "OPENSKY.opensky_processing")
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(*result)
	fmt.Println("ShowStoredProcedureDDL", duration.Milliseconds(), " ms")
}

func runShowTableDDL(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	// string field
	result, err := gpudbInst.ShowTableDDL(context.TODO(), "demo.nyctaxi_shard")
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(*result)
	fmt.Println("ShowTableDDL", duration.Milliseconds(), " ms")
}

func runShowTable1(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	//result, err := gpudbInst.ShowTableRaw("MASTER.nyctaxi")
	result, err := gpudbInst.ShowTableRawWithOpts(context.TODO(), "otel.trace_span", &gpudb.ShowTableOptions{
		ForceSynchronous:   true,
		GetSizes:           true,
		ShowChildren:       false, // needs to be false for tables
		NoErrorIfNotExists: false,
		GetColumnInfo:      true,
	})
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(result)
	fmt.Println("ShowTable", duration.Milliseconds(), " ms")
}

func runShowSystemStatus1(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	result, err := gpudbInst.ShowSystemStatusRaw(context.TODO())
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(*result)
	fmt.Println("ShowSystemStatusRaw", duration.Milliseconds(), " ms")
}

func runShowSystemProperties1(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	result, err := gpudbInst.ShowSystemPropertiesRaw(context.TODO())
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(*result)
	fmt.Println("ShowSystemPropertiesRaw", duration.Milliseconds(), " ms")
}

// Insert SQL test - Works
func runExecuteSql4(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	sql := `INSERT INTO otel.trace_span
	(id, resource_id, scope_id, event_id, link_id, trace_id, span_id, parent_span_id, trace_state, name, span_kind, start_time_unix_nano, end_time_unix_nano, dropped_attributes_count, dropped_events_count, dropped_links_count, message, status_code)
	VALUES('3f02a130-726f-4fda-a115-18a12e7c5884', 'ba3bd6c7-afdf-43aa-9287-375b956720db', 'b4006169-85d3-4874-bcff-49f2af13d9a0', '6d820f15-99f2-4e05-803c-e80c8be8748a', '56e36b10-a184-4fc1-8aa6-d3105e954972', '5b8aa5a2d2c872e8321cf37308d69df2', '5fb397be34d26b51', '051581bf3cb55c13', '', 'Hello-Greetings', 0, 1651258378114000000, 1651272778114000000, 0, 0, 0, '', 0)`
	result, err := gpudbInst.ExecuteSqlRaw(context.TODO(), sql, 0, 0, "", nil)
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println("Records inserted =", result.CountAffected)
	fmt.Println("ExecuteSqlRaw -Insert", duration.Milliseconds(), " ms")
}

func runExecuteSql5(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	uuid := uuid.New().String()
	fmt.Println(uuid)
	fmt.Println(len(uuid))

	result, err := gpudb.ExecuteSqlStructWithOpts[example.SampleExecuteSql](context.TODO(), *gpudbInst, "select trace_id from ki_home.sample", 0, 10, gpudb.NewDefaultExecuteSqlOptions())
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Printf("%+v\n", result)

	fmt.Println("ExecuteSqlStruct", duration.Milliseconds(), " ms")
}

// This fails
func runExecuteSql3(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	uuid := uuid.New().String()
	fmt.Println(uuid)
	fmt.Println(len(uuid))

	result, err := gpudbInst.ExecuteSqlStruct(context.TODO(), "select max(trace_id) as max_trace_id from ki_home.sample", 0, 10, func() interface{} { return example.SampleExecuteSql{} })
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Printf("%+v\n", *result.ResultsStruct)

	for i := 0; i < len(*result.ResultsStruct); i++ {
		span := (*result.ResultsStruct)[i].(example.SampleExecuteSql)
		fmt.Println(span)
	}
	fmt.Println("ExecuteSqlStruct", duration.Milliseconds(), " ms")
}

func runExecuteSql2(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	result, err := gpudbInst.ExecuteSqlMap(context.TODO(), "select max(trace_id) as max_trace_id from ki_home.sample", 0, 10)
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(*result.ResultsMap)
	fmt.Println("ExecuteSqlMap", duration.Milliseconds(), " ms")
}

// This works
func runExecuteSql1(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	result, err := gpudbInst.ExecuteSqlRaw(context.TODO(), "select max(trace_id) as max_trace_id from ki_home.sample", 0, 10, "", nil)
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(result)
	fmt.Println("ExecuteSqlRaw", duration.Milliseconds(), " ms")
}

func runGetRecords1(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	result, err := gpudbInst.GetRecordsStruct(context.TODO(), "ki_home.sample", 0, 10, func() interface{} { return example.SampleGetRecords{} })
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(result)
	fmt.Println("GetRecordsStruct", duration.Milliseconds(), " ms")
}

func runGetRecords2(gpudbInst *gpudb.Gpudb) {
	messages := make(chan interface{}, 1000)

	start := time.Now()
	gpudbInst.GetRecordsStructSendChannel(context.TODO(), "MASTER.nyctaxi", 0, 1000, messages, func() interface{} { return example.NycTaxi{} })
	for elem := range messages {
		noop(elem)
	}
	duration := time.Since(start)
	fmt.Println("GetRecordsStruct", duration.Milliseconds(), " ms")
}

func runGetRecords3(gpudbInst *gpudb.Gpudb) {

	start := time.Now()
	response, err := gpudbInst.GetRecordsRaw(context.TODO(), "otel.trace_span", 0, 50)
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println("GetRecordsRaw", duration.Milliseconds(), " ms")
	fmt.Println(response)
}

func runGetRecords4(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	result, err := gpudb.GetRecordsStructWithOpts[example.SampleGetRecords](context.TODO(), *gpudbInst, "ki_home.sample", 0, 10, gpudb.NewDefaultGetRecordsOptions())
	if err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println(result)
	fmt.Println("GetRecordsStruct", duration.Milliseconds(), " ms")
}

// ChunkBySize - Splits a slice into multiple slices of the given size
//
//	@param items
//	@param chunkSize
//	@return [][]T
func ChunkBySize[T any](items []T, chunkSize int) [][]T {
	var _chunks = make([][]T, 0, (len(items)/chunkSize)+1)
	for chunkSize < len(items) {
		items, _chunks = items[chunkSize:], append(_chunks, items[0:chunkSize:chunkSize])
	}
	return append(_chunks, items)
}

func runInsertRecords(logger *zap.Logger, gpudbInst *gpudb.Gpudb) {
	in, err := os.Open("./example/trace_span.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	spans := []example.Span{}

	if err := gocsv.UnmarshalFile(in, &spans); err != nil {
		panic(err)
	}

	// conversion from []example.Span to []interface{} is necessary
	// so that InsertRecords can deal with any structure and not
	// just a specific one
	spanRecords := make([]interface{}, len(spans))
	spanAttributeRecords := make([]interface{}, len(spans))

	for i := 0; i < len(spans); i++ {
		spanRecords[i] = spans[i]
		attribValue := example.AttributeValue{
			IntValue:    i*100 + 1,
			StringValue: "",
			BoolValue:   0,
			DoubleValue: 0,
			BytesValue:  []byte{},
		}
		spanAttributeRecords[i] = example.SpanAttribute{SpanID: spans[i].ID, Key: fmt.Sprintf("Key%d", i), AttributeValue: attribValue}
	}

	logger.Info("Span records : ", zap.Int("Count", len(spanRecords)))
	logger.Info("Span attribute records : ", zap.Int("Count", len(spanAttributeRecords)))

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(data []any, logger *zap.Logger, wg *sync.WaitGroup) {
		logger.Info("Inserting Span records")
		err = doChunkedInsert(context.TODO(), logger, gpudbInst, "otel.trace_span", data)
		if err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}(spanRecords, logger, wg)

	wg.Add(1)
	go func(data []any, logger *zap.Logger, wg *sync.WaitGroup) {
		logger.Info("Inserting Span attribute records")
		err = doChunkedInsert(context.TODO(), logger, gpudbInst, "otel.trace_span_attribute", data)
		if err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}(spanAttributeRecords, logger, wg)

	wg.Wait()
	// response, err := gpudbInst.InsertRecordsRaw(context.TODO(), "otel.trace_span_attribute", spanAttributeRecords)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("All done")
}

func runCreateJob(gpudbInst *gpudb.Gpudb) {
	start := time.Now()
	// result, err := gpudbInst.CreateJobRaw(
	// 	"OPENSKY")
	// if err != nil {
	// 	panic(err)
	// }
	duration := time.Since(start)
	// fmt.Println(result)
	fmt.Println("ShowSchema", duration.Milliseconds(), " ms")
}

func noop(elem interface{}) {

}

func doChunkedInsert(ctx context.Context, logger *zap.Logger, gpudbInst *gpudb.Gpudb, tableName string, records []any) error {

	logger.Info("Writing to - ", zap.String("Table", tableName), zap.Int("Record count", len(records)))

	recordChunks := ChunkBySize(records, 10000)

	errsChan := make(chan error, len(recordChunks))
	respChan := make(chan int, len(recordChunks))

	wg := &sync.WaitGroup{}
	// var mutex = &sync.Mutex{}

	for _, recordChunk := range recordChunks {
		wg.Add(1)
		go func(data []any, wg *sync.WaitGroup) {

			// mutex.Lock()
			resp, err := gpudbInst.InsertRecordsRaw(context.TODO(), tableName, data)
			errsChan <- err
			respChan <- resp.CountInserted
			// mutex.Unlock()

			wg.Done()
		}(recordChunk, wg)
	}
	wg.Wait()
	close(errsChan)
	close(respChan)
	var errs error
	for err := range errsChan {
		errs = multierr.Append(errs, err)
	}

	for resp := range respChan {
		fmt.Println("Count Inserted = ", resp)
	}
	return errs
}
