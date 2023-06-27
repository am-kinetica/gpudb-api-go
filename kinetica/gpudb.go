package kinetica

import (
	"bytes"
	"context"
	"crypto/tls"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
	"github.com/golang/snappy"
	"github.com/hamba/avro"
	"github.com/ztrue/tracerr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Kinetica struct {
	url     string
	options *KineticaOptions
	client  *resty.Client
	tracer  trace.Tracer
	mutex   sync.Mutex
}

type KineticaOptions struct {
	Username           string
	Password           string
	UseSnappy          bool
	ByPassSslCertCheck bool
	TraceHTTP          bool
	Timeout            string
}

var (
	//go:embed avsc/*
	schemaFS embed.FS
)

func New(ctx context.Context, url string) *Kinetica {
	return NewWithOptions(ctx, url, &KineticaOptions{})
}

func NewWithOptions(ctx context.Context, url string, options *KineticaOptions) *Kinetica {
	var (
		// childCtx context.Context
		childSpan trace.Span
	)

	client := resty.New()
	tracer := otel.GetTracerProvider().Tracer("kinetica-golang-api")
	client.DisableWarn = true

	_, childSpan = tracer.Start(ctx, "NewWithOptions()")
	defer childSpan.End()

	if options.Timeout != "" {
		duration, err := time.ParseDuration(options.Timeout)
		if err != nil {
			fmt.Println("Error parsing timout, no timeout set.", err)
		} else {
			client.SetTimeout(duration)
		}
	}

	return &Kinetica{url: url, options: options, client: client, tracer: tracer, mutex: sync.Mutex{}}
}

func parseSchema(asset string) *avro.Schema {
	var (
		err          error
		schemaByte   []byte
		schemaString string
		schema       avro.Schema
	)

	schemaByte, err = schemaFS.ReadFile(asset)
	if err != nil {
		panic(err)
	}

	schemaString = strings.ReplaceAll(string(schemaByte), "\n", "")
	schema, err = avro.Parse(schemaString)
	if err != nil {
		panic(err)
	}

	return &schema
}

// GetOptions
//
//	@receiver kinetica
//	@return *GpudbOptions
func (kinetica *Kinetica) GetOptions() *KineticaOptions {
	return kinetica.options
}

func (kinetica *Kinetica) submitRawRequest(
	ctx context.Context,
	restUrl string, requestSchema *avro.Schema, responseSchema *avro.Schema,
	requestStruct interface{}, responseStruct interface{}) error {
	var (
		childCtx  context.Context
		childSpan trace.Span
	)

	childCtx, childSpan = kinetica.tracer.Start(ctx, "kinetica.submitRawRequest()")
	defer childSpan.End()

	requestBody, err := avro.Marshal(*requestSchema, requestStruct)
	if err != nil {
		err = tracerr.Wrap(err)
		childSpan.RecordError(err)
		childSpan.SetStatus(codes.Error, err.Error())

		return err
	}

	httpResponse, err := kinetica.buildHTTPRequest(childCtx, &requestBody).Post(kinetica.url + restUrl)
	if err != nil {
		err = tracerr.Wrap(err)
		childSpan.RecordError(err)
		childSpan.SetStatus(codes.Error, err.Error())

		return err
	}

	childSpan.AddEvent("Response Info:", trace.WithAttributes(
		attribute.Int("Status Code: ", httpResponse.StatusCode()),
		attribute.String("Proto: ", httpResponse.Proto()),
		attribute.String("Time: ", httpResponse.Time().String()),
		attribute.String("Received At: ", httpResponse.ReceivedAt().String()),
	))

	ti := httpResponse.Request.TraceInfo()
	childSpan.AddEvent("Request Trace Info:", trace.WithAttributes(
		attribute.String("DNSLookup: ", ti.DNSLookup.String()),
		attribute.String("ConnTime: ", ti.ConnTime.String()),
		attribute.String("TCPConnTime: ", ti.TCPConnTime.String()),
		attribute.String("TLSHandshake: ", ti.TLSHandshake.String()),
		attribute.String("ServerTime: ", ti.ServerTime.String()),
		attribute.String("ResponseTime: ", ti.ResponseTime.String()),
		attribute.String("TotalTime: ", ti.TotalTime.String()),
		attribute.Bool("IsConnReused: ", ti.IsConnReused),
		attribute.Bool("IsConnWasIdle: ", ti.IsConnWasIdle),
		attribute.String("ConnIdleTime: ", ti.ConnIdleTime.String()),
		attribute.Int("RequestAttempt: ", ti.RequestAttempt),
		// attribute.String("RemoteAddr: ", ti.RemoteAddr.String()),
	))

	body := httpResponse.Body()
	reader := avro.NewReader(bytes.NewBuffer(body), len(body))

	status := reader.ReadString()
	if reader.Error != nil {
		err = errors.New(reader.Error.Error())
		childSpan.RecordError(err)
		childSpan.SetStatus(codes.Error, err.Error())

		return err
	}

	message := reader.ReadString()
	if reader.Error != nil {
		err = errors.New(reader.Error.Error())
		childSpan.RecordError(err)
		childSpan.SetStatus(codes.Error, err.Error())

		return err
	}

	childSpan.AddEvent("Response Message:", trace.WithAttributes(
		attribute.String("Message: ", message),
	))

	if status == "ERROR" {
		err = errors.New(message)
		childSpan.RecordError(err)
		childSpan.SetStatus(codes.Error, err.Error())

		return err
	}

	reader.SkipString()
	reader.SkipInt()
	reader.ReadVal(*responseSchema, &responseStruct)
	if reader.Error != nil {
		err = errors.New(reader.Error.Error())
		childSpan.RecordError(err)
		childSpan.SetStatus(codes.Error, err.Error())

		return err
	}

	childSpan.AddEvent("Response Struct:", trace.WithAttributes(
		attribute.String("Spew: ", spew.Sdump(responseStruct)),
	))

	return nil
}

func (kinetica *Kinetica) buildHTTPRequest(ctx context.Context, requestBody *[]byte) *resty.Request {

	var (
	// childCtx context.Context
	// childSpan trace.Span
	)

	if kinetica.options.ByPassSslCertCheck {
		kinetica.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	request := kinetica.client.R().
		SetBasicAuth(kinetica.options.Username, kinetica.options.Password)

	if kinetica.options.UseSnappy {
		snappyRequestBody := snappy.Encode(nil, *requestBody)

		request = request.SetHeader("Content-type", "application/x-snappy").SetBody(snappyRequestBody)
	} else {
		request = request.SetHeader("Content-type", "application/octet-stream").SetBody(*requestBody)
	}

	if kinetica.options.TraceHTTP {
		request = request.EnableTrace()
	}

	return request
}
