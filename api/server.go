package api

import (
	"context"
	"io"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/redis/go-redis/v9"

	"github.com/arvan/qoute/config"
	pkg "github.com/arvan/qoute/pkg/redisclient"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	jaegerClientConfig "github.com/uber/jaeger-client-go/config"
)

type appServer struct {
	cfg            *config.Config
	app            *fiber.App
	TracingCloser  io.Closer
	RedisClientPtr *redis.Ring
}

var (
	appSrv  *appServer
	ctx     = context.Background()
	zlogger = log.With().Str("service", "Arvan-Qoute").Logger()
)

type jaegerLog struct {
}

func (j jaegerLog) Error(msg string) {
	zlogger.Error().Msg(msg)
}

func (j jaegerLog) Infof(msg string, args ...interface{}) {
	zlogger.Info().Msgf(msg, args)
}

func (j jaegerLog) Debugf(msg string, args ...interface{}) {
	zlogger.Debug().Msgf(msg, args)
}

func NewAppServer(cfg *config.Config) *appServer {

	appSrv = &appServer{
		cfg: cfg,
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	if cfg.Observability.Prometheus == true {
		prometheus := fiberprometheus.New("Arvan-Qoute-Service")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	if cfg.Observability.Jaeger == true {
		defcfg := jaegerClientConfig.Configuration{
			ServiceName: "Arvan-Qoute-Service",
			Reporter: &jaegerClientConfig.ReporterConfig{
				LocalAgentHostPort: cfg.Observability.Addr,
				LogSpans:           true,
			},
		}
		_, err := defcfg.FromEnv()
		if err != nil {
			zlogger.Info().Msgf("Could not parse Jaeger env vars: " + err.Error())
		}

		jlogger := jaegerLog{}
		closer, err := defcfg.InitGlobalTracer(
			"Courier",
			jaegerClientConfig.Logger(jlogger),
		)

		if err != nil {
			zlogger.Info().Msgf("Could not initialize jaeger tracer: " + err.Error())
		}

		appSrv.TracingCloser = closer
	}

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})
	api := app.Group("/api", middleware)
	v1 := api.Group("/v1", middleware)
	qouta := v1.Group("/qoute", middleware)

	v1.Get("/health", healthCheck)
	v1.Get("/info", info)
	qouta.Post("/add/:uuid/:username", add)

	appSrv.RedisClientPtr = pkg.Connect(ctx, cfg)

	if appSrv.RedisClientPtr == nil {
		zlogger.Info().Msgf("Could not connect to redis server ")
	}

	appSrv.app = app

	return appSrv
}

func (appSrv *appServer) ListenAndServe() chan error {
	errCh := make(chan error)
	go func() {
		zlogger.Info().Msgf("Started listening addr http://" + appSrv.cfg.Server.Addr)
		errCh <- appSrv.app.Listen(appSrv.cfg.Server.Addr)
	}()
	return errCh
}
