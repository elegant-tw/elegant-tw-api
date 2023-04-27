package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_sentenceHandlerHttpDelivery "elegant-tw-api/sentence/delivery/http"
	_sentenceRepo "elegant-tw-api/sentence/repository/postgresql"
	_sentenceUsecase "elegant-tw-api/sentence/usecase"
	"elegant-tw-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	mredis "github.com/ulule/limiter/v3/drivers/store/redis"

	_ "github.com/lib/pq"
)

//go:embed db/migrations/*.sql
var fs embed.FS

func main() {

	cfg, err := utils.Read()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Connecting to database...")
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUsername, cfg.DBPassword, cfg.DBName,
		),
	)

	if err != nil {
		logrus.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Connected to database.")

	d, err := iofs.New(fs, "db/migrations")
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Starting migration.")
	m, err := migrate.NewWithSourceInstance("iofs", d,
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
		),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logrus.Info("Nothing to migrate.")
		} else {
			logrus.Fatal(err)
		}
	} else {
		logrus.Info("Migration completed.")
	}

	gin.SetMode(cfg.GinMode)

	router := gin.Default()

	utils.BuildCORSConfig(router, cfg)

	if cfg.RateLimitEnabled {
		logrus.Info("Rate limit is enable.")
		rateLimitInit(router, cfg)
	} else {
		logrus.Info("Rate limit is disabled.")
	}

	sentenceRepo := _sentenceRepo.NewpostgresqlSentenceRepoistory(db)
	sentenceUsecase := _sentenceUsecase.NewSentenceUsecase(sentenceRepo)
	_sentenceHandlerHttpDelivery.NewSentenceHandler(router, sentenceUsecase)

	logrus.Info("HTTP server started.")
	srvStart(router, *cfg)
}

func rateLimitInit(router *gin.Engine, cfg *utils.Config) {
	rate, err := limiter.NewRateFromFormatted(cfg.RateLimitFormatted)

	if err != nil {
		logrus.Fatal(err)
	}

	var store limiter.Store

	if cfg.RateLimitStoreRedis {
		logrus.Info("Rate limit stores in redis.")
		client := redisInit(cfg)
		mstore, err := mredis.NewStoreWithOptions(client, limiter.StoreOptions{
			Prefix: cfg.RateLimitRedisPrefix,
		})
		if err != nil {
			logrus.Panic(err)
		}
		store = mstore
	} else {
		logrus.Info("Rate limit stores in memory.")
		store = memory.NewStore()
	}
	instance := limiter.New(store, rate)
	middleware := mgin.NewMiddleware(instance)
	router.Use(middleware)
	logrus.Infof("Request limit is %d reqs / %s", rate.Limit, rate.Period)
}

func redisInit(cfg *utils.Config) *redis.Client {
	logrus.Info("Connecting redis...")
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RateLimitRedisHost + ":" + cfg.RateLimitRedisPort,
	})
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Redis message: %+v", pong)
	logrus.Infof("Connected redis.")
	return rdb
}

func srvStart(router *gin.Engine, cfg utils.Config) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    cfg.ServerAddr + ":" + cfg.ServerPort,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logrus.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown: ", err)
	}

	logrus.Println("Server exiting")
}
