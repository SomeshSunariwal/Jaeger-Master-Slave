package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
	span := tracer.StartSpan("Say-Hello")
	span.SetTag("Hello-to", "Somesh")
	helloStr := fmt.Sprintf("Hello, %s", "Somesh")
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	span.LogKV("event", "Println")

	url := "http://localhost:8082/publish"
	req, _ := http.NewRequest("GET", url, nil)

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	e.GET("/", func(context echo.Context) error {
		resp, _ := http.DefaultClient.Do(req)
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println("http.StatusOK", string(data))
		span.LogFields(
			log.String("event", "string-format"),
			log.String("value", "helloStr"),
		)
		span.Finish()
		return context.String(http.StatusOK, "Done")
	})

	e.Logger.Fatal((e.Start(":8080")))
}
