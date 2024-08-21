package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderapp"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderstore/orderdb"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productapp"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productstore/productdb"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userapp"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userstore/userdb"
	"github.com/nhannguyenacademy/ecommerce/pkg/keystore"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	apiV1GroupName = "api/v1"
)

var build = "develop"

type config struct {
	conf.Version
	Server struct {
		ReadTimeout        time.Duration `conf:"default:5s"`
		WriteTimeout       time.Duration `conf:"default:10s"`
		IdleTimeout        time.Duration `conf:"default:120s"`
		ShutdownTimeout    time.Duration `conf:"default:20s"`
		Host               string        `conf:"default:0.0.0.0:8080"`
		CORSAllowedOrigins []string      `conf:"default:*"`
	}
	Auth struct {
		KeysFolder string `conf:"default:configs/keys/"`
		ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
		Issuer     string `conf:"default:Nhan Nguyen Academy"`
	}
	DB struct {
		User            string        `conf:"default:postgres"`
		Password        string        `conf:"default:postgres,mask"`
		Host            string        `conf:"default:database"`
		Name            string        `conf:"default:postgres"`
		MaxOpenConns    int           `conf:"default:25"`
		MaxIdleConns    int           `conf:"default:25"`
		MaxConnLifeTime time.Duration `conf:"default:5m"`
		DisableTLS      bool          `conf:"default:true"`
	}
	Tempo struct {
		Host        string  `conf:"default:tempo:4317"`
		ServiceName string  `conf:"default:ecommerce"`
		Probability float64 `conf:"default:0.05"`
		// Shouldn't use a high Probability value in non-developer systems.
		// 0.05 should be enough for most systems. Some might want to have
		// this even lower.
	}
}

func main() {
	// todo: implement this
	traceIDFn := func(ctx context.Context) string {
		return ""
	}

	log := logger.New(os.Stdout, logger.LevelInfo, "ecommerce", traceIDFn)

	// -------------------------------------------------------------------------

	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// Configuration

	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "Ecommerce",
		},
	}

	const prefix = "ECOMMERCE"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// app Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.Host)

	db, err := sqldb.Open(sqldb.Config{
		User:            cfg.DB.User,
		Password:        cfg.DB.Password,
		Host:            cfg.DB.Host,
		Name:            cfg.DB.Name,
		MaxIdleConns:    cfg.DB.MaxIdleConns,
		MaxOpenConns:    cfg.DB.MaxOpenConns,
		MaxConnLifeTime: cfg.DB.MaxConnLifeTime,
		DisableTLS:      cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	defer db.Close()

	// -------------------------------------------------------------------------
	// Init auth
	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	ath, err := auth.New(auth.Config{Log: log, DB: db, KeyLookup: ks, Issuer: "ecommerce"})
	if err != nil {
		return fmt.Errorf("creating auth: %w", err)
	}

	// -------------------------------------------------------------------------
	// Init businesses

	userBus := userbus.NewBusiness(log, userdb.NewStore(log, db))

	productBus := productbus.NewBusiness(log, productdb.NewStore(log, db))

	orderBus := orderbus.NewBusiness(log, orderdb.NewStore(log, db))

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	// Shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Setup routes
	// todo: cors, tracing, metrics
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	apiV1Router := ginEngine.Group(apiV1GroupName)
	apiV1Router.Use(mid.Logging(log, []string{}), mid.Panic(log))
	userapp.New(log, ath, cfg.Auth.ActiveKID, userBus).Routes(apiV1Router)
	productapp.New(log, ath, productBus).Routes(apiV1Router)
	orderapp.New(log, ath, sqldb.NewBeginner(db), orderBus, productBus, userBus).Routes(apiV1Router)

	// Construct API server
	api := http.Server{
		Addr:         cfg.Server.Host,
		Handler:      ginEngine.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
