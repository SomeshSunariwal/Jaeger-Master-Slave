package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	DefaultComponentName = "tracer-demo"
)

func Init(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("Error : Cannot Init jaeger : %v\n"))
	}
	return tracer, closer
}

func main() {
	tracer, closer := Init(DefaultComponentName)
	defer closer.Close()

	e := echo.New()
	e.GET("/publish", func(context echo.Context) error {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(context.Request().Header))
		serverSpan := tracer.StartSpan("server", ext.RPCServerOption(spanCtx))
		serverSpan.Finish()
		serverSpan.LogFields(
			log.String("event", "string-format"),
			log.String("value", "helloStr"),
		)
		return context.String(http.StatusOK, "Hello From master")
	})

	e.Logger.Fatal(e.Start(":8082"))
}
