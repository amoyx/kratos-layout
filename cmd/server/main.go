package main

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos-layout/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/contrib/config/apollo/v2"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software. use eureka kratos auto set eureka.metadata["Version"] = service.Version
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&Name, "name", "helloworld", "name of the application, eg: -name myapp")
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, rs registrySet) *kratos.App {
	var opts = []kratos.Option{
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	}

	for _, v := range rs {
		opts = append(opts, kratos.Registrar(v))
	}

	return kratos.New(opts...)
}

func configComplete() config.Config {
	lf := config.WithSource(
		file.NewSource(flagconf),
	)

	//1.load local file config
	lc := config.New(lf)
	if err := lc.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := lc.Scan(&bc); err != nil {
		lc.Close()
		panic(err)
	}

	if bc.Apollo == nil {
		return lc
	}

	//2.load apollo config
	sr := config.WithSource(apollo.NewSource(
		apollo.WithAppID(bc.Apollo.AppId),
		apollo.WithCluster(bc.Apollo.Cluster),
		apollo.WithEndpoint(bc.Apollo.Ip),
		apollo.WithEnableBackup(),
		apollo.WithNamespace(bc.Apollo.NamespaceName),
	))

	ac := config.New(sr)

	// If the Apollo configuration loading fails, the local configuration will be used directly.
	//In Kratos, when multiple configuration sources are provided, the data from the last source
	//will take precedence and override any conflicting values from previous sources.
	if err := ac.Load(); err != nil {
		lc.Close()
		panic(err)
	}
	return ac
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	c := configComplete()
	defer c.Close()

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Application, bc.Application.Server, bc.Application.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	initTracer(app, bc.Application.Tracing)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
