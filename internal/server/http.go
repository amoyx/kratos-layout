package server

import (
	gohttp "net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	v2 "github.com/go-kratos/kratos-layout/api/helloworld/v2"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Server,
	greeter *service.GreeterService,
	greeterV2 *service.GreeterServiceV2,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ErrorEncoder(CustomErrorEncoder),
		http.ResponseEncoder(CustomResponseEncoder),
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			validate.Validator(),
			ratelimit.Server(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	openApiHandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openApiHandler)
	v1.RegisterUserServiceHTTPServer(srv, greeter)
	v2.RegisterUserServiceHTTPServer(srv, greeterV2)
	return srv
}

// CustomErrorEncoder encodes the error to the HTTP response.
func CustomErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(gohttp.StatusOK)
		return
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	w.WriteHeader(gohttp.StatusOK)
	_, _ = w.Write(body)
}

// CustomResponseEncoder encodes the object to the HTTP response.
func CustomResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	resp := &v1.Response{
		Code:    200,
		Message: "success",
	}
	if m, ok := v.(proto.Message); ok {
		any, err := anypb.New(m)
		if err != nil {
			return err
		}
		resp.Data = any
	}

	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{"application", subtype}, "/")
}
