package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

func main() {
	logrus.Info("HTTP server started")
	cfg, err := utils.Read()

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

	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	router := gin.Default()

	sentenceRepo := _sentenceRepo.NewpostgresqlSentenceRepoistory(db)
	sentenceUsecase := _sentenceUsecase.NewSentenceUsecase(sentenceRepo)
	_sentenceHandlerHttpDelivery.NewSentenceHandler(router, sentenceUsecase)

	srvStart(router, *cfg)
}

func srvStart(router *gin.Engine, cfg utils.Config) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
