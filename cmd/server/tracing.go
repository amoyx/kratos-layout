package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/go-kratos/kratos-layout/internal/conf"
)

func initTracer(app *kratos.App, conf *conf.Tracing) {
	opts := make([]tracesdk.TracerProviderOption, 0, 2)
	kvs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(app.Name()),
		attribute.String("exporter", "jaeger"),
		attribute.Float64("float", 312.23),
	}

	//未配置收集器则不上传数据
	if conf != nil && conf.Endpoint != "" {
		exp, err := jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(conf.Endpoint),
			),
		)

		if err != nil {
			log.Error(err)
			return
		}
		opts = append(opts, tracesdk.WithBatcher(exp))
	} else {
		opts = append(opts, tracesdk.WithSampler(tracesdk.AlwaysSample()))
	}

	if conf != nil && conf.Token != "" {
		kvs = append(kvs, attribute.String("token", conf.Token))
	}

	opts = append(opts, tracesdk.WithResource(resource.NewSchemaless(kvs...)))
	tp := tracesdk.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
}
