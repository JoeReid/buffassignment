package api

import (
	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/apiutils/jsoncodec"
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
func Versioned() (*chi.Mux, error) {
	r := chi.NewRouter()

	// Configure middleware
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

	routerV1, err := v1()
	if err != nil {
		return nil, err
	}

	r.Mount("/v1", routerV1)
	return r, nil
}

func v1() (*chi.Mux, error) {
	r := chi.NewRouter()

	// configure all the codec options
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
		return nil, err
	}

	// TODO: how do we shut this down?
	// Do we need to? can it just follow the lifecycle of the service?
	dc, err := config.DBConfig()
	if err != nil {
		return nil, err
	}

	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
	)
	if err != nil {
		return nil, err
	}

	// video_stream endpoint
	r.Method("GET", "/video_streams", apiutils.HandlerWithSelector(codecSelector, videostream.NewListHandler(store)))
	r.Method("GET", "/video_streams/{uuid}", apiutils.HandlerWithSelector(codecSelector, videostream.NewGetHandler(store)))
	r.Method("GET", "/video_streams/{uuid}/buffs", apiutils.HandlerWithSelector(codecSelector, buff.NewListForStreamHandler(store)))

	// buffs endpoint
	r.Method("GET", "/buffs", apiutils.HandlerWithSelector(codecSelector, buff.NewListHandler(store)))
	r.Method("GET", "/buffs/{uuid}", apiutils.HandlerWithSelector(codecSelector, buff.NewGetHandler(store)))

	return r, nil
}
