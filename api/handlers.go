package api

import (
	"context"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/apiutils/jsoncodec"
	"github.com/JoeReid/apiutils/tracer"
	"github.com/JoeReid/apiutils/yamlcodec"
	"github.com/JoeReid/buffassignment/api/buff"
	"github.com/JoeReid/buffassignment/api/videostream"
	"github.com/JoeReid/buffassignment/internal/config"
	"github.com/JoeReid/buffassignment/internal/model/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httptracer"
	"github.com/opentracing/opentracing-go"
)

// Versioned builds the full versioned api for the buff service
// complete with internal db connections.
//
// The api is returned as a chi router, allowing for easy use as a
// subrouter, if desired.
func Versioned(ctx context.Context) (*chi.Mux, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "build versioned api")
	defer sp.Finish()

	r := chi.NewRouter()

	// Configure middleware
	tracer.Log(sp, "configuring router middleware")
	r.Use(
		httptracer.Tracer(
			opentracing.GlobalTracer(),
			httptracer.Config{
				ServiceName:    "Buff API Service",
				ServiceVersion: "unversioned",
			},
		),
		middleware.Logger, // TODO: replace with tracing?
		middleware.RedirectSlashes,
	)

	tracer.Log(sp, "building api v1")
	routerV1, err := v1(ctx)
	if err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	r.Mount("/v1", routerV1)
	return r, nil
}

func v1(ctx context.Context) (*chi.Mux, error) {
	// linter thinks we want to assign here
	// realy we only want shadowing, we just dont call anything with the context yet
	// nolint:ineffassign,staticcheck
	sp, ctx := opentracing.StartSpanFromContext(ctx, "build versioned api")
	defer sp.Finish()

	r := chi.NewRouter()

	// configure all the codec options
	tracer.Log(sp, "configuring codec selectors")
	codecSelector, err := apiutils.NewRequestSelector(
		apiutils.RegisterCodec(
			jsoncodec.New(), "json", "application/json"),

		apiutils.RegisterCodec(
			jsoncodec.New(jsoncodec.SetIndent("", "\t")),
			"json,pretty", "application/json,pretty"),

		apiutils.RegisterCodec(
			yamlcodec.New(), "yaml", "application/x-yaml"),
	)
	if err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	tracer.Log(sp, "getting database config")
	dc, err := config.DBConfig()
	if err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	// TODO: how do we shut this down?
	// Do we need to? can it just follow the lifecycle of the service?
	tracer.Log(sp, "building store instance")
	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
	)
	if err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	// video_stream endpoint
	tracer.Log(sp, "configuring /video_streams handlers")
	r.Method("GET", "/video_streams", apiutils.HandlerWithSelector(codecSelector, videostream.NewListHandler(store)))
	r.Method("GET", "/video_streams/{uuid}", apiutils.HandlerWithSelector(codecSelector, videostream.NewGetHandler(store)))
	r.Method("GET", "/video_streams/{uuid}/buffs", apiutils.HandlerWithSelector(codecSelector, buff.NewListForStreamHandler(store)))

	// buffs endpoint
	tracer.Log(sp, "configuring /buffs handlers")
	r.Method("GET", "/buffs", apiutils.HandlerWithSelector(codecSelector, buff.NewListHandler(store)))
	r.Method("GET", "/buffs/{uuid}", apiutils.HandlerWithSelector(codecSelector, buff.NewGetHandler(store)))

	return r, nil
}
