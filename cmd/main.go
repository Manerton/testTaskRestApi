package main

import (
	"context"
	"log"
	"main/cmd/test"
	"main/internal/config"
	"main/internal/storage/postgresql"
	"main/internal/walletHandler"
	"main/internal/walletRepository"
	"main/internal/walletService"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

const configPath = "config.env"

func main() {
	// init config
	cfg := config.GetConfig(configPath)
	// init storage
	storage, err := postgresql.NewPosgreSQL(cfg.GetDataSourceName())
	if err != nil {
		log.Fatalf("faile to init storage")
	} else {
		log.Println("storage are enabled")
	}

	// init repository
	repository := walletRepository.NewWalletRepository()
	// init service
	service := walletService.NewWalletService(storage, repository)
	// init handler
	handler := walletHandler.NewWalletHandler(service)
	// init router
	router := chi.NewRouter()
	router.Get("/api/v1/wallets/{WALLET_UUID}", handler.GetWalletBalance)
	router.Post("/api/v1/wallets/", handler.CreateWallet)
	router.Put("/api/v1/wallet", handler.UpdateWallet)

	// init server
	server := &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: router,
	}
	// init gorutine for start server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("failed to start server")
		}
	}()
	log.Println("server start")
	test.StartTest()
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("failed to stop server")
		return
	}
	log.Println("server stopped")
}
