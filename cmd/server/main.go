package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/JoeReid/apiutils/tracer"
	"github.com/JoeReid/buffassignment/api"
	"github.com/JoeReid/buffassignment/internal/config"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
)

func main() {
	if err := tracer.InitTracer("Buff API Service"); err != nil {
		tracer.UntracedLogf("failed to configure tracer: %e", err)
		os.Exit(1)
	}

	srv, err := genServer()
	if err != nil {
		tracer.UntracedLogf("failed to setup api server: %e", err)
		os.Exit(1)
	}

	if err := srv.ListenAndServe(); err != nil {
		tracer.UntracedLogf("ListenAndServe error: %e", err)
		os.Exit(1)
	}
}

func genServer() (*http.Server, error) {
	sp, ctx := opentracing.StartSpanFromContext(context.Background(), "configure api server")
	defer sp.Finish()

	tracer.Log(sp, "build versioned api enpoints")
	r, err := api.Versioned(ctx)
	if err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	walkFunc := func(
		method string,
		route string,
		handler http.Handler,
		middlewares ...func(http.Handler) http.Handler,
	) error {
		tracer.Logf(sp, "discover route %s %s", method, route)
		return nil
	}

	tracer.Log(sp, "walk the router to print endpoints")
	if err := chi.Walk(r, walkFunc); err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	tracer.Log(sp, "read server config from environment")
	serverConfig, err := config.ServerConfig()
	if err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	return &http.Server{
		Handler: r,

		Addr:         fmt.Sprintf("%s:%d", serverConfig.ServeIP, serverConfig.ServePort),
		WriteTimeout: serverConfig.ServerWiteTimeout,
		ReadTimeout:  serverConfig.ServerReadTimeout,
	}, nil
}
